// Falcon RDP Brute-Force Tool
// Advanced multi-protocol brute-force attack framework
// Version: 1.0.0

package main

import (
	"fmt"
)

const (
	Version = "1.0.0"
	Author  = "Falcon Project"
	Github  = "https://github.com/drjyhpaz/new_falcon"
)

var Banner = fmt.Sprintf(`
╔═══════════════════════════════════════════════════════════════╗
║                                                               ║
║     🦅 FALCON - RDP Brute-Force Attack Framework 🦅          ║
║                                                               ║
║                     Version: %s                          ║
║               Author: %s                                ║
║          Repository: %s                                ║
║                                                               ║
║  Advanced Features:                                           ║
║  • Multi-protocol support (RDP, SSH, FTP)                    ║
║  • Password spraying & credential stuffing                   ║
║  • Proxy rotation & stealth mode                             ║
║  • Automatic reconnaissance                                  ║
║  • Post-login automation                                     ║
║  • Session management & persistence                          ║
║                                                               ║
║  ⚠️  WARNING: For authorized testing only!                   ║
║                                                               ║
╚═══════════════════════════════════════════════════════════════╝
`, Version, Author, Github)

func init() {
	fmt.Println(Banner)
}
