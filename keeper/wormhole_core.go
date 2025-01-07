package keeper

import (
	"context"
	"encoding/binary"
	"fmt"
	"strconv"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/event"
	"cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/common"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"

	"github.com/noble-assets/wormhole/types"
)

func (k *Keeper) ParseAndVerifyVAA(ctx context.Context, bz []byte) (*vaautils.VAA, error) {
	vaa, err := vaautils.Unmarshal(bz)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidVAA, "failed to unmarshal: %v", err)
	}

	hash := vaa.SigningDigest().Bytes()
	if has, err := k.VAAArchive.Has(ctx, hash); err != nil || has {
		return nil, types.ErrAlreadyExecutedVAA
	}

	guardianSet, err := k.GuardianSets.Get(ctx, vaa.GuardianSetIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to get guardian set %d from state", vaa.GuardianSetIndex)
	}

	blockTime := uint64(k.headerService.GetHeaderInfo(ctx).Time.Unix())
	if guardianSet.ExpirationTime != 0 && guardianSet.ExpirationTime < blockTime {
		return nil, fmt.Errorf("guardian set %d is expired", vaa.GuardianSetIndex)
	}

	var addresses []common.Address
	for _, address := range guardianSet.Addresses {
		addresses = append(addresses, common.BytesToAddress(address))
	}
	if !vaa.VerifySignatures(addresses) {
		return nil, fmt.Errorf("invalid guardian set signatures")
	}

	if err := k.VAAArchive.Set(ctx, hash, collections.Join(vaa.MessageID(), true)); err != nil {
		return nil, errors.Wrap(err, "failed to set vaa in state")
	}

	return vaa, nil
}

func (k *Keeper) HandleCoreGovernancePacket(ctx context.Context, pkt types.GovernancePacket) error {
	switch pkt.Action {
	case 2:
		if len(pkt.Payload) < 5 {
			return types.ErrMalformedPayload
		}

		index := binary.BigEndian.Uint32(pkt.Payload[0:4])

		length := int(pkt.Payload[4:5][0])
		if len(pkt.Payload[5:]) != 20*length {
			return types.ErrMalformedPayload
		}

		offset := 5
		addresses := make([][]byte, length)
		for i := range length {
			addresses[i] = pkt.Payload[offset : offset+20]
			offset += 20
		}

		cfg, err := k.Config.Get(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to get config from state")
		}
		blockTime := uint64(k.headerService.GetHeaderInfo(ctx).Time.Unix())

		oldIndex := cfg.GuardianSetIndex
		if index != oldIndex+1 {
			return fmt.Errorf("invalid guardian set index: expected %d, got %d", oldIndex+1, index)
		}

		err = k.GuardianSets.Set(ctx, index, types.GuardianSet{
			Addresses:      addresses,
			ExpirationTime: blockTime + cfg.GuardianSetExpiry,
		})
		if err != nil {
			return errors.Wrap(err, "failed to set guardian set in state")
		}
		cfg.GuardianSetIndex += 1
		err = k.Config.Set(ctx, cfg)
		if err != nil {
			return errors.Wrap(err, "failed to set config in state")
		}

		return k.eventService.EventManager(ctx).EmitKV(ctx, "GuardianSetUpgrade",
			event.Attribute{Key: "old", Value: strconv.Itoa(int(oldIndex))},
			event.Attribute{Key: "new", Value: strconv.Itoa(int(index))},
		)
	default:
		return errors.Wrapf(types.ErrUnsupportedGovernanceAction, "module: %s, type: %d", pkt.Module, pkt.Action)
	}
}
