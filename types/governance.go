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

const AddressLenght = 20

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

type GuardianSetUpgrade struct {
	NewGuardianSetIndex uint32
	NewGuardianSet      GuardianSet
}

func (p *GuardianSetUpgrade) Parse(payload []byte) error {
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

	p.NewGuardianSetIndex = newGuardianSetIndex

	return nil
}
