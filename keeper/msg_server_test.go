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

	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"github.com/stretchr/testify/require"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"

	"github.com/noble-assets/wormhole/keeper"
	"github.com/noble-assets/wormhole/types"
	"github.com/noble-assets/wormhole/utils"
	"github.com/noble-assets/wormhole/utils/mocks"
)

func TestSubmitVAA(t *testing.T) {
	// ARRANGE
	pk := mocks.PortKeeper{
		Ports: make(map[string]bool),
	}
	sk := mocks.ScopedKeeper{
		Capabilities: make(map[string]*capabilitytypes.Capability),
	}
	ics4w := mocks.ICS4Wrapper{}

	ctx, k := mocks.NewWormholeKeeper(t, ics4w, pk, sk)

	ms := keeper.NewMsgServer(k)

	msg := types.MsgSubmitVAA{}

	// ACT
	resp, err := ms.SubmitVAA(ctx, &msg)

	// ASSERT
	require.Error(t, err, "expected an error with an invalid message during parsing")
	require.ErrorContains(t, err, "failed to parse and verify vaa")
	require.Nil(t, resp, "expected nil response")

	// ARRANGE: Set the test to pass the parse and verification of the VAA
	signer := utils.TestAddress()
	guardian := utils.GuardianSigner()

	vaaBody := utils.VAABody{
		GuardianSetIndex: 0,
		Payload:          []byte("test vaa"),
		Sequence:         1,
		EmitterChain:     0,
		EmitterAddress:   [32]byte{},
	}
	vaa := utils.CreateVAA(t, []utils.Guardian{guardian}, vaaBody)
	bzVaa, err := vaa.Marshal()
	require.NoError(t, err, "expected no error marshaling the vaa")

	msg.Vaa = bzVaa
	msg.Signer = signer.Bech32

	guardianSet := types.GuardianSet{
		Addresses:      [][]byte{guardian.Address[:]},
		ExpirationTime: 0,
	}
	err = k.GuardianSets.Set(ctx, 0, guardianSet)
	require.NoError(t, err, "expected no error setting the guardian set")

	// ACT
	_, err = ms.SubmitVAA(ctx, &msg)

	// ASSERT
	require.Error(t, err, "expected an error when the config is not set")
	require.ErrorContains(t, err, "failed to get config from state")

	// ARRANGE: Set a valid config and clean the archive from previous VAA.
	cfg := types.Config{ChainId: uint16(2)}
	err = k.Config.Set(ctx, cfg)
	require.NoError(t, err, "expected no error setting the config")

	hash := vaa.SigningDigest().Bytes()
	err = k.VAAArchive.Remove(ctx, hash)
	require.NoError(t, err, "expected no error resetting vaa archive to empty")

	// ACT
	_, err = ms.SubmitVAA(ctx, &msg)

	// ASSERT
	require.Error(t, err, "expected an error when the config does not have a valid gov chain")
	require.ErrorIs(t, err, types.ErrNotGovernanceVAA)

	// ARRANGE: Set a valid ChainID but not a valid emitter address.
	vaaBody.EmitterChain = vaautils.ChainID(3)

	vaa = utils.CreateVAA(t, []utils.Guardian{guardian}, vaaBody)
	bzVaa, err = vaa.Marshal()
	require.NoError(t, err, "expected no error marshaling the vaa")
	msg.Vaa = bzVaa

	// ACT
	_, err = ms.SubmitVAA(ctx, &msg)

	// ASSERT
	require.Error(t, err, "expected an error when emitter address is not gov address")
	require.ErrorIs(t, err, types.ErrNotGovernanceVAA)

	// ARRANGE: Set a valid emitter address but different than the gov one.
	cfg.GovAddress = []byte("address") // not padded with zero
	cfg.GovChain = 3
	err = k.Config.Set(ctx, cfg)
	require.NoError(t, err, "expected no error setting the config")

	vaaBody.EmitterAddress = vaautils.Address([]byte("0000000000000000000000000address"))
	vaa = utils.CreateVAA(t, []utils.Guardian{guardian}, vaaBody)
	bzVaa, err = vaa.Marshal()
	require.NoError(t, err, "expected no error marshaling the vaa")
	msg.Vaa = bzVaa

	// ACT
	_, err = ms.SubmitVAA(ctx, &msg)

	// ASSERT
	require.Error(t, err, "expected an error when emitter address is not the gov address")
	require.ErrorIs(t, err, types.ErrNotGovernanceVAA)

	// ARRANGE: Set a guardian set index in the config different than the one in the VAA and
	// clean the vaa archive.
	hash = vaa.SigningDigest().Bytes()
	err = k.VAAArchive.Remove(ctx, hash)
	require.NoError(t, err, "expected no error resetting vaa archive to empty")

	cfg.GovAddress = []byte("0000000000000000000000000address")
	cfg.GuardianSetIndex = 99
	err = k.Config.Set(ctx, cfg)
	require.NoError(t, err, "expected no error setting the config")

	// ACT
	_, err = ms.SubmitVAA(ctx, &msg)

	// ASSERT
	require.Error(t, err, "expected an error when guardian set index is not the same")
	require.ErrorContains(t, err, "must be signed by current guardian set")

	// ARRANG: Reset the vaa archive and set the correct guardian set index.
	err = k.VAAArchive.Remove(ctx, hash)
	require.NoError(t, err, "expected no error resetting vaa archive to empty")

	cfg.GuardianSetIndex = 0
	err = k.Config.Set(ctx, cfg)
	require.NoError(t, err, "expected no error setting the config")

	// ACT
	_, err = ms.SubmitVAA(ctx, &msg)

	// ASSERT
	require.Error(t, err, "expected an error when the payload is not a valid governance packet")
	require.ErrorContains(t, err, "failed parsing the vaa payload")

	// ARRANGE
	err = k.VAAArchive.Remove(ctx, hash)
	require.NoError(t, err, "expected no error resetting vaa archive to empty")

	packet := types.GovernancePacket{
		Action: 3,
		Module: "Core",
		Chain:  3,
	}
	packetBz := packet.Serialize()

	vaaBody.Payload = packetBz
	vaa = utils.CreateVAA(t, []utils.Guardian{guardian}, vaaBody)
	bzVaa, err = vaa.Marshal()
	require.NoError(t, err, "expected no error marshaling the vaa")
	msg.Vaa = bzVaa

	// ACT
	_, err = ms.SubmitVAA(ctx, &msg)

	// ASSERT
	require.Error(t, err, "expected an error when the config chain id is not the same of vaa")
	require.ErrorContains(t, err, "packet not meant for this chain")

	// ARRANGE: Set the Chain id to zero to have the chain id check to always pass
	hash = vaa.SigningDigest().Bytes()
	err = k.VAAArchive.Remove(ctx, hash)
	require.NoError(t, err, "expected no error resetting vaa archive to empty")

	packet = types.GovernancePacket{
		Action: 3,
		Module: "Invalid",
		Chain:  0,
	}
	packetBz = packet.Serialize()

	vaaBody.Payload = packetBz
	vaa = utils.CreateVAA(t, []utils.Guardian{guardian}, vaaBody)
	bzVaa, err = vaa.Marshal()
	require.NoError(t, err, "expected no error marshaling the vaa")
	msg.Vaa = bzVaa

	// ACT
	resp, err = ms.SubmitVAA(ctx, &msg)

	// ASSERT
	require.Error(t, err, "expected an error when the governance action is not valid")
	require.NotContains(t, err.Error(), "packet not meant for this chain", "expected passing chain check when packet chain is 0")
	require.ErrorContains(t, err, "unsupported governance action")
	require.Equal(t, &types.MsgSubmitVAAResponse{}, resp, "expected a different response")

	// ARRANGE: Trigger a failure in the execution of the core gov.
	hash = vaa.SigningDigest().Bytes()
	err = k.VAAArchive.Remove(ctx, hash)
	require.NoError(t, err, "expected no error resetting vaa archive to empty")

	packet = types.GovernancePacket{
		Action: 3,
		Module: "Core",
		Chain:  0,
	}
	packetBz = packet.Serialize()

	vaaBody.Payload = packetBz
	vaa = utils.CreateVAA(t, []utils.Guardian{guardian}, vaaBody)
	bzVaa, err = vaa.Marshal()
	require.NoError(t, err, "expected no error marshaling the vaa")
	msg.Vaa = bzVaa

	// ACT
	resp, err = ms.SubmitVAA(ctx, &msg)

	// ASSERT
	require.Error(t, err, "expected an error in the core handler")
	require.ErrorContains(t, err, "failed handling the core", "expected failing handling the core governance packet")
	require.Equal(t, &types.MsgSubmitVAAResponse{}, resp, "expected a different response")

	// ARRANGE: Valid VAA submission via Core module
	hash = vaa.SigningDigest().Bytes()
	err = k.VAAArchive.Remove(ctx, hash)
	require.NoError(t, err, "expected no error resetting vaa archive to empty")

	payload := []byte{}
	payload = append(payload, []byte{0x00, 0x00, 0x00, 0x01}...)
	payload = append(payload, 0x01)
	payload = append(payload, guardian.Address.Bytes()...)

	packet = types.GovernancePacket{
		Action:  uint8(vaautils.ActionGuardianSetUpdate),
		Module:  "Core",
		Chain:   0,
		Payload: payload,
	}
	packetBz = packet.Serialize()

	vaaBody.Payload = packetBz
	vaa = utils.CreateVAA(t, []utils.Guardian{guardian}, vaaBody)
	bzVaa, err = vaa.Marshal()
	require.NoError(t, err, "expected no error marshaling the vaa")
	msg.Vaa = bzVaa

	// ACT
	resp, err = ms.SubmitVAA(ctx, &msg)

	// ASSERT
	require.NoError(t, err, "expected no error in the core handler")
	require.Equal(t, &types.MsgSubmitVAAResponse{}, resp, "expected a different response")

	// ARRANGE: Trigger a failure in the execution of the ibc gov. Here we have to update the
	// guardian set index to reflect previous vaa execution.
	hash = vaa.SigningDigest().Bytes()
	err = k.VAAArchive.Remove(ctx, hash)
	require.NoError(t, err, "expected no error resetting vaa archive to empty")

	packet = types.GovernancePacket{
		Action: 3,
		Module: "IbcReceiver",
		Chain:  0,
	}
	packetBz = packet.Serialize()

	vaaBody.Payload = packetBz
	vaaBody.GuardianSetIndex = 1
	vaa = utils.CreateVAA(t, []utils.Guardian{guardian}, vaaBody)
	bzVaa, err = vaa.Marshal()
	require.NoError(t, err, "expected no error marshaling the vaa")
	msg.Vaa = bzVaa

	// ACT
	resp, err = ms.SubmitVAA(ctx, &msg)

	// ASSERT
	require.Error(t, err, "expected an error in the ibc receiver handler")
	require.ErrorContains(t, err, "failed handling the ibc", "expected failing handling the ibc receive governance packet")
	require.Equal(t, &types.MsgSubmitVAAResponse{}, resp, "expected a different response")

	// ARRANGE: Valid VAA submission via IBC module
	hash = vaa.SigningDigest().Bytes()
	err = k.VAAArchive.Remove(ctx, hash)
	require.NoError(t, err, "expected no error resetting vaa archive to empty")

	channelBz, err := vaautils.LeftPadIbcChannelId("channel-0")
	require.NoError(t, err)
	chainIDBz := []byte{
		byte(uint16(vaautils.ChainIDWormchain) >> 8),
		byte(uint16(vaautils.ChainIDWormchain) & 0xFF),
	}

	payload = make([]byte, 66)
	copy(payload[:64], channelBz[:])
	copy(payload[64:], chainIDBz)

	packet = types.GovernancePacket{
		Action:  uint8(vaautils.IbcReceiverActionUpdateChannelChain),
		Module:  "IbcReceiver",
		Chain:   0,
		Payload: payload,
	}
	packetBz = packet.Serialize()

	vaaBody.Payload = packetBz
	vaa = utils.CreateVAA(t, []utils.Guardian{guardian}, vaaBody)
	bzVaa, err = vaa.Marshal()
	require.NoError(t, err, "expected no error marshaling the vaa")
	msg.Vaa = bzVaa

	// ACT
	resp, err = ms.SubmitVAA(ctx, &msg)

	// ASSERT
	require.NoError(t, err, "expected an error in the ibc receiver handler")
	require.Equal(t, &types.MsgSubmitVAAResponse{}, resp, "expected a different response")
}

