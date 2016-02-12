package clientpackets

import (
	"../../packets"
)

type ProtocolVersion struct {
	Version uint32
}

func NewProtocolVersion(request []byte) ProtocolVersion {
	var packet = packets.NewReader(request)
	var p ProtocolVersion

	p.Version = packet.ReadUInt32()

	return p
}
