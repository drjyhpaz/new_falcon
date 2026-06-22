package ui

import (
	"falcon/attack"
	"falcon/config"
	"falcon/credentials"
	"falcon/logger"
	"fmt"
	"image/color"
	"sync"

	"github.com/fyne-io/fyne/v2"
	app "github.com/fyne-io/fyne/v2/app"
	canvas "github.com/fyne-io/fyne/v2/canvas"
	container "github.com/fyne-io/fyne/v2/container"
	layout "github.com/fyne-io/fyne/v2/layout"
	widget "github.com/fyne-io/fyne/v2/widget"
)

// Dashboard represents the main dashboard
type Dashboard struct {
	App           app.App
	Window        fyne.Window
	Config        *config.Config
	Engine        *attack.AttackEngine
	IsRunning     bool
	Mutex         sync.RWMutex
	StatusLabel   *widget.Label
	StatsLabel    *widget.Label
	ProgressBar   *widget.ProgressBar
	LogsContainer *container.AppScroll
	ResultsTable  *widget.Table
}

// NewDashboard creates a new dashboard
func NewDashboard(cfg *config.Config) *Dashboard {
	return &Dashboard{
		App:    app.NewApp(),
		Config: cfg,
	}
}

// BuildUI builds the UI
func (d *Dashboard) BuildUI() {
	d.Window = d.App.NewWindow()
	d.Window.SetTitle("🦅 Falcon RDP Brute-Force Tool v1.0")
	d.Window.Resize(fyne.NewSize(1400, 900))

	// Create tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("🎯 Dashboard", d.buildDashboardTab()),
		container.NewTabItem("📁 Files", d.buildFilesTab()),
		container.NewTabItem("⚙️ Settings", d.buildSettingsTab()),
		container.NewTabItem("📊 Results", d.buildResultsTab()),
	)

	d.Window.SetContent(tabs)
}

// buildDashboardTab builds the dashboard tab
func (d *Dashboard) buildDashboardTab() *fyne.Container {
	// Title
	title := canvas.NewText("Attack Dashboard", color.White)
	title.TextSize = 24
	title.TextStyle.Bold = true

	// Status
	d.StatusLabel = widget.NewLabel("Status: Ready")
	d.StatsLabel = widget.NewLabel("Total: 0 | Success: 0 | Failed: 0 | PPS: 0")
	d.ProgressBar = widget.NewProgressBar()

	// Control buttons
	startBtn := widget.NewButton("▶️ Start Attack", func() {
		d.startAttack()
	})
	startBtn.Importance = widget.HighImportance

	stopBtn := widget.NewButton("⏹️ Stop Attack", func() {
		d.stopAttack()
	})
	stopBtn.Importance = widget.DangerImportance

	controlBox := container.NewHBox(
		startBtn,
		stopBtn,
	)

	// Results table
	d.ResultsTable = widget.NewTable(
		func() (int, int) {
			if d.Engine == nil {
				return 0, 5
			}
			return len(d.Engine.GetSuccessfulResults()), 5
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if d.Engine == nil {
				return
			}

			results := d.Engine.GetSuccessfulResults()
			if id.Row < len(results) {
				result := results[id.Row]
				switch id.Col {
				case 0:
					label.SetText(result.IP)
				case 1:
					label.SetText(fmt.Sprintf("%d", result.Port))
				case 2:
					label.SetText(result.Username)
				case 3:
					label.SetText(result.Password)
				case 4:
					label.SetText(result.Domain)
				}
			}
		},
	)

	d.ResultsTable.SetColumnWidth(0, 120)
	d.ResultsTable.SetColumnWidth(1, 80)
	d.ResultsTable.SetColumnWidth(2, 120)
	d.ResultsTable.SetColumnWidth(3, 150)
	d.ResultsTable.SetColumnWidth(4, 100)

	// Logs
	d.LogsContainer = container.NewVScroll(
		widget.NewLabel("Logs will appear here..."),
	)

	// Layout
	return container.New(
		layout.NewVBoxLayout(),
		title,
		d.StatusLabel,
		d.StatsLabel,
		d.ProgressBar,
		controlBox,
		widget.NewCard("Successful Logins", "", d.ResultsTable, nil),
		widget.NewCard("Activity Log", "", d.LogsContainer, nil),
	)
}

// buildFilesTab builds the files tab
func (d *Dashboard) buildFilesTab() *fyne.Container {
	title := canvas.NewText("File Management", color.White)
	title.TextSize = 24
	title.TextStyle.Bold = true

	serversLabel := widget.NewLabel("servers.txt:")
	serversEntry := widget.NewEntry()
	serversEntry.PlaceHolder = "/path/to/servers.txt"

	usersLabel := widget.NewLabel("users.txt:")
	usersEntry := widget.NewEntry()
	usersEntry.PlaceHolder = "/path/to/users.txt"

	passwordsLabel := widget.NewLabel("passwords.txt:")
	passwordsEntry := widget.NewEntry()
	passwordsEntry.PlaceHolder = "/path/to/passwords.txt"

	generateBtn := widget.NewButton("🔗 Generate Credentials", func() {
		d.generateCredentials(serversEntry.Text, usersEntry.Text, passwordsEntry.Text)
	})
	generateBtn.Importance = widget.HighImportance

	return container.New(
		layout.NewVBoxLayout(),
		title,
		widget.NewCard("Servers", "", container.New(layout.NewVBoxLayout(), serversLabel, serversEntry), nil),
		widget.NewCard("Users", "", container.New(layout.NewVBoxLayout(), usersLabel, usersEntry), nil),
		widget.NewCard("Passwords", "", container.New(layout.NewVBoxLayout(), passwordsLabel, passwordsEntry), nil),
		generateBtn,
	)
}

