package wormhole

import (
	"context"
	"encoding/json"
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/header"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	modulev1 "github.com/noble-assets/wormhole/api/module/v1"
	wormholev1 "github.com/noble-assets/wormhole/api/v1"
	"github.com/noble-assets/wormhole/keeper"
	"github.com/noble-assets/wormhole/types"
)

// ConsensusVersion defines the current x/wormhole module consensus version.
const ConsensusVersion = 1

var (
	_ module.AppModuleBasic      = AppModule{}
	_ appmodule.AppModule        = AppModule{}
	_ module.HasConsensusVersion = AppModule{}
	_ module.HasGenesis          = AppModule{}
	_ module.HasGenesisBasics    = AppModuleBasic{}
	_ module.HasServices         = AppModule{}
)

//

type AppModuleBasic struct{}

func NewAppModuleBasic() AppModuleBasic {
	return AppModuleBasic{}
}

func (AppModuleBasic) Name() string { return types.ModuleName }

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

func (AppModuleBasic) RegisterInterfaces(reg codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var genesis types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genesis); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return genesis.Validate()
}

//

type AppModule struct {
	AppModuleBasic

	keeper       *keeper.Keeper
	addressCodec address.Codec
}

func NewAppModule(keeper *keeper.Keeper, addressCodec address.Codec) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(),
		keeper:         keeper,
		addressCodec:   addressCodec,
	}
}

func (AppModule) IsOnePerModuleType() {}

func (AppModule) IsAppModule() {}

func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

func (m AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, bz json.RawMessage) {
	var genesis types.GenesisState
	cdc.MustUnmarshalJSON(bz, &genesis)

	InitGenesis(ctx, m.keeper, m.addressCodec, genesis)
}

func (m AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genesis := ExportGenesis(ctx, m.keeper)
	return cdc.MustMarshalJSON(genesis)
}

func (m AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(m.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(m.keeper))
}

//

func (AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: wormholev1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod:      "SubmitVAA",
					Use:            "submit-vaa [vaa]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "vaa"}},
				},
				{
					RpcMethod: "PostMessage",
					Use:       "post-message [message] [nonce]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "message"},
						{ProtoField: "nonce"},
					},
				},
			},
		},
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: wormholev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Config",
					Use:       "config",
					Short:     "Query the current module configuration",
				},
				{
					RpcMethod: "WormchainChannel",
					Use:       "wormchain-channel",
					Short:     "Query the current channel opened to Wormchain",
				},
				{
					RpcMethod: "GuardianSets",
					Use:       "guardian-sets",
					Short:     "Query all guardian sets registered in the module",
				},
				{
					RpcMethod:      "GuardianSet",
					Use:            "guardian-set [index]",
					Short:          "Query a specific guardian set given an index",
					Example:        "  nobled query wormhole guardian-set current\n  nobled query wormhole guardian-set 0",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				{
					RpcMethod:      "ExecutedVAA",
					Use:            "executed-vaa [digest]",
					Short:          "Query if a specific VAA has been executed on Noble",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "input"}},
				},
			},
		},
	}
}

//

func init() {
	appmodule.Register(&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Config        *modulev1.Module
	StoreService  store.KVStoreService
	HeaderService header.Service
	EventService  event.Service

	Codec        codec.Codec
	AddressCodec address.Codec
}

type ModuleOutputs struct {
	depinject.Out

	Keeper *keeper.Keeper
	Module appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	k := keeper.NewKeeper(
		in.Codec,
		in.StoreService,
		in.HeaderService,
		in.EventService,
		in.AddressCodec,
		nil,
		nil,
		nil,
	)
	m := NewAppModule(k, in.AddressCodec)

	return ModuleOutputs{Keeper: k, Module: m}
}
