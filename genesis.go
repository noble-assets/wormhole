package wormhole

import (
	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/noble-assets/wormhole/keeper"
	"github.com/noble-assets/wormhole/types"
)

func InitGenesis(ctx sdk.Context, k *keeper.Keeper, cdc address.Codec, genesis types.GenesisState) {
	if err := k.Config.Set(ctx, genesis.Config); err != nil {
		panic(err)
	}

	if err := k.WormchainChannel.Set(ctx, genesis.WormchainChannel); err != nil {
		panic(err)
	}

	for index, guardianSet := range genesis.GuardianSets {
		if err := k.GuardianSets.Set(ctx, index, guardianSet); err != nil {
			panic(err)
		}
	}

	for address, sequence := range genesis.Sequences {
		sender := make([]byte, 32)
		bz, err := cdc.StringToBytes(address)
		if err != nil {
			panic(err)
		}
		copy(sender[12:], bz)

		if err := k.Sequences.Set(ctx, sender, sequence); err != nil {
			panic(err)
		}
	}

	for _, hash := range genesis.VaaArchive {
		if err := k.VAAArchive.Set(ctx, hash, true); err != nil {
			panic(err)
		}
	}

	if err := k.BindPort(ctx); err != nil {
		panic(err)
	}
}

func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	config, err := k.Config.Get(ctx)
	if err != nil {
		panic(err)
	}

	wormchainChannel, err := k.WormchainChannel.Get(ctx)
	if err != nil {
		panic(err)
	}

	guardianSets, err := k.GetGuardianSets(ctx)
	if err != nil {
		panic(err)
	}

	sequences, err := k.GetSequences(ctx)
	if err != nil {
		panic(err)
	}

	vaaArchive, err := k.GetVAAArchive(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Config:           config,
		WormchainChannel: wormchainChannel,
		GuardianSets:     guardianSets,
		Sequences:        sequences,
		VaaArchive:       vaaArchive,
	}
}
