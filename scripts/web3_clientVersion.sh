#!/bin/bash
curl -X POST \
    -H "Content-Type: application/json" \
    --data '{"jsonrpc":"2.0","method":"web3_clientVersion","params":[],"id":67}' \
    https://ethereum-mainnet.g.allthatnode.com/full/evm/77fc5541d71a47a4b2b3a31fb2bb7abc

