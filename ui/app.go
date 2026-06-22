package ui

import (
	"falcon/config"
	fyne "github.com/fyne-io/fyne/v2/app"
	container "github.com/fyne-io/fyne/v2/container"
	widget "github.com/fyne-io/fyne/v2/widget"
)

// App represents the main application
type App struct {
	FyneApp fyne.App
	Config  *config.Config
	MainWindow *widget.TabContainer
}

// NewApp creates a new application
func NewApp(cfg *config.Config) *App {
	app := &App{
		FyneApp: fyne.NewApp(),
		Config:  cfg,
	}

	return app
}

// Run starts the application
func (a *App) Run() {
	window := a.FyneApp.NewWindow()
	window.SetTitle("Falcon RDP Brute-Force Tool")
	window.Resize(fyne.NewSize(1200, 800))

	// Create tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Dashboard", a.createDashboardTab()),
		container.NewTabItem("Files", a.createFilesTab()),
		container.NewTabItem("Settings", a.createSettingsTab()),
		container.NewTabItem("Results", a.createResultsTab()),
	)

	window.SetContent(tabs)
	window.ShowAndRun()
}

func (a *App) createDashboardTab() interface{} {
	label := widget.NewLabel("Dashboard - Coming Soon")
	return label
}

func (a *App) createFilesTab() interface{} {
	label := widget.NewLabel("Files Manager - Coming Soon")
	return label
}

func (a *App) createSettingsTab() interface{} {
	label := widget.NewLabel("Settings - Coming Soon")
	return label
}

func (a *App) createResultsTab() interface{} {
	label := widget.NewLabel("Results - Coming Soon")
	return label
}
