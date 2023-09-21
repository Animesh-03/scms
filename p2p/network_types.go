package p2p

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

type Network interface {
	AddStream(p string, handler func(stream network.Stream))
	Broadcast(topic string, msg string)
	GetHost() host.Host
	GetNumberOfPeers() int
	GetPeers() map[peer.ID]*peer.AddrInfo
	Init(config NetworkConfig)
	ListenBroadcast(topic string, handler func(sub *pubsub.Subscription, self peer.ID))
	SendTo(proto string, p peer.ID, msg string)
}

type Discoverer interface {
	StartDiscovery()
}
