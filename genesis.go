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

package wormhole

import (
	"encoding/hex"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/noble-assets/wormhole/keeper"
	"github.com/noble-assets/wormhole/types"
)

func InitGenesis(ctx sdk.Context, k *keeper.Keeper, cdc address.Codec, genesis types.GenesisState) {
	if err := k.Config.Set(ctx, genesis.Config); err != nil {
		panic(err)
	}

	if err := k.WormchainChannel.Set(ctx, genesis.WormchainChannel); err != nil {
		panic(err)
	}

	for index, guardianSet := range genesis.GuardianSets {
		if err := k.GuardianSets.Set(ctx, index, guardianSet); err != nil {
			panic(err)
		}
	}

	for address, sequence := range genesis.Sequences {
		sender := make([]byte, 32)
		bz, err := cdc.StringToBytes(address)
		if err != nil {
			panic(err)
		}
		copy(sender[12:], bz)

		if err := k.Sequences.Set(ctx, sender, sequence); err != nil {
			panic(err)
		}
	}

	for hash, id := range genesis.VaaArchive {
		bz, err := hex.DecodeString(hash)
		if err != nil {
			panic(err)
		}

		if err := k.VAAArchive.Set(ctx, bz, collections.Join(id, true)); err != nil {
			panic(err)
		}
	}
}

func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	config, err := k.Config.Get(ctx)
	if err != nil {
		panic(err)
	}

	wormchainChannel, err := k.WormchainChannel.Get(ctx)
	if err != nil {
		panic(err)
	}

	guardianSets, err := k.GetGuardianSets(ctx)
	if err != nil {
		panic(err)
	}

	sequences, err := k.GetSequences(ctx)
	if err != nil {
		panic(err)
	}

	vaaArchive, err := k.GetVAAArchive(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Config:           config,
		WormchainChannel: wormchainChannel,
		GuardianSets:     guardianSets,
		Sequences:        sequences,
		VaaArchive:       vaaArchive,
	}
}
