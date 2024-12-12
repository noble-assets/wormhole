package keeper

import (
	"context"

	"github.com/noble-assets/wormhole/types"
)

// GetGuardianSets is a helper function for retrieving all guardian sets from state.
func (k *Keeper) GetGuardianSets(ctx context.Context) (map[uint32]types.GuardianSet, error) {
	guardianSets := make(map[uint32]types.GuardianSet)

	err := k.GuardianSets.Walk(ctx, nil, func(index uint32, guardianSet types.GuardianSet) (stop bool, err error) {
		guardianSets[index] = guardianSet
		return false, nil
	})

	return guardianSets, err
}

// GetSequences is a helper function for retrieving all sequence entries from
// state. It encodes senders as Bech32 addresses for use in a genesis export.
func (k *Keeper) GetSequences(ctx context.Context) (map[string]uint64, error) {
	sequences := make(map[string]uint64)

	err := k.Sequences.Walk(ctx, nil, func(sender []byte, sequence uint64) (stop bool, err error) {
		// NOTE: This assumes that addresses contain 20 bytes.
		address, err := k.addressCodec.BytesToString(sender[12:])
		if err != nil {
			// NOTE: We continue in the case of an encoding error.
			return false, err
		}

		sequences[address] = sequence
		return false, nil
	})

	return sequences, err
}

// GetVAAArchive is a helper function for retrieving all executed VAAs from
// state. Note that this should only be used when exporting genesis.
func (k *Keeper) GetVAAArchive(ctx context.Context) ([][]byte, error) {
	var vaaArchive [][]byte

	err := k.VAAArchive.Walk(ctx, nil, func(hash []byte, executed bool) (stop bool, err error) {
		if executed {
			vaaArchive = append(vaaArchive, hash)
		}

		return false, nil
	})

	return vaaArchive, err
}
