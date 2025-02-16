package simapp

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	ibc "github.com/cosmos/ibc-go/v10/modules/core"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v10/modules/core/03-connection/types"
	porttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v10/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"
	solomachine "github.com/cosmos/ibc-go/v10/modules/light-clients/06-solomachine"
	tendermint "github.com/cosmos/ibc-go/v10/modules/light-clients/07-tendermint"

	"github.com/noble-assets/wormhole"
	"github.com/noble-assets/wormhole/types"
)

func (app *SimApp) RegisterLegacyModules() error {
	if err := app.RegisterStores(
		storetypes.NewKVStoreKey(exported.StoreKey),
	); err != nil {
		return err
	}

	app.ParamsKeeper.Subspace(exported.ModuleName).WithKeyTable(clienttypes.ParamKeyTable().RegisterParamSet(&connectiontypes.Params{}))

	app.IBCKeeper = ibckeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.GetKey(exported.StoreKey)),
		app.GetSubspace(exported.ModuleName),
		app.UpgradeKeeper,
		"noble1vvn7s88yj02ktwwckdzvtz64fvengfsjtwejck",
	)

	app.WormholeKeeper.SetIBCKeepers(app.IBCKeeper.ChannelKeeper)

	router := porttypes.NewRouter()
	router.AddRoute(types.ModuleName, wormhole.NewIBCModule(app.WormholeKeeper))
	app.IBCKeeper.SetRouter(router)

	clientKeeper := app.IBCKeeper.ClientKeeper
	storeProvider := clientKeeper.GetStoreProvider()

	tmLightClientModule := tendermint.NewLightClientModule(app.appCodec, storeProvider)
	clientKeeper.AddRoute(tendermint.ModuleName, &tmLightClientModule)

	smLightClientModule := solomachine.NewLightClientModule(app.appCodec, storeProvider)
	clientKeeper.AddRoute(solomachine.ModuleName, &smLightClientModule)

	return app.RegisterModules(
		ibc.NewAppModule(app.IBCKeeper),
		tendermint.NewAppModule(tmLightClientModule),
		solomachine.NewAppModule(smLightClientModule),
	)
}
