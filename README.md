[![codecov](https://codecov.io/gh/zarbchain/zarb-go/branch/main/graph/badge.svg?token=8N6N60D5UI)](https://codecov.io/gh/zarbchain/zarb-go)
[![CI](https://github.com/zarbchain/zarb-go/workflows/CI/badge.svg)](https://github.com/zarbchain/zarb-go/actions?query=workflow%3ACI+branch%3Amain+)
[![Go Report Card](https://goreportcard.com/badge/github.com/zarbchain/zarb-go)](https://goreportcard.com/report/github.com/zarbchain/zarb-go)
[![Discord](https://badgen.net/badge/icon/discord?icon=discord&label)](https://discord.gg/zPqWqV85ch)
------
# Zarb

*Zarb blockchain* (https://zarb.network)

## Compiling the code

You need to make sure you have installed [(Go1.15 or higher)](https://golang.org/).
Follow these steps to compile and build Zarb blockchain:

```bash
git clone https://github.com/zarbchain/zarb-go.git
cd zarb-go
make
```

Run `zarb version` to make sure Zarb is properly compiled and installed in your machine.

## Running Zarb


### Testnet

To join the TestNet, first you need to create a working directory
and then start the node:

```bash
zarb init  -w=<working_dir> --testnet
zarb start -w=<working_dir>
```

### Local net

You can create a local node with one validator to test Zerb in your machine:

 ```bash
 zarb init -w=<working_dir>
 zarb start -w=<working_dir>
 ```

## Usage of Docker

You can run the Zarb using docker file.
Please make sure you have installed [docker](https://docs.docker.com/engine/install/) in your machine.

Pull the docker from docker hub.

```bash
docker pull zarb/zarb
```

Let's create a worksapce at `~/zarb/testnet` for the testnet:

```bash
docker run -it --rm -v ~/zarb/testnet:/zarb zarb/zarb init -w /zarb --testnet
```

Now we can run the zarb and join the testnet:

```bash
docker run -it -v ~/zarb/testnet:/zarb -p 8080:8080 --name zarb-testnet zarb/zarb start -w /zarb
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

## Contribution

 Any ideas are welcome. Feel free to submit any issues or pull requests.

## License

The Zarb blockchain is under MIT license.
