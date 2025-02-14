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

	"github.com/noble-assets/wormhole/types"
)

func TestGuardianSetUpgrade_Parse(t *testing.T) {
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
			var guardianSetUpgrade types.GuardianSetUpgrade
			err := guardianSetUpgrade.Parse(tC.payload)

			if tC.errorContains != "" {
				require.Error(t, err, "expected an error")
				require.ErrorContains(t, err, tC.errorContains, "expected a different error")
			} else {
				require.NoError(t, err, "expected no error")
				require.Equal(t, tC.expectedIndex, guardianSetUpgrade.NewGuardianSetIndex, "expected different guardian set index")
			}
		})
	}
}
