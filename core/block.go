package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"time"

	"github.com/Animesh-03/scms/logger"
)

type Block struct {
	Height            uint          `json:"height"`
	Hash              []byte        `json:"hash"`
	Timestamp         int64         `json:"timestamp"`
	MerkleRoot        []byte        `json:"merkleroot"`
	PreviousBlockHash []byte        `json:"previousblockhash"`
	Transactions      []Transaction `json:"transactions"`
}

// Creates a new block with given transactions and height
func NewBlock(txs []Transaction, previousBlockHash []byte, height uint) *Block {
	block := &Block{
		Height:            uint(height),
		Timestamp:         time.Now().UnixMilli(),
		PreviousBlockHash: previousBlockHash,
		Transactions:      txs,
	}

	var txHashes [][]byte
	for _, tx := range txs {
		txHashes = append(txHashes, tx.ID)
	}

	block.MerkleRoot = NewMerkleTree(txHashes).Root.Hash
	block.Hash = block.ComputeHash()

	return block
}

// Generates the genesis block with hardcoded hashes
func CreateGenesisBlock() *Block {
	block := &Block{
		Height:            1,
		Timestamp:         time.Now().UnixMilli(),
		MerkleRoot:        []byte("0"),
		PreviousBlockHash: []byte("0"),
		Hash:              []byte("0"),
		Transactions:      []Transaction{},
	}
	block.Hash = block.ComputeHash()

	return block
}

func (b *Block) ComputeHash() []byte {
	data := bytes.Join([][]byte{
		ToByte(int64(b.Height)),
		ToByte(b.Timestamp),
		b.PreviousBlockHash,
		b.MerkleRoot,
	}, []byte{})

	hash := sha256.Sum256(data)

	return hash[:]
}

func (b *Block) Verify(prevBlock *Block) bool {
	// Check if block hash or height are invalid
	if !bytes.Equal(b.PreviousBlockHash, prevBlock.Hash) || b.Height != prevBlock.Height+1 {
		return false
	}

	// Verify all the transactions in the block
	for _, tx := range b.Transactions {
		if !tx.Verify() {
			return false
		}
	}

	// Check the MerkleRoot
	var txHashes [][]byte
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	return bytes.Equal(b.MerkleRoot, NewMerkleTree(txHashes).Root.Hash)
}

func ToByte(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		logger.LogError("Error converting to Bytes: %s", err.Error())
		panic(err)
	}

	return buff.Bytes()
}
