# Simple Ethereum Parse
## How to Run
### Build
#### Go Version
```shell
go version
```
```text
go version go1.22.1 darwin/arm64
```
```shell
go build -o build/server
```

### Run
```shell
./build/server
```

### Run and Build
```shell
go build -o build/server && build/server
```

## Concepts
I used channel in few places to also illustrate the possibility of building this into microservice and event-driven architecture.

### Parser
The parser is just an interface to the storage since a lot of heavy lifting (JSONRPC calls) happens in observers.

### Storage
Simple in memory storage with mutex

### HTTP Router
Simple HTTP Router setup to interact with Parser:
- `GET /block`: get current block number from storage
- `POST /subscribe-address`: add address to the storage, this address will be observed for transactions
- `GET /transactions?address=XXX`: get transactions history for an address since the address added to observed list

### Block Observer
- Every 5 seconds will call jsonrpc (`method=eth_blockNumber`) to get latest block number
- If the block is larger than we have in store, update it and send the block number to channel for transaction processing

### New Block Observer
- Observer will listen to channel for block number
- On every message on the channel, will call jsonrpc (`method=eth_getBlockByNumber`) and look for address in observed addresses list
- Send this to transaction channel to fetch transaction detail

### Transaction Observer
- Observer will listen to channel for transaction hash
- On every message on the channel, will call jsonrpc (`method=eth_getTransactionReceipt`)
- Add to transaction list for address
- Print receipt on console

## Examples
For the sake of simplicity, I added address `0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD` as observed address (Uniswap Router)

### Get Transactions
```shell
curl --location 'localhost:8090/transactions?address=0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD'
```

### Get Latest Block Number
```shell
curl --location 'localhost:8090/block'
```

### Subscribe Address
```shell
curl --location 'localhost:8090/subscribe-address' \
--header 'Content-Type: application/json' \
--data '{
    "address": "0x974CaA59e49682CdA0AD2bbe82983419A2ECC400"
}'
```

### Example Run Log
```shell
2024/05/02 11:31:04 fetching txs on new block: 0x12dd1e2
2024/05/02 11:31:04 saving last block number to mem store: 19780066
2024/05/02 11:31:04 fetching latest block number
2024/05/02 11:31:04 getting last block number from mem store: 19780066
2024/05/02 11:31:05 added new transaction 0xb5e41d0d991f3196e3634cefd5ada86d15bd0a80d4084c35300208e982820e0f to address [to:0x3fc91a3afd70395cd496c647d5a6cc9d4b2b7fad, from:0x62fd715041c1810cb9f93415bf76295a89c7fd91]
2024/05/02 11:31:05 observed address, adding txs
2024/05/02 11:31:05 observed address, adding txs
2024/05/02 11:31:05 fetching txs receipt: 0xc9f5ecdcf2da073c4660253132c1dc5b8601989732ec337c3a8d864ba48012ce
2024/05/02 11:31:05 added new transaction 0xc9f5ecdcf2da073c4660253132c1dc5b8601989732ec337c3a8d864ba48012ce to address [to:0x3fc91a3afd70395cd496c647d5a6cc9d4b2b7fad, from:0x8b7dc7ca12a40aeb6f0575af1d040ab7321d61c9]
2024/05/02 11:31:05 fetching txs receipt: 0x85e8da7c993fc3d5c02de40231d750ed84543b5e209a5e038eb197e9692fc895
2024/05/02 11:31:05 observed address, adding txs
2024/05/02 11:31:06 added new transaction 0x85e8da7c993fc3d5c02de40231d750ed84543b5e209a5e038eb197e9692fc895 to address [to:0x3fc91a3afd70395cd496c647d5a6cc9d4b2b7fad, from:0x9d15716d1de9b28ef1462163257891b37a74f091]
2024/05/02 11:31:06 fetching txs receipt: 0x23ae7c9397b02ee672279c5ee1936102da9cada3bc52ad3ac422fd0ead81064e
2024/05/02 11:31:06 observed address, adding txs
2024/05/02 11:31:07 added new transaction 0x23ae7c9397b02ee672279c5ee1936102da9cada3bc52ad3ac422fd0ead81064e to address [to:0x3fc91a3afd70395cd496c647d5a6cc9d4b2b7fad, from:0x0280650a232c48f86adb1cc9d6ab70f50363c45e]
2024/05/02 11:31:07 fetching txs receipt: 0x7bdb0f1d99e8e8c7c09c23ad6eaea420d451bae57e8e573462bb7749858cb670
2024/05/02 11:31:07 observed address, adding txs
2024/05/02 11:31:07 added new transaction 0x7bdb0f1d99e8e8c7c09c23ad6eaea420d451bae57e8e573462bb7749858cb670 to address [to:0x3fc91a3afd70395cd496c647d5a6cc9d4b2b7fad, from:0x6e3b2b6c08080f1aba2515fefcdd68817865072b]
2024/05/02 11:31:07 fetching txs receipt: 0x640475a56d32e49246d02ef7b778898162258ce6403c7bbf930b6758ab96139e
2024/05/02 11:31:07 observed address, adding txs
2024/05/02 11:31:08 added new transaction 0x640475a56d32e49246d02ef7b778898162258ce6403c7bbf930b6758ab96139e to address [to:0x3fc91a3afd70395cd496c647d5a6cc9d4b2b7fad, from:0xf2e7d327b10b0ab060557840e7d84a7b3564dffb]
2024/05/02 11:31:08 fetching txs receipt: 0x1aad4d9f7f83e5d92eabdad10e143d7026f0623703124f783449770658237a06
2024/05/02 11:31:08 observed address, adding txs
2024/05/02 11:31:08 added new transaction 0x1aad4d9f7f83e5d92eabdad10e143d7026f0623703124f783449770658237a06 to address [to:0x3fc91a3afd70395cd496c647d5a6cc9d4b2b7fad, from:0xfcbbba89859202e61554587a75a0ea97fb364dcf]
2024/05/02 11:31:08 fetching txs receipt: 0x1080d9ec7deed2b9ee53fc5fcaca203aaec155dfe83f55ce9f7327e0af3a8bb5
2024/05/02 11:31:08 observed address, adding txs
2024/05/02 11:31:08 added new transaction 0x1080d9ec7deed2b9ee53fc5fcaca203aaec155dfe83f55ce9f7327e0af3a8bb5 to address [to:0x3fc91a3afd70395cd496c647d5a6cc9d4b2b7fad, from:0xf3f3485f85951030865a8f80b597cf5ec24f5a9b]
2024/05/02 11:31:08 fetching txs receipt: 0x45aefd03cf6f4e8dbf50999ea45a20b93f935f71e213e5536380c8d3d705b24d
2024/05/02 11:31:08 observed address, adding txs
```