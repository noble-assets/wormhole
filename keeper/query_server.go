package keeper

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/noble-assets/wormhole/types"
)

var _ types.QueryServer = &queryServer{}

type queryServer struct {
	*Keeper
}

func NewQueryServer(keeper *Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

func (k queryServer) Config(ctx context.Context, req *types.QueryConfig) (*types.QueryConfigResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	config, err := k.Keeper.Config.Get(ctx)
	if err != nil {
		return nil, errors.New("unable to get config from state")
	}

	return &types.QueryConfigResponse{Config: config}, nil
}

func (k queryServer) WormchainChannel(ctx context.Context, req *types.QueryWormchainChannel) (*types.QueryWormchainChannelResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	wormchainChannel, err := k.Keeper.WormchainChannel.Get(ctx)
	if err != nil {
		return nil, errors.New("wormchain channel not configured in state")
	}

	return &types.QueryWormchainChannelResponse{WormchainChannel: wormchainChannel}, nil
}

func (k queryServer) GuardianSets(ctx context.Context, req *types.QueryGuardianSets) (*types.QueryGuardianSetsResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	guardianSets, err := k.GetGuardianSets(ctx)
	if err != nil {
		return nil, errors.New("unable to get guardian sets from state")
	}

	return &types.QueryGuardianSetsResponse{GuardianSets: guardianSets}, nil
}

func (k queryServer) GuardianSet(ctx context.Context, req *types.QueryGuardianSet) (*types.QueryGuardianSetResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	index := uint32(0)
	if req.Index == "current" {
		config, err := k.Keeper.Config.Get(ctx)
		if err != nil {
			return nil, errors.New("unable to get config from state")
		}

		index = config.GuardianSetIndex
	} else {
		raw, err := strconv.Atoi(req.Index)
		if err != nil {
			return nil, fmt.Errorf("invalid guardian set index %s", req.Index)
		}

		index = uint32(raw)
	}

	guardianSet, err := k.Keeper.GuardianSets.Get(ctx, index)
	if err != nil {
		return nil, fmt.Errorf("unable to get guardian set %d from state", index)
	}

	return &types.QueryGuardianSetResponse{GuardianSet: guardianSet}, nil
}

func (k queryServer) ExecutedVAA(ctx context.Context, req *types.QueryExecutedVAA) (*types.QueryExecutedVAAResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	// TODO: k.VAAArchive.Indexes.ByID.MatchExact(ctx, req.Input)

	digest, err := hex.DecodeString(req.Input)
	if err != nil {
		return nil, fmt.Errorf("unable to decode digest %s", req.Input)
	}

	executed, _ := k.VAAArchive.Has(ctx, digest)

	return &types.QueryExecutedVAAResponse{Executed: executed}, nil
}
