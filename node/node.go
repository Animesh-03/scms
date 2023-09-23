package node

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Animesh-03/scms/core"
	"github.com/Animesh-03/scms/logger"
	"github.com/Animesh-03/scms/p2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/gin-gonic/gin"
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
	node.MemPool = core.NewMemPool()
	node.Blockchain = make([]core.Block, 0)

	node.SetupListeners()
	node.SetupRPCs(uint(config.ListenPort + 1000))

	// Wait until terminated
	termCh := make(chan os.Signal, 1)
	signal.Notify(termCh, os.Interrupt, syscall.SIGTERM)
	<-termCh
	logger.LogInfo("Shutting Down Node...\n")
}

// Setup the listeners based on the type of node
func (node *Node) SetupListeners() {
	// Role Specific Listeners
	switch node.Type {
	case Manufacturer:

	case Distribtor:

	case Consumer:

	}
	// General Listeners
	node.Network.ListenBroadcast("transaction", func(sub *pubsub.Subscription, self peer.ID) { TransactionHandler(sub, self, node) })

}

func (node *Node) SetupRPCs(port uint) {
	router := gin.Default()

	router.POST("/transaction", func(ctx *gin.Context) { SendTransaction(ctx, node) })
	router.GET("/info", func(ctx *gin.Context) { GetNodeInfo(ctx, node) })

	router.Run(fmt.Sprintf("0.0.0.0:%d", port))
}
