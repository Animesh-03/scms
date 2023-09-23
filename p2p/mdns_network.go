package p2p

import (
	"bufio"
	"context"
	"fmt"

	"github.com/Animesh-03/scms/logger"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type discoveryNotifee struct {
	h host.Host
	n *MDNSNetwork
}

func (d *discoveryNotifee) HandlePeerFound(p peer.AddrInfo) {
	// Connect to the peer found
	err := d.h.Connect(context.Background(), p)

	if err != nil {
		logger.LogError("Error Connecting to peer %s: %s\n", p.ID.Pretty(), err)
		return
	}
	logger.LogConnectEvent("Connected Successfully To %s\n", p.ID.Pretty())
	d.n.peers[p.ID] = &p
}

type NetworkConfig struct {
	ListenAddr          string
	ListenPort          uint16
	DiscoveryServiceTag string
}

type MDNSNetwork struct {
	h      host.Host
	config NetworkConfig
	ps     *pubsub.PubSub
	subs   map[string]*pubsub.Subscription
	topics map[string]*pubsub.Topic
	peers  map[peer.ID]*peer.AddrInfo
}

// Start the MDNS discovery service
func (n *MDNSNetwork) StartDiscovery() {
	d := mdns.NewMdnsService(n.h, n.config.DiscoveryServiceTag, &discoveryNotifee{h: n.h, n: n})
	d.Start()
}

// Initialize the network based on the config and start the discovery service
func (n *MDNSNetwork) Init(config NetworkConfig) {
	n.config = config
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%d", n.config.ListenAddr, n.config.ListenPort)),
	)
	if err != nil {
		logger.LogError("Error occured while initializing network: %s", err)
		return
	}
	n.h = h

	n.subs = make(map[string]*pubsub.Subscription)
	n.topics = make(map[string]*pubsub.Topic)
	n.peers = make(map[peer.ID]*peer.AddrInfo)

	logger.LogInfo("Started node on %s\n", n.h.Addrs())
	logger.LogInfo("Self ID: %s\n", n.h.ID().String())

	ps, err := pubsub.NewGossipSub(context.Background(), n.h)
	if err != nil {
		logger.LogError("Error creating PubSub: %s\n", err)
	}
	n.ps = ps

	// Handle Peer Disconnections
	n.h.Network().Notify(&network.NotifyBundle{
		DisconnectedF: func(net network.Network, conn network.Conn) {
			_, ok := n.peers[conn.RemotePeer()]
			if ok {
				logger.LogDisconnectEvent("%s disconnected\n", conn.RemotePeer())
				delete(n.peers, conn.RemotePeer())
			}
		},
	})

	n.StartDiscovery()
}

func (n *MDNSNetwork) ListenBroadcast(topic string, handler func(sub *pubsub.Subscription, self peer.ID)) {
	_, ok := n.topics[topic]
	if !ok {
		t, err := n.ps.Join(topic)
		if err != nil {
			logger.LogError("Error Joining the topic %s: %s\n", topic, err)
		}
		n.topics[topic] = t
	}

	sub, err := n.topics[topic].Subscribe()
	if err != nil {
		logger.LogError("Error Subscribing to topic %s: %s\n", topic, err)
	}
	n.subs[topic] = sub

	logger.LogInfo("Listening to %s\n", sub.Topic())

	go handler(sub, n.h.ID())
}

func (n *MDNSNetwork) Broadcast(topic string, msg []byte) {
	logger.LogInfo("Broadcasting %s\n", string(msg))
	_, ok := n.topics[topic]

	if !ok {
		topicHandle, err := n.ps.Join(topic)
		if err != nil {
			logger.LogError("Error Joining the topic %s: %s", topic, err)
			return
		}

		n.topics[topic] = topicHandle
	}

	n.topics[topic].Publish(context.Background(), msg)
}

// Add a new stream to the network which is handled by the handler function
func (n *MDNSNetwork) AddStream(p string, handler func(stream network.Stream)) {
	n.h.SetStreamHandler(protocol.ID(p), handler)
	logger.LogInfo("Listening to Stream: %s", p)
}

func (n *MDNSNetwork) GetHost() host.Host {
	return n.h
}

// Send a message privately to the peer
func (n *MDNSNetwork) SendTo(proto string, p peer.ID, msg string) {
	stream, err := n.h.NewStream(context.Background(), p, protocol.ID(proto))
	if err != nil {
		logger.LogError("Error creating new stream: %s, Peer: %s", err, p.String())
		return
	}
	defer stream.Close()

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	_, err = rw.WriteString(msg + "\n")
	if err != nil {
		logger.LogError("Error writing to buffer: %s", err)
	}

	logger.LogInfo("Sending %s to %s\n", msg, p)

	err = rw.Flush()
	if err != nil {
		logger.LogError("Error flushing buffer: %s", err)
	}
}

func (n *MDNSNetwork) GetNumberOfPeers() int {
	return len(n.peers)
}

func (n *MDNSNetwork) GetPeers() map[peer.ID]*peer.AddrInfo {
	return n.peers
}
