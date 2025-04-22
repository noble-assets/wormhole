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

	"github.com/noble-assets/wormhole/types"
	"github.com/noble-assets/wormhole/utils/mocks"
)

func TestBindPort(t *testing.T) {
	// ARRANGE
	pk := mocks.PortKeeper{
		Ports: make(map[string]bool),
	}
	sk := mocks.ScopedKeeper{
		Capabilities: make(map[string]*capabilitytypes.Capability),
	}
	ics4w := mocks.ICS4Wrapper{}

	ctx, k := mocks.NewWormholeKeeper(t, ics4w, pk, sk)

	// ACT: Bind to port when no capability are stored.
	err := k.BindPort(ctx)

	// ASSERT
	require.NoError(t, err, "expected no error when the capability associated with the port does not exist")

	c, found := sk.Capabilities[host.PortPath(types.Port)]
	require.True(t, found, "expected the capability to be in the state after binding")
	require.Equal(t, &capabilitytypes.Capability{Index: uint64(3)}, c, "expected a different index for capability") // 3 is hardcoded in the mock

	// ACT
	err = k.BindPort(ctx)

	// ASSERT
	require.NoError(t, err, "expected no error when the capability is already registered")

	// ARRANGE: Reset to initial conditions.
	delete(sk.Capabilities, host.PortPath(types.Port))
	// Setting the entry in the ports instructs the mock port to return a nil capability object.
	pk.Ports[types.Port] = true

	// ACT
	err = k.BindPort(ctx)

	// ASSERT
	require.Error(t, err, "expected when claim capability returns an error")
	require.ErrorContains(t, err, "could not claim port capability", "expected a different error")
}

func TestClaimCapability(t *testing.T) {
	// ARRANGE
	sk := mocks.ScopedKeeper{
		Capabilities: make(map[string]*capabilitytypes.Capability),
	}
	pk := mocks.PortKeeper{}
	ics4w := mocks.ICS4Wrapper{}

	// ARRANGE
	ctx, k := mocks.NewWormholeKeeper(t, ics4w, pk, sk)

	// ACT
	err := k.ClaimCapability(ctx, nil, "name")

	// ASSERT
	require.Error(t, err, "expected an error when the capability is nil")

	// ARRANGE: Create a valid capability object.
	capability := capabilitytypes.Capability{
		Index: 3,
	}

	// ACT
	err = k.ClaimCapability(ctx, &capability, "name")

	c, found := sk.Capabilities["name"]
	require.True(t, found, "expected the capability to be in the state")
	require.Equal(t, &capability, c, "expected a different index for capability")

	// ASSERT
	require.NoError(t, err)
}
