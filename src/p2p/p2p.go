package p2p

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/Sirupsen/logrus"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"net"
	"time"
	"math/rand"
	"os"
	"strconv"
)

const messageId = 0

type Message struct {
	Content string
}

var log *logrus.Entry = logrus.WithField("package", "p2p")
var srv *p2p.Server

func MyProtocol() p2p.Protocol {
	return p2p.Protocol{
		Name:    "MyProtocol",
		Version: 1,
		Length:  1,
		Run:     msgHandler,
	}
}

func randomID() (id discover.NodeID) {
	for i := range id {
		id[i] = byte(rand.Intn(255))
	}
	return id
}

func AddPeer(id, dest string, sourcePort int) error {

	nodeID := randomID()

	log.Info(nodeID, net.ParseIP(dest), uint16(sourcePort), uint16(sourcePort))
	node := discover.NewNode(nodeID, net.ParseIP(dest), uint16(sourcePort), uint16(sourcePort))

	// add it as a peer
	// the connection and crypto handshake will be performed automatically
	srv.AddPeer(node)

	// wait for the connection to complete
	time.Sleep(time.Millisecond * 100)

	// inspect the results
	log.Info("node peers", srv.Peers(), srv.PeerCount())

	return nil
}

func msgHandler(peer *p2p.Peer, ws p2p.MsgReadWriter) error {
	outmsg := &Message{
		Content: "testing 123",
	}

	// send the message
	err := p2p.Send(ws, 0, outmsg)
	if err != nil {
		return fmt.Errorf("Send p2p message fail: %v", err)
	}
	log.Info("sending message", "peer", peer, "msg", outmsg)

	for {
		msg, err := ws.ReadMsg()
		if err != nil {
			return err
		}

		var myMessage Message
		err = msg.Decode(&myMessage)
		if err != nil {
			// handle decode error
			continue
		}

		switch myMessage.Content {
		case "foo":
			err := p2p.SendItems(ws, messageId, "bar")
			if err != nil {
				return err
			}
		default:
			fmt.Println("recv:", myMessage)
		}
	}

	return nil
}

type Peer struct {
	ID         string
	Dest       string
	SourcePort int
}

func GetPeers(id, dest string, sourcePort int) (peers []Peer) {
	//TODO: open database

	//TODO: get peers

	//TODO: connect to all peers

	return []Peer{
		{
			ID:         id,
			Dest:       dest,
			SourcePort: sourcePort,
		},
	}
}

func InitP2PPeers(peers []Peer) {
	for _, peer := range peers {

		err := AddPeer(peer.ID, peer.Dest, peer.SourcePort)
		if err != nil {
			log.WithField("method", "InitP2PPeers").
				WithField("peer.ID", peer.ID).
				WithField("peer.Dest", peer.Dest).
				WithField("peer.SourcePort", peer.SourcePort).
				WithError(err).
				Error("Failed to connect to peer")
		}

	}

}

var protocols []p2p.Protocol

func InitP2PServer(dest string, sourcePort int) {

	strCon := fmt.Sprintf("%s:%d", dest, sourcePort)
	log.Info("Starting P2P ", strCon)

	protocols = []p2p.Protocol{MyProtocol()}

	nodekey, _ := crypto.GenerateKey()
	srv = &p2p.Server{
		Config: p2p.Config{
			MaxPeers:        10,
			PrivateKey:      nodekey,
			Name:            "my node name "+strconv.Itoa(sourcePort),
			ListenAddr:      strCon,
			Protocols:       protocols,
			EnableMsgEvents: true,
		},
	}



	// the Err() on the Subscription object returns when subscription is closed
	eventOneC := make(chan *p2p.PeerEvent)
	sub_one := srv.SubscribeEvents(eventOneC)
	go func() {
		for {
			select {
			case peerevent := <-eventOneC:
				if peerevent.Type == "add" {
					log.Info("Received peer add notification on node #1", "peer", peerevent.Peer)
				} else if peerevent.Type == "msgsend" {
					log.Info("Received message send notification on node #1", "event", peerevent)
				}
			case <-sub_one.Err():
				return
			}
		}
	}()


	if err := srv.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
