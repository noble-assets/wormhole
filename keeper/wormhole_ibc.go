package keeper

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"strconv"

	"cosmossdk.io/core/event"
	"cosmossdk.io/errors"
	"github.com/wormhole-foundation/wormhole/sdk/vaa"

	"github.com/noble-assets/wormhole/types"
)

func (k *Keeper) GetPacketData(ctx context.Context, message []byte, nonce uint32, signer string) (*types.PacketData, error) {
	config, err := k.Config.Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config from state")
	}

	emitter := make([]byte, 32)
	bz, err := k.addressCodec.StringToBytes(signer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode signer address")
	}
	copy(emitter[12:], bz)

	sequence, _ := k.Sequences.Get(ctx, emitter)
	err = k.Sequences.Set(ctx, emitter, sequence+1)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set sequence in state")
	}

	return &types.PacketData{
		Publish: struct {
			Msg []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"msg"`
		}(struct {
			Msg []struct {
				Key   string
				Value string
			}
		}{
			Msg: []struct {
				Key   string
				Value string
			}{
				{
					Key:   "message.message",
					Value: hex.EncodeToString(message),
				},
				{
					Key:   "message.sender",
					Value: hex.EncodeToString(emitter),
				},
				{
					Key:   "message.chain_id",
					Value: strconv.Itoa(int(config.ChainId)),
				},
				{
					Key:   "message.nonce",
					Value: strconv.Itoa(int(nonce)),
				},
				{
					Key:   "message.sequence",
					Value: strconv.Itoa(int(sequence)),
				},
				{
					Key:   "message.block_time",
					Value: strconv.Itoa(int(k.headerService.GetHeaderInfo(ctx).Time.Unix())),
				},
			},
		}),
	}, nil
}

func (k *Keeper) HandleIBCReceiverGovernancePacket(ctx context.Context, pkt types.GovernancePacket) error {
	switch pkt.Action {
	case 1:
		if len(pkt.Payload) != 66 {
			return types.ErrMalformedPayload
		}

		channel := string(bytes.TrimLeft(pkt.Payload[0:64], "\x00"))
		chain := binary.BigEndian.Uint16(pkt.Payload[64:66])

		if chain != uint16(vaa.ChainIDWormchain) {
			return types.ErrInvalidChannel
		}

		if err := k.WormchainChannel.Set(ctx, channel); err != nil {
			return errors.Wrap(err, "failed to set wormchain channel in state")
		}

		return k.eventService.EventManager(ctx).EmitKV(ctx, "UpdateChannelChain",
			event.Attribute{Key: "chain_id", Value: strconv.Itoa(int(chain))},
			event.Attribute{Key: "channel_id", Value: channel},
		)
	default:
		return errors.Wrapf(types.ErrUnsupportedGovernanceAction, "module: %s, type: %d", pkt.Module, pkt.Action)
	}
}
