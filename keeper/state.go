package keeper

import (
	"context"
	"encoding/hex"

	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"

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
func (k *Keeper) GetVAAArchive(ctx context.Context) (map[string]string, error) {
	vaaArchive := make(map[string]string)

	err := k.VAAArchive.Walk(ctx, nil, func(hash []byte, value collections.Pair[string, bool]) (stop bool, err error) {
		if value.K2() {
			vaaArchive[hex.EncodeToString(hash)] = value.K1()
		}

		return false, nil
	})

	return vaaArchive, err
}

//

type VAAArchiveIndexes struct {
	ByID *indexes.Unique[string, []byte, collections.Pair[string, bool]]
}

func (i VAAArchiveIndexes) IndexesList() []collections.Index[[]byte, collections.Pair[string, bool]] {
	return []collections.Index[[]byte, collections.Pair[string, bool]]{i.ByID}
}

func NewVAAArchiveIndexes(builder *collections.SchemaBuilder) VAAArchiveIndexes {
	return VAAArchiveIndexes{
		ByID: indexes.NewUnique(
			builder, types.VAAByIDPrefix, "vaa_by_id",
			collections.StringKey, collections.BytesKey,
			func(_ []byte, value collections.Pair[string, bool]) (string, error) {
				return value.K1(), nil
			},
		),
	}
}
