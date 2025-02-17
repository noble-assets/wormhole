// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, NASD Inc. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

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

func (pkt *GovernancePacket) Serialize() []byte {
	buf := make([]byte, 35+len(pkt.Payload))

	moduleBz := []byte(pkt.Module)
	copy(buf[0:32], make([]byte, 32))
	if len(moduleBz) < 32 {
		copy(buf[32-len(moduleBz):], moduleBz) // Right-pad with 0x00
	}

	buf[32] = pkt.Action

	binary.BigEndian.PutUint16(buf[33:35], pkt.Chain)

	copy(buf[35:], pkt.Payload)

	return buf
}

// GuardianSetUpdate represents the governance action to update the guardian set.
type GuardianSetUpdate struct {
	NewGuardianSetIndex uint32
	NewGuardianSet      GuardianSet
}

func (a *GuardianSetUpdate) Parse(payload []byte) error {
	if len(payload) < 5 {
		return ErrMalformedPayload
	}

	newGuardianSetIndex := binary.BigEndian.Uint32(payload[0:4])

	newGuardianSetLength := int(payload[4:5][0])

	if len(payload[5:]) != AddressLenght*newGuardianSetLength {
		return ErrMalformedPayload
	}

	// Offset is given by the 4 bytes of the new index + the single byte of the guardian set lenght.
	offset := 5
	addresses := make([][]byte, newGuardianSetLength)
	for i := range newGuardianSetLength {
		addresses[i] = payload[offset : offset+20]
		offset += 20
	}

	a.NewGuardianSetIndex = newGuardianSetIndex
	a.NewGuardianSet.Addresses = addresses
	a.NewGuardianSet.ExpirationTime = 0

	return nil
}

// UpdateChannelChain represents the governance action to update the IBC
// channel associated with a chain.
type UpdateChannelChain struct {
	ChannelID []byte
	Chain     uint16
}

func (a *UpdateChannelChain) Parse(payload []byte) error {
	if len(payload) != 66 {
		return ErrMalformedPayload
	}

	a.ChannelID = bytes.TrimLeft(payload[0:64], "\x00")
	a.Chain = binary.BigEndian.Uint16(payload[64:66])

	return nil
}
