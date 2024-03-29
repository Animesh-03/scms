# Group Details

Shivansh Shrivastava - 2020A7PS2095H

Ishan Chhangani - 2020A7PS0230H

Kartikeya Dubey - 2020A7PS0031H

Aditha Venkata Animesh - 2020A7PS0193H

Sriram Balasubramanian - 2020A7PS0002H

# DPoS Implementation

## Phase 1 - Election of Group of Verifiers

1. Nodes initally register themselves in the network by broadcasting a stake amount. This amount is randomised for simulation purposes and is a value between 0 and 30. These stakes are stored by every node in the object `node.Dpos.Stakes`. The high level code for this is in the 
[node.go](node/node.go#L83) which sends the broadcast and [dpos.go](node/dpos.go#L53) which contains the handler function for the broadcast received.

2. Once all the nodes are registered to the network, the nodes start voting for the group of verifiers which in real world applications are decided on various factors like reputation but in this implementation the votes are randomised. The code for this can be found in [node.go](node/node.go#L90) which broadcasts the vote to the other nodes and is handled in [dpos.go](node/dpos.go#L73) which stores the votes.

3. After the votes from all the nodes are received then the nodes compute the top nodes by summing up all the nodes' votes. The top `n` nodes are selected to form the group of verifiers. The value of `n` in this implementation is 2. The code for this can be found in [node.go](node/node.go#L96)

## Phase 2 - Consensus on Blocks Generated

1. In real world applications the group of verifiers take turns creating the blocks and verifying the blocks but in this implmentation the top elected node creates the blocks and the rest verify the blocks. After all the verifiers are done verifying the block, the block is broadcast to all the other nodes who then add it to their copy of the blockchain. The code for this can be found [node.go](node/node.go#L112) and [dpos.go](node/dpos.go#L90).

2. The block is broadcasted to all the other nodes only when all the verifiers verify the block. The code for this can be found in [dpos.go](node/dpos.go#L112)

## Implementation with no P2P

There is also an implementation of DPoS with no P2P in [main.go](main.go) but is in the git branch `nop2p`.

To switch to the branch:
```bash
git checkout nop2p
```
# RPCs

The nodes also have a set of RPCs included with them which can be used to interact with them

## GET /info

This returns the current state of the node including DPoS details like the stake and votes

The code for the RPC is located in [rpc.go](node/rpc.go#L31)

Sample Response:
```json
{
    "ID": "3000",
    "Type": 3,
    "Network": {},
    "Blockchain": [
        {
            "height": 1,
            "hash": "OHzp6Xt9nQC3iWibqrJ2ZYhmmR8VBNsDQvuhJGXUqBM=",
            "timestamp": 0,
            "merkleroot": "MA==",
            "previousblockhash": "MA==",
            "transactions": []
        },
        {
            "height": 2,
            "hash": "sQeoAynviqd5TZ6bnjRRVRObL+5kFkFayjfev2Cf9aE=",
            "timestamp": 1696083580346,
            "merkleroot": "3A2J+RNoYGkGwHGhz8vdAtWg8opjYW9eX0VHqdYQLhg=",
            "previousblockhash": "OHzp6Xt9nQC3iWibqrJ2ZYhmmR8VBNsDQvuhJGXUqBM=",
            "transactions": [
                {
                    "id": "3A2J+RNoYGkGwHGhz8vdAtWg8opjYW9eX0VHqdYQLhg=",
                    "sender": "3002",
                    "receiver": "abc",
                    "productid": "123",
                    "status": 1,
                    "signature": "MEUCIHNb2cKZdXoIOh9Hzu81xZ+bafYp2yacYl1jnZkcEMfOAiEA/MbWYbvyLIpSpx7okoXOBrdeSJHVbePzZX5q/KJ7Flg="
                }
            ]
        },
        {
            "height": 3,
            "hash": "KryxO2CVZ02UsQKCdzWXxqeJsU/p1a23u1ftntCBQKg=",
            "timestamp": 1696083590347,
            "merkleroot": "j7Fk2Phdp8nH4CDvknbPxaEBpbfrzNRADvFwtT79BZc=",
            "previousblockhash": "sQeoAynviqd5TZ6bnjRRVRObL+5kFkFayjfev2Cf9aE=",
            "transactions": [
                {
                    "id": "j7Fk2Phdp8nH4CDvknbPxaEBpbfrzNRADvFwtT79BZc=",
                    "sender": "3001",
                    "receiver": "abc",
                    "productid": "123",
                    "status": 2,
                    "signature": "MEYCIQDex4lsSSS2iMNT9MrfML52uHvEfrwHr53nw3FWQHUahQIhALSmPucmfSgZWr6ncqqSLqaUcPGd2mwrXRIaoK/948cZ"
                }
            ]
        },
        {
            "height": 4,
            "hash": "ohxhPe9NqteXTP1lg6d1eujqOKxpX852zLOdNIv/9+8=",
            "timestamp": 1696083600348,
            "merkleroot": "MA==",
            "previousblockhash": "KryxO2CVZ02UsQKCdzWXxqeJsU/p1a23u1ftntCBQKg=",
            "transactions": []
        }
    ],
    "MemPool": {
        "pool": {}
    },
    "CurrentProduct": "",
    "PubKeyMap": {
        "3000": {
            "Curve": null,
            "X": 30246794350822519054778452508393483556851050147495679395194914196969419200222,
            "Y": 22806389722793802419181077681160860018805169143432979971790518932972790042990
        },
        "3001": {
            "Curve": null,
            "X": 111699090184077551379246621370075291095205079480239697690512798089746456225519,
            "Y": 23745705857280546853113833094637348877448314857391126289603691139401930801198
        },
        "3002": {
            "Curve": null,
            "X": 102267057708492298109409727589472747383437496181344795389041675352874742788062,
            "Y": 34500421244774085528174091607038487313492071036499885878458852534495372449814
        }
    },
    "PeerMap": {
        "3000": "12D3KooWRGbei5QaeFq3hDYbR73zHtahQPEFvX3Q8LhiFWBd3638",
        "3001": "12D3KooWP9GCyXCHkfatEC5hEHo64x5zv8NTNWx9Eyc8Nf22sRmn",
        "3002": "12D3KooWJpV4r4AT6ggD6EQfPnPr4CcYsjbusCwrrrNLCbFViR1V"
    },
    "IDMap": {
        "id1": "3002",
        "id2": "3001",
        "id3": "3000"
    }
    "PrivKey": {
        "Curve": {},
        "X": 30246794350822519054778452508393483556851050147495679395194914196969419200222,
        "Y": 22806389722793802419181077681160860018805169143432979971790518932972790042990,
        "D": 21811584513617742247549216491386892127787873856301732411581725258711867224663
    },
    "PubKey": {
        "Curve": {},
        "X": 30246794350822519054778452508393483556851050147495679395194914196969419200222,
        "Y": 22806389722793802419181077681160860018805169143432979971790518932972790042990
    },
    "Dpos": {
        "stakes": {
            "3000": 6,
            "3001": 12,
            "3002": 22
        },
        "votes": {
            "3000": 16,
            "3001": 12,
            "3002": 22
        },
        "verifiers": [
            "3002",
            "3000"
        ],
        "blockvotes": {}
    }
}
```

## POST /transaction

Depending on the type of node i.e, Manufacturer, Distributor and Consumer the behavior of this RPC changes.

For a manufacturer node, it takes in the `productid` which is the product that it creates and the `receiver` which is the id of the distributor node that it wants to send the product to.

For a distributor node, it takes in the `productid` which is the product that it wants to dispatch to the `receiver` which is the ID of a consumer node

For a consumer node, it takes in the `productid` which is the product that is received by the consumer node.

All the above behaviours result in the generation of a transaction that is broadcast over the network.

The code for the RPC is located in [rpc.go](node/rpc.go#L17)

Sample Request:
```json
{
    "receiver": "abc",
    "productid": "123"
}
```

Sample Response:
```json
{
    "id": "piS0YGuKlQ8LTuB9O/l1AUN/C9sJOka+tQWtlQruPZI=",
    "sender": "3000",
    "receiver": "abc",
    "productid": "123",
    "status": 3,
    "signature": "MEUCIQDTP9NUhHOCE9Hnx6G63K/9ARgDL0bGFq1V3Wv/ywfNaQIgWOIo4+pzovlv/MMkLMwgsvNzmaXZQXeqJpxqT7hKiKE="
}
```

## POST /product_status

This returns a QR code that contains the status of the `productid` that is passed in the request body.

The code for the RPC is located in [rpc.go](node/rpc.go#L39)

Sample Request:
```json
{
    "productid": "123"
}
```

## POST /dispute

This RPC can be called when a consumer node wants to raise a dispute on the delivery of a `productid` that is passed in the request body.

If the dispatcher truly dispatched then the consumer is making a false claim and his stake in the network is penalised.

If the distributor node did not dispatch the product but claims to have dispatched it then the distributor node's stake in the network is penalised.

The code for the RPC is located in [rpc.go](node/rpc.go#L69)

Sample Request:
```json
{
    "productid": "123"
}
```

Sample Response:
```json
{
    "error": "Consumer is wrong, stake is being deducted"
}
```

# Requirements

golang >= go1.20.0

Installation instructions for golang can be found [here](https://go.dev/dl/)

# How to run

After cloning for the first time run the following in the terminal at the root directory
```bash
go get .
```

This will install the necessary libraries to run the node

To run a consumer node
```bash
./run-consumer.sh
```

To run a distributor node
```bash
./run-distributor.sh
```

To run a manufacturer node
```bash
./run-manufacturer.sh
```

To run more nodes use the help text that can be obtained by using
```bash
./bin/scms -h
```

The nodes can discover other nodes on the network automatically and connect to them.
