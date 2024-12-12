package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/header"
	"cosmossdk.io/core/store"
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	"github.com/cosmos/ibc-go/v8/modules/core/24-host"

	"github.com/noble-assets/wormhole/types"
)

type Keeper struct {
	schema        collections.Schema
	headerService header.Service
	eventService  event.Service
	addressCodec  address.Codec

	Config           collections.Item[types.Config]
	WormchainChannel collections.Item[string]
	GuardianSets     collections.Map[uint32, types.GuardianSet]
	Sequences        collections.Map[[]byte, uint64]
	VAAArchive       collections.Map[[]byte, bool]

	ics4Wrapper  types.ICS4Wrapper
	portKeeper   types.PortKeeper
	scopedKeeper types.ScopedKeeper
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	headerService header.Service,
	eventService event.Service,
	addressCodec address.Codec,
	ics4Wrapper types.ICS4Wrapper,
	portKeeper types.PortKeeper,
	scopedKeeper types.ScopedKeeper,
) *Keeper {
	builder := collections.NewSchemaBuilder(storeService)

	keeper := &Keeper{
		headerService: headerService,
		eventService:  eventService,
		addressCodec:  addressCodec,

		Config:           collections.NewItem(builder, types.ConfigKey, "config", codec.CollValue[types.Config](cdc)),
		WormchainChannel: collections.NewItem(builder, types.WormchainChannelKey, "wormchain_channel", collections.StringValue),
		GuardianSets:     collections.NewMap(builder, types.GuardianSetPrefix, "guardian_sets", collections.Uint32Key, codec.CollValue[types.GuardianSet](cdc)),
		Sequences:        collections.NewMap(builder, types.SequencePrefix, "sequences", collections.BytesKey, collections.Uint64Value),
		VAAArchive:       collections.NewMap(builder, types.VAAArchivePrefix, "vaa_archive", collections.BytesKey, collections.BoolValue),

		ics4Wrapper:  ics4Wrapper,
		portKeeper:   portKeeper,
		scopedKeeper: scopedKeeper,
	}

	schema, err := builder.Build()
	if err != nil {
		panic(err)
	}

	keeper.schema = schema
	return keeper
}

// SetIBCKeepers overrides the provided IBC specific keepers for this module.
// This exists because IBC doesn't support dependency injection.
func (k *Keeper) SetIBCKeepers(ics4Wrapper types.ICS4Wrapper, portKeeper types.PortKeeper, scopedKeeper types.ScopedKeeper) {
	k.ics4Wrapper = ics4Wrapper
	k.portKeeper = portKeeper
	k.scopedKeeper = scopedKeeper
}

// BindPort allows the module to bind a specific port on initialization.
func (k *Keeper) BindPort(ctx sdk.Context) error {
	if _, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(types.Port)); !ok {
		capability := k.portKeeper.BindPort(ctx, types.Port)
		err := k.ClaimCapability(ctx, capability, host.PortPath(types.Port))
		if err != nil {
			return errors.Wrap(err, "could not claim port capability")
		}
	}

	return nil
}

// ClaimCapability allows the module to claim port or channel capabilities.
func (k *Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}
