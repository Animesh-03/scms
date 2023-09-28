package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/Animesh-03/scms/core"
	"github.com/Animesh-03/scms/logger"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

// Check the status of the given product by the product ID
func (n *Node) GetStatusOfProduct(productId string) (core.TransactionStatus, error) {
	if productId == "" {
		return 0, errors.New("product not found")
	}

	status := uint16(0)
	for _, block := range n.Blockchain {
		for _, tx := range block.Transactions {
			if tx.ProductID == productId {
				status = uint16(math.Max(float64(status), float64(tx.Status)))
			}
		}
	}

	if status == 0 {
		return 0, errors.New("product not found")
	}
	return core.TransactionStatus(status), nil
}

// Broadcast a transaction with status based on the type of node
func (n *Node) MakeTransaction(receiver, productId string) (*core.Transaction, error) {
	var transaction *core.Transaction

	switch n.Type {
	case Manufacturer:
		status, _ := n.GetStatusOfProduct(productId)
		if status == 0 {
			transaction = core.NewTransaction(n.ID, receiver, productId, core.Manufactured)
		} else {
			return nil, fmt.Errorf("product %s already exists", productId)
		}

	case Distribtor:
		// Check if a previous product is still progress
		status, _ := n.GetStatusOfProduct(n.CurrentProduct)
		status1, _ := n.GetStatusOfProduct(productId)
		if status1 == core.Manufactured && (status == core.Received || n.CurrentProduct == "") {
			transaction = core.NewTransaction(n.ID, receiver, productId, core.Dispatched)
			n.CurrentProduct = productId
		} else if status1 != core.Manufactured {
			return nil, fmt.Errorf("product %s not yet manufactured", productId)
		} else {
			return nil, fmt.Errorf("previous product %s still in progress", n.CurrentProduct)
		}

	case Consumer:
		status, _ := n.GetStatusOfProduct(productId)

		if status == 0 {
			return nil, fmt.Errorf("product %s not found", productId)
		}

		if status != core.Dispatched {
			return nil, fmt.Errorf("product %s not dispatched", productId)
		}

		transaction = core.NewTransaction(n.ID, receiver, productId, core.Received)
	}

	n.SignTransaction(transaction)

	transactionBytes, err := json.Marshal(transaction)
	if err != nil {
		logger.LogError("error marshalling transaction: %s", err.Error())
		return nil, err
	}

	n.Network.Broadcast("transaction", transactionBytes)

	return transaction, nil
}

func TransactionHandler(sub *pubsub.Subscription, self peer.ID, node *Node) {
	for {
		msg, err := sub.Next(context.Background())
		if err != nil {
			logger.LogError("Error reading from %s\n", sub.Topic())
			return
		}
		logger.LogInfo("Received Transaction from %s:\n%s\n", msg.ReceivedFrom.String(), msg.GetData())

		var transaction core.Transaction
		json.Unmarshal(msg.Data, &transaction)

		if transaction.Verify(node.PubKeyMap[transaction.Sender]) {
			node.MemPool.AddToPool(&transaction)
		} else {
			logger.LogWarn("Transaction Invalid: %s", transaction.Stringify())
		}
	}
}
