# zarb tx

`zarb tx` is a simple command-line tool to make a raw transaction for Zarb offline.
After creating the raw transaction, you can sign it using `zarb key sign -t <raw_transaction>`
The signed transaction can broadcast through `grpc/send_raw_transaction` API.


## Usage

### Send transaction

To create a send transaction you use `zarb tx send` command.

> if the seq, stamp and endpoint were not specified then it will prompt for endpoint to pull these from gRPC server

> if endpoint was supplied then it will publish the transaction otherwise it will just print signed transaction and exit

Example:
```bash
$ zarb tx send --seq=10 --sender=zrb1x8qy6v8lr0x5uxn0lp4aygxh44wtdrz6y82jxd --receiver=zrb1team0xhxarezhy96z6yt9kkpztrn8f8kmpndm0 -k ./build/6/validator_key.json --amount=123000 --stamp=17913ea30d60133f0cc65f69bced0547ed89318b78994e9b1f3fc8c2bf33e067 --fee=123 -e=localhost:9010

 ███████╗  █████╗  ██████╗  ██████╗
 ╚══███╔╝ ██╔══██╗ ██╔══██╗ ██╔══██╗
   ███╔╝  ███████║ ██████╔╝ ██████╔╝
  ███╔╝   ██╔══██║ ██╔══██╗ ██╔══██╗
 ███████╗ ██║  ██║ ██║  ██║ ██████╔╝
 ╚══════╝ ╚═╝  ╚═╝ ╚═╝  ╚═╝ ╚═════╝

Passphrase: 
[WARN] Your transaction:
{
   "Version": 1,
   "Stamp": "17913ea30d60133f0cc65f69bced0547ed89318b78994e9b1f3fc8c2bf33e067",
   "Sequence": 10,
   "Fee": 123,
   "Type": 1,
   "Payload": {
      "Sender": "zrb1x8qy6v8lr0x5uxn0lp4aygxh44wtdrz6y82jxd",
      "Receiver": "zrb1team0xhxarezhy96z6yt9kkpztrn8f8kmpndm0",
      "Amount": 123000
   },
   "Memo": "",
   "PublicKey": "6df01b4b4f49b26692d83add4bf9a47c8a3b5db2f5000b80a399a9b1b6afe04f8afd6749354e3f766c877b2837351004a279f4834dd532018766c0446cec1d1903735d52cafdb5ad2c61764fe89da05d139f7efe5f049d8ec92727ba93c74595",
   "Signature": "f7f9a5c54dc9ed6248732401b8d38ae6039264f8bfe3a5cc94b25c053d1d53fd6bc5e92c8874e123006b84c2fa38eb10"
}

This operation is "not reversible". Are you sure [yes/no]? 

```

```bash
$ zarb tx send --seq=[Senders Sequance Number] --sender=[Senders Address] --receiver=[Recivers Address] -k=[Senders Key File Path] --amount=[Amount To Send] --stamp=[BlockChains Latest Blocks Hash] --fee=[Fee Willing To Pay For This Transaction ]

 ███████╗  █████╗  ██████╗  ██████╗ 
 ╚══███╔╝ ██╔══██╗ ██╔══██╗ ██╔══██╗
   ███╔╝  ███████║ ██████╔╝ ██████╔╝
  ███╔╝   ██╔══██║ ██╔══██╗ ██╔══██╗
 ███████╗ ██║  ██║ ██║  ██║ ██████╔╝
 ╚══════╝ ╚═╝  ╚═╝ ╚═╝  ╚═╝ ╚═════╝ 

Passphrase: 
[WARN] Your transaction:
{
   "Version": 1,
   "Stamp": "17913ea30d60133f0cc65f69bced0547ed89318b78994e9b1f3fc8c2bf33e067",
   "Sequence": 10,
   "Fee": 123,
   "Type": 1,
   "Payload": {
      "Sender": "zrb1x8qy6v8lr0x5uxn0lp4aygxh44wtdrz6y82jxd",
      "Receiver": "zrb1team0xhxarezhy96z6yt9kkpztrn8f8kmpndm0",
      "Amount": 123000
   },
   "Memo": "",
   "PublicKey": "[Senders Public Key]",
   "Signature": "[Payloads Signature Generated From Senders Private Key]"
}

raw signed transaction payload:
a8010102582017913ea30d60133f0cc65f69bced0547ed89318b78994e9b1f3fc8c2bf33e067
[..snip...]
1d53fd6bc5e92c8874e123006b84c2fa38eb10
```

