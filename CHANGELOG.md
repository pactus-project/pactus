# Changelog

## [1.7.5](https://github.com/pactus-project/pactus/compare/v1.7.4...v1.7.5) (2025-06-06)

## [1.7.4](https://github.com/pactus-project/pactus/compare/v1.7.3...v1.7.4) (2025-05-04)

### Feat

- **packager**: add python package ([#1775](https://github.com/pactus-project/pactus/pull/1775))
- **release**: add packager jsonrpc for javascript ([#1773](https://github.com/pactus-project/pactus/pull/1773))
- generate grpc-web files

### Fix

- fix openrpc issues

## [1.7.3](https://github.com/pactus-project/pactus/compare/v1.7.2...v1.7.3) (2025-04-16)

### Fix

- **grpc**: update gRPC client to use @grpc/grpc-js ([#1751](https://github.com/pactus-project/pactus/pull/1751))

## [1.7.2](https://github.com/pactus-project/pactus/compare/v1.7.1...v1.7.2) (2025-04-12)

### Feat

- **other**: add packager ([#1738](https://github.com/pactus-project/pactus/pull/1738))

## [1.7.1](https://github.com/pactus-project/pactus/compare/v1.7.0...v1.7.1) (2025-02-09)

### Fix

- **state**: check the xeggex account hash on committing new block  ([#1688](https://github.com/pactus-project/pactus/pull/1688))
- **state, sandbox**: implement the PIP-38 ([#1687](https://github.com/pactus-project/pactus/pull/1687))

## [1.7.0](https://github.com/pactus-project/pactus/compare/v1.6.0...v1.7.0) (2025-01-23)

### Feat

- **zeromq**: add ZMQ Publishers to NodeInfo API ([#1674](https://github.com/pactus-project/pactus/pull/1674))
- **zeromq**: add publisher raw tx ([#1672](https://github.com/pactus-project/pactus/pull/1672))
- **zeromq**: add publisher raw block ([#1670](https://github.com/pactus-project/pactus/pull/1670))
- **zeromq**: add publisher transaction info ([#1669](https://github.com/pactus-project/pactus/pull/1669))
- **zeromq**: add block info publisher ([#1666](https://github.com/pactus-project/pactus/pull/1666))
- **grpc**: support Ed25519 message signing and verification ([#1667](https://github.com/pactus-project/pactus/pull/1667))
- **other**: add zeromq server with configuration ([#1660](https://github.com/pactus-project/pactus/pull/1660))
- **cmd**: read password from file ([#1653](https://github.com/pactus-project/pactus/pull/1653))
- **network**: evaluate propagate policy for gossip messages ([#1647](https://github.com/pactus-project/pactus/pull/1647))
- **config**: add firewall module to logger config ([#1637](https://github.com/pactus-project/pactus/pull/1637))

### Fix

- **gtk**: panic on windows 11 and mac ([#1650](https://github.com/pactus-project/pactus/pull/1650))
- **consensus**: refactor strong termination for change-proposer phase ([#1643](https://github.com/pactus-project/pactus/pull/1643))
- **network**: restore default message queue size in PubSub module ([#1642](https://github.com/pactus-project/pactus/pull/1642))
- **consensus**: refactor strong termination for change-proposer phase ([#1641](https://github.com/pactus-project/pactus/pull/1641))
- **grpc**: define address for the aggregated public key ([#1608](https://github.com/pactus-project/pactus/pull/1608))
- **gtk**: last block height icon ([#1611](https://github.com/pactus-project/pactus/pull/1611))
- **crypto, state**: resolve panic on 32-bit OSes ([#1604](https://github.com/pactus-project/pactus/pull/1604))
- **cmd**: parse withdraw fee using transaction options ([#1602](https://github.com/pactus-project/pactus/pull/1602))

## [1.6.0](https://github.com/pactus-project/pactus/compare/v1.5.0...v1.6.0) (2024-11-14)

### Feat

- **grpc**: add bls public key and signature aggregate methods ([#1587](https://github.com/pactus-project/pactus/pull/1587))
- **wallet**: create single ed25519 reward address for all validators ([#1570](https://github.com/pactus-project/pactus/pull/1570))
- **gtk**: add fee entry for transfer, bond and withdraw ([#1575](https://github.com/pactus-project/pactus/pull/1575))
- **txpool**: add consumptional fee model ([#1572](https://github.com/pactus-project/pactus/pull/1572))
- **txpool**: calculate consumption when committing a new block ([#1554](https://github.com/pactus-project/pactus/pull/1554))
- **sync**: add metric to track the network activity ([#1552](https://github.com/pactus-project/pactus/pull/1552))
- **wallet**: add wallet service API ([#1548](https://github.com/pactus-project/pactus/pull/1548))
- **config**: add consumption fee configs ([#1547](https://github.com/pactus-project/pactus/pull/1547))

### Fix

- **config**: update TOML parser ([#1592](https://github.com/pactus-project/pactus/pull/1592))
- **gtk**: prevent duplicate address on enter signal in create modal ([#1590](https://github.com/pactus-project/pactus/pull/1590))
- **txpool, cmd, gtk**: broadcast transactions with zero fee ([#1589](https://github.com/pactus-project/pactus/pull/1589))
- **consensus**: send decided vote for previous round on query vote ([#1567](https://github.com/pactus-project/pactus/pull/1567))
- **grpc**: get tx pool content filter by payload type ([#1581](https://github.com/pactus-project/pactus/pull/1581))
- **wallet, cmd**: add support for importing Ed25519 private keys ([#1584](https://github.com/pactus-project/pactus/pull/1584))
- **gtk**: change transactions to transaction in tx link ([#1580](https://github.com/pactus-project/pactus/pull/1580))
- **grpc**: set Bond public key for decoded transaction ([#1577](https://github.com/pactus-project/pactus/pull/1577))
- **other**: add varnamelen linter to improve name convention ([#1568](https://github.com/pactus-project/pactus/pull/1568))
- **grpc**: encode data and signature properly ([#1538](https://github.com/pactus-project/pactus/pull/1538))
- **gtk**: change some text in GUI and pruned position ([#1536](https://github.com/pactus-project/pactus/pull/1536))

### Refactor

- **cmd**: get first account address from wallet as reward address ([#1594](https://github.com/pactus-project/pactus/pull/1594))
- **grpc**: revert GetRawTransfer method and undo deprecation ([#1560](https://github.com/pactus-project/pactus/pull/1560))
- **crypto**: define SerializeSize for PublicKey and Signature ([#1534](https://github.com/pactus-project/pactus/pull/1534))

## [1.5.0](https://github.com/pactus-project/pactus/compare/v1.4.0...v1.5.0) (2024-10-08)

### Feat

- **cmd**: pactus-wallet add info commands ([#1496](https://github.com/pactus-project/pactus/pull/1496))
- **state**: enable Ed25519 for the Testnet ([#1497](https://github.com/pactus-project/pactus/pull/1497))
- **gtk**: support create Ed25519 in gtk ([#1489](https://github.com/pactus-project/pactus/pull/1489))
- **grpc**: add Ed25519 to AddressType proto ([#1492](https://github.com/pactus-project/pactus/pull/1492))
- **wallet**: upgrade wallet ([#1491](https://github.com/pactus-project/pactus/pull/1491))
- **wallet**: supporting Ed25519 curve in wallet ([#1484](https://github.com/pactus-project/pactus/pull/1484))
- **grpc**: add `Proposal` to `ConsensusInfo` API ([#1469](https://github.com/pactus-project/pactus/pull/1469))
- **crypto**: supporting ed25519 ([#1481](https://github.com/pactus-project/pactus/pull/1481))
- **gtk**: adding IsPrune to node widget ([#1470](https://github.com/pactus-project/pactus/pull/1470))
- **daemon**: warn at pruning a prune node attempt ([#1471](https://github.com/pactus-project/pactus/pull/1471))
- **genesis**: separating chain param from genesis param ([#1463](https://github.com/pactus-project/pactus/pull/1463))
- **cmd**: pactus-shell support interactive shell ([#1460](https://github.com/pactus-project/pactus/pull/1460))

### Fix

- **gtk**: increase window width to show availability score ([#1529](https://github.com/pactus-project/pactus/pull/1529))
- **state**: set hard-fork height for the mainnet ([#1528](https://github.com/pactus-project/pactus/pull/1528))
- **wallet**: change to prompt password for masking ([#1527](https://github.com/pactus-project/pactus/pull/1527))
- **deps**: go version docker image to build go v1.23.2 ([#1522](https://github.com/pactus-project/pactus/pull/1522))
- **network**: close stream on timeout ([#1520](https://github.com/pactus-project/pactus/pull/1520))
- **http**: add pprof link in http web interface ([#1518](https://github.com/pactus-project/pactus/pull/1518))
- **sync**: close stream on read error ([#1519](https://github.com/pactus-project/pactus/pull/1519))
- **sync**: set last support version to 1.5.0 ([#1517](https://github.com/pactus-project/pactus/pull/1517))
- **http**: pprof in http server ([#1515](https://github.com/pactus-project/pactus/pull/1515))
- **cmd**: add flag debug to enable pprof ([#1512](https://github.com/pactus-project/pactus/pull/1512))
- **cmd**: add pprof as default in http server ([#1511](https://github.com/pactus-project/pactus/pull/1511))
- **grpc**: merge raw transaction methods to one rpc method ([#1500](https://github.com/pactus-project/pactus/pull/1500))
- **wallet, cmd**: adding ed25519_account in help and set as default ([#1485](https://github.com/pactus-project/pactus/pull/1485))
- **wallet**: add memo in confirmation wallet CLI ([#1499](https://github.com/pactus-project/pactus/pull/1499))
- **store**: cache Ed25519 Public Keys ([#1495](https://github.com/pactus-project/pactus/pull/1495))
- **grpc**: adding pyi files for python generated files ([#1479](https://github.com/pactus-project/pactus/pull/1479))
- **grpc**: change enum type to numeric for documentation ([#1474](https://github.com/pactus-project/pactus/pull/1474))
- **shell**: stop showing usage on error ([#1467](https://github.com/pactus-project/pactus/pull/1467))
- **util**: chunked download to improve download speed ([#1459](https://github.com/pactus-project/pactus/pull/1459))
- **gtk**: width size of listbox and download button ([#1434](https://github.com/pactus-project/pactus/pull/1434))
- **grpc**: add example json-rpc in generated doc ([#1461](https://github.com/pactus-project/pactus/pull/1461))
- **grpc**: add basic check for grpc configuration to check basic auth ([#1455](https://github.com/pactus-project/pactus/pull/1455))
- **util**: remove util.Now helper function ([#1442](https://github.com/pactus-project/pactus/pull/1442))

### Refactor

- **crypto**: replace bls12-381 kilic with gnark ([#1510](https://github.com/pactus-project/pactus/pull/1510))
- **crypto**: define errors for crypto package ([#1507](https://github.com/pactus-project/pactus/pull/1507))
- **sync**: define errors for sync package ([#1504](https://github.com/pactus-project/pactus/pull/1504))
- **types**: define errors for vote package ([#1503](https://github.com/pactus-project/pactus/pull/1503))
- **state**: define errors for state package ([#1457](https://github.com/pactus-project/pactus/pull/1457))
- **util**: remove GenericError code ([#1454](https://github.com/pactus-project/pactus/pull/1454))
- **types**: using options pattern for memo parameter on new tx functions ([#1443](https://github.com/pactus-project/pactus/pull/1443))

## [1.4.0](https://github.com/pactus-project/pactus/compare/v1.3.0...v1.4.0) (2024-08-01)

### Feat

- **cmd**: add node type page to the startup assistant  ([#1431](https://github.com/pactus-project/pactus/pull/1431))
- **grpc**: adding is-pruned and pruning-height to blockchain info API ([#1420](https://github.com/pactus-project/pactus/pull/1420))
- **daemon**: add import command to download pruned snapshots ([#1424](https://github.com/pactus-project/pactus/pull/1424))
- **util**: file downloader with verify sha256 hash ([#1422](https://github.com/pactus-project/pactus/pull/1422))
- **sync**: define full and prune service ([#1412](https://github.com/pactus-project/pactus/pull/1412))
- **pip**: implement PIP-23 ([#1397](https://github.com/pactus-project/pactus/pull/1397))
- **firewall**: check valid gossip and stream messages ([#1402](https://github.com/pactus-project/pactus/pull/1402))
- **state**: prune block on commit ([#1404](https://github.com/pactus-project/pactus/pull/1404))
- **core**: pruning client by prune command ([#1400](https://github.com/pactus-project/pactus/pull/1400))
- **store**: prune block function ([#1399](https://github.com/pactus-project/pactus/pull/1399))
- **wallet**: add timeout client connection ([#1396](https://github.com/pactus-project/pactus/pull/1396))
- add backup tool script ([#1373](https://github.com/pactus-project/pactus/pull/1373))

### Fix

- **consensus**: handle query for decided proposal ([#1438](https://github.com/pactus-project/pactus/pull/1438))
- **gtk**: solve dynamic library dependencies and import path on macOS ([#1435](https://github.com/pactus-project/pactus/pull/1435))
- **cmd**: prevent sudden crash on download error ([#1432](https://github.com/pactus-project/pactus/pull/1432))
- **store**: pruning height returns zero when store is not in prune mode ([#1430](https://github.com/pactus-project/pactus/pull/1430))
- **grpc**: add last-block-time to blockchain-info API ([#1428](https://github.com/pactus-project/pactus/pull/1428))
- **grpc**: show negative pruning height when is pruned false ([#1429](https://github.com/pactus-project/pactus/pull/1429))
- **sync**: fix syncing issue on prune mode ([#1415](https://github.com/pactus-project/pactus/pull/1415))
- **grpc**: return error on invalid arguments for VerifyMessage ([#1411](https://github.com/pactus-project/pactus/pull/1411))
- **network**: accept messages originating from self ([#1408](https://github.com/pactus-project/pactus/pull/1408))
- change wallet rpc ip to dns address ([#1398](https://github.com/pactus-project/pactus/pull/1398))
- **pactus-shell**: pactus shell support basic auth ([#1384](https://github.com/pactus-project/pactus/pull/1384))
- **gui**: support ctrl+c for interrupt gui ([#1385](https://github.com/pactus-project/pactus/pull/1385))
- **grpc**: add basic auth in swagger header ([#1383](https://github.com/pactus-project/pactus/pull/1383))

### Refactor

- **execution**: simplify executors and tests ([#1425](https://github.com/pactus-project/pactus/pull/1425))

## [1.3.0](https://github.com/pactus-project/pactus/compare/v1.2.0...v1.3.0) (2024-06-27)

### Feat

- **grpc**: get txpool content API ([#1364](https://github.com/pactus-project/pactus/pull/1364))
- **network**: permanent peer store ([#1354](https://github.com/pactus-project/pactus/pull/1354))

### Fix

- **grpc**: change bytes type to hex string ([#1371](https://github.com/pactus-project/pactus/pull/1371))
- **http**: add basic auth middleware for http server ([#1372](https://github.com/pactus-project/pactus/pull/1372))
- **network**: use goroutines for sending streams ([#1365](https://github.com/pactus-project/pactus/pull/1365))

## [1.2.0](https://github.com/pactus-project/pactus/compare/v1.1.0...v1.2.0) (2024-06-20)

### Feat

- **config**: make minimum fee configurable ([#1349](https://github.com/pactus-project/pactus/pull/1349))
- apply rate limit for the network topics ([#1332](https://github.com/pactus-project/pactus/pull/1332))
- add ipblocker package ([#1323](https://github.com/pactus-project/pactus/pull/1323))
- **consensus**: fast consensus path implementation ([#1253](https://github.com/pactus-project/pactus/pull/1253))
- **version**: add alias to node version ([#1281](https://github.com/pactus-project/pactus/pull/1281))
- **ntp**: add ntp util ([#1274](https://github.com/pactus-project/pactus/pull/1274))
- **gRPC**: add connection info to node info ([#1273](https://github.com/pactus-project/pactus/pull/1273))
- **gRPC**: add only_connected parameter to getNetworkInfo API ([#1264](https://github.com/pactus-project/pactus/pull/1264))
- **grpc**: refactor CreateWallet and add RestoreWallet API endpoint ([#1256](https://github.com/pactus-project/pactus/pull/1256))
- add wallet service ([#1241](https://github.com/pactus-project/pactus/pull/1241))
- ban attacker validators ([#1235](https://github.com/pactus-project/pactus/pull/1235))
- **txpool**: prevent spamming transactions by defining a minimum value ([#1233](https://github.com/pactus-project/pactus/pull/1233))
- reject direct message from non-supporting agents ([#1225](https://github.com/pactus-project/pactus/pull/1225))

### Fix

- **wallet**: add public key on get new address ([#1350](https://github.com/pactus-project/pactus/pull/1350))
- **sync**: add IsBannedAddress check in processing connect event ([#1347](https://github.com/pactus-project/pactus/pull/1347))
- **sync**: update latest supporting version ([#1336](https://github.com/pactus-project/pactus/pull/1336))
- **state**: improve node startup by optimizing availability score calculation ([#1338](https://github.com/pactus-project/pactus/pull/1338))
- **HTTP**: add clock offset and connection info to node-info API ([#1334](https://github.com/pactus-project/pactus/pull/1334))
- **grpc**: add stacktrace to locate panic ([#1333](https://github.com/pactus-project/pactus/pull/1333))
- **consensus**: implement PIP-26 ([#1331](https://github.com/pactus-project/pactus/pull/1331))
- **grpc**: improve grpc server and client ([#1330](https://github.com/pactus-project/pactus/pull/1330))
- **util**: add more ntp pool ([#1328](https://github.com/pactus-project/pactus/pull/1328))
- **jsonrpc**: update JSON-RPC Gateway to support headers and improve client registry ([#1327](https://github.com/pactus-project/pactus/pull/1327))
- **consensus**: improve consensus alghorithm ([#1329](https://github.com/pactus-project/pactus/pull/1329))
- **txpool**: set fix fee of 0.1 PAC for transactions ([#1320](https://github.com/pactus-project/pactus/pull/1320))
- **network**: add block and transaction topics ([#1319](https://github.com/pactus-project/pactus/pull/1319))
- **gRPC**: prevent concurrent map iteration and map write ([#1279](https://github.com/pactus-project/pactus/pull/1279))
- **api**: add swagger schemes ([#1270](https://github.com/pactus-project/pactus/pull/1270))
- **network**: set infinite limit for resource manager  ([#1261](https://github.com/pactus-project/pactus/pull/1261))
- **sync**: introduce session manager ([#1257](https://github.com/pactus-project/pactus/pull/1257))
- **HTTP**: using amount type for fee in transaction details ([#1255](https://github.com/pactus-project/pactus/pull/1255))
- **network**: disconnect from peers that has no protocol ([#1243](https://github.com/pactus-project/pactus/pull/1243))
- **wallet**: saving wallet file after generating new address in gRPC ([#1236](https://github.com/pactus-project/pactus/pull/1236))
- prevent zero stake for bond transactions ([#1227](https://github.com/pactus-project/pactus/pull/1227))
- set bounding interval for first boudning tx only ([#1224](https://github.com/pactus-project/pactus/pull/1224))

### Refactor

- **wallet**: set server address on loading wallet ([#1348](https://github.com/pactus-project/pactus/pull/1348))
- removed deprecated LockWallet and UnLockWallet from WalletService ([#1343](https://github.com/pactus-project/pactus/pull/1343))
- **crypto**: decode data to point on verification ([#1339](https://github.com/pactus-project/pactus/pull/1339))
- **network**: define connection info in network proto ([#1297](https://github.com/pactus-project/pactus/pull/1297))
- **sync**: define peer package ([#1271](https://github.com/pactus-project/pactus/pull/1271))
- **network**: refactor peer manager and redefine the min cons ([#1259](https://github.com/pactus-project/pactus/pull/1259))

## [1.1.0](https://github.com/pactus-project/pactus/compare/v1.0.0...v1.1.0) (2024-04-14)

### Feat

- **gRPC**: add get address history method ([#1206](https://github.com/pactus-project/pactus/pull/1206))
- **grpc**: Add GetNewAddress/GetTotalBalance endpoint to gateway ([#1197](https://github.com/pactus-project/pactus/pull/1197))
- **GUI**: adding total balance to wallet widget ([#1194](https://github.com/pactus-project/pactus/pull/1194))
- Add GetNewAddress gRPC API ([#1193](https://github.com/pactus-project/pactus/pull/1193))
- **gRPC**: add new API to get the total balance of wallet ([#1190](https://github.com/pactus-project/pactus/pull/1190))
- **GUI**: showing transaction hash after broadcasting transaction ([#1187](https://github.com/pactus-project/pactus/pull/1187))
- add jsonrpc gateway support ([#1183](https://github.com/pactus-project/pactus/pull/1183))
- **config**: one reward address in config for all validators ([#1178](https://github.com/pactus-project/pactus/pull/1178))
- **GUI**: memo field for transactions on GUI wallet ([#1182](https://github.com/pactus-project/pactus/pull/1182))
- implement basic auth for pactus shell ([#1177](https://github.com/pactus-project/pactus/pull/1177))
- **grpc**: add rust code gen for proto ([#1151](https://github.com/pactus-project/pactus/pull/1151))
- **testnet**: define permanent Testent genesis ([#1173](https://github.com/pactus-project/pactus/pull/1173))
- add basic auth authentication for securing grpc ([#1162](https://github.com/pactus-project/pactus/pull/1162))
- **grpc**: calculate fee for create-raw-transaction APIs ([#1159](https://github.com/pactus-project/pactus/pull/1159))
- **grpc**: add fixed-amount to calc-fee API ([#1146](https://github.com/pactus-project/pactus/pull/1146))
- **wallet**: adding all account address functions ([#1128](https://github.com/pactus-project/pactus/pull/1128))
- **grpc**: update swagger API to version 1.1 ([#1106](https://github.com/pactus-project/pactus/pull/1106))
- **GUI**: adding availability score in wallet ([#1118](https://github.com/pactus-project/pactus/pull/1118))
- **logger**: adding log target ([#1122](https://github.com/pactus-project/pactus/pull/1122))
- **logger**: adding file_only option ([#1117](https://github.com/pactus-project/pactus/pull/1117))
- **gui**: add connections and moniker fields to main windows ([#1090](https://github.com/pactus-project/pactus/pull/1090))
- implementation for PIP-22 ([#1067](https://github.com/pactus-project/pactus/pull/1067))
- generate documentation for proto files ([#1064](https://github.com/pactus-project/pactus/pull/1064))
- pactus-ctl ([#946](https://github.com/pactus-project/pactus/pull/946))

### Fix

- **cmd**: ignore error on balance query ([#1220](https://github.com/pactus-project/pactus/pull/1220))
- **gRPC**: add basic auth option in header ([#1217](https://github.com/pactus-project/pactus/pull/1217))
- **gRPC**: not return block data on information verbosity ([#1212](https://github.com/pactus-project/pactus/pull/1212))
- **wallet**: fix wallet conn issue ([#1211](https://github.com/pactus-project/pactus/pull/1211))
- **GUI**: update total balance on wallet timeout ([#1204](https://github.com/pactus-project/pactus/pull/1204))
- accept small bond to existing validator ([#1152](https://github.com/pactus-project/pactus/pull/1152))
- **GUI**: make transaction hash selactable ([#1196](https://github.com/pactus-project/pactus/pull/1196))
- close connections with peers that have no supported protocol ([#1181](https://github.com/pactus-project/pactus/pull/1181))
- **sync**: check the start block request height ([#1176](https://github.com/pactus-project/pactus/pull/1176))
- **config**: load logger levels in Mainnet config ([#1168](https://github.com/pactus-project/pactus/pull/1168))
- **gRPC**: pactus swagger not found ([#1163](https://github.com/pactus-project/pactus/pull/1163))
- add error type for invalid configs ([#1153](https://github.com/pactus-project/pactus/pull/1153))
- save Mainnet config with inline comments ([#1143](https://github.com/pactus-project/pactus/pull/1143))
- **network**: set deadline for streams ([#1149](https://github.com/pactus-project/pactus/pull/1149))
- **grpc**: fix error 404 on grpc gateway ([#1144](https://github.com/pactus-project/pactus/pull/1144))
- **wallet**: checking args in history add ([#1141](https://github.com/pactus-project/pactus/pull/1141))
- **gRPC**: adding sign raw transaction API to gateway ([#1116](https://github.com/pactus-project/pactus/pull/1116))
- **sync**: fix concurrent map read-write crash ([#1112](https://github.com/pactus-project/pactus/pull/1112))
- **network**: remove disconnected peers from peerMgr ([#1110](https://github.com/pactus-project/pactus/pull/1110))
- **network**: set dial and accept limit in connection gater ([#1089](https://github.com/pactus-project/pactus/pull/1089))
- stderr logger in windows os ([#1081](https://github.com/pactus-project/pactus/pull/1081))
- **sync**: improve syncing process ([#1087](https://github.com/pactus-project/pactus/pull/1087))
- **network**: redefine resource limits ([#1086](https://github.com/pactus-project/pactus/pull/1086))

### Refactor

- **sync**: improve syncing process ([#1207](https://github.com/pactus-project/pactus/pull/1207))
- move fee calculation logic to execution package  ([#1195](https://github.com/pactus-project/pactus/pull/1195))
- introduce Amount data type for converting PAC units ([#1174](https://github.com/pactus-project/pactus/pull/1174))
- using PAC instead of atomic units for external input/outputs ([#1161](https://github.com/pactus-project/pactus/pull/1161))
- change func() to cancel func type ([#1142](https://github.com/pactus-project/pactus/pull/1142))

## [1.0.0](https://github.com/pactus-project/pactus/compare/v0.20.0...v1.0.0) (2024-01-31)

### Feat

- implement get validator address for grpc ([#975](https://github.com/pactus-project/pactus/pull/975))
- add bootstrap.json and load in config on build ([#964](https://github.com/pactus-project/pactus/pull/964))
- add mainnet config and genesis files ([#951](https://github.com/pactus-project/pactus/pull/951))
- add Consensus-address to network info ([#952](https://github.com/pactus-project/pactus/pull/952))
- **grpc**: sign transaction using wallet client ([#945](https://github.com/pactus-project/pactus/pull/945))
- pactus gui lock support ([#947](https://github.com/pactus-project/pactus/pull/947))
- **consensus**: turning consensus to a zero-config module ([#942](https://github.com/pactus-project/pactus/pull/942))

### Fix

- localnet wallet issue ([#970](https://github.com/pactus-project/pactus/pull/970))
- **sync**: remove ReachabilityStatus from agent info ([#956](https://github.com/pactus-project/pactus/pull/956))
- **daemon**: keeping previous behavior for password flag, linting CLI messages ([#950](https://github.com/pactus-project/pactus/pull/950))
- **consensus**: detect if the system time is behind the network ([#939](https://github.com/pactus-project/pactus/pull/939))

## [0.20.0](https://github.com/pactus-project/pactus/compare/v0.19.0...v0.20.0) (2024-01-11)

### Feat

- implement relay service ([#931](https://github.com/pactus-project/pactus/pull/931))
- **HTTP**: Integrate AddRowDouble and update tests ([#926](https://github.com/pactus-project/pactus/pull/926))
- **network**: making listen address private in config ([#921](https://github.com/pactus-project/pactus/pull/921))
- **http**: adding AvailabilityScore to http module ([#917](https://github.com/pactus-project/pactus/pull/917))
- **network**: adding 'enable_udp' config ([#918](https://github.com/pactus-project/pactus/pull/918))
- **network**: removing gossip node service ([#916](https://github.com/pactus-project/pactus/pull/916))
- **gRPC**: adding AvailabilityScore to gRPC ([#910](https://github.com/pactus-project/pactus/pull/910))
- **GUI**: unbond and withdraw transaction dialogs ([#908](https://github.com/pactus-project/pactus/pull/908))

### Fix

- **gRPC**: adding missing get raw transaction APIs to gRPC gateway ([#925](https://github.com/pactus-project/pactus/pull/925))
- **network**: preventing self dial ([#924](https://github.com/pactus-project/pactus/pull/924))
- fixing time lag on starting node ([#923](https://github.com/pactus-project/pactus/pull/923))
- **network**: fixing network deadlock on linux arm64 ([#922](https://github.com/pactus-project/pactus/pull/922))
- **GUI**: updating unbond and withdraw dialogs ([#911](https://github.com/pactus-project/pactus/pull/911))
- fixing gRPC node info issue ([#906](https://github.com/pactus-project/pactus/pull/906))

## [0.19.0](https://github.com/pactus-project/pactus/compare/v0.18.0...v0.19.0) (2024-01-04)

### Feat

- **gRPC**: defining network and peers info response's properly ([#898](https://github.com/pactus-project/pactus/pull/898))
- implementing pip-19 ([#899](https://github.com/pactus-project/pactus/pull/899))
- **network**: disabling GosipSub, only FloodSub ([#895](https://github.com/pactus-project/pactus/pull/895))
- **www**: adding change proposer round and value to consensus info votes ([#892](https://github.com/pactus-project/pactus/pull/892))
- **network**: adding relay service to dial relay nodes ([#887](https://github.com/pactus-project/pactus/pull/887))
- implementing pip-15 ([#843](https://github.com/pactus-project/pactus/pull/843))
- check already running by lock file ([#871](https://github.com/pactus-project/pactus/pull/871))

### Fix

- **store**: use cache to check if public key exists ([#902](https://github.com/pactus-project/pactus/pull/902))
- **executor**: not rejecting bond transaction for bootstrap validator ([#901](https://github.com/pactus-project/pactus/pull/901))
- **GUI**: removing unnecessary tags in transaction confirm dialog ([#893](https://github.com/pactus-project/pactus/pull/893))
- **network**: close relay connection for public node ([#891](https://github.com/pactus-project/pactus/pull/891))
- **network**: refining GossipSubParams for Gossiper node ([#882](https://github.com/pactus-project/pactus/pull/882))
- **sync**: adding sequence number to the bundle ([#881](https://github.com/pactus-project/pactus/pull/881))
- **network**: turn off mesh for gossiper node ([#880](https://github.com/pactus-project/pactus/pull/880))
- **consensus**: check voteset for CP strong termination ([#879](https://github.com/pactus-project/pactus/pull/879))
- adding querier to query messages ([#878](https://github.com/pactus-project/pactus/pull/878))
- **execution**: fixing issue #869 ([#870](https://github.com/pactus-project/pactus/pull/870))
- fixing logger issue on rotating log file ([#859](https://github.com/pactus-project/pactus/pull/859))

## [0.18.0](https://github.com/pactus-project/pactus/compare/v0.17.0...v0.18.0) (2023-12-11)

### Feat

- implement pip-14 ([#841](https://github.com/pactus-project/pactus/pull/841))
- sort wallet addresses ([#836](https://github.com/pactus-project/pactus/pull/836))
- **grpc**: endpoints for creating raw transaction ([#838](https://github.com/pactus-project/pactus/pull/838))
- network reachability API ([#834](https://github.com/pactus-project/pactus/pull/834))
- implement pip-13 ([#835](https://github.com/pactus-project/pactus/pull/835))
- subscribing to libp2p eventbus ([#831](https://github.com/pactus-project/pactus/pull/831))
- implement helper methods for wallet address path ([#830](https://github.com/pactus-project/pactus/pull/830))
- **logger**: adding rotate log file after days, compress and max backups for logger config ([#822](https://github.com/pactus-project/pactus/pull/822))
- enable bandwidth router metric ([#819](https://github.com/pactus-project/pactus/pull/819))

### Fix

- **network**: refining the connection limit ([#849](https://github.com/pactus-project/pactus/pull/849))
- corrected mistake when retrieving the reward address ([#848](https://github.com/pactus-project/pactus/pull/848))
- **config**: restore default config when it is deleted ([#847](https://github.com/pactus-project/pactus/pull/847))
- **cmd**: changing home directory for root users ([#846](https://github.com/pactus-project/pactus/pull/846))
- removing BasicCheck for hash ([#845](https://github.com/pactus-project/pactus/pull/845))
- disabling libp2p ping protocol ([#844](https://github.com/pactus-project/pactus/pull/844))
- build docker file ([#839](https://github.com/pactus-project/pactus/pull/839))
- **sync**: ignore publishing a block if it is received before ([#829](https://github.com/pactus-project/pactus/pull/829))
- **network**: subscribing to the Libp2p event bus ([#828](https://github.com/pactus-project/pactus/pull/828))
- **sync**: ignore block request if blocks are already inside the cache ([#817](https://github.com/pactus-project/pactus/pull/817))

## [0.17.0](https://github.com/pactus-project/pactus/compare/v0.16.0...v0.17.0) (2023-11-12)

### Feat

- **network**: default configs for bootstrap and relay peers ([#812](https://github.com/pactus-project/pactus/pull/812))
- introducing node gossip type ([#811](https://github.com/pactus-project/pactus/pull/811))
- **sync**: adding remote address to the peer info ([#804](https://github.com/pactus-project/pactus/pull/804))
- **network**: adding public address to factory ([#795](https://github.com/pactus-project/pactus/pull/795))
- **network**: filter private ips ([#793](https://github.com/pactus-project/pactus/pull/793))

### Fix

- upgrading Testnet ([#814](https://github.com/pactus-project/pactus/pull/814))
- **sync**: prevent opening sessions indefinitely ([#813](https://github.com/pactus-project/pactus/pull/813))
- **execution**: fixing mistake on calculating unbonded power ([#806](https://github.com/pactus-project/pactus/pull/806))
- **network**: check connection threshold on gater ([#803](https://github.com/pactus-project/pactus/pull/803))
- **network**: no transient connection ([#799](https://github.com/pactus-project/pactus/pull/799))
- not close connection for bootstrap nodes ([#792](https://github.com/pactus-project/pactus/pull/792))

### Refactor

- **sync**: refactoring sync process ([#807](https://github.com/pactus-project/pactus/pull/807))

## [0.16.0](https://github.com/pactus-project/pactus/compare/v0.15.0...v0.16.0) (2023-10-29)

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

## [0.15.0](https://github.com/pactus-project/pactus/compare/v0.13.0...v0.15.0) (2023-10-15)

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

## [0.13.0](https://github.com/pactus-project/pactus/compare/v0.12.0...v0.13.0) (2023-06-30)

### Fix

- implemented restore wallet base on input seed ([#553](https://github.com/pactus-project/pactus/pull/553))
- measuring total sent and received bytes ([#552](https://github.com/pactus-project/pactus/pull/552))
- add validate seed and restore wallet ([#533](https://github.com/pactus-project/pactus/pull/533))
- **HTTP**: Missing handlers ([#549](https://github.com/pactus-project/pactus/pull/549))
- **gui**: update about dialog and menu items in help ([#532](https://github.com/pactus-project/pactus/pull/532))

### Refactor

- implementing TestSuite ([#535](https://github.com/pactus-project/pactus/pull/535))

## [0.12.0](https://github.com/pactus-project/pactus/compare/v0.11.0...v0.12.0) (2023-06-19)

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
