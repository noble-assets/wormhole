# E2E Testing

This directory contains everything necessary for locally testing Noble's integration with Wormhole.

Please note that the instructions below are incomplete!

## Running the IBC Relayer

- `rly config init --home relayer`
- `rly chains add-dir relayer/chains --home relayer`
- `rly paths new duke-1 wormchain wormhole --home relayer`
- `rly transact link wormhole --version "ibc-wormhole-v1" --src-port "wormhole" --dst-port "wasm.wormhole1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrq0kdhcj" --home relayer`

## Generating the VAAs

- `/guardiand template ibc-receiver-update-channel-chain --chain-id 4009 --channel-id channel-0 --target-chain-id 3104 --idx 0 > vaa`
- `/guardiand admin governance-vaa-inject vaa --socket /tmp/admin.sock`
- `/guardiand template ibc-receiver-update-channel-chain --chain-id 3104 --channel-id channel-0 --target-chain-id 4009 --idx 0 > vaa`
- `/guardiand admin governance-vaa-inject vaa --socket /tmp/admin.sock`

## Register the IBC Channels

- `/simapp/build/simd tx wormhole submit-vaa [VAA] --chain-id duke-1 --from validator --keyring-backend test --home .duke`
- `./build/wormchaind tx wasm execute wormhole1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrq0kdhcj '{"submit_update_channel_chain":{"vaas":["TODO"]}}' --from tiltRelayer --home build`