func TestPostMessage_MsgServer(t *testing.T) {
	// ARRANGE
	pk := mocks.PortKeeper{
		Ports: make(map[string]bool),
	}
	sk := mocks.ScopedKeeper{
		Capabilities: make(map[string]*capabilitytypes.Capability),
	}
	ics4w := mocks.ICS4Wrapper{}

	ctx, k := mocks.NewWormholeKeeper(t, ics4w, pk, sk)

	ms := keeper.NewMsgServer(k)

	msg := types.MsgPostMessage{}

	// ACT
	resp, err := ms.PostMessage(ctx, &msg)

	// ASSERT
	require.Error(t, err, "expected an error the message has invalid values")
	require.Equal(t, &types.MsgPostMessageResponse{}, resp)

	// ARRANGE: Setting the values required to make the keeper method called under the hood working.
	err = k.WormchainChannelId.Set(ctx, "channel-0")
	require.NoError(t, err, "expecting no error setting the wormhole channel")

	capability := capabilitytypes.Capability{Index: uint64(3)}
	sk.Capabilities[host.ChannelCapabilityPath(types.Port, "channel-0")] = &capability

	cfg := types.Config{ChainId: uint16(3)}
	err = k.Config.Set(ctx, cfg)
	require.NoError(t, err, "expected no error setting the config")

	msg.Message = []byte("Hello from Noble")
	msg.Nonce = 0
	msg.Signer = utils.TestAddress().Bech32

	// ACT
	resp, err = ms.PostMessage(ctx, &msg)

	// ASSERT
	require.NoError(t, err)
	require.Equal(t, &types.MsgPostMessageResponse{}, resp)
}
