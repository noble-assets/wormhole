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
	"errors"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"

	"github.com/noble-assets/wormhole/types"
)

var _ types.QueryServer = &queryServer{}

type queryServer struct {
	*Keeper
}

func NewQueryServer(keeper *Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

func (k queryServer) Config(ctx context.Context, req *types.QueryConfig) (*types.QueryConfigResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	config, err := k.Keeper.Config.Get(ctx)
	if err != nil {
		return nil, errors.New("unable to get config from state")
	}

	return &types.QueryConfigResponse{Config: config}, nil
}

func (k queryServer) WormchainChannel(ctx context.Context, req *types.QueryWormchainChannel) (*types.QueryWormchainChannelResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	wormchainChannel, err := k.Keeper.WormchainChannel.Get(ctx)
	if err != nil {
		return nil, errors.New("wormchain channel not configured in state")
	}

	return &types.QueryWormchainChannelResponse{WormchainChannel: wormchainChannel}, nil
}

func (k queryServer) GuardianSets(ctx context.Context, req *types.QueryGuardianSets) (*types.QueryGuardianSetsResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	guardianSets, err := k.GetGuardianSets(ctx)
	if err != nil {
		return nil, errors.New("unable to get guardian sets from state")
	}

	return &types.QueryGuardianSetsResponse{GuardianSets: guardianSets}, nil
}

func (k queryServer) GuardianSet(ctx context.Context, req *types.QueryGuardianSet) (*types.QueryGuardianSetResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	index := uint32(0)
	if req.Index == "current" {
		config, err := k.Keeper.Config.Get(ctx)
		if err != nil {
			return nil, errors.New("unable to get config from state")
		}

		index = config.GuardianSetIndex
	} else {
		raw, err := strconv.Atoi(req.Index)
		if err != nil {
			return nil, fmt.Errorf("invalid guardian set index %s", req.Index)
		}

		index = uint32(raw)
	}

	guardianSet, err := k.Keeper.GuardianSets.Get(ctx, index)
	if err != nil {
		return nil, fmt.Errorf("unable to get guardian set %d from state", index)
	}

	return &types.QueryGuardianSetResponse{GuardianSet: guardianSet}, nil
}

func (k queryServer) ExecutedVAA(ctx context.Context, req *types.QueryExecutedVAA) (*types.QueryExecutedVAAResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	switch req.InputType {
	case "", "digest":
		digest := common.FromHex(req.Input)
		executed, _ := k.VAAArchive.Has(ctx, digest)

		return &types.QueryExecutedVAAResponse{Executed: executed}, nil
	case "id":
		digest, _ := k.VAAArchive.Indexes.ByID.MatchExact(ctx, req.Input)

		return &types.QueryExecutedVAAResponse{Executed: digest != nil}, nil
	default:
		return nil, fmt.Errorf("invalid input type '%s', expected 'digest' or 'id'", req.InputType)
	}
}
