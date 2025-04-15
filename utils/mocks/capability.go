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

package mocks

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"

	"github.com/noble-assets/wormhole/types"
)

var _ types.ScopedKeeper = ScopedKeeper{}

type ScopedKeeper struct {
	Capabilities map[string]*capabilitytypes.Capability
}

// ClaimCapability implements types.ScopedKeeper.
func (s ScopedKeeper) ClaimCapability(ctx sdk.Context, capability *capabilitytypes.Capability, name string) error {
	if capability == nil {
		return errors.New("empty capability")
	}

	if name == "" {
		return errors.New("empty name")
	}

	s.Capabilities[name] = capability
	return nil
}

// GetCapability implements types.ScopedKeeper
func (s ScopedKeeper) GetCapability(ctx sdk.Context, name string) (*capabilitytypes.Capability, bool) {
	capability, found := s.Capabilities[name]

	if !found {
		return nil, false
	}

	return capability, true
}
