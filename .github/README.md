# `x/wormhole`

This repository contains a native Cosmos SDK module that enables Noble's connection to Wormhole. It is essentially a reimplementation of their [`wormholeIbc`] CosmWasm smart contract that is already deployed on Cosmos zones like Neutron, Injective, and Sei.

## Architecture Overview

![](design.png)

Noble (and specifically this module) connects to Wormhole via Wormchain, their purpose built Cosmos SDK blockchain. Wormchain implements a standard, as described [here][whitepaper], which allows certain Cosmos zones to emit messages for verification without the need for the guardian set to support each zone directly.

We have tackled this by allowing a custom channel with version `ibc-wormhole-v1` to be opened between our `wormhole` port on Noble and a port exposed by their [`wormchainIbcReceiver`] CosmWasm smart contract.

Once opened, this channel then allows a new custom packet to be sent from Noble to Wormchain. This packet allows us to emit messages to Wormhole's guardian set, who are listening to these packets. In general, these packets take the following format:

```json
{
  "publish": {
    "msg": [
      {
        "key": "...",
        "value": "..."
      },
      ...
    ]
  }
}
```

Once these packets have been relayed to Wormchain via IBC, the guardian set emits a VAA once consensus is reached!

### Governance Packets

On top of sending packets, the module additionally exposes a custom transaction that allows Wormhole governance actions to be submitted on Noble. This message is permissionless, allowing any user to submit these packets once they are emitted by Wormhole's guardian set.

Wormhole takes a module / action approach to governance, allowing them to easily distinguish different governance actions across their various implementations and infrastructure. The module currently implements and can execute the following packets:

| Module        | Action               |
|---------------|----------------------|
| `Core`        | `GuardianSetUpgrade` |
| `IbcReceiver` | `UpdateChannelChain` |

[`wormholeIbc`]: https://github.com/wormhole-foundation/wormhole/tree/main/cosmwasm/contracts/wormhole-ibc
[`wormchainIbcReceiver`]: https://github.com/wormhole-foundation/wormhole/tree/main/cosmwasm/contracts/wormchain-ibc-receiver
[whitepaper]: https://github.com/wormhole-foundation/wormhole/blob/main/whitepapers/0012_ibc_generic_messaging.md
