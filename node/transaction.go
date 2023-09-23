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

func (n *Node) MakeTransaction(receiver, productId string) error {
	var transaction *core.Transaction
	switch n.Type {
	case Manufacturer:
		transaction = core.NewTransaction(n.Network.GetHost().ID().String(), receiver, productId, core.Manufactured)

	case Distribtor:
		// Check if a previous product is still progress
		status, _ := n.GetStatusOfProduct(n.CurrentProduct)
		if n.CurrentProduct == "" || status == core.Received {
			transaction = core.NewTransaction(n.Network.GetHost().ID().String(), receiver, productId, core.Dispatched)
			n.CurrentProduct = productId
		} else {
			return fmt.Errorf("previous product %s still in progress", n.CurrentProduct)
		}

	case Consumer:
		transaction = core.NewTransaction(n.Network.GetHost().ID().String(), receiver, productId, core.Received)
	}

	transactionBytes, err := json.Marshal(*transaction)
	if err != nil {
		logger.LogError("error marshalling transaction: %s", err.Error())
		return err
	}

	n.Network.Broadcast("transaction", transactionBytes)

	return nil
}

func TransactionHandler(sub *pubsub.Subscription, self peer.ID) {
	for {
		msg, err := sub.Next(context.Background())
		if err != nil {
			logger.LogError("Error reading from %s\n", sub.Topic())
			return
		}

		var transaction core.Transaction
		json.Unmarshal(msg.Data, &transaction)
		logger.LogInfo("Received Transaction from %s:\n%+v\n", msg.ReceivedFrom.String(), transaction)
	}
}