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

	"github.com/stretchr/testify/require"

	"github.com/noble-assets/wormhole/types"
	"github.com/noble-assets/wormhole/utils"
	"github.com/noble-assets/wormhole/utils/mocks"
)

func TestGetChain(t *testing.T) {
	// ARRANGE: Create environment
	ctx, k := mocks.WormholeKeeper(t)

	// ACT
	_, err := k.GetChain(ctx)

	// ASSERT
	require.Error(t, err, "expected error when config not set")
	require.ErrorContains(t, err, "unable to get", "expected a different error")

	// ARRANGE: Add empty config data to state
	config := types.Config{
		ChainId:           0,
		GuardianSetIndex:  0,
		GuardianSetExpiry: 0,
		GovChain:          0,
		GovAddress:        []byte{},
	}
	err = k.Config.Set(ctx, config)
	require.NoError(t, err, "expected no error setting the config")

	// ACT
	chain, err := k.GetChain(ctx)

	// ASSERT
	require.NoError(t, err, "expected no error when config is set")
	require.Equal(t, config.ChainId, chain, "expected different chain id")
}

func TestGetGuardianSets(t *testing.T) {
	// ARRANGE: Create environment
	ctx, k := mocks.WormholeKeeper(t)

	// ACT
	guardianSets, err := k.GetGuardianSets(ctx)

	// ASSERT
	require.NoError(t, err, "expected no error when the set is empty")
	require.Empty(t, guardianSets, "expected empty map")

	// ARRANGE: Add guardian sets
	set1 := types.GuardianSet{
		Addresses: [][]byte{
			utils.GenerateRandomBytes(20),
			utils.GenerateRandomBytes(20),
		},
		ExpirationTime: uint64(0),
	}
	set2 := types.GuardianSet{
		Addresses: [][]byte{
			utils.GenerateRandomBytes(20),
		},
		ExpirationTime: uint64(1),
	}

	err = k.GuardianSets.Set(ctx, 1, set1)
	require.NoError(t, err, "expected no error setting the guardian")
	err = k.GuardianSets.Set(ctx, 2, set2)
	require.NoError(t, err, "expected no error setting the guardian")

	// ACT
	guardianSets, err = k.GetGuardianSets(ctx)

	// ASSERT
	require.NoError(t, err, "expected no error when the set is empty")
	require.Len(t, guardianSets, 2, "expected two sets")
	require.Equal(t, set1, guardianSets[1], "expected different values in first set")
	require.Equal(t, set2, guardianSets[2], "expected different values in second set")
}

func TestGetSequences(t *testing.T) {
	// ARRANGE: Create environment
	ctx, k := mocks.WormholeKeeper(t)

	// ACT
	sequences, err := k.GetSequences(ctx)

	// ASSERT
	require.NoError(t, err, "expected no error when no sequences are registered")
	require.Empty(t, sequences, "expected empty map")

	// ARRANGE: Add addresses with less than 20 bytes
	adddress1 := utils.GenerateRandomBytes(10)

	err = k.Sequences.Set(ctx, adddress1, 0)
	require.NoError(t, err, "expected no error setting the sequence")

	// ACT
	sequences, err = k.GetSequences(ctx)

	// ASSERT
	require.Error(t, err, "expected an error when an address is less than 20 bytes")
	require.ErrorContains(t, err, "address with less than 20 bytes", "expected a different error")
	require.Empty(t, sequences, "expected empty map")

	// ARRANGE: Add another valid address
	adddress2 := utils.GenerateRandomBytes(20)

	err = k.Sequences.Set(ctx, adddress2, 1)
	require.NoError(t, err, "expected no error setting the sequence")

	// ACT
	sequences, err = k.GetSequences(ctx)

	// ASSERT
	require.Error(t, err, "expected an error when an address is less than 20 bytes")
	require.ErrorContains(t, err, "address with less than 20 bytes", "expected a different error")
	require.Len(t, sequences, 1, "expected the valid address to be in the map")

	// TODO: how to make address codec fail?
}
