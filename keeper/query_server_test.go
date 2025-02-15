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
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/noble-assets/wormhole/keeper"
	"github.com/noble-assets/wormhole/types"
	"github.com/noble-assets/wormhole/utils"
	"github.com/noble-assets/wormhole/utils/mocks"
)

func TestConfig(t *testing.T) {
	// ARRANGE
	ctx, k := mocks.WormholeKeeper(t)
	qs := keeper.NewQueryServer(k)

	// ACT
	resp, err := qs.Config(ctx, nil)

	// ASSERT
	require.Error(t, err, "expected error when the request is nil")
	require.ErrorIs(t, err, types.ErrInvalidRequest, "expected a different error")
	require.Nil(t, resp, "expected nil response")

	// ARRANGE
	req := types.QueryConfig{}

	// ACT
	resp, err = qs.Config(ctx, &req)

	// ASSERT
	require.Error(t, err, "expected error when the request is nil")
	require.ErrorContains(t, err, "unable to get config", "expected a different error")
	require.Nil(t, resp, "expected nil response")

	// ARRANGE
	cfg := types.Config{
		ChainId:          1,
		GuardianSetIndex: 2,
		GovChain:         3,
		GovAddress:       []byte("4"),
	}
	err = k.Config.Set(ctx, cfg)
	require.NoError(t, err, "expected no error setting the config")

	// ACT
	resp, err = qs.Config(ctx, &req)

	// ASSERT
	require.NoError(t, err, "expected no error when the request is valid and config exists")
	require.Equal(t, cfg.ChainId, resp.Config.ChainId, "expected a different chain id")
	require.Equal(t, cfg.GuardianSetIndex, resp.Config.GuardianSetIndex, "expected a different guardian set index")
	require.Equal(t, cfg.GovChain, resp.Config.GovChain, "expected a different gov chain")
	require.Equal(t, cfg.GovAddress, resp.Config.GovAddress, "expected a different gov address")
}

func TestWormchainChannel(t *testing.T) {
	// ARRANGE
	ctx, k := mocks.WormholeKeeper(t)
	qs := keeper.NewQueryServer(k)

	// ACT
	resp, err := qs.WormchainChannel(ctx, nil)

	// ASSERT
	require.Error(t, err, "expected error when the request is nil")
	require.ErrorIs(t, err, types.ErrInvalidRequest, "expected a different error")
	require.Nil(t, resp, "expected nil response")

	// ARRANGE
	req := types.QueryWormchainChannel{}

	// ACT
	resp, err = qs.WormchainChannel(ctx, &req)

	// ASSERT
	require.Error(t, err)
	require.ErrorContains(t, err, "wormchain channel not configured in state", "expected a different error")
	require.Nil(t, resp, "expected nil response")

	// ARRANGE
	wormchainChannel := "channel-0"
	err = k.WormchainChannel.Set(ctx, wormchainChannel)
	require.NoError(t, err, "expected no error setting the wormchain channel")

	// ACT
	resp, err = qs.WormchainChannel(ctx, &req)

	// ASSERT
	require.NoError(t, err, "expected no error when the request is valid and channel is set")
	require.Equal(t, wormchainChannel, resp.WormchainChannel, "expected a different channel")
}

