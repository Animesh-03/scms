package node

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	ID      string
	Type    NodeType
	Network *p2p.MDNSNetwork

	Blockchain     []core.Block
	MemPool        *core.MemPool
	CurrentProduct string
	PubKeyMap      map[string]ecdsa.PublicKey
	PeerMap        map[string]peer.ID
	IDMap          map[peer.ID]string

	PrivKey *ecdsa.PrivateKey
	PubKey  *ecdsa.PublicKey

	Dpos DposClient
}

// Initialize the node by joining the network
func (node *Node) Start(config *p2p.NetworkConfig) {
	// Initialize Node
	node.CurrentProduct = ""

	// Initialize the network
	net := p2p.MDNSNetwork{}
	net.Init(*config)
	defer net.GetHost().Close()

	node.ID = fmt.Sprintf("%d", config.ListenPort)
	node.Network = &net
	node.MemPool = core.NewMemPool()
	node.Blockchain = make([]core.Block, 0)
	node.Dpos = NewDposClient()

	node.PubKeyMap = make(map[string]ecdsa.PublicKey)
	node.PeerMap = make(map[string]peer.ID)
	node.IDMap = make(map[peer.ID]string)

	privKey, err := ecdsa.GenerateKey(elliptic.P256(), crand.Reader)
	if err != nil {
		logger.LogError("Error initializing node: %s\n", err.Error())
		return
	}
	node.PrivKey = privKey
	node.PubKey = &privKey.PublicKey

	node.SetupListeners()
	go node.SetupRPCs(uint(config.ListenPort + 1000))

	// Register the node after a delay
	// Stake a random amount of tokens < 30
	go func() {
		time.Sleep(10 * time.Second)
		node.Register(uint(rand.Int()) % 30)
	}()

	// Vote for a random node after a delay
	go func() {
		time.Sleep(15 * time.Second)
		node.VoteRandomNode()
	}()

	// Compute the final
	go func() {
		time.Sleep(20 * time.Second)
		node.Dpos.ComputeVerfiers(2)
		logger.LogInfo("Final Votes are: %+v\n", node.Dpos.Votes)
		logger.LogInfo("Verifiers are: %+v\n", node.Dpos.Verifiers)
	}()

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
	node.Network.ListenBroadcast("register", func(sub *pubsub.Subscription, self peer.ID) { RegistrationHandler(sub, self, node) })
	node.Network.ListenBroadcast("vote", func(sub *pubsub.Subscription, self peer.ID) { VotingHandler(sub, self, node) })

	logger.LogInfo("Listeners Setup Successfully\n")
}

func (node *Node) SetupRPCs(port uint) {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.POST("/transaction", func(ctx *gin.Context) { SendTransaction(ctx, node) })
	router.GET("/info", func(ctx *gin.Context) { GetNodeInfo(ctx, node) })
	router.GET("/product_status", func(ctx *gin.Context) { GetProductStatus(ctx, node) })

	router.Run(fmt.Sprintf("0.0.0.0:%d", port))
}

// Sign A Transaction
func (node *Node) SignTransaction(tx *core.Transaction) {
	signature, err := ecdsa.SignASN1(crand.Reader, node.PrivKey, tx.Bytes())
	if err != nil {
		logger.LogError("Error signing tx: %s\n", err.Error())
		return
	}

	tx.Signature = signature
}

func (node *Node) CreateBlock() *core.Block {
	block := core.NewBlock(node.MemPool.GetTransactions(5), node.Blockchain[len(node.Blockchain)-1].Hash, node.Blockchain[len(node.Blockchain)-1].Height+1)
	return block
}

func (node *Node) VerifyBlock(block *core.Block) bool {
	return block.Verify(&node.Blockchain[len(node.Blockchain)-1], node.PubKeyMap)
}

func (node *Node) AddBlockToBlockChain(block *core.Block) {
	node.Blockchain = append(node.Blockchain, *block)
	node.MemPool.RemoveAll(block.Transactions)
}

// Broadcast the stake to register
func (node *Node) Register(stakeAmount uint) {
	logger.LogInfo("Registering self with amount: %d\n", stakeAmount)
	stake := RegistrationData{
		PeerId: node.ID,
		Amount: stakeAmount,
	}
	stakeBytes, err := json.Marshal(stake)
	if err != nil {
		logger.LogError("error marshalling stake\n")
		return
	}

	node.Network.Broadcast("register", stakeBytes)
}

func (node *Node) VoteRandomNode() {
	keys := make([]string, 0, len(node.PubKeyMap))
	for k := range node.PubKeyMap {
		keys = append(keys, k)
	}

	voteNode := keys[rand.Int()%len(node.PeerMap)]
	node.Network.Broadcast("vote", []byte(voteNode))
}
