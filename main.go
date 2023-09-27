package main

import (
	"crypto/ecdsa"
	"math/rand"
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
	// distributor1 := node.NewNode("d1", 2, pubKeyMap, nodeMap)
	// distributor2 := node.NewNode("d2", 2, pubKeyMap, nodeMap)
	// consumer1 := node.NewNode("c1", 3, pubKeyMap, nodeMap)
	// consumer2 := node.NewNode("c2", 3, pubKeyMap, nodeMap)
	// consumer3 := node.NewNode("c3", 3, pubKeyMap, nodeMap)
	// consumer4 := node.NewNode("c4", 3, pubKeyMap, nodeMap)
	// consumer5 := node.NewNode("c5", 3, pubKeyMap, nodeMap)
	// consumer6 := node.NewNode("c6", 3, pubKeyMap, nodeMap)
	node.NewNode("d1", 2, pubKeyMap, nodeMap)
	node.NewNode("d2", 2, pubKeyMap, nodeMap)
	node.NewNode("c1", 3, pubKeyMap, nodeMap)
	node.NewNode("c2", 3, pubKeyMap, nodeMap)
	node.NewNode("c3", 3, pubKeyMap, nodeMap)
	node.NewNode("c4", 3, pubKeyMap, nodeMap)
	node.NewNode("c5", 3, pubKeyMap, nodeMap)
	node.NewNode("c6", 3, pubKeyMap, nodeMap)

	// Register Stakes
	for _, node := range nodeMap {
		node.RegisterStake(stakes, uint(rand.Int()%30))
	}

	// Start Voting for Verifier Selection
	// The voting is random in this implementation
	for _, node := range nodeMap {
		node.AddRandomVote(stakes, votes)
	}

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
	for _, node := range nodeMap {
		node.MemPool.AddToPool(createProductTx)
	}

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
	for _, node := range nodeMap {
		node.AddBlockToChain(block)
	}
}
