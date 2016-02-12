package serverpackets

import (
	"../../packets"
)

func NewCharCreateOkPacket() []byte {
	buffer := packets.NewBuffer()
	buffer.WriteByte(0x25)          // Packet type: CharCreateOk
	buffer.WriteUInt32(0x01)        // Everything went like expected

	return buffer.Bytes()
}

func NewCharListPacket() []byte {
	buffer := packets.NewBuffer()
	buffer.WriteByte(0x1f)                       // Packet type: CharList
	buffer.Write([]byte{0x00, 0x00, 0x00, 0x00}) // TODO

	return buffer.Bytes()
}

func NewCharTemplatePacket() []byte {
	buffer := packets.NewBuffer()
	buffer.WriteByte(0x23)          // Packet type: CharTemplate
	buffer.WriteUInt32(0x00)        // We don't actually need to send the template to the client

	return buffer.Bytes()
}
