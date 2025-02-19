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
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"

	"github.com/noble-assets/wormhole/types"
)

// SetIBCKeepers overrides the provided IBC specific keepers for this module.
// This exists because IBC doesn't support dependency injection.
func (k *Keeper) SetIBCKeepers(ics4Wrapper types.ICS4Wrapper, portKeeper types.PortKeeper, scopedKeeper types.ScopedKeeper) {
	k.ics4Wrapper = ics4Wrapper
	k.portKeeper = portKeeper
	k.scopedKeeper = scopedKeeper
}

// BindPort allows the module to bind a specific port on initialization.
func (k *Keeper) BindPort(ctx sdk.Context) error {
	if _, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(types.Port)); !ok {
		capability := k.portKeeper.BindPort(ctx, types.Port)
		err := k.ClaimCapability(ctx, capability, host.PortPath(types.Port))
		if err != nil {
			return errors.Wrap(err, "could not claim port capability")
		}
	}

	return nil
}

// ClaimCapability allows the module to claim port or channel capabilities.
func (k *Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}
