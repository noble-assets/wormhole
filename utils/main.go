package main

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec/address"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func main() {
	addressCodec := address.NewBech32Codec("noble")

	transceiverAddress := authtypes.NewModuleAddress("dollar/portal/transceiver")
	transceiver, _ := addressCodec.BytesToString(transceiverAddress)
	fmt.Println("TRANSCEIVER:", transceiver)

	paddedTransceiverAddress := make([]byte, 32)
	copy(paddedTransceiverAddress[12:], transceiverAddress)
	fmt.Println("PADDED:     ", hex.EncodeToString(paddedTransceiverAddress))

	fmt.Println()

	managerAddress := authtypes.NewModuleAddress("dollar/portal/manager")
	manager, _ := addressCodec.BytesToString(managerAddress)
	fmt.Println("MANAGER:    ", manager)

	paddedManagerAddress := make([]byte, 32)
	copy(paddedManagerAddress[12:], managerAddress)
	fmt.Println("PADDED:     ", hex.EncodeToString(paddedManagerAddress))
}
