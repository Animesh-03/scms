package node

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"math/rand"

	"github.com/Animesh-03/scms/core"
	"github.com/Animesh-03/scms/logger"
)

type NodeType uint

const (
	Manufacturer NodeType = 1
	Distribtor   NodeType = 2
	Consumer     NodeType = 3
)

type Node struct {
	Type NodeType
	// Network *p2p.MDNSNetwork

	ID string

	Blockchain     []core.Block
	MemPool        *core.MemPool
	CurrentProduct string
	Stake          uint

	PrivKey *ecdsa.PrivateKey
	PubKey  *ecdsa.PublicKey
}

func NewNode(id string, t NodeType, pubKeyMap map[string]ecdsa.PublicKey, nodeMap map[string]*Node) *Node {
	node := Node{}

	node.ID = id
	node.Blockchain = make([]core.Block, 0)
	node.Blockchain = append(node.Blockchain, *core.CreateGenesisBlock())
	node.MemPool = core.NewMemPool()
	node.CurrentProduct = ""

	privKey, err := ecdsa.GenerateKey(elliptic.P256(), crand.Reader)
	if err != nil {
		logger.LogError("Error initializing node: %s\n", err.Error())
		return nil
	}
	node.PrivKey = privKey
	node.PubKey = &privKey.PublicKey

	pubKeyMap[id] = *node.PubKey
	nodeMap[id] = &node

	return &node
}

func (node *Node) RegisterStake(stakes map[string]uint, stake uint) {
	stakes[node.ID] = stake
	node.Stake = stake
}

func (node *Node) AddRandomVote(stakes map[string]uint, votes map[string]uint) {
	keys := make([]string, 0, len(stakes))
	for k := range stakes {
		keys = append(keys, k)
	}

	votes[keys[rand.Int()%len(keys)]] += node.Stake
}

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

func (node *Node) VerifyBlock(block *core.Block, pubKeyMap map[string]ecdsa.PublicKey) bool {
	return block.Verify(&node.Blockchain[len(node.Blockchain)-1], pubKeyMap)
}

func (node *Node) AddBlockToChain(block *core.Block) {
	node.Blockchain = append(node.Blockchain, *block)
	node.MemPool.RemoveAll(block.Transactions)
}
