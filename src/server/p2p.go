package server

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"time"
	"math/rand"
	"strconv"
	"github.com/ethereum/go-ethereum/common"
	"pscoin/src/model"
)

type Message struct {
	Content string
}

var srv *p2p.Server

var (
	proto = p2p.Protocol{
		Name:    "foo",
		Version: 42,
		Length:  1,
		Run: func(peer *p2p.Peer, rw p2p.MsgReadWriter) error {

			go func() {

				for {

					select {

					case content, ok := <-msgC:

						if ok {
							outmsg := &Message{
								Content: content,
							}

							// send the message
							log.WithField("peer", peer).WithField("msg", outmsg).Info("p2p.Send, sending message")

							err := p2p.Send(rw, 0, outmsg)
							if err != nil {
								log.WithField("peer", peer).WithField("msg", outmsg).WithError(err).Error("p2p.Send, Could not send message")
								//return fmt.Errorf("send p2p message fail: %v", err)
							}
						} else {
							log.WithField("peer", peer).Error("p2p.Send, Could not send message")
							//return fmt.Errorf("send p2p message fail: %v")
						}
					}
				}
			}()

			for {
				msg, err := rw.ReadMsg()
				if err != nil {
					log.WithError(err).Error("GOT ERROR")
				}

				var myMessage Message
				err = msg.Decode(&myMessage)
				if err != nil {
					log.WithError(err).Error("GOT ERROR")
				}

				switch myMessage.Content {
				case "foobar":
					log.Info("AEEEEEEEEEEEEEEEEEEEEEE!!!!!!!!!!!!!!!!!!!!!!!!!!")
					//err := p2p.SendItems(rw, messageId, "bar")
					if err != nil {
						log.WithError(err).Error("GOT ERROR")
					}
					continue
				default:
					fmt.Println("GOT WEIRD MSG!!!!!!!!!!!!!!!!!!!!:", myMessage)
					continue
				}
			}
			log.Info("terminate the protocol ??")

			return nil
		},
	}
)

type Peer struct {
	ID         string
	Dest       string
	SourcePort int
}

func randomID() (id discover.NodeID) {
	for i := range id {
		id[i] = byte(rand.Intn(255))
	}
	return id
}

func AddPeer(node *discover.Node) error {
	log.Info("Going to connect ", node)

	// add it as a peer
	// the connection and crypto handshake will be performed automatically
	srv.AddPeer(node)

	// inspect the results
	log.Info("node peers", srv.Peers(), srv.PeerCount())

	return nil
}

func InitP2PPeers(peers []model.PeerNode) {
	log.Infof("InitP2PPeers, total peers %d", len(peers))
	for i, peer := range peers {

		if &peer.ID != nil {

			peer, err := discover.ParseNode(peer.ID)

			if err != nil {
				log.WithField("method", "InitP2PPeers").
					WithField("peer.ID", peer.ID).
					WithError(err).
					Error("Failed to connect to peer")
			} else {

				err = AddPeer(peer)
				if err != nil {
					log.WithField("method", "InitP2PPeers").
						WithField("peer.ID", peer.ID).
						WithError(err).
						Error("Failed to connect to peer")
				} else {

					eventOneC := make(chan *p2p.PeerEvent)
					sub_one := srv.SubscribeEvents(eventOneC)
					go func() {

						for {
							select {
							case peerevent := <-eventOneC:
								if peerevent.Type == "add" {
									log.WithField("Type", peerevent.Type).WithField("Peer", peerevent.Peer).WithField("node", i).Info("Received peer add notification")
								} else if peerevent.Type == "msgsend" {
									log.WithField("Type", peerevent.Type).WithField("Event", peerevent).WithField("node", i).Info("Received message send notification")
								} else if peerevent.Type == "drop" {
									log.WithField("Type", peerevent.Type).WithField("Event", peerevent).WithField("node", i).Error("Received DROP message")
								} else {
									log.WithField("Type", peerevent.Type).WithField("Event", peerevent).WithField("node", i).Info("Received message")
								}
							case <-sub_one.Err():
								return
							}
						}
					}()
				}
			}
		} else {
			log.WithField("method", "InitP2PPeers").
				WithField("peer.ID", peer.ID).Error("Could not get a Peer instance from this peer")
		}

	}

	// wait for the connections to complete
	time.Sleep(time.Millisecond * 1000)

}

var protocols []p2p.Protocol

func InitP2PServer(dest string, sourcePort int) (*p2p.Server, error) {

	strCon := fmt.Sprintf("%s:%d", dest, sourcePort)
	log.WithField("strCon", strCon).Info("Starting P2P ")

	nodekey, err := crypto.GenerateKey()
	if err != nil {
		log.WithError(err).Error("Generate private key failed")
		return nil, err
	}

	srv = &p2p.Server{
		Config: p2p.Config{
			MaxPeers:        10,
			PrivateKey:      nodekey,
			Name:            common.MakeName("pscoin node "+strconv.Itoa(sourcePort), "1"),
			ListenAddr:      strCon,
			Protocols:       []p2p.Protocol{proto},
			EnableMsgEvents: true,
		},
	}

	// the Err() on the Subscription object returns when subscription is closed

	if err := srv.Start(); err != nil {
		log.WithError(err).Error("Could not start p2p Server")
		return srv, err
	}

	return srv, nil

}

func InitP2p(host string, port int) {

	srv, err := InitP2PServer(host, port)
	if err != nil {
		log.WithError(err).Fatal("Could not start P2P Server")
	}

	srv_node := srv.Self()

	//verify if this node is already on our database
	//save to db

	targetObject := model.PeerNode{
		ID:       srv_node.String(),
		LastSeen: time.Now(),
	}

	CockroachClient.Save(&targetObject)

	var peerNodes []model.PeerNode
	//get all peers except yourself
	err = CockroachClient.Model(model.PeerNode{}).Where("id <> ?",targetObject.ID).Order("last_seen DESC").Limit(100).Find(&peerNodes).Error

	if err != nil {
		log.WithError(err).Error("Could not get peer list")
		return
	}

	InitP2PPeers(peerNodes)

}
