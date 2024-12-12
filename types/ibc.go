package types

import "time"

const (
	Port           = "wormhole"
	Version        = "ibc-wormhole-v1"
	PacketLifetime = 365 * 24 * time.Hour
)

// PacketData defines the data included in an IBC packet sent to Wormchain.
// There is currently only one packet type, publish, which allows messages to
// be sent between Noble and other integrated Wormhole chains.
type PacketData struct {
	Publish struct {
		Msg []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"msg"`
	} `json:"publish"`
}
