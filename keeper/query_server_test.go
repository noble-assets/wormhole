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

	"github.com/stretchr/testify/require"

	"github.com/noble-assets/wormhole/keeper"
	"github.com/noble-assets/wormhole/types"
	"github.com/noble-assets/wormhole/utils/mocks"
)

func TestConfig(t *testing.T) {
	// ARRANGE: Create environment and nil request.
	ctx, k := mocks.WormholeKeeper(t)
	qs := keeper.NewQueryServer(k)

	// ACT
	resp, err := qs.Config(ctx, nil)

	// ASSERT
	require.Error(t, err, "expected error when the request is nil")
	require.ErrorIs(t, err, types.ErrInvalidRequest, "expected a different error")
	require.Nil(t, resp, "expected nil response")

	// ARRANGE: Request not nil, but empty state
	req := types.QueryConfig{}

	// ACT
	resp, err = qs.Config(ctx, &req)

	// ASSERT
	require.Error(t, err, "expected error when the request is nil")
	require.ErrorContains(t, err, "unable to get config", "expected a different error")
	require.Nil(t, resp, "expected nil response")

	// ARRANGE: Add config to the state.
	cfg := types.Config{
		ChainId:          1,
		GuardianSetIndex: 2,
		GovChain:         3,
		GovAddress:       []byte("4"),
	}
	err = k.Config.Set(ctx, cfg)
	require.NoError(t, err, "expected no error setting the config")

	// ACT
	resp, err = qs.Config(ctx, &req)

	// ASSERT
	require.NoError(t, err, "expected no error when the request is valid and config exists")
	require.Equal(t, cfg.ChainId, resp.Config.ChainId, "expected a different chain id")
	require.Equal(t, cfg.GuardianSetIndex, resp.Config.GuardianSetIndex, "expected a different guardian set index")
	require.Equal(t, cfg.GovChain, resp.Config.GovChain, "expected a different gov chain")
	require.Equal(t, cfg.GovAddress, resp.Config.GovAddress, "expected a different gov address")
}
