# 🦅 Falcon - Advanced RDP Brute-Force Tool

## 📋 Table of Contents

1. [Features](#-features)
2. [Installation](#-installation)
3. [Usage](#-usage)
4. [Configuration](#-configuration)
5. [Examples](#-examples)
6. [Architecture](#-architecture)
7. [API Reference](#-api-reference)
8. [Troubleshooting](#-troubleshooting)

## ✨ Features

### Core Capabilities
- **Multi-Protocol Support**: RDP, SSH, FTP, SMB
- **Worker Pool**: Configurable concurrent threads with auto CPU detection
- **Advanced Concurrency**: Goroutines, channels, and context-based timeout
- **High Performance**: Multi-core optimization and efficient resource management

### Attack Modes
- **Password Spraying**: Single password against multiple users (safe)
- **Credential Stuffing**: Multiple passwords against single user
- **Hybrid Mode**: Combined attack strategy
- **Real Authentication**: Full protocol stack implementation

### Evasion & Stealth
- **Stealth Mode**: Random jitter between requests
- **Adaptive Rate Limiting**: Auto-adjust based on errors
- **Low & Slow**: Extended duration capability
- **Proxy Support**: SOCKS5, HTTP, TOR with IP rotation

### Reconnaissance
- **Automatic Port Scanning**: Check open ports
- **Service Detection**: Identify service and version
- **NLA Detection**: Network Level Authentication check
- **SSL/TLS Detection**: Protocol variant detection
- **Latency Measurement**: Response time analysis

### Lockout Prevention
- **Failed Attempt Tracking**: Per-account failure counting
- **Automatic Cooldown**: Prevent account lockout
- **Error Classification**: Smart error handling
- **Strategy Adjustment**: Auto-switch to safer modes

### State Management
- **Checkpoint System**: Periodic state saving
- **Resume Support**: Continue from last checkpoint
- **Progress Tracking**: Real-time statistics
- **Result Persistence**: All results logged

### Post-Login Automation
- **Command Execution**: Run arbitrary commands
- **System Discovery**: Auto system info gathering
- **Admin Detection**: Check privilege level
- **Session Management**: Save and reuse sessions
- **Lateral Movement**: Pivot through network

### Reporting
- **Real-time Logging**: Console and file logs
- **JSON Reports**: Machine-readable output
- **CSV Export**: Spreadsheet-compatible
- **Session Storage**: Reusable session data

## 🚀 Installation

### Requirements
- Go 1.21+
- Linux/macOS/Windows

### Build from Source

```bash
git clone https://github.com/drjyhpaz/new_falcon.git
cd new_falcon
go mod download

# Simple build
go build -o falcon main.go

# Or use build script
chmod +x build.sh
./build.sh
```

## 💻 Usage

### CLI Mode

```bash
# Basic usage with default settings
./falcon --servers servers.txt --users users.txt --passwords passwords.txt

# With custom threads and stealth
./falcon --servers targets.txt --users users.txt --passwords pass.txt --threads 64 --stealth

# With proxy support
./falcon --servers targets.txt --users users.txt --passwords pass.txt --proxy --proxy-file proxies.txt

# Enable post-login automation
./falcon --servers targets.txt --users users.txt --passwords pass.txt --postlogin

# Resume from checkpoint
./falcon --resume

# Generate credentials
./falcon --generate --users users.txt --passwords passwords.txt
```

### GUI Mode

```bash
./falcon --ui
```

### Command Line Options

```
  -servers string
    	Path to servers.txt (default "servers.txt")
  -users string
    	Path to users.txt (default "users.txt")
  -passwords string
    	Path to passwords.txt (default "passwords.txt")
  -threads int
    	Number of threads (0 for auto) (default 0)
  -timeout int
    	Timeout in seconds (default 10)
  -stealth
    	Enable stealth mode
  -proxy
    	Enable proxy support
  -proxy-file string
    	Path to proxies.txt (default "proxies.txt")
  -resume
    	Resume from checkpoint
  -postlogin
    	Enable post-login automation
  -generate
    	Generate credentials from files
  -version
    	Show version
  -help
    	Show this help message
```

## 📁 Configuration

### Input Files

#### servers.txt
One target per line in `IP:Port` format:
```
192.168.1.100:3389
192.168.1.101:3389
10.0.0.50:22
```

#### users.txt
One username per line:
```
administrator
admin
guest
root
```

#### passwords.txt
One password per line:
```
Password123
Admin@123
Guest123
```

#### proxies.txt (optional)
Proxy format: `type://[user:pass@]host:port`
```
socks5://proxy.example.com:1080
http://user:pass@proxy.example.com:8080
tor://localhost:9050
```

## 📚 Examples

### Example 1: Basic RDP Brute-Force

```bash
# Create input files
echo "192.168.1.100:3389" > servers.txt
echo -e "admin\nadministrator\nguest" > users.txt
echo -e "Password123\nAdmin@123" > passwords.txt

# Run attack
./falcon --servers servers.txt --users users.txt --passwords passwords.txt --threads 32
```

### Example 2: Stealth Attack with Proxy

```bash
# Create proxy file
echo "socks5://proxy1.com:1080" > proxies.txt
echo "socks5://proxy2.com:1080" >> proxies.txt

# Run with stealth and proxy
./falcon --servers servers.txt --users users.txt --passwords passwords.txt \\
  --stealth --proxy --proxy-file proxies.txt --threads 16
```

### Example 3: SSH Attack with Post-Login

```bash
# Modify servers.txt for SSH
echo "10.0.0.1:22" > servers.txt

# Run with post-login commands
./falcon --servers servers.txt --users users.txt --passwords passwords.txt \\
  --postlogin --threads 20
```

## 🏗️ Architecture

### Directory Structure

```
falcon/
├── main.go                 # Entry point
├── go.mod/go.sum          # Dependencies
├── config/                # Configuration management
│   ├── config.go
│   └── types.go
├── attack/                # Core attack engine
│   ├── engine.go
│   ├── worker.go
│   ├── rate_limiter.go
│   ├── lockout.go
│   └── strategies.go
├── credentials/           # Credential handling
│   ├── loader.go
│   └── generator.go
├── recon/                 # Reconnaissance
│   └── scanner.go
├── auth/                  # Authentication
│   └── authenticator.go
├── evasion/               # Evasion techniques
│   └── stealth.go
├── proxy/                 # Proxy management
│   └── manager.go
├── state/                 # State management
│   └── manager.go
├── postlogin/             # Post-login automation
│   └── session.go
├── report/                # Report generation
│   └── generator.go
├── logger/                # Logging
│   └── log.go
├── utils/                 # Utilities
│   ├── helpers.go
│   ├── errors.go
│   └── nmap_parser.go
├── ui/                    # GUI components
│   ├── app.go
│   └── dashboard.go
├── constants/             # Constants
│   └── banner.go
├── README.md              # This file
├── LICENSE                # MIT License
└── .gitignore
```

### Component Overview

**AttackEngine**: Main orchestrator managing:
- Worker pool coordination
- Rate limiting
- Lockout detection
- Result collection

**WorkerPool**: Manages concurrent:
- Job distribution
- Worker threads
- Result aggregation

**RateLimiter**: Token bucket implementation:
- Packets per second control
- Adaptive adjustment

**LockoutManager**: Account safety:
- Failure tracking
- Cooldown periods
- Automatic prevention

**Stealth**: Evasion techniques:
- Random jitter
- Adaptive rates
- Low & slow modes

**ProxyRotator**: IP rotation:
- Proxy validation
- Round-robin rotation
- Multi-protocol support

## 🔌 API Reference

### AttackEngine

```go
// Create engine
engine := attack.NewAttackEngine(config, targets, credentials)

// Start attack
engine.Start()

// Stop attack
engine.Stop()

// Get results
results := engine.GetResults()
successful := engine.GetSuccessfulResults()
```

### Credentials

```go
// Load files
targets, _ := credentials.LoadServers("servers.txt")
users, _ := credentials.LoadUsers("users.txt")
passwords, _ := credentials.LoadPasswords("passwords.txt")

// Generate combinations
creds := credentials.GenerateCredentials(users, passwords, domain)

// Save to file
credentials.SaveCredentials(creds, "credentials.txt")
```

### State Management

```go
// Create state manager
sm := state.NewStateManager("state.json")

// Save checkpoint
sm.SaveCheckpoint(targets, results, stats)

// Load checkpoint
checkpoint, _ := sm.LoadCheckpoint()

// Check if checkpoint exists
if sm.HasCheckpoint() {
    // Resume from checkpoint
}
```

## 🐛 Troubleshooting

### Common Issues

**Q: Connection timeout errors**
A: Increase timeout with `--timeout 30` or check network connectivity

**Q: Out of memory errors**
A: Reduce thread count with `--threads 8` and enable streaming

**Q: Permission denied on files**
A: Ensure read permissions on input files: `chmod 644 servers.txt`

**Q: Proxy not working**
A: Validate proxy format and connectivity: `curl -x socks5://proxy:1080 http://example.com`

## 📝 License

MIT License - See LICENSE file for details

## ⚠️ Disclaimer

**This tool is for authorized security testing only.** Unauthorized access to computer systems is illegal. Ensure you have written permission before using this tool.

---

**Report issues**: https://github.com/drjyhpaz/new_falcon/issues