// buildSettingsTab builds the settings tab
func (d *Dashboard) buildSettingsTab() *fyne.Container {
	title := canvas.NewText("Settings", color.White)
	title.TextSize = 24
	title.TextStyle.Bold = true

	// Thread count
	threadsLabel := widget.NewLabel(fmt.Sprintf("Threads: %d", d.Config.Attack.Threads))
	threadsSlider := widget.NewSlider(1, 128)
	threadsSlider.Value = float64(d.Config.Attack.Threads)
	threadsSlider.OnChanged = func(v float64) {
		d.Config.Attack.Threads = int(v)
		threadsLabel.SetText(fmt.Sprintf("Threads: %d", int(v)))
	}

	// Timeout
	timeoutLabel := widget.NewLabel(fmt.Sprintf("Timeout: %dms", d.Config.Attack.Timeout.Milliseconds()))
	timeoutSlider := widget.NewSlider(1000, 30000)
	timeoutSlider.Value = float64(d.Config.Attack.Timeout.Milliseconds())
	timeoutSlider.OnChanged = func(v float64) {
		d.Config.Attack.Timeout = fyne.Duration(int(v)) // Placeholder
		timeoutLabel.SetText(fmt.Sprintf("Timeout: %dms", int(v)))
	}

	// Stealth mode
	stealthCheck := widget.NewCheck("Enable Stealth Mode", func(b bool) {
		d.Config.Attack.StealthMode = b
	})
	stealthCheck.Checked = d.Config.Attack.StealthMode

	// Proxy
	proxyCheck := widget.NewCheck("Enable Proxy", func(b bool) {
		d.Config.Attack.ProxyEnabled = b
	})
	proxyCheck.Checked = d.Config.Attack.ProxyEnabled

	return container.New(
		layout.NewVBoxLayout(),
		title,
		widget.NewCard("Performance", "", container.New(layout.NewVBoxLayout(), threadsLabel, threadsSlider, timeoutLabel, timeoutSlider), nil),
		widget.NewCard("Evasion", "", container.New(layout.NewVBoxLayout(), stealthCheck, proxyCheck), nil),
	)
}

// buildResultsTab builds the results tab
func (d *Dashboard) buildResultsTab() *fyne.Container {
	title := canvas.NewText("Results & Reports", color.White)
	title.TextSize = 24
	title.TextStyle.Bold = true

	exportJSONBtn := widget.NewButton("📥 Export JSON", func() {
		logger.Info("Exporting JSON...")
	})

	exportCSVBtn := widget.NewButton("📥 Export CSV", func() {
		logger.Info("Exporting CSV...")
	})

	return container.New(
		layout.NewVBoxLayout(),
		title,
		container.NewHBox(exportJSONBtn, exportCSVBtn),
	)
}

// startAttack starts the attack
func (d *Dashboard) startAttack() {
	d.Mutex.Lock()
	if d.IsRunning {
		d.Mutex.Unlock()
		return
	}
	d.IsRunning = true
	d.Mutex.Unlock()

	d.StatusLabel.SetText("Status: Running...")
	logger.Info("Attack started from UI")
}

// stopAttack stops the attack
func (d *Dashboard) stopAttack() {
	d.Mutex.Lock()
	if !d.IsRunning {
		d.Mutex.Unlock()
		return
	}
	d.IsRunning = false
	d.Mutex.Unlock()

	d.StatusLabel.SetText("Status: Stopped")
	logger.Info("Attack stopped from UI")
}

// generateCredentials generates credentials
func (d *Dashboard) generateCredentials(serversFile, usersFile, passwordsFile string) {
	if serversFile == "" || usersFile == "" || passwordsFile == "" {
		logger.Error("Please provide all file paths")
		return
	}

	// Load files
	users, err := credentials.LoadUsers(usersFile)
	if err != nil {
		logger.Error("Failed to load users: %v", err)
		return
	}

	passwords, err := credentials.LoadPasswords(passwordsFile)
	if err != nil {
		logger.Error("Failed to load passwords: %v", err)
		return
	}

	// Generate
	creds := credentials.GenerateCredentials(users, passwords, d.Config.Attack.DefaultDomain)

	// Save
	err = credentials.SaveCredentials(creds, "credentials.txt")
	if err != nil {
		logger.Error("Failed to save credentials: %v", err)
		return
	}

	logger.Success("Generated %d credentials", len(creds))
}

// Run starts the application
func (d *Dashboard) Run() {
	d.BuildUI()
	d.Window.ShowAndRun()
}
