# Installing Zarb

## Requirements

You need to make sure you have installed [Git](https://git-scm.com/downloads)
and [Go 1.19 or higher](https://golang.org/) in your machine.
If you want to install GUI application, make sure you have install
[GTK+3](https://www.gtk.org/docs/getting-started/) as well.

## Compiling the code

Follow these steps to compile and build Zarb:

```bash
git clone https://github.com/zarbchain/zarb-go.git
cd zarb-go
make build
```

This will be compile `zarb-daemon` and `zarb-wallet` in your machine.
Make sure Zarb is properly compiled and installed in your machine:

```
cd build
./zarb-daemon version
```

If you want to compile the GUI application run this command in the root folder:

```bash
make build_gui
```

To run the tests use this command:

```bash
make test
```

This may takes several minutes to finish.


## What is zarb-daemon

`zarb-daemon` is a full node implementation of Zarb blockchain.
You can use `zarb-daemon` to run a full node on Zarb blockchain.

```bash
./zarb-daemon init  -w=<working_dir>
./zarb-daemon start -w=<working_dir>
```

### Testnet

To join the TestNet, first you need to create a working directory
and then start the node:

```bash
./zarb-daemon init  -w=<working_dir> --testnet
./zarb-daemon start -w=<working_dir>
```

### Local net

You can create a local node with one validator to test Zarb in your machine:

 ```bash
 ./zarb-daemon init -w=<working_dir> --localnet
 ./zarb-daemon start -w=<working_dir>
 ```

## What is zarb-wallet

Zarb wallet is a native wallet in Zarb blockchain and let user to easily manage
their accounts on Zarb blockchain.

### Getting started

To create a new wallet run this command. The wallet will be encrypted by the
provided password.

```bash
./zarb-wallet -w=<wallet_path> create
```

You can create new address like this:

```bash
./zarb-wallet -w=<wallet_path> address new
```

List of addresses are available by this command:

```bash
./zarb-wallet -w=<wallet_path> address balance all
```

You can check the balance:

```bash
./zarb-wallet -w=<wallet_path> address balance <ADDRESS>
```

To publish a transactions use tx subcommand:

```bash
./zarb-wallet -w=<wallet_path> tx
```

For example, to publish a bond transaction:

```bash
./zarb-wallet -w=<wallet_path> tx bond <FROM> <TO> <STAKE>
```

You can recover a wallet if you have the seed phrase.

```bash
./zarb-wallet -w=<wallet_path> recover
```


## Usage of Docker

You can run the Zarb using docker file. Please make sure you have installed
[docker](https://docs.docker.com/engine/install/) in your machine.

Pull the docker from docker hub.

```bash
docker pull zarb/zarb
```

Let's create a working directory at `~/zarb/testnet` for the testnet:

```bash
docker run -it --rm -v ~/zarb/testnet:/zarb zarb/zarb init -w /zarb --testnet
```

Now we can run the zarb and join the testnet:

```bash
docker run -it -v ~/zarb/testnet:/zarb -p 8080:8080 -p 21777:21777 --name zarb-testnet zarb/zarb start -w /zarb
```

check "[http://localhost:8080](http://localhost:8080)" for the list of APIs.

Also you can stop/start docker:

```
docker stop zarb-testnet
docker start zarb-testnet
```

Or check the logs:
```
docker logs zarb-testnet --tail 1000 -f
```
