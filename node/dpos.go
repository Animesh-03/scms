package node

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"sort"

	"github.com/Animesh-03/scms/core"
	"github.com/Animesh-03/scms/logger"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

type DposClient struct {
	Stakes     map[string]uint `json:"stakes"`
	Votes      map[string]uint `json:"votes"`
	Verifiers  []string        `json:"verifiers"`
	BlockVotes map[string]uint `json:"blockvotes"`
}

func NewDposClient() DposClient {
	return DposClient{
		Stakes:     make(map[string]uint),
		Votes:      make(map[string]uint),
		BlockVotes: make(map[string]uint),
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

func BlockVerificationHandler(sub *pubsub.Subscription, self peer.ID, node *Node) {
	for {
		msg, err := sub.Next(context.Background())
		if err != nil {
			logger.LogError("Error reading from %s\n", sub.Topic())
			return
		}

		var block core.Block
		json.Unmarshal(msg.Data, &block)

		logger.LogInfo("Received block to verify: %+v\n", block.Stringify())

		if node.VerifyBlock(&block) {
			node.Network.Broadcast("block.verified", []byte(block.Stringify()))
		} else {
			logger.LogWarn("Received Invalid block to verify: %+v\n", block)
		}

	}
}

func BlockVerifiedHandler(sub *pubsub.Subscription, self peer.ID, node *Node) {
	for {
		msg, err := sub.Next(context.Background())
		if err != nil {
			logger.LogError("Error reading from %s\n", sub.Topic())
			return
		}

		var block core.Block
		json.Unmarshal(msg.Data, &block)

		logger.LogInfo("Block Verified with hash: %+v\n", block.Stringify())

		node.Dpos.BlockVotes[string(block.Hash)]++

		blockBytes, err := json.Marshal(block)
		if err != nil {
			logger.LogError("Error Marhsalling block: %+v\n", block.Stringify())
		}

		if node.Dpos.BlockVotes[string(block.Hash)] == uint(len(node.Dpos.Verifiers)) {
			node.Network.Broadcast("block.add", blockBytes)
		}
	}
}
