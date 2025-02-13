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

package keeper

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"

	"github.com/noble-assets/wormhole/types"
)

// GetChain is a helper function for retrieving the local Wormhole Chain ID.
func (k *Keeper) GetChain(ctx context.Context) (uint16, error) {
	config, err := k.Config.Get(ctx)
	if err != nil {
		return 0, errors.New("unable to get config from state")
	}

	return config.ChainId, nil
}

// GetGuardianSets is a helper function for retrieving all guardian sets from state.
func (k *Keeper) GetGuardianSets(ctx context.Context) (map[uint32]types.GuardianSet, error) {
	guardianSets := make(map[uint32]types.GuardianSet)

	err := k.GuardianSets.Walk(ctx, nil, func(index uint32, guardianSet types.GuardianSet) (stop bool, err error) {
		guardianSets[index] = guardianSet
		return false, nil
	})

	return guardianSets, err
}

// GetSequences is a helper function for retrieving all sequence entries from
// state. It encodes senders as Bech32 addresses for use in a genesis export.
func (k *Keeper) GetSequences(ctx context.Context) (map[string]uint64, error) {
	sequences := make(map[string]uint64)

	err := k.Sequences.Walk(ctx, nil, func(sender []byte, sequence uint64) (stop bool, err error) {
		if len(sender) < 20 {
			return false, fmt.Errorf("address with less than 20 bytes: %s", sender)
		}

		address, err := k.addressCodec.BytesToString(sender[12:])
		if err != nil {
			// NOTE: We continue in the case of an encoding error.
			return false, err
		}

		sequences[address] = sequence
		return false, nil
	})

	return sequences, err
}

// GetVAAArchive is a helper function for retrieving all executed VAAs from
// state. Note that this should only be used when exporting genesis.
func (k *Keeper) GetVAAArchive(ctx context.Context) (map[string]string, error) {
	vaaArchive := make(map[string]string)

	err := k.VAAArchive.Walk(ctx, nil, func(hash []byte, value collections.Pair[string, bool]) (stop bool, err error) {
		if value.K2() {
			vaaArchive[hex.EncodeToString(hash)] = value.K1()
		}

		return false, nil
	})

	return vaaArchive, err
}

//

type VAAArchiveIndexes struct {
	ByID *indexes.Unique[string, []byte, collections.Pair[string, bool]]
}

func (i VAAArchiveIndexes) IndexesList() []collections.Index[[]byte, collections.Pair[string, bool]] {
	return []collections.Index[[]byte, collections.Pair[string, bool]]{i.ByID}
}

func NewVAAArchiveIndexes(builder *collections.SchemaBuilder) VAAArchiveIndexes {
	return VAAArchiveIndexes{
		ByID: indexes.NewUnique(
			builder, types.VAAByIDPrefix, "vaa_by_id",
			collections.StringKey, collections.BytesKey,
			func(_ []byte, value collections.Pair[string, bool]) (string, error) {
				return value.K1(), nil
			},
		),
	}
}
