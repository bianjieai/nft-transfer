<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.
Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## [Unreleased]

## [v1.1.3]

### Improvements

* [\#48](https://github.com/bianjieai/nft-transfer/pull/48) apply audit suggestion.
* [\#17](https://github.com/bianjieai/nft-transfer/pull/17) replace param proposal with MsgUpdateParams.
* [\#15](https://github.com/bianjieai/nft-transfer/pull/15) solve the problem of "/" parsing error in classID.
* [\#55](https://github.com/bianjieai/nft-transfer/pull/55) rename variable name
* [\#40](https://github.com/bianjieai/nft-transfer/pull/40) bump up IBC-Go.


## [v1.1.2]

### API Breaking

* [\#16](https://github.com/bianjieai/nft-transfer/pull/16) return the sequence of packet in `MsgTransferResponse`.

### Improvements

* [\#11](https://github.com/bianjieai/nft-transfer/pull/11) adjust the verification order of nft.

### Features

* [\#13](https://github.com/bianjieai/nft-transfer/pull/13) add params to control whether the module is enabled
* [\#18](https://github.com/bianjieai/nft-transfer/pull/13) add params query api

### Bug Fixes

* [\#12](https://github.com/bianjieai/nft-transfer/pull/12) fix `critical vulnerability allows attacker to take control of any NFT`.
* [\#19](https://github.com/bianjieai/nft-transfer/pull/19) fix `delete redundant parameters of the params command`.
* [\#27](https://github.com/bianjieai/nft-transfer/pull/27) fix `allow value with slashes in URL template`.

## [v1.1.2-beta]

### API Breaking

* [\#16](https://github.com/bianjieai/nft-transfer/pull/16) return the sequence of packet in `MsgTransferResponse`.

### Improvements

* [\#11](https://github.com/bianjieai/nft-transfer/pull/11) adjust the verification order of nft.

### Features

* [\#13](https://github.com/bianjieai/nft-transfer/pull/13) add params to control whether the module is enabled

### Bug Fixes

* [\#12](https://github.com/bianjieai/nft-transfer/pull/12) fix `critical vulnerability allows attacker to take control of any NFT`.

## [v1.1.1-beta]

### Dependencies

### API Breaking

### State Machine Breaking

### Improvements

* [\#7](https://github.com/bianjieai/nft-transfer/pull/7) modify JSON encoding rules

### Features

### Bug Fixes

## [v1.1.0-beta]

### Dependencies

### API Breaking

### State Machine Breaking

* (proto) [\#6](https://github.com/bianjieai/nft-transfer/pull/6) add `class_data` & `token_data` field for `NonFungibleTokenPacketData`, add `memo` field for `MsgTransfer`

### Improvements

### Features

### Bug Fixes

## [v1.0.0-beta]

### Dependencies

* [\#1](https://github.com/bianjieai/nft-transfer/pull/1) Bump ibc-go to v5.0.1.

### API Breaking

* (types/codec) [\#2](https://github.com/bianjieai/nft-transfer/pull/2) `NonFungibleTokenPacketData` uses camel case json encoding.

### State Machine Breaking

### Improvements

### Features

### Bug Fixes

* (types/packet) [\#3](https://github.com/bianjieai/nft-transfer/pull/3) It should not verify whether the address of the original chain is legal.
