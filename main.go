package main

import (
	"context"
	"encoding/json"
	"flag"

	"github.com/Animesh-03/scms/logger"
	"github.com/Animesh-03/scms/node"
	"github.com/Animesh-03/scms/p2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

func testHandler(sub *pubsub.Subscription, self peer.ID) {
	for {
		msg, err := sub.Next(context.Background())
		if err != nil {
			logger.LogError("Error reading from %s\n", sub.Topic())
			return
		}

		if msg.ReceivedFrom == self {
			continue
		}

		var msgString string
		json.Unmarshal(msg.Data, &msgString)
		logger.LogInfo("Message: %s\n", msgString)
	}
}

func main() {
	// Get command line args
	addr := flag.String("addr", "0.0.0.0", "Address of Network Inteface to be used")
	port := flag.Uint("p", 3000, "Port to be used to run the node")
	discoveryTag := flag.String("t", "mdns-discovery-tag", "Discovery tag")
	nodeType := flag.Uint("n", 3, "Enter the following: Manufacturer - 1, Distribtor - 2, Consumer - 3\n Default is Consumer")

	flag.Parse()

	//Create the config Object
	cfg := p2p.NetworkConfig{
		ListenAddr:          *addr,
		ListenPort:          uint16(*port),
		DiscoveryServiceTag: *discoveryTag,
	}

	node := &node.Node{
		Type: node.NodeType(*nodeType),
	}
	node.Start(&cfg)
}
