# Installing Pactus

## Requirements

You need to make sure you have installed [Git](https://git-scm.com/downloads)
and [Go 1.21 or higher](https://golang.org/) on your machine.
If you want to install a GUI application, make sure you have installed
[GTK+3](https://www.gtk.org/docs/getting-started/) as well.

## Compiling the code

Follow these steps to compile and build Pactus:

```bash
git clone https://github.com/pactus-project/pactus.git
cd pactus
make build
```

This will compile `pactus-daemon`, `pactus-wallet` and `pactus-shell` on your machine.

```bash
cd build
./pactus-daemon version
```

If you want to compile the GUI application, run this command in the root folder:

```bash
make build_gui
```

To run the tests, use this command:

```bash
make test
```

This may take several minutes to finish.

## What is `pactus-daemon`?

`pactus-daemon` is a full node implementation of the Pactus blockchain.
You can use `pactus-daemon` to run a full node:

```bash
./pactus-daemon init  -w=<working_dir>
./pactus-daemon start -w=<working_dir>
```

### Testnet

To join the Testnet, first you need to initialize your node
and then start the node:

```bash
./pactus-daemon init  -w=<working_dir> --testnet
./pactus-daemon start -w=<working_dir>
```

### Localnet

You can create a local node to set up a local network for development purposes on your machine:

```bash
./pactus-daemon init  -w=<working_dir> --localnet
./pactus-daemon start -w=<working_dir>
```

## What is `pactus-wallet`?

`pactus-wallet` is the CLI tool for creating and managing wallets on the Pactus blockchain.

### Getting started

To create a new wallet, run this command. The wallet will be encrypted by the
provided password.

```bash
./pactus-wallet --path ~/pactus/wallets/wallet_1 create
```

You can create a new address like this:

```bash
./pactus-wallet --path ~/pactus/wallets/wallet_1 address new
```

A list of addresses is available with this command:

```bash
./pactus-wallet --path ~/pactus/wallets/wallet_1 address all
```

To obtain the public key of an address, run this command:

```bash
./pactus-wallet --path ~/pactus/wallets/wallet_1 address pub <ADDRESS>
```

To send a transaction, use the `send` subcommand.
For example, to send a bond transaction:

```bash
./pactus-wallet --path ~/pactus/wallets/wallet_1 send bond <FROM> <TO> <AMOUNT>
```

You can recover a wallet if you have the seed phrase.

```bash
./pactus-wallet --path ~/pactus/wallets/wallet_2 recover
```

## What is `pactus-shell`?

`pactus-shell` is an interactive command-line client for exploring and calling the Pactus gRPC APIs.

Start it against your node (default gRPC address is `localhost:50051`):

```bash
./pactus-shell interactive --server-addr localhost:50051
```

## Docker

You can run Pactus using a Docker image. Please make sure you have installed
[Docker](https://docs.docker.com/engine/install/) on your machine.

Pull the image from Docker Hub:

```bash
docker pull pactus/pactus:main
```

Let's create a working directory at `~/pactus/testnet` for the testnet:

```bash
docker run -it --rm -v ~/pactus/testnet:/root/pactus pactus/pactus:main pactus-daemon init --testnet
```

Now we can run Pactus and join the testnet:

```bash
docker run -it -v ~/pactus/testnet:/root/pactus -p 8080:8080 -p 21777:21777 --name pactus-testnet pactus/pactus:main pactus-daemon start
```

Check "[http://localhost:8080](http://localhost:8080)" for the list of APIs.

You can stop or start the container:

```bash
docker start pactus-testnet
docker stop pactus-testnet
```

Or check the logs:

```bash
docker logs pactus-testnet --tail 1000 -f
```

## Profiling with pprof

If you need runtime profiling, enable the HTML server with pprof in your node configuration and run the node.

Once running, you can collect and explore a CPU profile with the pprof web UI (replace the host and port with your HTML server address):

```bash
go tool pprof -http :3000 "http://localhost:8080/debug/pprof/profile?debug=1"
```