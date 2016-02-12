package serverpackets

import "../../packets"

func NewPlayFailPacket(reason uint32) []byte {
	buffer := new(packets.Buffer)
	buffer.WriteByte(0x06) // Packet type: PlayFail
	buffer.WriteUInt32(reason)

	return buffer.Bytes()
}

func NewPlayOkPacket() []byte {
	buffer := new(packets.Buffer)
	buffer.WriteByte(0x07)
	buffer.Write([]byte{0x34, 0x0b, 0x00, 0x01}) // Session Key
	buffer.Write([]byte{0x55, 0x66, 0x77, 0x88}) // Session Key 2?

	return buffer.Bytes()
}