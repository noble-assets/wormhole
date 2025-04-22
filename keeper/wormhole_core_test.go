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

	"github.com/stretchr/testify/require"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"

	"github.com/noble-assets/wormhole/keeper"
	"github.com/noble-assets/wormhole/types"
	"github.com/noble-assets/wormhole/utils/mocks"
)

func TestHandleCoreGovernancePacket(t *testing.T) {
	// ARRANGE: Create environment.
	ctx, k := mocks.WormholeKeeper(t)
	packet := types.GovernancePacket{
		Action: uint8(vaautils.ActionGuardianSetUpdate) + 1, // action guardian set update is the only supported action
	}

	// ACT
	err := k.HandleCoreGovernancePacket(ctx, packet)

	// ASSERT
	require.Error(t, err, "expected error when governance action is not 2")
	require.ErrorContains(t, err, "unsupported governance action", "expected a different error")

	// ARRANGE: The action is valid but the payload is empty.
	packet.Action = uint8(vaautils.ActionGuardianSetUpdate)

	// ACT
	err = k.HandleCoreGovernancePacket(ctx, packet)

	// ASSERT
	require.Error(t, err, "expected error when the payload is malformed")
	require.ErrorIs(t, err, types.ErrMalformedPayload, "expected a different error")

	// ARRANGE: The payload is valid but config is not set.
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
	require.Error(t, err, "expected error when the payload is valid but the config is not set")
	require.ErrorContains(t, err, "failed to get config", "expected a different error")

	// ARRANGE: Set the config with a guardian set index different that makes the payload invalid
	cfg := types.Config{
		GuardianSetIndex: 1, // same index of the payload (the valid should be payload - 1)
	}
	err = k.Config.Set(ctx, cfg)
	require.NoError(t, err, "expected no error setting the config")

	// ACT
	err = k.HandleCoreGovernancePacket(ctx, packet)

	// ASSERT
	require.Error(t, err, "expected error when guardian set index is not valid")
	require.ErrorContains(t, err, "invalid guardian set index", "expected a different error")

	// ARRANGE: Set valid guardian set index but don't populate the state with old guardian set.
	cfg.GuardianSetIndex = 0
	err = k.Config.Set(ctx, cfg)
	require.NoError(t, err, "expected no error setting the config")

	// ACT
	err = k.HandleCoreGovernancePacket(ctx, packet)

	// ASSERT
	require.Error(t, err, "expected an error when no old guardian set is present")
	require.ErrorContains(t, err, "failed to get old guardian set", "expected a different error")

	// ARRANGE: Set empty guardian state to the store.
	err = k.GuardianSets.Set(ctx, 0, types.GuardianSet{})
	require.NoError(t, err, "expected no error setting the old guardian set")

	// ACT
	err = k.HandleCoreGovernancePacket(ctx, packet)

	// ASSERT
	require.NoError(t, err, "expected no error when the payload and the state are valid")

	respCfg, err := k.Config.Get(ctx)
	require.NoError(t, err, "expected no error reading the config")
	require.Equal(t, uint32(1), respCfg.GuardianSetIndex, "expected a different guardian set index after the update")

	// Check old guardian set.
	respOldGuardianSet, err := k.GuardianSets.Get(ctx, 0)
	require.NoError(t, err, "expected no error reading the old guardian set")
	expTime := uint64(time.Now().Truncate(time.Second).Unix()) + keeper.GuardianSetExpiry
	require.Equal(t, expTime, respOldGuardianSet.ExpirationTime, "expected a different timestamp for old guardian set")

	// Check updated guardian set.
	respGuardianSet, err := k.GuardianSets.Get(ctx, 1)
	require.NoError(t, err, "expected no error reading the guardian set")
	require.Len(t, respGuardianSet.Addresses, 2, "expected a different number of addresses in the set")
}
