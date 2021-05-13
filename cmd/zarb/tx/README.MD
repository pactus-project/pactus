# zarb tx

`zarb tx` is a simple command-line tool to make a raw transaction for Zarb offline.
After creating the raw transaction, you can sign it using `zarb key sign -t <raw_transaction>`
The signed transaction can broadcast through `grpc/send_raw_transaction` API.


## Usage

### Send transaction

To create a send transaction you use `zarb tx send` command.

Example:
```
zarb tx send --seq=10 --sender=zrb1x8qy6v8lr0x5uxn0lp4aygxh44wtdrz6y82jxd --receiver=zrb1team0xhxarezhy96z6yt9kkpztrn8f8kmpndm0 --amount=123000 --stamp=17913ea30d60133f0cc65f69bced0547ed89318b78994e9b1f3fc8c2bf33e067 --fee=123
```

### Bond transaction

To create a bond transaction you use `zarb tx bond` command.

Example:
```
zarb tx bond --seq=10 --bonder=zrb1x8qy6v8lr0x5uxn0lp4aygxh44wtdrz6y82jxd --pub=806140b991a8aac9e9b0271e8f0c20c59bd1f101a71a26203299627e8f5a57d71bda69d7c97b1b71ec3d8c1e0458b8087724dbcba7fd83c620cd2e6b4da189f3e623ad14fcfab845f08fc70c85e7518336359f87d9fb8efd874b3a2752d47698 --stake=123000 --stamp=17913ea30d60133f0cc65f69bced0547ed89318b78994e9b1f3fc8c2bf33e067 --fee 123
```