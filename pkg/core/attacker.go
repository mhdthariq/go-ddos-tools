package core

import (
	"context"

	"github.com/go-ddos-tools/pkg/proxy"
	"github.com/go-ddos-tools/pkg/utils"
)

// AttackConfig holds common configuration for all attacks
type AttackConfig struct {
	Target       string
	Duration     int
	Threads      int
	RPC          int
	Proxies      []proxy.Proxy
	Reflectors   []string
	ProtocolID   int
	UserAgents   []string
	Referers     []string
	RequestsSent *utils.Counter
	BytesSent    *utils.Counter
}

// Attacker is the interface that all attack methods must implement
type Attacker interface {
	// Attack performs a single iteration of the attack
	// It should respect the context cancellation for stopping
	Attack(ctx context.Context) error
}
