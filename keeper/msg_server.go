package keeper

import (
	"bytes"
	"context"
	"encoding/json"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"

	"github.com/noble-assets/wormhole/types"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) SubmitVAA(ctx context.Context, msg *types.MsgSubmitVAA) (*types.MsgSubmitVAAResponse, error) {
	vaa, err := k.ParseAndVerifyVAA(ctx, msg.Vaa)
	if err != nil {
		return nil, err
	}

	config, err := k.Config.Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config from state")
	}

	if !(vaa.EmitterChain == vaautils.ChainID(config.GovChain) && bytes.Equal(vaa.EmitterAddress.Bytes(), config.GovAddress)) {
		return nil, types.ErrNotGovernanceVAA
	}
	if vaa.GuardianSetIndex != config.GuardianSetIndex {
		return nil, errors.Wrap(types.ErrInvalidGovernanceVAA, "must be signed by current guardian set")
	}

	var pkt types.GovernancePacket
	err = pkt.Parse(vaa.Payload)
	if err != nil {
		return nil, err
	}

	switch pkt.Module {
	case "Core":
		err = k.HandleCoreGovernancePacket(ctx, pkt)
		return &types.MsgSubmitVAAResponse{}, err
	case "IbcReceiver":
		err = k.HandleIBCReceiverGovernancePacket(ctx, pkt)
		return &types.MsgSubmitVAAResponse{}, err
	default:
		return &types.MsgSubmitVAAResponse{}, errors.Wrapf(types.ErrUnsupportedGovernanceAction, "module: %s, type: %d", pkt.Module, pkt.Action)
	}
}

func (k msgServer) PostMessage(ctx context.Context, msg *types.MsgPostMessage) (*types.MsgPostMessageResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	info := k.headerService.GetHeaderInfo(ctx)
	channel, err := k.WormchainChannel.Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get wormchain channel from state")
	}

	data, err := k.GetPacketData(ctx, msg.Message, msg.Nonce, msg.Signer)
	if err != nil {
		return nil, err
	}
	bz, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	capability, _ := k.scopedKeeper.GetCapability(sdkCtx, host.ChannelCapabilityPath(types.Port, channel))
	_, err = k.ics4Wrapper.SendPacket(
		sdkCtx, capability, types.Port, channel,
		clienttypes.ZeroHeight(),
		uint64(info.Time.Add(types.PacketLifetime).UnixNano()),
		bz,
	)

	return &types.MsgPostMessageResponse{}, err
}
