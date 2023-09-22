package core

import "crypto/sha256"

type MerkleTree struct {
	Root *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Hash  []byte
}

// Creates a new Merkle Node
func NewMerkleNode(left *MerkleNode, right *MerkleNode, hash []byte) *MerkleNode {
	node := MerkleNode{}

	if left == nil && right == nil {
		// If leaf node then assign the hash directly
		node.Hash = hash
	} else {
		// Compute the hash by hashing the child nodes' hashes
		childHash := append(left.Hash, right.Hash...)
		currentHash := sha256.Sum256(childHash)
		node.Hash = currentHash[:]
	}

	node.Left = left
	node.Right = right
	return &node
}

// Construct the merkle tree given the array of hashes of transactions
func NewMerkleTree(txHashes [][]byte) *MerkleTree {
	nodes := make([]*MerkleNode, 0)

	// Create the leaf nodes of the tree
	for _, hash := range txHashes {
		node := NewMerkleNode(nil, nil, hash)
		nodes = append(nodes, node)
	}

	// Make the leaf nodes even by duplicating the last node
	if len(nodes)%2 == 0 {
		nodes = append(nodes, nodes[len(nodes)-1])
	}

	// Recursively form the merkle tree from the leaves
	for len(nodes) > 1 {
		var nextLevel []*MerkleNode
		for i := 0; i < len(nodes); i += 2 {
			var node *MerkleNode
			// If there are odd number of nodes in the current level
			// then duplicate the last node
			if i+1 == len(nodes) {
				node = NewMerkleNode(nodes[i], nodes[i], nil)
			} else {
				node = NewMerkleNode(nodes[i], nodes[i+1], nil)
			}
			nextLevel = append(nextLevel, node)
		}
		nodes = nextLevel
	}

	return &MerkleTree{nodes[0]}
}
