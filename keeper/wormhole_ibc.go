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
	"encoding/json"
	"strconv"

	"cosmossdk.io/core/event"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"

	"github.com/noble-assets/wormhole/types"
)

func (k *Keeper) HandleIBCReceiverGovernancePacket(ctx context.Context, pkt types.GovernancePacket) error {
	switch pkt.Action {
	case uint8(vaautils.IbcReceiverActionUpdateChannelChain):
		var updateChannelChain types.UpdateChannelChain
		err := updateChannelChain.Parse(pkt.Payload)
		if err != nil {
			return err
		}

		if updateChannelChain.Chain != uint16(vaautils.ChainIDWormchain) {
			return types.ErrInvalidChain
		}

		if err := k.WormchainChannelId.Set(ctx, string(updateChannelChain.ChannelID)); err != nil {
			return errors.Wrap(err, "failed to set wormchain channel in state")
		}

		return k.eventService.EventManager(ctx).EmitKV(ctx, "UpdateChannelChain",
			event.Attribute{Key: "chain_id", Value: strconv.Itoa(int(updateChannelChain.Chain))},
			event.Attribute{Key: "channel_id", Value: string(updateChannelChain.ChannelID)},
		)
	default:
		return errors.Wrapf(types.ErrUnsupportedGovernanceAction, "module: %s, type: %d", pkt.Module, pkt.Action)
	}
}

// PostMessage allows the module to send messages from Noble via Wormhole.
func (k *Keeper) PostMessage(ctx context.Context, signer string, message []byte, nonce uint32) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	info := k.headerService.GetHeaderInfo(ctx)

	channel, err := k.WormchainChannelId.Get(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get wormchain channel from state")
	}

	data, err := k.GetPacketData(ctx, message, nonce, signer)
	if err != nil {
		return errors.Wrap(err, "failed to get packet data")
	}
	bz, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshad packed data")
	}

	capability, _ := k.scopedKeeper.GetCapability(sdkCtx, host.ChannelCapabilityPath(types.Port, channel))

	_, err = k.ics4Wrapper.SendPacket(
		sdkCtx, capability, types.Port, channel,
		clienttypes.ZeroHeight(),
		uint64(info.Time.Add(types.PacketLifetime).UnixNano()),
		bz,
	)

	return errors.Wrap(err, "failed to send packet")
}

func (k *Keeper) GetPacketData(ctx context.Context, message []byte, nonce uint32, signer string) (*types.PacketData, error) {
	config, err := k.Config.Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config from state")
	}

	bz, err := k.addressCodec.StringToBytes(signer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode signer address")
	}
	emitter := make([]byte, 32)
	copy(emitter[12:], bz)

	sequence, _ := k.Sequences.Get(ctx, emitter)
	if err := k.Sequences.Set(ctx, emitter, sequence+1); err != nil {
		return nil, errors.Wrap(err, "failed to set sequence in state")
	}

	pkt := CreatePacket(message, emitter, config.ChainId, nonce, sequence, k.headerService.GetHeaderInfo(ctx).Time.Unix())

	return pkt, nil
}

func CreatePacket(message, sender []byte, chainId uint16, nonce uint32, sequence uint64, timestamp int64) *types.PacketData {
	return &types.PacketData{
		Publish: struct {
			Msg []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"msg"`
		}(struct {
			Msg []struct {
				Key   string
				Value string
			}
		}{
			Msg: []struct {
				Key   string
				Value string
			}{
				{"message.message", hex.EncodeToString(message)},
				{"message.sender", hex.EncodeToString(sender)},
				{"message.chain_id", strconv.Itoa(int(chainId))},
				{"message.nonce", strconv.Itoa(int(nonce))},
				{"message.sequence", strconv.Itoa(int(sequence))},
				{"message.block_time", strconv.Itoa(int(timestamp))},
			},
		}),
	}
}
