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

package keeper_test

import (
	"testing"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/noble-assets/wormhole/types"
	"github.com/noble-assets/wormhole/utils"
	"github.com/noble-assets/wormhole/utils/mocks"
)

func TestParseAndVerifyVAA(t *testing.T) {
	// ARRANGE
	ctx, k := mocks.WormholeKeeper(t)

	// ACT: Query with empty bytes as vaa.
	_, err := k.ParseAndVerifyVAA(ctx, []byte{})

	// ASSERT
	require.Error(t, err, "expected an error when the vaa cannot be unmarshaled")
	require.ErrorContains(t, err, "failed to unmarshal", "expected a different error")

	// ARRANGE: Create a VAA and register it in the archive simulating a previous execution.
	guardian1 := utils.GuardianSigner()
	vaaBody := utils.VAABody{
		GuardianSetIndex: 0,
		Payload:          []byte("first test vaa"),
		Sequence:         1,
	}
	vaa1 := utils.CreateVAA(t, []utils.Guardian{guardian1}, vaaBody)
	hash1 := vaa1.SigningDigest().Bytes()
	err = k.VAAArchive.Set(ctx, hash1, collections.Join(vaa1.MessageID(), true))
	require.NoError(t, err, "expected no error setting the vaa in the archive")

	bzVaa, err := vaa1.Marshal()
	require.NoError(t, err, "expected no error marshaling the vaa")

	// ACT
	_, err = k.ParseAndVerifyVAA(ctx, bzVaa)

	// ASSERT
	require.Error(t, err, "expected error when vaa is already in the archive")
	require.ErrorIs(t, err, types.ErrAlreadyExecutedVAA, "expected a different error")

	// ARRANGE: VAA is valid (not in the achive) but the guardian set is not registered.
	err = k.VAAArchive.Remove(ctx, hash1)
	require.NoError(t, err, "expected no error removing the vaa from archive")

	// ACT
	_, err = k.ParseAndVerifyVAA(ctx, bzVaa)

	// ASSERT
	require.Error(t, err, "expected error when guardian set is not stored")
	require.ErrorContains(t, err, "failed to get guardian set", "expected a different error")

	// ARRANGE: VAA is valid but the guardian set expired.
	headerTime := uint64(sdk.UnwrapSDKContext(ctx).HeaderInfo().Time.Unix())
	err = k.GuardianSets.Set(ctx, 0, types.GuardianSet{ExpirationTime: headerTime - 1})
	require.NoError(t, err, "expected no error setting the guardian set")

	// ACT
	_, err = k.ParseAndVerifyVAA(ctx, bzVaa)

	// ASSERT
	require.Error(t, err, "expected error when guardian set is expired")
	require.ErrorContains(t, err, "expired", "expected a different error")

	// ARRANGE: VAA is valid but the guardian set is different than the one who signed vaa. This
	// makes signature verification failing.
	guardianSet := types.GuardianSet{
		Addresses: [][]byte{
			utils.GuardianSigner().Address.Bytes(),
		},
		ExpirationTime: headerTime,
	}
	err = k.GuardianSets.Set(ctx, 0, guardianSet)
	require.NoError(t, err, "expected no error setting the guardian set")

	// ACT
	_, err = k.ParseAndVerifyVAA(ctx, bzVaa)

	// ASSERT
	require.Error(t, err, "expected error when the signing guardian set is not the correct one")
	require.ErrorContains(t, err, "failed to verify", "expected a different error")

	// ARRANGE: VAA is valid.
	vaaBody = utils.VAABody{
		GuardianSetIndex: 0,
		Payload:          []byte("second test vaa"),
		Sequence:         1,
		EmitterChain:     0,
		EmitterAddress:   [32]byte{},
	}
	guardian2 := utils.GuardianSigner()
	vaa2 := utils.CreateVAA(t, []utils.Guardian{guardian1, guardian2}, vaaBody)
	bzVaa2, err := vaa2.Marshal()
	require.NoError(t, err, "expected no error marshaling the vaa")

	guardianSet = types.GuardianSet{
		ExpirationTime: 0,
		Addresses: [][]byte{
			guardian1.Address[:],
			guardian2.Address[:],
		},
	}
	err = k.GuardianSets.Set(ctx, 0, guardianSet)
	require.NoError(t, err, "expected no error setting the guardian set")

	// ACT
	vaaResp, err := k.ParseAndVerifyVAA(ctx, bzVaa2)

	// ASSERT
	require.NoError(t, err, "expected no error when VAA is valid and signed by the current guardian set")
	require.NoError(t, err, "expected no error retrieving an archived vaa")
	require.Equal(t, vaa2.ConsistencyLevel, vaaResp.ConsistencyLevel, "expected a different ConsistencyLevel")
	require.Equal(t, vaa2.EmitterAddress, vaaResp.EmitterAddress, "expected a different EmitterAddress")
	require.Equal(t, vaa2.EmitterChain, vaaResp.EmitterChain, "expected a different EmitterChain")
	require.Equal(t, vaa2.GuardianSetIndex, vaaResp.GuardianSetIndex, "expected a different GuardianSetIndex")
	require.Equal(t, vaa2.Nonce, vaaResp.Nonce, "expected a different Nonce")
	require.Equal(t, vaa2.Payload, vaaResp.Payload, "expected a different Payload")
	require.Equal(t, vaa2.Sequence, vaaResp.Sequence, "expected a different Sequence")
	require.Equal(t, vaa2.Signatures, vaaResp.Signatures, "expected a different Signatures")
	require.Equal(t, vaa2.Timestamp, vaaResp.Timestamp, "expected a different Timestamp")
	require.Equal(t, vaa2.Version, vaaResp.Version, "expected a different Version")
}
