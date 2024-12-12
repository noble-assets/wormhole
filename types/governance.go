package types

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// GovernancePacket defines the expected payload of a Governance VAA.
type GovernancePacket struct {
	Module  string
	Action  uint8
	Chain   uint16
	Payload []byte
}

func (pkt *GovernancePacket) Parse(bz []byte) error {
	if len(bz) < 35 {
		return errors.New("governance packet is malformed")
	}

	pkt.Module = string(bytes.TrimLeft(bz[0:32], "\x00"))
	pkt.Action = bz[32:33][0]
	pkt.Chain = binary.BigEndian.Uint16(bz[33:35])
	pkt.Payload = bz[35:]

	return nil
}
