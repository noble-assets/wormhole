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
	"bytes"
	"context"

	"cosmossdk.io/errors"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"

	"github.com/noble-assets/wormhole/types"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) SubmitVAA(ctx context.Context, msg *types.MsgSubmitVAA) (*types.MsgSubmitVAAResponse, error) {
	vaa, err := k.ParseAndVerifyVAA(ctx, msg.Vaa)
	if err != nil {
		return nil, errors.Wrap(err, "failed during vaa parsing and verification")
	}

	config, err := k.Config.Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config from state")
	}

	if vaa.EmitterChain != vaautils.ChainID(config.GovChain) || !bytes.Equal(vaa.EmitterAddress.Bytes(), config.GovAddress) {
		// TODO: is this the error we want to return?
		return nil, types.ErrNotGovernanceVAA
	}
	if vaa.GuardianSetIndex != config.GuardianSetIndex {
		return nil, errors.Wrap(types.ErrInvalidGovernanceVAA, "must be signed by current guardian set")
	}

	var pkt types.GovernancePacket
	err = pkt.Parse(vaa.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed parsing the vaa payload")
	}

	if pkt.Chain != config.ChainId && pkt.Chain != 0 {
		return nil, errors.Wrap(types.ErrInvalidGovernanceVAA, "packet not meant for this chain")
	}

	switch pkt.Module {
	case "Core":
		err = k.HandleCoreGovernancePacket(ctx, pkt)
		if err != nil {
			err = errors.Wrap(err, "failed handling the core governance packet")
		}
		return &types.MsgSubmitVAAResponse{}, err
	case "IbcReceiver":
		err = k.HandleIBCReceiverGovernancePacket(ctx, pkt)
		if err != nil {
			err = errors.Wrap(err, "failed handling the ibc receive governance packet")
		}
		return &types.MsgSubmitVAAResponse{}, err
	default:
		return &types.MsgSubmitVAAResponse{}, errors.Wrapf(types.ErrUnsupportedGovernanceAction, "module: %s, type: %d", pkt.Module, pkt.Action)
	}
}

func (k msgServer) PostMessage(ctx context.Context, msg *types.MsgPostMessage) (*types.MsgPostMessageResponse, error) {
	err := k.Keeper.PostMessage(ctx, msg.Signer, msg.Message, msg.Nonce)

	return &types.MsgPostMessageResponse{}, err
}
