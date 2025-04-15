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

package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"

	"github.com/noble-assets/wormhole/types"
	"github.com/noble-assets/wormhole/utils"
)

func TestUpdateChannelChain_Parse(t *testing.T) {
	channelBz, err := vaautils.LeftPadIbcChannelId("channel-0")
	require.NoError(t, err)
	// Shift left by eight for most significant byte and mask for less significant ones.
	chainIDBz := []byte{
		byte(uint16(vaautils.ChainIDWormchain) >> 8),
		byte(uint16(vaautils.ChainIDWormchain) & 0xFF),
	}
	validPayload := make([]byte, 66)
	copy(validPayload[:64], channelBz[:])
	copy(validPayload[64:], chainIDBz)

	require.NoError(t, err, "expected no error padding the channel")

	testCases := []struct {
		name              string
		payload           []byte
		expectedChannelID []byte
		expectedChain     uint16
		errorContains     string
	}{
		{
			"fail when payload is empty",
			[]byte{},
			[]byte{},
			0,
			"payload is malformed",
		},
		{
			"fail when payload is less than 66",
			utils.GenerateRandomBytes(65),
			[]byte{},
			0,
			"payload is malformed",
		},
		{
			"successful when payload is valid",
			validPayload,
			[]byte("channel-0"),
			uint16(vaautils.ChainIDWormchain),
			"",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			var updateChannelChain types.UpdateChannelChain
			err := updateChannelChain.Parse(tC.payload)

			if tC.errorContains != "" {
				require.Error(t, err, "expected an error")
				require.ErrorContains(t, err, tC.errorContains, "expected a different error")
			} else {
				require.NoError(t, err, "expected no error")
				require.Equal(t, tC.expectedChain, updateChannelChain.Chain, "expected different chain")
				require.Equal(t, tC.expectedChannelID, updateChannelChain.ChannelID, "expected different channel ID")
			}
		})
	}
}

func TestGuardianSetUpdate_Parse(t *testing.T) {
	testCases := []struct {
		name                    string
		payload                 []byte
		expectedIndex           uint32
		expectedAddressLenght   uint32
		exppectedExpirationTime uint64
		errorContains           string
	}{
		{
			"fail when payload is empty",
			[]byte{},
			0,
			0,
			0,
			"payload is malformed",
		},
		{
			"fail when payload is too short",
			[]byte{0x00, 0x00, 0x00, 0x01},
			0,
			0,
			0,
			"payload is malformed",
		},
		{
			"fail when addresses length doesn't match",
			[]byte{
				0x00, 0x00, 0x00, 0x01, // index
				0x02,                   // length of 2
				0x01, 0x02, 0x03, 0x04, // incomplete address data
			},
			0,
			0,
			0,
			"payload is malformed",
		},
		{
			"fails when only one address is present and two are expected",
			[]byte{
				0x00, 0x00, 0x00, 0x01, // index = 1
				0x02, // length of 2
				// Only address (20 bytes)
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
				0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14,
			},
			0,
			0,
			0,
			"payload is malformed",
		},
		{
			"successful when valid payload",
			[]byte{
				0x00, 0x00, 0x00, 0x01, // index = 1
				0x02, // length of 2
				// First address (20 bytes)
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
				0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14,
				// Second address (20 bytes)
				0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e,
				0x1f, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
			},
			1,
			2,
			0,
			"",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			var guardianSetUpgrade types.GuardianSetUpdate
			err := guardianSetUpgrade.Parse(tC.payload)

			if tC.errorContains != "" {
				require.Error(t, err, "expected an error")
				require.ErrorContains(t, err, tC.errorContains, "expected a different error")
			} else {
				require.NoError(t, err, "expected no error")
				require.Equal(t, tC.expectedIndex, guardianSetUpgrade.NewGuardianSetIndex, "expected different guardian set index")
				require.Len(t, guardianSetUpgrade.NewGuardianSet.Addresses, int(tC.expectedAddressLenght), "expected a different number of addresses")
				require.Equal(t, tC.exppectedExpirationTime, guardianSetUpgrade.NewGuardianSet.ExpirationTime, "expected a different expiration time")
			}
		})
	}
}
