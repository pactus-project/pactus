[![codecov](https://codecov.io/gh/zarbchain/zarb-go/branch/main/graph/badge.svg?token=8N6N60D5UI)](https://codecov.io/gh/zarbchain/zarb-go)
[![license](https://img.shields.io/github/license/zarbchain/zarb-go)](https://tldrlegal.com/license/mit-license)
[![CI](https://github.com/zarbchain/zarb-go/workflows/CI/badge.svg)](https://github.com/zarbchain/zarb-go/actions?query=workflow%3ACI+branch%3Amain+)
[![Go Report Card](https://goreportcard.com/badge/github.com/zarbchain/zarb-go)](https://goreportcard.com/report/github.com/zarbchain/zarb-go)
------
# Zarb

*Zarb blockchain*

## Compiling the code

You need to make sure you have install [Go](https://golang.org/).
Follow these steps to compile and build Zarb blockchain:

```bash
mkdir -p $GOPATH/src/github.com/zarbchain/zarb-go
cd $GOPATH/src/github.com/zarbchain/zarb-go
git clone https://github.com/zarbchain/zarb-go.git .
make
```

Run `zarb version` to make sure Zarb is properly compiled and installed in your machine.

## Running Zarb


### Testnet

To join test nest, run this command:

```bash
zarb start --wizard
```

And follow the instruction.

### Local net

Initialize the working directory by running:

 ```bash
 zarb init -w=<workspace_directory>
 zarb start -w=<workspace_directory>
 ```

 This command will create config.toml, genesis.json and private key for the validator.

## Usage of Docker

Install [Docker](https://www.docker.com/) and run the following commands to build the docker file:

```bash
make docker
```

Then you can execute the Zarb blockchain, using the docker:

```bash
docker pull zarb/zarb
docker run -it zarb/zarb start --wizard
```

## Contribution

 Any ideas are welcome. Feel free to submit any issues or pull requests.

## License

The Zarb blockchain is under MIT license.
