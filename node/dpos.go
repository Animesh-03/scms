package node

import (
	"context"
	"encoding/json"

	"github.com/Animesh-03/scms/logger"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

type DposClient struct {
	Stakes map[string]uint `json:"stakes"`
}

type StakeData struct {
	PeerId string `json:"peerId"`
	Amount uint   `json:"amount"`
}

func (d *DposClient) RegisterStake(stake StakeData) {
	d.Stakes[stake.PeerId] = stake.Amount
}

func RegistrationHandler(sub *pubsub.Subscription, self peer.ID, node *Node) {
	for {
		msg, err := sub.Next(context.Background())
		if err != nil {
			logger.LogError("Error reading from %s\n", sub.Topic())
			return
		}

		var stake StakeData
		json.Unmarshal(msg.Data, &stake)

		node.Dpos.RegisterStake(stake)

		logger.LogInfo("Registered node %s with stake amount: %d\n", stake.PeerId, stake.Amount)
	}
}
