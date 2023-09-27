package main

import (
	"crypto/ecdsa"
	"sort"

	"github.com/Animesh-03/scms/core"
	"github.com/Animesh-03/scms/logger"
	"github.com/Animesh-03/scms/node"
)

func main() {
	stakes := make(map[string]uint)
	votes := make(map[string]uint)
	pubKeyMap := make(map[string]ecdsa.PublicKey)
	nodeMap := make(map[string]*node.Node)

	// Initialize nodes
	manufacturer := node.NewNode("m1", 1, pubKeyMap, nodeMap)
	distributor1 := node.NewNode("d1", 2, pubKeyMap, nodeMap)
	distributor2 := node.NewNode("d2", 2, pubKeyMap, nodeMap)
	consumer1 := node.NewNode("c1", 3, pubKeyMap, nodeMap)
	consumer2 := node.NewNode("c2", 3, pubKeyMap, nodeMap)
	consumer3 := node.NewNode("c3", 3, pubKeyMap, nodeMap)
	consumer4 := node.NewNode("c4", 3, pubKeyMap, nodeMap)
	consumer5 := node.NewNode("c5", 3, pubKeyMap, nodeMap)
	consumer6 := node.NewNode("c6", 3, pubKeyMap, nodeMap)

	// Register Stakes
	manufacturer.RegisterStake(stakes, 20)
	distributor1.RegisterStake(stakes, 25)
	distributor2.RegisterStake(stakes, 25)
	consumer1.RegisterStake(stakes, 15)
	consumer2.RegisterStake(stakes, 15)
	consumer3.RegisterStake(stakes, 15)
	consumer4.RegisterStake(stakes, 15)
	consumer5.RegisterStake(stakes, 15)
	consumer6.RegisterStake(stakes, 15)

	// Start Voting for Verifier Selection
	// The voting is random in this implementation
	manufacturer.AddRandomVote(stakes, votes)
	distributor1.AddRandomVote(stakes, votes)
	distributor2.AddRandomVote(stakes, votes)
	consumer1.AddRandomVote(stakes, votes)
	consumer2.AddRandomVote(stakes, votes)
	consumer3.AddRandomVote(stakes, votes)
	consumer4.AddRandomVote(stakes, votes)
	consumer5.AddRandomVote(stakes, votes)
	consumer6.AddRandomVote(stakes, votes)

	logger.LogInfo("Votes are: %+v\n", votes)

	// Get the top elected nodes
	verifiers := make([]string, 0, len(votes))
	for k := range votes {
		verifiers = append(verifiers, k)
	}
	sort.SliceStable(verifiers, func(i, j int) bool {
		return votes[verifiers[i]] > votes[verifiers[j]]
	})

	verifiers = verifiers[:4]

	logger.LogInfo("Top Nodes are: %+v\n", verifiers)

	// Manufacturer Creates a product
	// Create a transaction to create a product
	createProductTx := core.NewTransaction(manufacturer.ID, "", "prod1", core.Manufactured)

	// Manufacturer Signs Tx
	manufacturer.SignTransaction(createProductTx)
	logger.LogInfo("Manufacturer Created Transaction: %+v\n", createProductTx.Stringify())

	// In a p2p network the transaction is broadcasted

	// Verify Transaction
	if !createProductTx.Verify(pubKeyMap[createProductTx.Sender]) {
		logger.LogWarn("Transaction Not Verified: %+v", createProductTx.Stringify())
		return
	}
	// Add transaction to mempool of all the nodes
	manufacturer.MemPool.AddToPool(createProductTx)
	distributor1.MemPool.AddToPool(createProductTx)
	distributor2.MemPool.AddToPool(createProductTx)
	consumer1.MemPool.AddToPool(createProductTx)
	consumer2.MemPool.AddToPool(createProductTx)
	consumer3.MemPool.AddToPool(createProductTx)
	consumer4.MemPool.AddToPool(createProductTx)
	consumer5.MemPool.AddToPool(createProductTx)
	consumer6.MemPool.AddToPool(createProductTx)

	// Random Node from the verifiers creates a block
	block := nodeMap[verifiers[0]].CreateBlock()
	logger.LogInfo("%s created block: %+v\n", verifiers[0], block.Stringify())

	// In a p2p network the created block is broadcasted

	// The block is verified by all the verifiers
	for _, v := range verifiers {
		if !nodeMap[v].VerifyBlock(block, pubKeyMap) {
			logger.LogWarn("Block not verified: %+v\n", block.Stringify())
			return
		}
	}

	// In p2p the block is broadcasted to the entire network after the verifiers verify the block

	// Block is added to blockchain by all the nodes
	manufacturer.AddBlockToChain(block)
	distributor1.AddBlockToChain(block)
	distributor2.AddBlockToChain(block)
	consumer1.AddBlockToChain(block)
	consumer2.AddBlockToChain(block)
	consumer3.AddBlockToChain(block)
	consumer4.AddBlockToChain(block)
	consumer5.AddBlockToChain(block)
	consumer6.AddBlockToChain(block)
}
