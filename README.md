# Zarb

*Zarb blockchain*

## Compiling the code

You need to make sure you have install [Go](https://golang.org/) and [rust](https://www.rust-lang.org). 
Follow these steps to compile and build Zarb blockchain:

```bash
mkdir -p $GOPATH/src/gitlab.com/zarb-chain/zarb-go
cd $GOPATH/src/gitlab.com/zarb-chain/zarb-go
git clone https://gitlab.com/zarb-chain/zarb-go.git .
make
```

Run `zarb version` to make sure Zarb is properly compiled and installed in your machine.

## Running Zarb

### Initialize

Initialize the working directory by running:

 ```bash
 zarb init -w=<workspace_directory>
 ```

 This command will create config.toml, genesis.json and private key for the validator.

### Run

For running a Zarb node, use:

```bash
zarb start -w=<workspace_directory>
```

The Zarb blockchain starts immediately.

## Usage of Docker

Install [Docker](https://www.docker.com/) and run the following commands to build the docker file:

```bash
make docker
```

Then you can execute the Zarb blockchain, using the docker:

```bash
# Initializing the working directory
docker run -it --rm -v "/zarb:/zarb" zarb init
# Starting the blockchain
docker run -it --rm -v "/zarb:/zarb" -p 1337:1337 -p 50051:50051 -p 46656:46656 zarb start
```

## Contribution

 Any ideas are welcome. Feel free to submit any issues or pull requests. 

## License

The Zarb blockchain is under MIT license.