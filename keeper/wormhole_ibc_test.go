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
	"encoding/hex"
	"strconv"
	"strings"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"github.com/stretchr/testify/require"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"

	"github.com/noble-assets/wormhole/types"
	"github.com/noble-assets/wormhole/utils"
	"github.com/noble-assets/wormhole/utils/mocks"
)

func TestHandleIBCReceiverGovernancePacket(t *testing.T) {
	// ARRANGE
	ctx, k := mocks.WormholeKeeper(t)
	packet := types.GovernancePacket{
		Action: uint8(vaautils.IbcReceiverActionUpdateChannelChain) + 1,
	}

	// ACT
	err := k.HandleIBCReceiverGovernancePacket(ctx, packet)

	// ASSERT
	require.Error(t, err, "expected error when governance action is not supported")
	require.ErrorContains(t, err, "unsupported governance action", "expected a different error")

	// ARRANGE: The action is valid but the payload is empty.
	packet.Action = uint8(vaautils.IbcReceiverActionUpdateChannelChain)

	// ACT
	err = k.HandleIBCReceiverGovernancePacket(ctx, packet)

	// ASSERT
	require.Error(t, err, "expected error when the payload is malformed")
	require.ErrorIs(t, err, types.ErrMalformedPayload, "expected a different error")

	// ARRANGE: Set a Chain different than wormchain in payload.
	channelBz, err := vaautils.LeftPadIbcChannelId("channel-0")
	require.NoError(t, err)
	// Shift left by eight for most significant byte and mask for less significant ones.
	chainIDBz := []byte{
		byte(uint16(vaautils.ChainIDNoble) >> 8), // NobleChainId
		byte(uint16(vaautils.ChainIDWormchain) & 0xFF),
	}
	invalidPayload := make([]byte, 66)
	copy(invalidPayload[:64], channelBz[:])
	copy(invalidPayload[64:], chainIDBz)
	packet.Payload = invalidPayload

	// ACT
	err = k.HandleIBCReceiverGovernancePacket(ctx, packet)

	// ASSERT
	require.Error(t, err, "expected error when the chain in the payload is not valid")
	require.ErrorIs(t, err, types.ErrInvalidChain, "expected a different error")

	// ARRANGE
	chainIDBz = []byte{
		byte(uint16(vaautils.ChainIDWormchain) >> 8),
		byte(uint16(vaautils.ChainIDWormchain) & 0xFF),
	}
	validPayload := make([]byte, 66)
	copy(validPayload[:64], channelBz[:])
	copy(validPayload[64:], chainIDBz)
	packet.Payload = validPayload

	// ACT
	err = k.HandleIBCReceiverGovernancePacket(ctx, packet)

	// ASSERT
	require.NoError(t, err, "expected no error when the payload is valid")
	channel, err := k.WormchainChannelId.Get(ctx)
	require.NoError(t, err, "expected non error retrieving the channel")
	require.Equal(t, "channel-0", channel, "expected a different channel")
}

