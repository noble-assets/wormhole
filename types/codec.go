package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSubmitVAA{}, "wormhole/SubmitVAA", nil)
	cdc.RegisterConcrete(&MsgPostMessage{}, "wormhole/PostMessage", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSubmitVAA{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgPostMessage{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
