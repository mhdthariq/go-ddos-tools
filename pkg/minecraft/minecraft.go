package minecraft

import (
	"bytes"
	"encoding/binary"
)

// VarInt encodes an integer as a Minecraft VarInt
func VarInt(value int) []byte {
	var buf bytes.Buffer
	for {
		temp := byte(value & 0x7F)
		value >>= 7
		if value != 0 {
			temp |= 0x80
		}
		buf.WriteByte(temp)
		if value == 0 {
			break
		}
	}
	return buf.Bytes()
}

// Data wraps data with VarInt length prefix
func Data(payload []byte) []byte {
	length := VarInt(len(payload))
	return append(length, payload...)
}

// Short encodes a 16-bit integer in big-endian format
func Short(value uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, value)
	return buf
}

// Long encodes a 64-bit integer in big-endian format
func Long(value int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(value))
	return buf
}

// Handshake creates a Minecraft handshake packet
func Handshake(host string, port uint16, protocolVersion int, state int) []byte {
	var packet bytes.Buffer

	// Packet ID (0x00 for handshake)
	packet.Write(VarInt(0x00))

	// Protocol version
	packet.Write(VarInt(protocolVersion))

	// Server address (with length prefix)
	packet.Write(VarInt(len(host)))
	packet.WriteString(host)

	// Server port
	packet.Write(Short(port))

	// Next state (1 for status, 2 for login)
	packet.Write(VarInt(state))

	return Data(packet.Bytes())
}

// Login creates a Minecraft login packet
func Login(protocolVersion int, username string) []byte {
	var packet bytes.Buffer

	// Packet ID depends on protocol version
	packetID := 0x00
	if protocolVersion >= 391 {
		packetID = 0x00
	} else if protocolVersion >= 385 {
		packetID = 0x01
	}

	packet.Write(VarInt(packetID))
	packet.Write(VarInt(len(username)))
	packet.WriteString(username)

	return Data(packet.Bytes())
}

// KeepAlive creates a Minecraft keep-alive packet
func KeepAlive(protocolVersion int, keepAliveID int64) []byte {
	var packet bytes.Buffer

	// Packet ID varies by protocol version
	packetID := 0x0F
	if protocolVersion >= 755 {
		packetID = 0x0F
	} else if protocolVersion >= 712 {
		packetID = 0x10
	}

	packet.Write(VarInt(packetID))

	if protocolVersion >= 339 {
		packet.Write(Long(keepAliveID))
	} else {
		packet.Write(VarInt(int(keepAliveID)))
	}

	return Data(packet.Bytes())
}

// Chat creates a Minecraft chat packet
func Chat(protocolVersion int, message string) []byte {
	var packet bytes.Buffer

	// Packet ID varies by protocol version
	packetID := 0x03
	if protocolVersion >= 755 {
		packetID = 0x03
	} else if protocolVersion >= 464 {
		packetID = 0x03
	} else if protocolVersion >= 389 {
		packetID = 0x02
	}

	packet.Write(VarInt(packetID))
	packet.Write(VarInt(len(message)))
	packet.WriteString(message)

	return Data(packet.Bytes())
}
