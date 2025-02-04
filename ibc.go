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

package wormhole

import (
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"

	"github.com/noble-assets/wormhole/keeper"
	"github.com/noble-assets/wormhole/types"
)

var _ porttypes.IBCModule = IBCModule{}

type IBCModule struct {
	*keeper.Keeper
}

func NewIBCModule(keeper *keeper.Keeper) IBCModule {
	return IBCModule{Keeper: keeper}
}

func (m IBCModule) OnChanOpenInit(ctx sdk.Context, _ channeltypes.Order, _ []string, port string, channel string, cap *capabilitytypes.Capability, _ channeltypes.Counterparty, version string) (string, error) {
	if port != types.Port {
		return "", errors.Wrapf(types.ErrInvalidPort, "expected port %s, got %s", types.Port, port)
	}
	if version != types.Version {
		return "", errors.Wrapf(types.ErrInvalidVersion, "expected version %s, got %s", types.Version, version)
	}

	err := m.ClaimCapability(ctx, cap, host.ChannelCapabilityPath(port, channel))

	return types.Version, err
}

func (m IBCModule) OnChanOpenTry(ctx sdk.Context, _ channeltypes.Order, _ []string, port string, channel string, cap *capabilitytypes.Capability, _ channeltypes.Counterparty, counterpartyVersion string) (string, error) {
	if port != types.Port {
		return "", errors.Wrapf(types.ErrInvalidPort, "expected port %s, got %s", types.Port, port)
	}
	if counterpartyVersion != types.Version {
		return "", errors.Wrapf(types.ErrInvalidVersion, "expected counterparty version %s, got %s", types.Version, counterpartyVersion)
	}

	err := m.ClaimCapability(ctx, cap, host.ChannelCapabilityPath(port, channel))

	return types.Version, err
}

func (m IBCModule) OnChanOpenAck(_ sdk.Context, port string, _ string, _ string, counterpartyVersion string) error {
	if port != types.Port {
		return errors.Wrapf(types.ErrInvalidPort, "expected port %s, got %s", types.Port, port)
	}
	if counterpartyVersion != types.Version {
		return errors.Wrapf(types.ErrInvalidVersion, "expected counterparty version %s, got %s", types.Version, counterpartyVersion)
	}

	return nil
}

func (m IBCModule) OnChanOpenConfirm(_ sdk.Context, port string, _ string) error {
	if port != types.Port {
		return errors.Wrapf(types.ErrInvalidPort, "expected port %s, got %s", types.Port, port)
	}

	return nil
}

func (m IBCModule) OnChanCloseInit(_ sdk.Context, _ string, _ string) error {
	return fmt.Errorf("channels with version %s cannot be closed", types.Version)
}

func (m IBCModule) OnChanCloseConfirm(_ sdk.Context, _ string, _ string) error {
	return fmt.Errorf("channels with version %s cannot be closed", types.Version)
}

func (m IBCModule) OnRecvPacket(_ sdk.Context, _ channeltypes.Packet, _ sdk.AccAddress) exported.Acknowledgement {
	return channeltypes.NewErrorAcknowledgement(fmt.Errorf("channels with version %s cannot receive packets", types.Version))
}

func (m IBCModule) OnAcknowledgementPacket(_ sdk.Context, _ channeltypes.Packet, _ []byte, _ sdk.AccAddress) error {
	return nil
}

func (m IBCModule) OnTimeoutPacket(_ sdk.Context, _ channeltypes.Packet, _ sdk.AccAddress) error {
	return nil
}
