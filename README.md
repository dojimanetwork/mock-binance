[![pipeline status](https://gitlab.com/thorchain/bepswap/mock-binance/badges/master/pipeline.svg)](https://gitlab.com/thorchain/bepswap/mock-binance/commits/master)

Mock Binance
============

This project mocks some of the API endpoints of a [binance node
rpc](https://docs.binance.org/api-reference/node-rpc.html) for
testing purposes.
When doing a full integration test of thornode, it helps to have a fake
Binance full node that allows us to easily push fake transactions, blocks, and
account balances, without needing to use actual tokens (mainnet or testnet).
Also helps against API rate limiting problems when running full integration
tests in our pull requests.

### Install
```bash
make install
mock-binance # for a custom port, use `PORT` env var
```

### Seeding funds
There is a hard coded specific wallet address that has the ability to send
funds WITHOUT actually holding those funds. This address is
`tbnb1ht7v08hv2lhtmk8y7szl2hjexqryc3hcldlztl`. This is referred to as the
"seed" address. Use this address to send funds to other addresses.

### Sending transactions
When you want to send a transaction on mock-binance, you can use any of the
broadcast endpoints. No signature is required to be able to send from any
address.

We've added an additional broadcast endpoint to make it easy to send
transaction from the CLI. This endpoint is `/broadcast/easy`. An example of
using this endpoint is...

```bash
curl -s -X POST -d "{
  \"from\": \"tbnbXXXX\",
  \"to\": \"tbnbYYYY\",
  \"coins\":[
      {\"denom\": \"BNB\", \"amount\": 101010101}
  ],
  \"memo\": \"this is my memo\"
}" localhost:26660/broadcast/easy
```

Keep in mind amounts are `1e8` notation. So `1 BNB` is actually `100000000`.

### Tests
```bash
make test
```

or

```bash
make test-watch # hot reloading
```

