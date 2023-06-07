# Installing Pactus

## Requirements

You need to make sure you have installed [Git](https://git-scm.com/downloads)
and [Go 1.19 or higher](https://golang.org/) on your machine.
If you want to install a GUI application, make sure you have installed
[GTK+3](https://www.gtk.org/docs/getting-started/) as well.

## Compiling the code

Follow these steps to compile and build Pactus:

```text
git clone https://github.com/pactus-project/pactus.git
cd pactus
make build
```

This will be compile `pactus-daemon` and `pactus-wallet` on your machine.
Make sure Pactus is properly compiled and installed on your machine:

```
cd build
./pactus-daemon version
```

If you want to compile the GUI application, run this command in the root folder:

```text
make build_gui
```

To run the tests, use this command:

```text
make test
```

This may take several minutes to finish.

## What is pactus-daemon?

`pactus-daemon` is a full node implementation of Pactus blockchain.
You can use `pactus-daemon` to run a full node:

```text
./pactus-daemon init  -w=<working_dir>
./pactus-daemon start -w=<working_dir>
```

### Testnet

To join the TestNet, first you need to initialize your node
and then start the node:

```text
./pactus-daemon init  -w=<working_dir> --testnet
./pactus-daemon start -w=<working_dir>
```

### Local net

You can create a local node to set up a local network for development purposes on your machine:

 ```text
 ./pactus-daemon init  -w=<working_dir> --localnet
 ./pactus-daemon start -w=<working_dir>
 ```

## What is pactus-wallet?

Pactus wallet is a native wallet in the Pactus blockchain that lets users easily manage
their accounts on the Pactus blockchain.

### Getting started

To create a new wallet, run this command. The wallet will be encrypted by the
provided password.

```text
./pactus-wallet ~/pactus/wallets/wallet_1 create
```

You can create a new address like this:

```text
./pactus-wallet ~/pactus/wallets/wallet_1 address new
```

A list of addresses is available with this command:

```text
./pactus-wallet ~/pactus/wallets/wallet_1 address all
```

To obtain the public key of an address, run this command:

```text
./pactus-wallet ~/pactus/wallets/wallet_1 address pub <ADDRESS>
```

To publish a transaction, use the tx subcommand.
For example, to publish a bond transaction:

```text
./pactus-wallet ~/pactus/wallets/wallet_1 tx bond <FROM> <TO> <AMOUNT>
```

You can recover a wallet if you have the seed phrase.

```text
./pactus-wallet ~/pactus/wallets/wallet_2 recover
```


## Docker

You can run Pactus using a Docker file. Please make sure you have installed
[docker](https://docs.docker.com/engine/install/) on your machine.

Pull the Docker from Docker Hub:

```text
docker pull pactus/pactus
```

Let's create a working directory at `~/pactus/testnet` for the testnet:

```text
docker run -it --rm -v ~/pactus/testnet:/pactus pactus/pactus init -w /pactus --testnet
```

Now we can run Pactus and join the testnet:

```text
docker run -it -v ~/pactus/testnet:/pactus -p 8080:8080 -p 21777:21777 --name pactus-testnet pactus/pactus start -w /pactus
```

check "[http://localhost:8080](http://localhost:8080)" for the list of APIs.

Also you can stop/start docker:

```
docker stop pactus-testnet
docker start pactus-testnet
```

Or check the logs:

```
docker logs pactus-testnet --tail 1000 -f
```
