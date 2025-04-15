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

package types

import "cosmossdk.io/errors"

var (
	ErrInvalidRequest = errors.Register(ModuleName, 0, "invalid request")

	ErrInvalidVAA                  = errors.Register(ModuleName, 1, "invalid vaa")
	ErrAlreadyExecutedVAA          = errors.Register(ModuleName, 2, "vaa already executed")
	ErrNotGovernanceVAA            = errors.Register(ModuleName, 3, "not governance vaa")
	ErrInvalidGovernanceVAA        = errors.Register(ModuleName, 4, "invalid governance vaa")
	ErrUnsupportedGovernanceAction = errors.Register(ModuleName, 5, "unsupported governance action")
	ErrMalformedPayload            = errors.New(ModuleName, 6, "payload is malformed")

	ErrInvalidPort    = errors.New(ModuleName, 101, "invalid port")
	ErrInvalidVersion = errors.New(ModuleName, 102, "invalid version")
	ErrInvalidChain   = errors.New(ModuleName, 103, "chain id is not associated with wormchain")
)
