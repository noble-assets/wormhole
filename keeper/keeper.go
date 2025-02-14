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
	"fmt"

	"cosmossdk.io/collections"
	collcodec "cosmossdk.io/collections/codec"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/header"
	"cosmossdk.io/core/store"
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"github.com/ethereum/go-ethereum/common"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"

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
	// TODO: what is the purpose of the bool?
	// VAAArchive is used to store the touple (VAA message ID, ???) indexed by the hash of the VAA.
	VAAArchive *collections.IndexedMap[[]byte, collections.Pair[string, bool], VAAArchiveIndexes]

	ics4Wrapper  types.ICS4Wrapper
	portKeeper   types.PortKeeper
	scopedKeeper types.ScopedKeeper
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	headerService header.Service,
	eventService event.Service,
	addressCodec address.Codec,
	ics4Wrapper types.ICS4Wrapper,
	portKeeper types.PortKeeper,
	scopedKeeper types.ScopedKeeper,
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

		ics4Wrapper:  ics4Wrapper,
		portKeeper:   portKeeper,
		scopedKeeper: scopedKeeper,
	}

	schema, err := builder.Build()
	if err != nil {
		panic(err)
	}

	keeper.schema = schema
	return keeper
}

// SetIBCKeepers overrides the provided IBC specific keepers for this module.
// This exists because IBC doesn't support dependency injection.
func (k *Keeper) SetIBCKeepers(ics4Wrapper types.ICS4Wrapper, portKeeper types.PortKeeper, scopedKeeper types.ScopedKeeper) {
	k.ics4Wrapper = ics4Wrapper
	k.portKeeper = portKeeper
	k.scopedKeeper = scopedKeeper
}

// BindPort allows the module to bind a specific port on initialization.
func (k *Keeper) BindPort(ctx sdk.Context) error {
	if _, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(types.Port)); !ok {
		capability := k.portKeeper.BindPort(ctx, types.Port)
		err := k.ClaimCapability(ctx, capability, host.PortPath(types.Port))
		if err != nil {
			return errors.Wrap(err, "could not claim port capability")
		}
	}

	return nil
}

// ClaimCapability allows the module to claim port or channel capabilities.
func (k *Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
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

	capability, _ := k.scopedKeeper.GetCapability(sdkCtx, host.ChannelCapabilityPath(types.Port, channel))
	_, err = k.ics4Wrapper.SendPacket(
		sdkCtx, capability, types.Port, channel,
		clienttypes.ZeroHeight(),
		uint64(info.Time.Add(types.PacketLifetime).UnixNano()),
		bz,
	)

	return err
}

func (k *Keeper) ParseAndVerifyVAA(ctx context.Context, bz []byte) (*vaautils.VAA, error) {
	vaa, err := vaautils.Unmarshal(bz)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidVAA, "failed to unmarshal: %v", err)
	}

	hash := vaa.SigningDigest().Bytes()
	if has, err := k.VAAArchive.Has(ctx, hash); err != nil || has {
		return nil, types.ErrAlreadyExecutedVAA
	}

	guardianSet, err := k.GuardianSets.Get(ctx, vaa.GuardianSetIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to get guardian set %d from state", vaa.GuardianSetIndex)
	}

	blockTime := uint64(k.headerService.GetHeaderInfo(ctx).Time.Unix())
	// TODO: is zero a no expiration?
	if guardianSet.ExpirationTime != 0 && guardianSet.ExpirationTime < blockTime {
		return nil, fmt.Errorf("guardian set %d is expired", vaa.GuardianSetIndex)
	}

	var addresses []common.Address
	for _, address := range guardianSet.Addresses {
		addresses = append(addresses, common.BytesToAddress(address))
	}
	if err := vaa.Verify(addresses); err != nil {
		return nil, errors.Wrap(err, "failed to verify vaa")
	}

	if err := k.VAAArchive.Set(ctx, hash, collections.Join(vaa.MessageID(), true)); err != nil {
		return nil, errors.Wrap(err, "failed to set vaa in state")
	}

	return vaa, nil
}
