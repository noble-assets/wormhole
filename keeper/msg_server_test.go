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

	"github.com/noble-assets/wormhole/keeper"
	"github.com/noble-assets/wormhole/types"
	"github.com/noble-assets/wormhole/utils"
	"github.com/noble-assets/wormhole/utils/mocks"
)

func TestPostMessage_MsgServer(t *testing.T) {
	// ARRANGE
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
	require.Error(t, err, "expected an error with an empty message")
	require.Equal(t, &types.MsgPostMessageResponse{}, resp)

	// ARRANGE
	err = k.WormchainChannel.Set(ctx, "channel-0")
	require.NoError(t, err, "expecting no error setting the wormhole channel")

	sk.Capabilities[host.ChannelCapabilityPath(types.Port, "channel-0")] = &capabilitytypes.Capability{Index: uint64(3)}

	cfg := types.Config{
		ChainId: uint16(3),
	}
	err = k.Config.Set(ctx, cfg)
	require.NoError(t, err, "expected no error setting the config")

	signer := utils.TestAddress()

	msg.Message = []byte("Hello from Noble")
	msg.Nonce = 0
	msg.Signer = signer.Bech32

	// ACT
	resp, err = ms.PostMessage(ctx, &msg)

	// ASSERT
	require.NoError(t, err)
	require.Equal(t, &types.MsgPostMessageResponse{}, resp)
}
