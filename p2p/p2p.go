package p2p

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Animesh-03/scms/logger"
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
	flag.Parse()
	//
	cfg := NetworkConfig{
		ListenAddr:          *addr,
		ListenPort:          uint16(*port),
		DiscoveryServiceTag: *discoveryTag,
	}

	// Initialize the network
	net := MDNSNetwork{}
	net.Init(cfg)
	defer net.GetHost().Close()

	net.ListenBroadcast("test", testHandler)

	// Wait until terminated
	termCh := make(chan os.Signal, 1)
	signal.Notify(termCh, os.Interrupt, syscall.SIGTERM)
	<-termCh
	logger.LogInfo("Shutting Down Node...")
}