### Bond transaction

To create a bond transaction you use `zarb tx bond` command.


> if the seq, stamp and endpoint were not specified then it will prompt for endpoint to pull these from gRPC server

> if endpoint was supplied then it will publish the transaction otherwise it will just print signed transaction and exit

Example:
```bash
$ zarb tx bond --stake=[Amount of Stack To be Sent] --fee=[Fee Willing To Pay For This Transaction ] --bonder=[Address Of Account Will Pay For Stack And Fee, And Will Sign This Transaction] --pub=[Public Key of Validator To Bond Stack To] -k=[Senders Key File Path] -e [gRPC Endpoint Address]


 ███████╗  █████╗  ██████╗  ██████╗
 ╚══███╔╝ ██╔══██╗ ██╔══██╗ ██╔══██╗
   ███╔╝  ███████║ ██████╔╝ ██████╔╝
  ███╔╝   ██╔══██║ ██╔══██╗ ██╔══██╗
 ███████╗ ██║  ██║ ██║  ██║ ██████╔╝
 ╚══════╝ ╚═╝  ╚═╝ ╚═╝  ╚═╝ ╚═════╝

Passphrase: 
[WARN] Your transaction:
{
   "Version": 1,
   "Stamp": "db5057350d920eaf855cfbdc8ce46e195d11384d364a528597dd1702c1aaad82",
   "Sequence": 0,
   "Fee": 2500000,
   "Type": 2,
   "Payload": {
      "Bonder": "zrb1h87hfkn3wa36xwypjz8aep3hu4ssdrt86chs3c",
      "Validator": "6df01b4b4f49b26692d83add4bf9a47c8a3[...snip...]7ba93c74595",
      "Stake": 2500000000
   },
   "Memo": "",
   "PublicKey": "6df01b4b4f49b26692d[...snip...]ba93c74595",
   "Signature": "b598a9d4e284eeb85714e8d638679af815885a24916f751465cef15af0c72cfe2c082103b477ad05ff401fbe3130c186"
}

This operation is "not reversible". Are you sure [yes/no]?
```
```bash
$ zarb tx bond ---stake=[Amount of Stack To be Sent] --fee=[Fee Willing To Pay For This Transaction ] --bonder=[Address Of Account Will Pay For Stack And Fee, And Will Sign This Transaction] --pub=[Public Key of Validator To Bond Stack To] -k=[Senders Key File Path] --seq=[Senders Sequance Number] --stamp=[BlockChains Latest Blocks Hash]                                                                    

 ███████╗  █████╗  ██████╗  ██████╗ 
 ╚══███╔╝ ██╔══██╗ ██╔══██╗ ██╔══██╗
   ███╔╝  ███████║ ██████╔╝ ██████╔╝
  ███╔╝   ██╔══██║ ██╔══██╗ ██╔══██╗
 ███████╗ ██║  ██║ ██║  ██║ ██████╔╝
 ╚══════╝ ╚═╝  ╚═╝ ╚═╝  ╚═╝ ╚═════╝ 

Passphrase: 
←[33m[WARN] Your transaction:←[0m
{
   "Version": 1,
   "Stamp": "db5057350d920eaf855cfbdc8ce46e195d11384d364a528597dd1702c1aaad82",
   "Sequence": 0,
   "Fee": 2500000,
   "Type": 2,
   "Payload": {
      "Bonder": "zrb1h87hfkn3wa36xwypjz8aep3hu4ssdrt86chs3c",
      "Validator": "6df01b4b4f49b[...snip...]c74595",
      "Stake": 2500000000
   },
   "Memo": "",
   "PublicKey": "6df01b4b4f49b266[...snip...]c92727ba93c74595",
   "Signature": "b598a9d4e284eeb85714e8d638679af815885a24916f751465cef15af0c72cfe2c082103b477ad05ff401fbe3130c186"
}

raw signed transaction payload:
a80101025820db5057350d920eaf855cfbdc8ce46e195d11384d364a528597dd1702c1aaad820300041a002625a0050206a30154b9fd74da717763a33881908fdc8637e561068d670258606df01b4b4f49b26692d
[...snip...]
1764fe89da05d139f7efe5f049d8ec92727ba93c74595155830b598a9d4e284eeb85714e8d638679af815885a24916f751465cef15af0c72cfe2c082103b477ad05ff401fbe3130c186
```