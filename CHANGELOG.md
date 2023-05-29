# Changelog

## [0.11.0](https://github.com/pactus-project/pactus/compare/v0.10.0...v0.11.0) (2023-05-29)

### Fix

- **gui**: showing the total number of validators ([#474](https://github.com/pactus-project/pactus/pull/474))
- **network**: fixing relay connection issue ([#475](https://github.com/pactus-project/pactus/pull/475))
- **consensus**: rejecting vote with round numbers exceeding the limit ([#466](https://github.com/pactus-project/pactus/pull/466))
- increase failed counter when stream got error ([#489](https://github.com/pactus-project/pactus/pull/489))
- boosting syncing process ([#482](https://github.com/pactus-project/pactus/pull/482))
- waiting for proposal in pre-commit phase ([#486](https://github.com/pactus-project/pactus/pull/486))
- retrieving public key from wallet for bond transactions ([#485](https://github.com/pactus-project/pactus/pull/485))
- restoring config file to the default ([#484](https://github.com/pactus-project/pactus/pull/484))
- defining ChainType in genesis to detect the type of network ([#483](https://github.com/pactus-project/pactus/pull/483))
- reducing the default Argon2d to consume less memory ([#480](https://github.com/pactus-project/pactus/pull/480))
- adding password option to the start commands ([#473](https://github.com/pactus-project/pactus/pull/473))

### Refactor

- rename send to transfer. ([#469](https://github.com/pactus-project/pactus/pull/469))

## [0.10.0](https://github.com/pactus-project/pactus/compare/v0.9.0...v0.10.0) (2023-05-09)

### Feat

- removing address from account ([#454](https://github.com/pactus-project/pactus/pull/454))
- supporting multiple consensus instances ([#450](https://github.com/pactus-project/pactus/pull/450))
- adding sortition interval to the parameters ([#442](https://github.com/pactus-project/pactus/pull/442))

### Fix

- `GetBlockchainInfo` API in gRPC now returns the total number of validators and accounts
- **gui**: check if the node has an active consensus instance ([#458](https://github.com/pactus-project/pactus/pull/458))
- wallet path as argument ([#455](https://github.com/pactus-project/pactus/pull/455))
- adding reward addresses to config ([#453](https://github.com/pactus-project/pactus/pull/453))
- removing committers from the certificate hash ([#444](https://github.com/pactus-project/pactus/pull/444))
- prevent data race conditions in committee  ([#452](https://github.com/pactus-project/pactus/pull/452))
- using 2^256 for the vrf denominator ([#445](https://github.com/pactus-project/pactus/pull/445))
- updating tla+ readme ([#443](https://github.com/pactus-project/pactus/pull/443))

## 0.9.0 (2022-09-05)

No changelog
