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
	"encoding/json"

	"cosmossdk.io/collections"
	collcodec "cosmossdk.io/collections/codec"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/header"
	"cosmossdk.io/core/store"
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"

	"github.com/noble-assets/wormhole/types"
)

type Keeper struct {
	schema        collections.Schema
	headerService header.Service
	eventService  event.Service
	addressCodec  address.Codec

	Config           collections.Item[types.Config]
	WormchainChannel collections.Item[string]
	GuardianSets     collections.Map[uint32, types.GuardianSet]
	Sequences        collections.Map[[]byte, uint64]
	VAAArchive       *collections.IndexedMap[[]byte, collections.Pair[string, bool], VAAArchiveIndexes]

	ics4Wrapper types.ICS4Wrapper
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	headerService header.Service,
	eventService event.Service,
	addressCodec address.Codec,
	ics4Wrapper types.ICS4Wrapper,
) *Keeper {
	builder := collections.NewSchemaBuilder(storeService)

	keeper := &Keeper{
		headerService: headerService,
		eventService:  eventService,
		addressCodec:  addressCodec,

		Config:           collections.NewItem(builder, types.ConfigKey, "config", codec.CollValue[types.Config](cdc)),
		WormchainChannel: collections.NewItem(builder, types.WormchainChannelKey, "wormchain_channel", collections.StringValue),
		GuardianSets:     collections.NewMap(builder, types.GuardianSetPrefix, "guardian_sets", collections.Uint32Key, codec.CollValue[types.GuardianSet](cdc)),
		Sequences:        collections.NewMap(builder, types.SequencePrefix, "sequences", collections.BytesKey, collections.Uint64Value),
		VAAArchive: collections.NewIndexedMap(
			builder, types.VAAArchivePrefix, "vaa_archive",
			collections.BytesKey,
			collcodec.KeyToValueCodec(collections.PairKeyCodec(
				collections.StringKey, collections.BoolKey,
			)),
			NewVAAArchiveIndexes(builder),
		),

		ics4Wrapper: ics4Wrapper,
	}

	schema, err := builder.Build()
	if err != nil {
		panic(err)
	}

	keeper.schema = schema
	return keeper
}

// SetICS4Wrapper overrides the provided ICS4 wrapper for this module.
// This exists because IBC doesn't support dependency injection.
func (k *Keeper) SetICS4Wrapper(ics4Wrapper types.ICS4Wrapper) {
	k.ics4Wrapper = ics4Wrapper
}

// PostMessage allows the module to send messages from Noble via Wormhole.
func (k *Keeper) PostMessage(ctx context.Context, signer string, message []byte, nonce uint32) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	info := k.headerService.GetHeaderInfo(ctx)
	channel, err := k.WormchainChannel.Get(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get wormchain channel from state")
	}

	data, err := k.GetPacketData(ctx, message, nonce, signer)
	if err != nil {
		return err
	}
	bz, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = k.ics4Wrapper.SendPacket(
		sdkCtx, types.Port, channel,
		clienttypes.ZeroHeight(),
		uint64(info.Time.Add(types.PacketLifetime).UnixNano()),
		bz,
	)

	return err
}
