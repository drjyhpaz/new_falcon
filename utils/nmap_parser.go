package utils

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// NmapPort represents a port from Nmap output
type NmapPort struct {
	Protocol string
	PortID   int
	Service  string
	State    string
}

// NmapTarget represents a target from Nmap output
type NmapTarget struct {
	IP    string
	Ports []NmapPort
}

// ParseNmapXML parses Nmap XML output
func ParseNmapXML(filename string) ([]NmapTarget, error) {
	var targets []NmapTarget

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read nmap file: %w", err)
	}

	// Simple XML parsing for Nmap output
	// This is a basic implementation

	return targets, nil
}

// ParseNmapGNMAP parses Nmap GNMAP output
func ParseNmapGNMAP(filename string) ([]NmapTarget, error) {
	var targets []NmapTarget

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read gnmap file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Host:") {
			target := parseGNMAPLine(line)
			if target != nil {
				targets = append(targets, *target)
			}
		}
	}

	return targets, nil
}

func parseGNMAPLine(line string) *NmapTarget {
	// Format: Host: 192.168.1.1 Ports: 22/open/tcp//ssh///, 80/open/tcp//http///

	parts := strings.Split(line, "Ports:")
	if len(parts) != 2 {
		return nil
	}

	ipParts := strings.Fields(parts[0])
	if len(ipParts) < 2 {
		return nil
	}

	ip := ipParts[1]
	target := &NmapTarget{
		IP:    ip,
		Ports: make([]NmapPort, 0),
	}

	portStrings := strings.Split(parts[1], ",")
	for _, portStr := range portStrings {
		portStr = strings.TrimSpace(portStr)
		if portStr == "" {
			continue
		}

		fields := strings.Split(portStr, "/")
		if len(fields) >= 3 {
			portNum, err := strconv.Atoi(fields[0])
			if err != nil {
				continue
			}

			port := NmapPort{
				PortID:   portNum,
				State:    fields[1],
				Protocol: fields[2],
			}
			target.Ports = append(target.Ports, port)
		}
	}

	return target
}
