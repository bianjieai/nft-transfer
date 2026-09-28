# Cosmos SDK v0.53.8 and IBC Go v10.5.0

This module requires Go 1.23.8 or later. It uses Cosmos SDK v0.53.8 and
`github.com/cosmos/ibc-go/v10` v10.5.0. CometBFT resolves to v0.38.23 through
the SDK dependency requirements.

The changes follow the [Cosmos SDK upgrade guide](https://github.com/cosmos/cosmos-sdk/blob/v0.53.8/UPGRADING.md)
and the [IBC v8.1 to v10 migration guide](https://github.com/cosmos/ibc-go/blob/v10.5.0/docs/docs/05-migrations/13-v8_1-to-v10.md).

## Application integration

- Update IBC imports to `github.com/cosmos/ibc-go/v10`.
- Construct the NFT transfer keeper without the port keeper and scoped
  capability keeper. Its store key, codec, authority, ICS4 wrapper, channel
  keeper, account keeper, and NFT keeper arguments retain their meaning.
- Remove capability and IBC fee middleware wiring. IBC v10 no longer has these
  modules or channel upgrade callbacks. Update ICS4 wrappers to omit channel
  capabilities and packet callbacks to accept `channelVersion`.
- Wire core IBC keepers with `runtime.NewKVStoreService` and register Tendermint
  and Solo Machine light client routes explicitly. See `testing/simapp/app.go`.
- Include `authtypes.ModuleName` after the upgrade module in the pre-block order.

IBC v10 routes callbacks using the port identifier. The default NFT port is
`nft-transfer`, whereas this module's name is `nonfungibletokentransfer`.
The v10 router only accepts alphanumeric route keys, so the test application
registers the default NFT route as follows:

```go
ibcRouter.AddRoute(nfttransfertypes.PortRouterKey, nfttransfer.NewIBCModule(nftTransferKeeper))
```

`PortRouterKey` is `nft`. IBC v10 first tries an exact route, then tries route
keys in alphabetical order using substring matching. This key matches
`nft-transfer` before `transfer`. Applications using a custom genesis port must
choose a matching route and check for overlap with other application routes.
NFT channel handshakes and packet callbacks enforce the configured local port;
sending on another application's port is rejected before NFT state changes.

The NFT packet format, `ics721-1` version, class tracing, escrow address derivation,
and NFT store keys are unchanged.

## CLI timeouts

The transfer CLI follows IBC v10's timeout behavior:

- Relative timeouts use the local clock and a nonzero timestamp duration.
  The default remains 10 minutes.
- A nonzero relative block height is rejected. Use `--absolute-timeouts` with
  `--packet-timeout-height` for an absolute counterparty height.
- `--absolute-timeouts` continues to accept absolute height and timestamp
  values. Zero disables the corresponding timeout.

## Existing chain state

This repository provides a library and a test application. The test application's
`v10` upgrade handler runs registered module migrations; it does not provide a
complete deployment upgrade from historical chain state. Obsolete v5-v8 test
handlers referencing removed IBC migration APIs have been removed.

A chain integrating this release must provide its own upgrade handler and store
loader, settle any legacy IBC fee escrow before removing its store, and remove
the capability store according to the IBC migration guide. Run the registered
module migrations using the previous module version map. IBC's core migrations
from consensus version 6 to 8 remove the old localhost client state and channel
upgrade state. Channels must not be in `FLUSHING` or `FLUSHCOMPLETE` when the
channel migration runs. Validate the chain's actual upgrade state separately;
fresh-genesis tests do not establish this deployment behavior.
