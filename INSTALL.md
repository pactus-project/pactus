# Installing Pactus

## Requirements

You need to make sure you have installed [Git](https://git-scm.com/downloads)
and [Go 1.19 or higher](https://golang.org/) in your machine.
If you want to install GUI application, make sure you have install
[GTK+3](https://www.gtk.org/docs/getting-started/) as well.

## Compiling the code

Follow these steps to compile and build Pactus:

```bash
git clone https://github.com/pactus-project/pactus.git
cd pactus
make build
```

This will be compile `pactus-daemon` and `pactus-wallet` in your machine.
Make sure Pactus is properly compiled and installed in your machine:

```
cd build
./pactus-daemon version
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


## What is pactus-daemon

`pactus-daemon` is a full node implementation of Pactus blockchain.
You can use `pactus-daemon` to run a full node on Pactus blockchain.

```bash
./pactus-daemon init  -w=<working_dir>
./pactus-daemon start -w=<working_dir>
```

### Testnet

To join the TestNet, first you need to create a working directory
and then start the node:

```bash
./pactus-daemon init  -w=<working_dir> --testnet
./pactus-daemon start -w=<working_dir>
```

### Local net

You can create a local node with one validator to test Pactus in your machine:

 ```bash
 ./pactus-daemon init -w=<working_dir> --localnet
 ./pactus-daemon start -w=<working_dir>
 ```

## What is pactus-wallet

Pactus wallet is a native wallet in Pactus blockchain and let user to easily manage
their accounts on Pactus blockchain.

### Getting started

To create a new wallet run this command. The wallet will be encrypted by the
provided password.

```bash
./pactus-wallet -w=<wallet_path> create
```

You can create new address like this:

```bash
./pactus-wallet -w=<wallet_path> address new
```

List of addresses are available by this command:

```bash
./pactus-wallet -w=<wallet_path> address balance all
```

You can check the balance:

```bash
./pactus-wallet -w=<wallet_path> address balance <ADDRESS>
```

To publish a transactions use tx subcommand:

```bash
./pactus-wallet -w=<wallet_path> tx
```

For example, to publish a bond transaction:

```bash
./pactus-wallet -w=<wallet_path> tx bond <FROM> <TO> <STAKE>
```

You can recover a wallet if you have the seed phrase.

```bash
./pactus-wallet -w=<wallet_path> recover
```


## Usage of Docker

You can run the Pactus using docker file. Please make sure you have installed
[docker](https://docs.docker.com/engine/install/) in your machine.

Pull the docker from docker hub.

```bash
docker pull pactus/pactus
```

Let's create a working directory at `~/pactus/testnet` for the testnet:

```bash
docker run -it --rm -v ~/pactus/testnet:/pactus pactus/pactus init -w /pactus --testnet
```

Now we can run the pactus and join the testnet:

```bash
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
