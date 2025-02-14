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
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/noble-assets/wormhole/keeper"
	"github.com/noble-assets/wormhole/types"
	"github.com/noble-assets/wormhole/utils/mocks"
)

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

	// ASSERT: Set valid guardian index but don't populate the state with old guardian set.
	require.Error(t, err, "expected error when guardian set index is not valid")
	require.ErrorContains(t, err, "invalid guardian set index", "expected a different error")

	// ARRANGE
	cfg.GuardianSetIndex = 0
	err = k.Config.Set(ctx, cfg)
	require.NoError(t, err, "expected no error setting the config")

	// ACT
	err = k.HandleCoreGovernancePacket(ctx, packet)

	// ASSERT
	require.Error(t, err, "expected expected an error when no old guardian set is present")
	require.ErrorContains(t, err, "failed to get old guardian set", "expected a different error")

	// ARRANGE: Set empty guardian state to the store
	err = k.GuardianSets.Set(ctx, 0, types.GuardianSet{})
	require.NoError(t, err, "expected no error setting the old guardian set")

	// ACT
	err = k.HandleCoreGovernancePacket(ctx, packet)

	// ASSERT
	require.NoError(t, err, "expected no error when the payload and the state are valid")

	respCfg, err := k.Config.Get(ctx)
	require.NoError(t, err, "expected no error reading the config")
	require.Equal(t, uint32(1), respCfg.GuardianSetIndex, "expected a different guardian set index")

	// Check old guardian set
	respOldGuardianSet, err := k.GuardianSets.Get(ctx, 0)
	require.NoError(t, err, "expected no error reading the old guardian set")
	println("Time from ctx :", uint64(sdk.UnwrapSDKContext(ctx).HeaderInfo().Time.Unix()))
	println("Time in tests: ", time.Now().Truncate(time.Second).Unix())
	expTime := uint64(time.Now().Truncate(time.Second).Unix()) + keeper.GuardianSetExpiry
	require.Equal(t, expTime, respOldGuardianSet.ExpirationTime, "expected a different timestamp for old guardian set")

	// Check updated guardian set.
	respGuardianSet, err := k.GuardianSets.Get(ctx, 1)
	require.NoError(t, err, "expected no error reading the guardian set")
	require.Len(t, respGuardianSet.Addresses, 2, "expected a different number of addresses in the set")
}
