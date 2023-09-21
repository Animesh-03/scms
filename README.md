# Requirements

golang >= go1.20.0

Installation instructions for golang can be found [here](https://go.dev/dl/)

# How to run

After cloning for the first time run the following in the terminal at the root directory
```bash
go get .
```

This will install the necessary libraries to run the node

To run a node
```bash
go run . -p <port> # Default port is 3000
```

To run multiple nodes, run the above command in different terminal and assign a different port

The nodes can discover other nodes on the network automatically and connect to them.
