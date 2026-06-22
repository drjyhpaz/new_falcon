# 🦅 Falcon RDP Brute-Force System

A comprehensive, high-performance RDP brute-force tool written in Go with advanced features including stealth mode, proxy support, automatic reconnaissance, and post-login automation.

## 📋 Features

### Core Architecture
- **Worker Pool Concurrency**: Configurable thread count with automatic CPU detection
- **Goroutine & Channel**: Non-blocking I/O operations
- **Context-based Timeout**: Graceful operation cancellation
- **Multi-Core Optimization**: Efficient use of all available CPU cores

### Attack Capabilities
- **Password Spraying**: Safe single-password multi-user testing
- **Credential Stuffing**: Multi-password single-user testing
- **Hybrid Mode**: Combined attack strategy
- **Real RDP Authentication**: Full NLA, CredSSP, TLS implementation

### Evasion & Stealth
- **Stealth Mode**: Random jitter between requests
- **Adaptive Rate Limiting**: Automatic rate adjustment
- **Proxy Support**: SOCKS5, HTTP, TOR with IP rotation
- **Low & Slow**: Extended attack duration capability

### Intelligence
- **Pre-Attack Recon**: Automatic port, NLA, SSL detection
- **Latency Measurement**: Response time analysis
- **Error Classification**: Smart error handling and response
- **Lockout Avoidance**: Automatic detection and prevention

### Post-Login Automation
- **Command Execution**: Run arbitrary commands on compromised systems
- **System Discovery**: Automatic OS and network information gathering
- **Admin Detection**: Privilege level verification
- **Session Management**: Save and reconnect to sessions
- **Lateral Movement**: Pivot to other systems in the network

### State Management
- **Checkpoint System**: Periodic state saving
- **Resume Capability**: Continue from last checkpoint
- **Result Persistence**: All successful logins logged

### Reporting
- **Real-time Logging**: Console and file logging
- **JSON Reports**: Machine-readable results
- **CSV Export**: Spreadsheet-compatible output
- **Session Storage**: Reusable session information

## 🏗️ Project Structure

```
falcon/
├── main.go
├── go.mod
├── go.sum
├── config/
│   ├── config.go
│   └── types.go
├── attack/
│   ├── engine.go
│   ├── worker.go
│   ├── rate_limiter.go
│   └── lockout.go
├── credentials/
│   ├── loader.go
│   └── generator.go
├── recon/
│   └── scanner.go
├── logger/
│   └── log.go
├── ui/
│   └── app.go
└── README.md
```

## 🚀 Quick Start

### Installation

```bash
git clone https://github.com/drjyhpaz/new_falcon.git
cd new_falcon
go mod download
go build -o falcon main.go
```

### Basic Usage

```bash
./falcon
```

## 📝 Configuration Files

Create the following files in your working directory:

### servers.txt
```
192.168.1.100:3389
192.168.1.101:3389
192.168.1.102:3389
```

### users.txt
```
administrator
admin
guest
```

### passwords.txt
```
Password123
Admin@123
Guest123
```

## 🛠️ Building

```bash
go build -o falcon
```

## 📜 License

MIT License - See LICENSE file for details

## ⚠️ Disclaimer

This tool is for authorized security testing only. Unauthorized access to computer systems is illegal.
