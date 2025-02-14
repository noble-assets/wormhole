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
	"strconv"

	"cosmossdk.io/core/event"
	"cosmossdk.io/errors"
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

		if err := k.WormchainChannel.Set(ctx, string(updateChannelChain.ChannelID)); err != nil {
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

func (k *Keeper) GetPacketData(ctx context.Context, message []byte, nonce uint32, signer string) (*types.PacketData, error) {
	config, err := k.Config.Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config from state")
	}

	emitter := make([]byte, 32)
	bz, err := k.addressCodec.StringToBytes(signer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode signer address")
	}
	copy(emitter[12:], bz)

	sequence, _ := k.Sequences.Get(ctx, emitter)
	err = k.Sequences.Set(ctx, emitter, sequence+1)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set sequence in state")
	}

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
				{
					Key:   "message.message",
					Value: hex.EncodeToString(message),
				},
				{
					Key:   "message.sender",
					Value: hex.EncodeToString(emitter),
				},
				{
					Key:   "message.chain_id",
					Value: strconv.Itoa(int(config.ChainId)),
				},
				{
					Key:   "message.nonce",
					Value: strconv.Itoa(int(nonce)),
				},
				{
					Key:   "message.sequence",
					Value: strconv.Itoa(int(sequence)),
				},
				{
					Key:   "message.block_time",
					Value: strconv.Itoa(int(k.headerService.GetHeaderInfo(ctx).Time.Unix())),
				},
			},
		}),
	}, nil
}
