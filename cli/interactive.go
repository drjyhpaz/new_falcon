package main

import (
	"bufio"
	"falcon/logger"
	"fmt"
	"os"
	"strings"
)

// InteractiveMode starts an interactive CLI
func InteractiveMode() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n" + Banner)
	fmt.Println("\n📝 Falcon Interactive Mode")
	fmt.Println("Type 'help' for available commands\n")

	for {
		fmt.Print("falcon> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		if !processCommand(input) {
			break
		}
	}
}

func processCommand(input string) bool {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return true
	}

	cmd := parts[0]

	switch cmd {
	case "help":
		printHelp()
	case "version":
		fmt.Printf("Falcon v%s\n", Version)
	case "exit", "quit":
		fmt.Println("👋 Goodbye!")
		return false
	case "show":
		if len(parts) > 1 {
			showCommand(parts[1])
		}
	case "load":
		if len(parts) > 2 {
			loadCommand(parts[1], parts[2])
		}
	case "set":
		if len(parts) > 2 {
			setCommand(parts[1], parts[2])
		}
	case "run":
		runCommand()
	case "clear":
		clearScreen()
	default:
		fmt.Println("❌ Unknown command. Type 'help' for more info.")
	}

	return true
}

func printHelp() {
	fmt.Println(`
╔══════════════════════════════════════════════════════════════════════════════╗
║                         FALCON COMMAND REFERENCE                             ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  Configuration Commands:                                                    ║
║  ├─ load servers <file>      Load targets from file                        ║
║  ├─ load users <file>        Load usernames from file                      ║
║  ├─ load passwords <file>    Load passwords from file                      ║
║  ├─ set threads <count>      Set thread count                              ║
║  ├─ set timeout <seconds>    Set connection timeout                        ║
║  ├─ set stealth <on|off>     Enable/disable stealth mode                   ║
║  ├─ set proxy <on|off>       Enable/disable proxy                          ║
║  └─ set pps <rate>           Set packets per second                        ║
║                                                                              ║
║  Information Commands:                                                      ║
║  ├─ show config              Display current configuration                  ║
║  ├─ show servers             Show loaded targets                            ║
║  ├─ show users               Show loaded usernames                          ║
║  ├─ show passwords           Show loaded passwords                          ║
║  └─ show stats               Display statistics                             ║
║                                                                              ║
║  Execution Commands:                                                        ║
║  ├─ run                      Start the attack                               ║
║  ├─ clear                    Clear screen                                   ║
║  ├─ version                  Show version                                   ║
║  ├─ help                     Show this help                                 ║
║  └─ exit / quit              Exit Falcon                                    ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
`)
}

func showCommand(target string) {
	switch target {
	case "config":
		fmt.Println("📋 Current Configuration:")
		fmt.Println("  [Not implemented yet]")
	case "servers":
		fmt.Println("🎯 Loaded Servers:")
		fmt.Println("  [No servers loaded]")
	case "users":
		fmt.Println("👥 Loaded Users:")
		fmt.Println("  [No users loaded]")
	case "passwords":
		fmt.Println("🔐 Loaded Passwords:")
		fmt.Println("  [No passwords loaded]")
	case "stats":
		fmt.Println("📊 Statistics:")
		fmt.Println("  Total Attempts: 0")
		fmt.Println("  Successful: 0")
		fmt.Println("  Failed: 0")
	default:
		fmt.Printf("❌ Unknown target: %s\n", target)
	}
}

func loadCommand(filetype, path string) {
	switch filetype {
	case "servers":
		logger.Info("Loading servers from %s", path)
	case "users":
		logger.Info("Loading users from %s", path)
	case "passwords":
		logger.Info("Loading passwords from %s", path)
	default:
		fmt.Printf("❌ Unknown file type: %s\n", filetype)
	}
}

func setCommand(key, value string) {
	switch key {
	case "threads":
		fmt.Printf("✅ Threads set to: %s\n", value)
	case "timeout":
		fmt.Printf("✅ Timeout set to: %s seconds\n", value)
	case "stealth":
		fmt.Printf("✅ Stealth mode: %s\n", value)
	case "proxy":
		fmt.Printf("✅ Proxy: %s\n", value)
	case "pps":
		fmt.Printf("✅ PPS set to: %s\n", value)
	default:
		fmt.Printf("❌ Unknown setting: %s\n", key)
	}
}

func runCommand() {
	logger.Info("Starting attack...")
	fmt.Println("🚀 Attack in progress...")
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}
