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

func TestParseAndVerifyVAA(t *testing.T) { // ARRANGE: Create environment
	ctx, k := mocks.WormholeKeeper(t)

	// ACT
	_, err := k.ParseAndVerifyVAA(ctx, []byte{})

	// ASSERT
	require.Error(t, err, "expected error when config not set")
	require.ErrorContains(t, err, "failed to unmarshal", "expected a different error")

	// ARRANGE: Create a VAA already registered in the archive.
	guardian := utils.GuardianSigner()
	vaa1 := utils.CreateVAA(t, []utils.Guardian{guardian}, "first test vaa", 1)
	hash1 := vaa1.SigningDigest().Bytes()
	err = k.VAAArchive.Set(ctx, hash1, collections.Join(vaa1.MessageID(), true))
	require.NoError(t, err, "expected no error setting the vaa in the archive")
	bzVaa, err := vaa1.Marshal()
	require.NoError(t, err, "expected no error marshaling the vaa")

	// ACT
	_, err = k.ParseAndVerifyVAA(ctx, bzVaa)

	// ASSERT
	require.Error(t, err, "expected error when vaa is in the archive")
	require.ErrorIs(t, err, types.ErrAlreadyExecutedVAA, "expected a different error")

	// ARRANGE: VAA is valid but the guardian set is not registered.
	err = k.VAAArchive.Remove(ctx, hash1)
	require.NoError(t, err, "expected no error removing the vaa from archive")

	// ACT
	_, err = k.ParseAndVerifyVAA(ctx, bzVaa)

	// ASSERT
	require.Error(t, err, "expected error when guardian set is not stored")
	require.ErrorContains(t, err, "failed to get guardian set", "expected a different error")

	// ARRANGE: VAA is valid but the guardian set expired.
	err = k.GuardianSets.Set(ctx, 0, types.GuardianSet{ExpirationTime: 1})
	require.NoError(t, err, "expected no error setting the guardian set")

	// ACT
	_, err = k.ParseAndVerifyVAA(ctx, bzVaa)

	// ASSERT
	require.Error(t, err, "expected error when guardian is expired")
	require.ErrorContains(t, err, "expired", "expected a different error")

	// ARRANGE: VAA is valid but no addresses in the VAA for the guardian set.
	guardianSet := types.GuardianSet{ExpirationTime: uint64(sdk.UnwrapSDKContext(ctx).HeaderInfo().Time.Unix())}
	err = k.GuardianSets.Set(ctx, 0, guardianSet)
	require.NoError(t, err, "expected no error setting the guardian set")

	// ACT
	_, err = k.ParseAndVerifyVAA(ctx, bzVaa)

	// ASSERT
	require.Error(t, err, "expected error when the addresses are not valid")
	require.ErrorContains(t, err, "failed to verify", "expected a different error")

	// ARRANGE: VAA is valid but no addresses in the VAA for the guardian set.
	invalidGuardian := utils.GuardianSigner()
	guardianSet.Addresses = [][]byte{invalidGuardian.Address[:]}
	err = k.GuardianSets.Set(ctx, 0, guardianSet)
	require.NoError(t, err, "expected no error setting the guardian set")

	// ACT
	_, err = k.ParseAndVerifyVAA(ctx, bzVaa)

	// ASSERT
	require.Error(t, err, "expected error when the guardian is set is different than signing set")
	require.ErrorContains(t, err, "failed to verify", "expected a different error")

	// ARRANGE: VAA is valid.
	vaa2 := utils.CreateVAA(t, []utils.Guardian{guardian}, "second test vaa", 1)
	bzVaa2, err := vaa2.Marshal()
	require.NoError(t, err, "expected no error marshaling the vaa")
	guardianSet = types.GuardianSet{
		ExpirationTime: 0,
		Addresses:      [][]byte{guardian.Address[:]},
	}
	err = k.GuardianSets.Set(ctx, 0, guardianSet)
	require.NoError(t, err, "expected no error setting the guardian set")

	// ACT
	vaaResp, err := k.ParseAndVerifyVAA(ctx, bzVaa2)

	// ASSERT
	require.NoError(t, err, "expected no error when VAA is valid")
	require.NoError(t, err, "expected no error retrieving an archived vaa")
	require.Equal(t, vaa2.ConsistencyLevel, vaaResp.ConsistencyLevel, "expected a different ConsistencyLevel")
	require.Equal(t, vaa2.EmitterAddress, vaaResp.EmitterAddress, "expected a different EmitterAddress")
	require.Equal(t, vaa2.EmitterChain, vaaResp.EmitterChain, "expected a different EmitterChain")
	require.Equal(t, vaa2.GuardianSetIndex, vaaResp.GuardianSetIndex, "expected a different GuardianSetIndex")
	require.Equal(t, vaa2.Nonce, vaaResp.Nonce, "expected a different Nonce")
	require.Equal(t, vaa2.Payload, vaaResp.Payload, "expected a different Payload")
	require.Equal(t, vaa2.Sequence, vaaResp.Sequence, "expected a different Sequence")
	require.Equal(t, vaa2.Signatures, vaaResp.Signatures, "expected a different Signatures")
	require.Equal(t, vaa2.Timestamp, vaaResp.Timestamp, "expected a different timestamp")
	require.Equal(t, vaa2.Version, vaaResp.Version, "expected a different version")
}