func TestGuardianSets(t *testing.T) {
	// ARRANGE
	ctx, k := mocks.WormholeKeeper(t)
	qs := keeper.NewQueryServer(k)

	// ACT
	resp, err := qs.GuardianSets(ctx, nil)

	// ASSERT
	require.Error(t, err, "expected error when the request is nil")
	require.ErrorIs(t, err, types.ErrInvalidRequest, "expected a different error")
	require.Nil(t, resp, "expected nil response")

	// ARRANGE
	req := types.QueryGuardianSets{}

	// ACT
	resp, err = qs.GuardianSets(ctx, &req)

	// ASSERT
	require.NoError(t, err, "expected no error when no set is stored")
	require.Equal(t, make(map[uint32]types.GuardianSet), resp.GuardianSets, "expected no sets")

	// ARRANGE
	key1 := uint32(0)
	set1 := types.GuardianSet{
		Addresses: [][]byte{
			[]byte("address1"),
			[]byte("address2"),
		},
		ExpirationTime: uint64(1),
	}
	err = k.GuardianSets.Set(ctx, key1, set1)
	require.NoError(t, err, "expected no error setting the guardian set")

	key2 := uint32(1)
	set2 := types.GuardianSet{
		Addresses: [][]byte{
			[]byte("address3"),
			[]byte("address4"),
			[]byte("address5"),
		},
		ExpirationTime: uint64(3),
	}
	err = k.GuardianSets.Set(ctx, key2, set2)
	require.NoError(t, err, "expected no error setting the guardian set")

	// ACT
	resp, err = qs.GuardianSets(ctx, &req)

	// ASSERT
	require.NoError(t, err, "expected no error when the request is valid and the a guardian set is in the store")
	require.Len(t, resp.GuardianSets, 2, "expected a different number of sets")
	require.Equal(t, set1, resp.GuardianSets[0], "expected a different first set")
	require.Equal(t, set2, resp.GuardianSets[1], "expected a different second set")
}

func TestGuardianSet(t *testing.T) {
	// ARRANGE
	ctx, k := mocks.WormholeKeeper(t)
	qs := keeper.NewQueryServer(k)

	// ACT
	resp, err := qs.GuardianSet(ctx, nil)

	// ASSERT
	require.Error(t, err, "expected error when the request is nil")
	require.ErrorIs(t, err, types.ErrInvalidRequest, "expected a different error")
	require.Nil(t, resp, "expected nil response")

	// ARRANGE
	req := types.QueryGuardianSet{}

	// ACT
	resp, err = qs.GuardianSet(ctx, &req)

	// ASSERT
	require.Error(t, err, "expected error when the request is nil")
	require.ErrorContains(t, err, "invalid guardian set index", "expected a different error")
	require.Nil(t, resp, "expected nil response")

	// ARRANGE
	req = types.QueryGuardianSet{Index: "invalid"}

	// ACT
	resp, err = qs.GuardianSet(ctx, &req)

	// ASSERT
	require.Error(t, err, "expected error when the request index is invalid")
	require.ErrorContains(t, err, "invalid guardian set index", "expected a different error")
	require.Nil(t, resp, "expected nil response")

	// ARRANGE
	req = types.QueryGuardianSet{Index: "current"}

	// ACT
	resp, err = qs.GuardianSet(ctx, &req)

	// ASSERT
	require.Error(t, err, "expected error when the request is valid but config not set")
	require.ErrorContains(t, err, "unable to get config from state", "expected a different error")
	require.Nil(t, resp, "expected nil response")

	// ARRANGE
	req = types.QueryGuardianSet{Index: "current"}

	// ACT
	cfg := types.Config{GuardianSetIndex: 33}
	err = k.Config.Set(ctx, cfg)
	require.NoError(t, err, "expected no error setting the config")
	resp, err = qs.GuardianSet(ctx, &req)

	// ASSERT
	require.Error(t, err, "expected error when the guardian set is not stored")
	require.ErrorContains(t, err, "unable to get guardian set", "expected a different error")
	require.Nil(t, resp, "expected nil response")

	// ARRANGE
	req = types.QueryGuardianSet{Index: "1"}

	// ACT
	resp, err = qs.GuardianSet(ctx, &req)

	// ASSERT
	require.Error(t, err, "expected error when the guardian set is not stored")
	require.ErrorContains(t, err, "unable to get guardian set", "expected a different error")
	require.Nil(t, resp, "expected nil response")

	// ARRANGE
	req = types.QueryGuardianSet{Index: "1"}

	key := uint32(1)
	set := types.GuardianSet{
		Addresses: [][]byte{
			[]byte("address1"),
			[]byte("address2"),
		},
		ExpirationTime: uint64(1),
	}
	err = k.GuardianSets.Set(ctx, key, set)
	require.NoError(t, err, "expected no error setting the guardian set")

	// ACT
	resp, err = qs.GuardianSet(ctx, &req)

	// ASSERT
	require.NoError(t, err, "expected no error when the guardian set associated with the index exists")
	require.Equal(t, set, resp.GuardianSet, "expected a different set")

	// ARRANGE
	req = types.QueryGuardianSet{Index: "current"}

	key = uint32(33)
	set = types.GuardianSet{
		Addresses: [][]byte{
			[]byte("address0"),
		},
		ExpirationTime: uint64(0),
	}
	err = k.GuardianSets.Set(ctx, key, set)
	require.NoError(t, err, "expected no error setting the guardian set")

	// ACT
	resp, err = qs.GuardianSet(ctx, &req)

	// ASSERT
	require.NoError(t, err, "expected no error when the guardian set associated with the current index exists")
	require.Equal(t, set, resp.GuardianSet, "expected a different set")
}

