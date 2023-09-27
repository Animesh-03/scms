package node

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"sort"

	"github.com/Animesh-03/scms/logger"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

type DposClient struct {
	Stakes    map[string]uint `json:"stakes"`
	Votes     map[string]uint `json:"votes"`
	Verifiers []string        `json:"verifiers"`
}

func NewDposClient() DposClient {
	return DposClient{
		Stakes: make(map[string]uint),
		Votes:  make(map[string]uint),
	}
}

type RegistrationData struct {
	PeerId    string          `json:"peerId"`
	Amount    uint            `json:"amount"`
	PublicKey ecdsa.PublicKey `json:"publickey"`
}

// Adds the stake to the respective node
func (d *DposClient) RegisterStake(stake RegistrationData) {
	d.Stakes[stake.PeerId] = stake.Amount
}

func (d *DposClient) ComputeVerfiers(n int) {
	verifiers := make([]string, 0, len(d.Votes))
	for k := range d.Votes {
		verifiers = append(verifiers, k)
	}
	sort.SliceStable(verifiers, func(i, j int) bool {
		return d.Votes[verifiers[i]] > d.Votes[verifiers[j]]
	})

	d.Verifiers = verifiers[:n]
}

func RegistrationHandler(sub *pubsub.Subscription, self peer.ID, node *Node) {
	for {
		msg, err := sub.Next(context.Background())
		if err != nil {
			logger.LogError("Error reading from %s\n", sub.Topic())
			return
		}

		var stake RegistrationData
		json.Unmarshal(msg.Data, &stake)

		node.PubKeyMap[stake.PeerId] = stake.PublicKey
		node.PeerMap[stake.PeerId] = msg.ReceivedFrom
		node.IDMap[msg.ReceivedFrom] = stake.PeerId
		node.Dpos.RegisterStake(stake)

		logger.LogInfo("Registered node %s with stake amount: %d\n", stake.PeerId, stake.Amount)
	}
}

func VotingHandler(sub *pubsub.Subscription, self peer.ID, node *Node) {
	for {
		msg, err := sub.Next(context.Background())
		if err != nil {
			logger.LogError("Error reading from %s\n", sub.Topic())
			return
		}

		var voteNode string
		json.Unmarshal(msg.Data, &voteNode)

		logger.LogInfo("Received vote from %s to %s\n", node.IDMap[msg.ReceivedFrom], voteNode)

		node.Dpos.Votes[node.IDMap[msg.ReceivedFrom]] += node.Dpos.Stakes[node.IDMap[msg.ReceivedFrom]]
	}
}
