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

	"github.com/stretchr/testify/require"

	"github.com/noble-assets/wormhole/types"
	"github.com/noble-assets/wormhole/utils"
	"github.com/noble-assets/wormhole/utils/mocks"
)

func TestGetPacketData(t *testing.T) {
	// ARRANGE: Set default variable and does not initialize the state.
	ctx, k := mocks.WormholeKeeper(t)
	message := []byte{}
	nonce := uint32(0)

	// ACT
	_, err := k.GetPacketData(ctx, message, nonce, "")

	// ASSERT
	require.Error(t, err, "expected an error")
	require.ErrorContains(t, err, "failed to get config", "expected a different error")

	// ARRANGE: Set empty config
	cfg := types.Config{}
	err = k.Config.Set(ctx, cfg)
	require.NoError(t, err, "expected no error setting the config")

	// ACT
	_, err = k.GetPacketData(ctx, message, nonce, "")

	// ASSERT
	require.Error(t, err, "expected an error when signer is not valid")
	require.ErrorContains(t, err, "failed to decode signer", "expected a different error")

	// ARRANGE: Create an invalid signer
	signer := utils.TestAddress()
	invalidSigner := strings.Join([]string{"cosmos", strings.Split(signer.Bech32, "noble")[1]}, "")

	// ACT
	_, err = k.GetPacketData(ctx, message, nonce, invalidSigner)

	// ASSERT
	require.Error(t, err, "expected an error when the address is not valid for the codec")
	require.ErrorContains(t, err, "failed to decode signer address", "expected a different error")

	// ARRANGE: Add more information to the state to better test the valid case
	cfg = types.Config{
		ChainId: uint16(3),
	}
	err = k.Config.Set(ctx, cfg)
	require.NoError(t, err, "expected no error setting the config")

	// ACT: Call with valid signer now
	message = []byte("Hello from Noble")
	resp, err := k.GetPacketData(ctx, message, nonce, signer.Bech32)

	// ASSERT
	require.NoError(t, err, "expected no error when the signer is valid")

	emitter := make([]byte, 32)
	copy(emitter[12:], signer.Bytes)
	s, err := k.Sequences.Get(ctx, emitter)
	require.NoError(t, err, "expected no error getting the updated sequence")
	require.Equal(t, uint64(1), s, "expected 1 for a previously not existent key")

	require.Len(t, resp.Publish.Msg, 6, "expected a different number of messages")
	require.Equal(t, hex.EncodeToString(message), resp.Publish.Msg[0].Value, "expected a different message")
	require.Equal(t, hex.EncodeToString(emitter), resp.Publish.Msg[1].Value, "expected a different emitter")
	require.Equal(t, "3", resp.Publish.Msg[2].Value, "expected a different chain ID")
	require.Equal(t, "0", resp.Publish.Msg[3].Value, "expected a different nonce")
	require.Equal(t, "0", resp.Publish.Msg[4].Value, "expected a different sequence")
	require.Equal(t, strconv.Itoa(int(time.Now().Truncate(time.Second).Unix())), resp.Publish.Msg[5].Value, "expected a different timestamp")
}