func TestExecutedVAA(t *testing.T) {
	// ARRANGE
	ctx, k := mocks.WormholeKeeper(t)
	qs := keeper.NewQueryServer(k)

	// ACT
	resp, err := qs.ExecutedVAA(ctx, nil)

	// ASSERT
	require.Error(t, err, "expected error when the request is nil")
	require.ErrorIs(t, err, types.ErrInvalidRequest, "expected a different error")
	require.Nil(t, resp, "expected nil response")

	// ARRANGE
	req := types.QueryExecutedVAA{InputType: "wrong"}

	// ACT
	resp, err = qs.ExecutedVAA(ctx, &req)

	// ASSERT
	require.Error(t, err, "expected error when the input type is not allowed")
	require.ErrorContains(t, err, "invalid input type", "expected a different error")
	require.Nil(t, resp, "expected nil response")

	// ARRANGE
	req = types.QueryExecutedVAA{InputType: ""}

	// ACT
	resp, err = qs.ExecutedVAA(ctx, &req)

	// ASSERT
	require.NoError(t, err, "expected no error when the input is not a valid digest")
	require.False(t, resp.Executed, "expected no vaa in executed")

	// ARRANGE
	vaa := utils.CreateVAA(t, []utils.Guardian{utils.GuardianSigner()}, "valid", 0)
	digest := vaa.SigningDigest().Bytes()
	req = types.QueryExecutedVAA{InputType: "", Input: common.Bytes2Hex(digest)}

	// ACT
	resp, err = qs.ExecutedVAA(ctx, &req)

	// ASSERT
	require.NoError(t, err, "expected no error when req is valid but archive is empty")
	require.False(t, resp.Executed, "expected no vaa in executed")

	// ARRANGE
	err = k.VAAArchive.Set(ctx, digest, collections.Join(vaa.MessageID(), true))
	require.NoError(t, err, "expected no error setting the vaa")

	// ACT
	resp, err = qs.ExecutedVAA(ctx, &req)

	// ASSERT
	require.NoError(t, err)
	require.True(t, resp.Executed, "expected the vaa to be found via digest")

	// ARRANGE
	req = types.QueryExecutedVAA{InputType: "id", Input: "id"}

	// ACT
	resp, err = qs.ExecutedVAA(ctx, &req)

	// ASSERT
	require.NoError(t, err, "expected no error when the id is not in the archive")
	require.False(t, resp.Executed, "expected no vaa")

	// ARRANGE
	req = types.QueryExecutedVAA{InputType: "id", Input: vaa.MessageID()}

	// ACT
	resp, err = qs.ExecutedVAA(ctx, &req)

	// ASSERT
	require.NoError(t, err)
	require.True(t, resp.Executed, "expected the vaa to be found via id")
}