func TestHandleCoreGovernancePacket(t *testing.T) {
	// ARRANGE: Create environment
	ctx, k := mocks.WormholeKeeper(t)
	packet := types.GovernancePacket{}

	// ACT
	err := k.HandleCoreGovernancePacket(ctx, packet)

	// ASSERT
	require.Error(t, err, "expected error when packet is empty")
	require.ErrorContains(t, err, "unsupported governance action", "expected a different error")

	// ARRANGE
	packet.Action = 3
	packet.Module = "Core"

	// ACT
	err = k.HandleCoreGovernancePacket(ctx, packet)

	// ASSERT
	require.Error(t, err, "expected error when governance action is not 2")
	require.ErrorIs(t, err, types.ErrUnsupportedGovernanceAction, "expected a different error")

	// ARRANGE: The action is valid but the payload is empty.
	packet.Action = 2

	// ACT
	err = k.HandleCoreGovernancePacket(ctx, packet)

	// ASSERT
	require.Error(t, err, "expected error when governance malformed payload")
	require.ErrorIs(t, err, types.ErrMalformedPayload, "expected a different error")

	// ARRANGE: The action is valid but the payload is too short.
	packet.Payload = []byte("shrt")

	// ACT
	err = k.HandleCoreGovernancePacket(ctx, packet)

	// ASSERT
	require.Error(t, err, "expected error when governance malformed payload")
	require.ErrorIs(t, err, types.ErrMalformedPayload, "expected a different error")

	// ARRANGE: Set an invalid payload to fail during parsing
	// This payload is valid and will be used for all tests below.
	packet.Payload = []byte{
		0x00, 0x00, 0x00, 0x01, // index = 1
		0x02, // length of 2
		// First address (20 bytes)
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
		0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14,
	}

	// ACT
	err = k.HandleCoreGovernancePacket(ctx, packet)

	// ASSERT
	require.Error(t, err, "expected error when payload is malformed ")
	require.ErrorIs(t, err, types.ErrMalformedPayload, "expected a different error")

	// ARRANGE: Config is not set
	packet.Payload = []byte{
		0x00, 0x00, 0x00, 0x01, // index = 1
		0x02, // length of 2
		// First address (20 bytes)
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
		0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14,
		// Second address (20 bytes)
		0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e,
		0x1f, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
	}

	// ACT
	err = k.HandleCoreGovernancePacket(ctx, packet)

	// ASSERT
	require.Error(t, err, "expected error when the config is not set")
	require.ErrorContains(t, err, "failed to get config", "expected a different error")

	// ARRANGE: Set the governance with an index that makes fail the payload
	cfg := types.Config{
		GuardianSetIndex: 1, // same index of the payload
	}
	err = k.Config.Set(ctx, cfg)
	require.NoError(t, err, "expected no error setting the config")

	// ACT
	err = k.HandleCoreGovernancePacket(ctx, packet)

	// ASSERT
	require.Error(t, err, "expected error when guardian set index is not valid")
	require.ErrorContains(t, err, "invalid guardian set index", "expected a different error")

	// ARRANGE: Set the governance with an index that makes fail the payload
	cfg.GuardianSetIndex = 0
	err = k.Config.Set(ctx, cfg)
	require.NoError(t, err, "expected no error setting the config")

	// ACT
	err = k.HandleCoreGovernancePacket(ctx, packet)

	// ASSERT
	require.NoError(t, err, "expected no error when the payload and the state are valid")

	respCfg, err := k.Config.Get(ctx)
	require.NoError(t, err, "expected no error reading the config")
	require.Equal(t, uint32(1), respCfg.GuardianSetIndex, "expected a different guardian set index")

	respGuardianSet, err := k.GuardianSets.Get(ctx, 1)
	require.NoError(t, err, "expected no error reading the guardian set")
	require.Len(t, respGuardianSet.Addresses, 2, "expected a different number of addresses in the set")
}
