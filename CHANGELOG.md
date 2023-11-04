# Changelog

## [0.16.1](https://github.com/pactus-project/pactus/compare/v0.16.0...v0.16.1)(2023-11-04)

### Feat

- **network**: adding public address to factory ([#795](https://github.com/pactus-project/pactus/pull/795))
- **network**: filter private ips ([#793](https://github.com/pactus-project/pactus/pull/793))

### Fix

- **network**: not close connection for bootstrap nodes ([#792](https://github.com/pactus-project/pactus/pull/792))

## [0.16.0](https://github.com/pactus-project/pactus/compare/v0.15.0...v0.16.0)(2023-10-29)

### Feat

- **gui**: display network ID ([#780](https://github.com/pactus-project/pactus/pull/780))
- create new validator address (CLI and GUI) ([#757](https://github.com/pactus-project/pactus/pull/757))
- add community bootstrap nodes to testnet config ([#764](https://github.com/pactus-project/pactus/pull/764))
- **network**: implementing connection manager ([#773](https://github.com/pactus-project/pactus/pull/773))
- **network**: adding bootstrapper mode to the network config ([#760](https://github.com/pactus-project/pactus/pull/760))

### Fix

- **network**: redefine the network limits ([#788](https://github.com/pactus-project/pactus/pull/788))
- **consensus**: not increase the vote-box power on duplicated votes ([#785](https://github.com/pactus-project/pactus/pull/785))
- **network**: close connection when unbale to get supported protocols ([#781](https://github.com/pactus-project/pactus/pull/781))
- **network**: enabling peer exchange for bootstrappers ([#779](https://github.com/pactus-project/pactus/pull/779))
- **network**: set connection limit for the resource manager ([#775](https://github.com/pactus-project/pactus/pull/775))
- **sync**: peer status set to known on sucessfull handshaking ([#774](https://github.com/pactus-project/pactus/pull/774))
- **consensus**: strong termination for the binary agreement ([#765](https://github.com/pactus-project/pactus/pull/765))
- **consensus**: not increase the voting power on duplicated binary votes ([#762](https://github.com/pactus-project/pactus/pull/762))

### Refactor

- **network**: refactoring peer manager ([#787](https://github.com/pactus-project/pactus/pull/787))

## [0.15.0](https://github.com/pactus-project/pactus/compare/v0.13.0...v0.15.0)(2023-10-15)

### Feat

- **gui**: adding the splash screen ([#743](https://github.com/pactus-project/pactus/pull/743))
- add absentees votes to the certificate ([#746](https://github.com/pactus-project/pactus/pull/746))
- **logger**: short stringer for loggers ([#732](https://github.com/pactus-project/pactus/pull/732))
- implementing pip-7 ([#731](https://github.com/pactus-project/pactus/pull/731))
- implementing pip-11 ([#712](https://github.com/pactus-project/pactus/pull/712))
- implementing pip-8 ([#711](https://github.com/pactus-project/pactus/pull/711))
- implementing pip-9 ([#706](https://github.com/pactus-project/pactus/pull/706))
- new API to get Public key by address ([#704](https://github.com/pactus-project/pactus/pull/704))
- Adding address field for AccountInfo ([#703](https://github.com/pactus-project/pactus/pull/703))
- CreateValidatorEvent and CreateAccountEvent for nanomsg ([#702](https://github.com/pactus-project/pactus/pull/702))
- implementing PIP-2 and PIP-3 ([#699](https://github.com/pactus-project/pactus/pull/699))
- Adding Hole Punching to network ([#697](https://github.com/pactus-project/pactus/pull/697))
- write logs into file ([#673](https://github.com/pactus-project/pactus/pull/673))
- check protocol support before sending connect/disconnect event ([#683](https://github.com/pactus-project/pactus/pull/683))
- updating genesis for pre-testnet-2 ([#679](https://github.com/pactus-project/pactus/pull/679))
- adding udp protocol for network ([#672](https://github.com/pactus-project/pactus/pull/672))
- implementing pip-4 ([#671](https://github.com/pactus-project/pactus/pull/671))
- Notifee service events ([#628](https://github.com/pactus-project/pactus/pull/628))
- adding MinimumStake parameter ([#574](https://github.com/pactus-project/pactus/pull/574))
- adding Sent and Received bytes per message metrics for peers ([#618](https://github.com/pactus-project/pactus/pull/618))
- add reason to BlockResponse messages ([#607](https://github.com/pactus-project/pactus/pull/607))
- Add CalcualteFee in GRPC ([#601](https://github.com/pactus-project/pactus/pull/601))
- add sent bytes and received bytes metrics to peerset plus update grpc ([#606](https://github.com/pactus-project/pactus/pull/606))
- added metrics of libp2p with supporting prometheus ([#588](https://github.com/pactus-project/pactus/pull/588))
- Check node address is valid ([#565](https://github.com/pactus-project/pactus/pull/565))
- add LastSent and LastReceived properties to peer ([#569](https://github.com/pactus-project/pactus/pull/569))

### Fix

- data race issue on updating certificate ([#747](https://github.com/pactus-project/pactus/pull/747))
- **network**: async connection ([#744](https://github.com/pactus-project/pactus/pull/744))
- adding query vote timer for CP phase ([#738](https://github.com/pactus-project/pactus/pull/738))
- trim transactions in proposed block ([#737](https://github.com/pactus-project/pactus/pull/737))
- fixing query votes and proposal issue ([#736](https://github.com/pactus-project/pactus/pull/736))
- fixing issue when a block has max transactions ([#735](https://github.com/pactus-project/pactus/pull/735))
- **consensus**: anti-entroy mechanism for the consensus ([#734](https://github.com/pactus-project/pactus/pull/734))
- **logger**: invalid level parsing error ([#733](https://github.com/pactus-project/pactus/pull/733))
- cache certificate by height ([#730](https://github.com/pactus-project/pactus/pull/730))
- fixing a crash on consensus ([#729](https://github.com/pactus-project/pactus/pull/729))
- **consensus**: prevent double entry in new height ([#728](https://github.com/pactus-project/pactus/pull/728))
- resolve consensus halt caused by time discrepancy in network. ([#727](https://github.com/pactus-project/pactus/pull/727))
- unsorted addresses in wallet listing ([#721](https://github.com/pactus-project/pactus/pull/721))
- send query votes message, if there is no proposal yet ([#723](https://github.com/pactus-project/pactus/pull/723))
- fixing logger level issue ([#722](https://github.com/pactus-project/pactus/pull/722))
- fixing syncing stuck issue ([#720](https://github.com/pactus-project/pactus/pull/720))
- fixing some minor issues on pre-testnet ([#719](https://github.com/pactus-project/pactus/pull/719))
- supporting go version 1.21 and higher ([#692](https://github.com/pactus-project/pactus/pull/692))
- ensure log rotation using tests ([#693](https://github.com/pactus-project/pactus/pull/693))
- restoring at the first block ([#691](https://github.com/pactus-project/pactus/pull/691))
- swagger doesn't work with multiple proto files ([#687](https://github.com/pactus-project/pactus/pull/687))
- fixing wallet-cli issues ([#686](https://github.com/pactus-project/pactus/pull/686))
- prevent stripping public key for subsidy transactions ([#678](https://github.com/pactus-project/pactus/pull/678))
- updating the consensus protocol ([#668](https://github.com/pactus-project/pactus/pull/668))
- aggregating signature for hello message ([#640](https://github.com/pactus-project/pactus/pull/640))
- error case for logger ([#634](https://github.com/pactus-project/pactus/pull/634))
- adding committers to the certificate ([#623](https://github.com/pactus-project/pactus/pull/623))
- updating sortition executor ([#608](https://github.com/pactus-project/pactus/pull/608))
- update buf and fixing proto generation issue   ([#600](https://github.com/pactus-project/pactus/pull/600))
- adding block hash to peer ([#584](https://github.com/pactus-project/pactus/pull/584))
- copy to clipboard option for address and pubkey ([#583](https://github.com/pactus-project/pactus/pull/583))
- public key aggregate ([#576](https://github.com/pactus-project/pactus/pull/576))
- remove GetValidators rpc method ([#573](https://github.com/pactus-project/pactus/pull/573))
- missing swagger ui for grpc get account by number ([#564](https://github.com/pactus-project/pactus/pull/564))
- incorrect handler for validator by number ([#563](https://github.com/pactus-project/pactus/pull/563))

### Refactor

- **sync**: refactoring syncing process ([#676](https://github.com/pactus-project/pactus/pull/676))
- remove payload prefix from payload transaction type ([#669](https://github.com/pactus-project/pactus/pull/669))
- change Hello message from broadcasting to direct messaging ([#665](https://github.com/pactus-project/pactus/pull/665))
- **committee**: using generic list for validators ([#667](https://github.com/pactus-project/pactus/pull/667))
- rename SanityCheck to BasicCheck ([#643](https://github.com/pactus-project/pactus/pull/643))
- **cli**: Migrating from mow.cli to cobra for wallet ([#629](https://github.com/pactus-project/pactus/pull/629))
- **cli**: replacing mow.cli with cobra for daemon ([#621](https://github.com/pactus-project/pactus/pull/621))
- **logger**: using fast JSON logger (zerolog) ([#613](https://github.com/pactus-project/pactus/pull/613))
- Using Generics for calculating Min and Max for numeric type #604 ([#609](https://github.com/pactus-project/pactus/pull/609))
- Updating LRU cache to version 2 #514 ([#602](https://github.com/pactus-project/pactus/pull/602))

## [0.13.0](https://github.com/pactus-project/pactus/compare/v0.12.0...v0.13.0)(2023-06-30)

### Fix

- implemented restore wallet base on input seed ([#553](https://github.com/pactus-project/pactus/pull/553))
- measuring total sent and received bytes ([#552](https://github.com/pactus-project/pactus/pull/552))
- add validate seed and restore wallet ([#533](https://github.com/pactus-project/pactus/pull/533))
- **HTTP**: Missing handlers ([#549](https://github.com/pactus-project/pactus/pull/549))
- **gui**: update about dialog and menu items in help ([#532](https://github.com/pactus-project/pactus/pull/532))

### Refactor

- implementing TestSuite ([#535](https://github.com/pactus-project/pactus/pull/535))

## [0.12.0](https://github.com/pactus-project/pactus/compare/v0.11.0...v0.12.0)(2023-06-19)

### Feat

- add GetAccountByNumber API to get account by number ([#511](https://github.com/pactus-project/pactus/pull/511))

### Fix

- caching account and validator in store ([#513](https://github.com/pactus-project/pactus/pull/513))
- get recent blocks by stamp ([#509](https://github.com/pactus-project/pactus/pull/509))
- closing the mDNS connection upon stopping the network ([#508](https://github.com/pactus-project/pactus/pull/508))
- updating linkedmap to use generic ([#507](https://github.com/pactus-project/pactus/pull/507))
- removing state from cache ([#506](https://github.com/pactus-project/pactus/pull/506))
- Typo in GUI ([#499](https://github.com/pactus-project/pactus/pull/499))
- supporting localnet ([#492](https://github.com/pactus-project/pactus/pull/492))

### Refactor

- update total power calculation based on power change(deltas) ([#518](https://github.com/pactus-project/pactus/pull/518))
- GetValidators return all validators in state validators map ([#512](https://github.com/pactus-project/pactus/pull/512))

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
