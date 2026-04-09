package scanner

import (
	"fmt"
	"net"
	"time"
)

// Port represents a network port with its state
type Port struct {
	Number   int
	Protocol string
	State    string
	Process  string
}

// Scanner handles port scanning operations
type Scanner struct {
	timeout time.Duration
}

// New creates a new Scanner with the specified timeout
func New(timeout time.Duration) *Scanner {
	return &Scanner{
		timeout: timeout,
	}
}

// ScanTCPPort checks if a TCP port is open on localhost
func (s *Scanner) ScanTCPPort(port int) (*Port, error) {
	address := fmt.Sprintf("127.0.0.1:%d", port)
	conn, err := net.DialTimeout("tcp", address, s.timeout)
	
	p := &Port{
		Number:   port,
		Protocol: "tcp",
		State:    "closed",
	}

	if err == nil {
		p.State = "open"
		conn.Close()
	}

	return p, nil
}

// ScanPortRange scans a range of TCP ports and returns open ports
func (s *Scanner) ScanPortRange(startPort, endPort int) ([]*Port, error) {
	if startPort < 1 || endPort > 65535 || startPort > endPort {
		return nil, fmt.Errorf("invalid port range: %d-%d", startPort, endPort)
	}

	var openPorts []*Port

	for port := startPort; port <= endPort; port++ {
		p, err := s.ScanTCPPort(port)
		if err != nil {
			continue
		}
		if p.State == "open" {
			openPorts = append(openPorts, p)
		}
	}

	return openPorts, nil
}

// ScanCommonPorts scans commonly used ports
func (s *Scanner) ScanCommonPorts() ([]*Port, error) {
	commonPorts := []int{20, 21, 22, 23, 25, 53, 80, 110, 143, 443, 3306, 5432, 6379, 8080, 8443, 27017}
	var openPorts []*Port

	for _, port := range commonPorts {
		p, err := s.ScanTCPPort(port)
		if err != nil {
			continue
		}
		if p.State == "open" {
			openPorts = append(openPorts, p)
		}
	}

	return openPorts, nil
}
