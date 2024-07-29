# NFT-Gated ETH-RPC Proxy

This is a gRPC service implemented in Golang with the help of Go-Ethereum.
The service serves as a proof-of-concept for a NFT-Gated ETH-RPC proxy.
It has a custom made in-memory caching system with the least recently used storing
strategy. Only holders of a specific NFT can submit transactions.

Before running the project, make sure you setup the environment by creating a
`.env` file. You can use the `.env.test` for reference.

### Runing the project

If you wish to run the project with Docker, you can do so by running: 
```bash
docker-compose up -d
```

You can also just run the database, and start the service locally.
If you don't have access to publicly available RPCs, you can run
your own fork of the chain with `ganache-cli`

Install it with `npm`:
```bash
npm install ganache -g
```
And start it with:
```bash
ganache
```

If you want to change the proto files and generate new code, you can do so with the 
prepared `Makefile` by running:
```bash
make generate
```

Just make sure you have `protoc` installed on your machine.

### Testing

For manual testing, please use the `./proto/ethproxy.proto` file for generating the request.
If you wish to run tests, you can do so with `goconvey`.

Runnig this command will start a local server on `http://127.0.0.1:8080/` where you can see
the test coverage.

```bash
goconvey -cover
```