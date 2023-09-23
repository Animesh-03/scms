package node

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Animesh-03/scms/core"
	"github.com/Animesh-03/scms/logger"
	"github.com/Animesh-03/scms/p2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

type NodeType uint

const (
	Manufacturer NodeType = 1
	Distribtor   NodeType = 2
	Consumer     NodeType = 3
)

type Node struct {
	Type    NodeType
	Network *p2p.MDNSNetwork

	Blockchain     []core.Block
	MemPool        *core.MemPool
	CurrentProduct string
}

// Initialize the node by joining the network
func (node *Node) Start(config *p2p.NetworkConfig) {
	// Initialize Node
	node.CurrentProduct = ""

	// Initialize the network
	net := p2p.MDNSNetwork{}
	net.Init(*config)
	defer net.GetHost().Close()

	node.Network = &net

	node.SetupListeners()

	// Wait until terminated
	termCh := make(chan os.Signal, 1)
	signal.Notify(termCh, os.Interrupt, syscall.SIGTERM)
	<-termCh
	logger.LogInfo("Shutting Down Node...\n")
}

// Setup the listeners based on the type of node
func (node *Node) SetupListeners() {
	switch node.Type {
	case Manufacturer:

	case Distribtor:

	case Consumer:

	}
	node.Network.ListenBroadcast("test", func(sub *pubsub.Subscription, self peer.ID) {})
	node.Network.ListenBroadcast("transaction", TransactionHandler)

}
