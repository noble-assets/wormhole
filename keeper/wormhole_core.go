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
	"fmt"
	"strconv"

	"cosmossdk.io/core/event"
	"cosmossdk.io/errors"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"

	"github.com/noble-assets/wormhole/types"
)

func (k *Keeper) HandleCoreGovernancePacket(ctx context.Context, pkt types.GovernancePacket) error {
	switch pkt.Action {
	case uint8(vaautils.ActionGuardianSetUpdate):

		var guardianSetUpdate types.GuardianSetUpdate
		err := guardianSetUpdate.Parse(pkt.Payload)
		if err != nil {
			return err
		}

		cfg, err := k.Config.Get(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to get config from state")
		}

		oldIndex := cfg.GuardianSetIndex
		if guardianSetUpdate.NewGuardianSetIndex != oldIndex+1 {
			return fmt.Errorf("invalid guardian set index: expected %d, got %d", oldIndex+1, guardianSetUpdate.NewGuardianSetIndex)
		}

		oldGuardianSet, err := k.GuardianSets.Get(ctx, oldIndex)
		if err != nil {
			return errors.Wrap(err, "failed to get old guardian set from state")
		}
		blockTime := uint64(k.headerService.GetHeaderInfo(ctx).Time.Unix())
		oldGuardianSet.ExpirationTime = blockTime + types.GuardianSetExpiry

		err = k.GuardianSets.Set(ctx, oldIndex, oldGuardianSet)
		if err != nil {
			return errors.Wrap(err, "failed to set old guardian set in state")
		}
		err = k.GuardianSets.Set(ctx, guardianSetUpdate.NewGuardianSetIndex, guardianSetUpdate.NewGuardianSet)
		if err != nil {
			return errors.Wrap(err, "failed to set new guardian set in state")
		}

		cfg.GuardianSetIndex += 1
		err = k.Config.Set(ctx, cfg)
		if err != nil {
			return errors.Wrap(err, "failed to set config in state")
		}

		return k.eventService.EventManager(ctx).EmitKV(ctx, "GuardianSetUpgrade",
			event.Attribute{Key: "old", Value: strconv.Itoa(int(oldIndex))},
			event.Attribute{Key: "new", Value: strconv.Itoa(int(cfg.GuardianSetIndex))},
		)
	default:
		return errors.Wrapf(types.ErrUnsupportedGovernanceAction, "module: %s, type: %d", pkt.Module, pkt.Action)
	}
}
