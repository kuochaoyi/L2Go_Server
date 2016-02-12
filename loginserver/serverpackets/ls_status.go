package serverpackets

import "../../packets"

func NewLoginFailPacket(reason uint32) []byte {
	buffer := new(packets.Buffer)
	buffer.WriteByte(0x01) // Packet type: LoginFail
	buffer.WriteUInt32(reason)

	return buffer.Bytes()
}

func NewLoginOkPacket(sessionID []byte) []byte {
	buffer := new(packets.Buffer)
	buffer.WriteByte(0x03) // Packet type: LoginOk
	buffer.Write(sessionID[:4])  // Session id 1/2
	buffer.Write(sessionID[4:8]) // Session id 2/2
	buffer.WriteUInt32(0x00)
	buffer.WriteUInt32(0x00)
	buffer.WriteUInt32(0x000003ea)
	buffer.WriteUInt32(0x00)
	buffer.WriteUInt32(0x00)
	buffer.WriteUInt32(0x02)

	return buffer.Bytes()
}