func TestPostMessage_Keeper(t *testing.T) {
	// ARRANGE
	pk := mocks.PortKeeper{
		Ports: make(map[string]bool),
	}
	sk := mocks.ScopedKeeper{
		Capabilities: make(map[string]*capabilitytypes.Capability),
	}
	ics4w := mocks.ICS4Wrapper{}

	ctx, k := mocks.NewWormholeKeeper(t, ics4w, pk, sk)

	// ACT
	err := k.PostMessage(ctx, "", []byte{}, 0)

	// ASSERT
	require.Error(t, err, "expected an error when the wormchain channel is not set")
	require.ErrorContains(t, err, "failed to get wormchain", "expected a different error")

	// ARRANGE: Set the channel and cause an error with the packet data.
	err = k.WormchainChannelId.Set(ctx, "channel-0")
	require.NoError(t, err, "expecting no error setting the wormhole channel")

	// ACT
	err = k.PostMessage(ctx, "", []byte{}, 0)

	// ASSERT
	require.Error(t, err, "expected an error because config is not set")
	require.ErrorContains(t, err, "failed to get packet data")

	// ARRANGE: Set the config but no capabilities. The GetCapability mock returns nil.
	cfg := types.Config{ChainId: uint16(3)}
	err = k.Config.Set(ctx, cfg)
	require.NoError(t, err, "expected no error setting the config")

	// ACT
	signer := utils.TestAddress()
	err = k.PostMessage(ctx, signer.Bech32, []byte{}, 0)

	// ASSERT
	require.Error(t, err, "expected an error when capability is nil and send packet is called")
	require.ErrorContains(t, err, "failed to send packet")

	// ARRANGE: Set a valid capability in the store.
	sk.Capabilities[host.ChannelCapabilityPath(types.Port, "channel-0")] = &capabilitytypes.Capability{Index: uint64(3)}

	// ACT
	err = k.PostMessage(ctx, signer.Bech32, []byte("Hello from Noble"), 0)

	// ASSERT
	require.NoError(t, err)
}

func TestGetPacketData(t *testing.T) {
	// ARRANGE
	ctx, k := mocks.WormholeKeeper(t)

	// Variables not relevant for the test.
	message := []byte("Hello from Noble")
	nonce := uint32(0)

	// ACT
	_, err := k.GetPacketData(ctx, message, nonce, "")

	// ASSERT
	require.Error(t, err, "expected an error when the config is not set")
	require.ErrorContains(t, err, "failed to get config", "expected a different error")

	// ARRANGE: Set config.
	cfg := types.Config{
		ChainId: uint16(3),
	}
	err = k.Config.Set(ctx, cfg)
	require.NoError(t, err, "expected no error setting the config")

	// ACT
	_, err = k.GetPacketData(ctx, message, nonce, "")

	// ASSERT
	require.Error(t, err, "expected an error when signer is empty")
	require.ErrorContains(t, err, "failed to decode signer", "expected a different error")

	// ARRANGE: Create a signer that makes fail the codec.
	signer := utils.TestAddress()
	invalidSigner := strings.Join([]string{"cosmos", strings.Split(signer.Bech32, "noble")[1]}, "")

	// ACT
	_, err = k.GetPacketData(ctx, message, nonce, invalidSigner)

	// ASSERT
	require.Error(t, err, "expected an error when the address is not valid for the codec")
	require.ErrorContains(t, err, "failed to decode signer address", "expected a different error")

	// ACT: Call with a valid signer now.
	resp, err := k.GetPacketData(ctx, message, nonce, signer.Bech32)

	// ASSERT
	require.NoError(t, err, "expected no error when the signer is valid")

	emitter := make([]byte, 32)
	copy(emitter[12:], signer.Bytes)
	s, err := k.Sequences.Get(ctx, emitter)
	require.NoError(t, err, "expected no error getting the updated sequence")
	require.Equal(t, uint64(1), s, "expected 1 when is first sender packet")

	require.Len(t, resp.Publish.Msg, 6, "expected a different number of messages in the packet")
	require.Equal(t, hex.EncodeToString(message), resp.Publish.Msg[0].Value, "expected a different message")
	require.Equal(t, hex.EncodeToString(emitter), resp.Publish.Msg[1].Value, "expected a different emitter")
	require.Equal(t, "3", resp.Publish.Msg[2].Value, "expected a different chain ID")
	require.Equal(t, "0", resp.Publish.Msg[3].Value, "expected a different nonce")
	require.Equal(t, "0", resp.Publish.Msg[4].Value, "expected a different sequence")
	headerTime := sdk.UnwrapSDKContext(ctx).HeaderInfo().Time.Truncate(time.Second).Unix()
	require.Equal(t, strconv.Itoa(int(headerTime)), resp.Publish.Msg[5].Value, "expected a different timestamp")
}
