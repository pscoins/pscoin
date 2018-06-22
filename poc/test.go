// trigger p2p message with RPC
package main

import (
	"crypto/ecdsa"
	"fmt"
	"net"
	"os"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/Sirupsen/logrus"
	"time"
)

var (
	//protoW                 = &sync.WaitGroup{}
	//messageW               = &sync.WaitGroup{}
	msgC                   = make(chan string)
	ipcpath                = ".demo.ipc"
	log      *logrus.Entry = logrus.WithField("package", "p2p")
)

// create a protocol that can take care of message sending
// the Run function is invoked upon connection
// it gets passed:
// * an instance of p2p.Peer, which represents the remote peer
// * an instance of p2p.MsgReadWriter, which is the io between the node and its peer

type FooMsg struct {
	Content string
}

var (
	proto = p2p.Protocol{
		Name:    "foo",
		Version: 42,
		Length:  1,
		Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {

			// only one of the peers will send this

			go func (){

			for {
				select {
				case content, ok := <-msgC:
					if ok {
						outmsg := &FooMsg{
							Content: content,
						}

						// send the message
						log.WithField("peer", p).WithField("msg", outmsg).Info("p2p.Send, sending message")

						err := p2p.Send(rw, 0, outmsg)
						if err != nil {
							log.WithField("peer", p).WithField("msg", outmsg).WithError(err).Error("p2p.Send, Could not send message")
							//return fmt.Errorf("send p2p message fail: %v", err)
						}
					} else {
						log.WithField("peer", p).Error("p2p.Send, Could not send message")
						//return fmt.Errorf("send p2p message fail: %v")
					}
				}
			}
			}()

			//go func(){

			for {
				msg, err := rw.ReadMsg()
				if err != nil {
					log.WithError(err).Error("GOT ERROR")
				}

				var myMessage FooMsg
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
			//}()

			//// wait for the subscriptions to end
			//messageW.Wait()
			//protoW.Done()

			log.Info("terminate the protocol ??")
			return nil
		},
	}
)

type FooAPI struct {
	sent bool
}

func (api *FooAPI) SendMsg(content string) error {
	log.Info("SendMsg, ", content)
	//if api.sent {
	//	return fmt.Errorf("Already sent")
	//}
	msgC <- content
	//close(msgC)
	//api.sent = true
	return nil
}

// create a server
func newP2pServer(privkey *ecdsa.PrivateKey, name string, version string, port int) *p2p.Server {
	// we need to explicitly allow at least one peer, otherwise the connection attempt will be refused
	// we also need to explicitly tell the server to generate events for messages
	cfg := p2p.Config{
		PrivateKey:      privkey,
		Name:            common.MakeName(name, version),
		MaxPeers:        1,
		Protocols:       []p2p.Protocol{proto},
		EnableMsgEvents: true,
	}
	if port > 0 {
		cfg.ListenAddr = fmt.Sprintf(":%d", port)
	}
	srv := &p2p.Server{
		Config: cfg,
	}
	return srv
}

func newRPCServer() (*rpc.Server, error) {
	// set up the RPC server
	rpcsrv := rpc.NewServer()
	err := rpcsrv.RegisterName("foo", &FooAPI{})
	if err != nil {
		return nil, fmt.Errorf("Register API method(s) fail: %v", err)
	}

	// create IPC endpoint
	ipclistener, err := net.Listen("unix", ipcpath)
	if err != nil {
		return nil, fmt.Errorf("IPC endpoint create fail: %v", err)
	}

	// mount RPC server on IPC endpoint
	// it will automatically detect and serve any valid methods
	go func() {
		err = rpcsrv.ServeListener(ipclistener)
		if err != nil {
			log.Error("Mount RPC on IPC fail", "err", err)
		}
	}()

	return rpcsrv, nil
}

func init(){
	os.Remove(ipcpath)
}

func main() {

	// we need private keys for both servers
	privkey_one, err := crypto.GenerateKey()
	if err != nil {
		log.Error("Generate private key #1 failed", "err", err)
	}
	privkey_two, err := crypto.GenerateKey()
	if err != nil {
		log.Error("Generate private key #2 failed", "err", err)
	}

	// set up the two servers
	srv_one := newP2pServer(privkey_one, "foo", "42", 0)
	err = srv_one.Start()
	if err != nil {
		log.Error("Start p2p.Server #1 failed", "err", err)
	}

	srv_two := newP2pServer(privkey_two, "bar", "666", 31234)
	err = srv_two.Start()
	if err != nil {
		log.Error("Start p2p.Server #2 failed", "err", err)
	}

	// set up the event subscriptions on both servers
	// the Err() on the Subscription object returns when subscription is closed
	eventOneC := make(chan *p2p.PeerEvent)
	sub_one := srv_one.SubscribeEvents(eventOneC)
	go func() {
		for {
			select {
			case peerevent := <-eventOneC:
				if peerevent.Type == "add" {
					log.WithField("Type",peerevent.Type).WithField("Peer", peerevent.Peer).WithField("node",1).Info("Received peer add notification")
				} else if peerevent.Type == "msgsend" {
					log.WithField("Type",peerevent.Type).WithField("Event", peerevent).WithField("node", 1).Info("Received message send notification")
				} else if peerevent.Type == "drop" {
					log.WithField("Type",peerevent.Type).WithField("Event", peerevent).WithField("node",1).Error("Received DROP message")
				} else {
					log.WithField("Type",peerevent.Type).WithField("Event", peerevent).WithField("node",1).Info("Received message")
				}
			case <-sub_one.Err():
				return
			}
		}
	}()

	eventTwoC := make(chan *p2p.PeerEvent)
	sub_two := srv_two.SubscribeEvents(eventTwoC)
	go func() {
		for {
			select {
			case peerevent := <-eventTwoC:
				if peerevent.Type == "add" {
					log.WithField("Type",peerevent.Type).WithField("Peer", peerevent.Peer).WithField("node",2).Info("Received peer add notification")
				} else if peerevent.Type == "msgsend" {
					log.WithField("Type",peerevent.Type).WithField("Event", peerevent).WithField("node",2).Info("Received message send notification")
				} else if peerevent.Type == "drop" {
					log.WithField("Type",peerevent.Type).WithField("Event", peerevent).WithField("node",2).Error("Received DROP message")
				} else {
					log.WithField("Type",peerevent.Type).WithField("Event", peerevent).WithField("node",2).Info("Received message")
				}
			case <-sub_two.Err():
				return
			}
		}
	}()

	// create and start RPC server
	rpcsrv, err := newRPCServer()
	if err != nil {
		log.Error(err.Error())
	}
	defer os.Remove(ipcpath)

	// get the node instance of the second server
	node_two := srv_two.Self()


	log.Info("Node", node_two)

	// add it as a peer to the first node
	// the connection and crypto handshake will be performed automatically
	srv_one.AddPeer(node_two)

	// create an IPC client
	rpcclient, err := rpc.Dial(ipcpath)
	if err != nil {
		log.Error("IPC dial fail", "err", err)
	}

	defer rpcsrv.Stop()
	//// wait for one message be sent, and both protocols to end
	//messageW.Add(1)
	//protoW.Add(2)

	// call the RPC method
	go func() {
		for {
			<-time.After(5 * time.Second)
			go func() {
				// call the RPC method
				err = rpcclient.Call(nil, "foo_sendMsg", "foobar")
				log.WithField("node", 1).WithField("Peers", srv_one.PeerCount()).Info("PEER INFO ")
				log.WithField("node", 2).WithField("Peers",srv_two.PeerCount()).Info("PEER INFO ")

				if err != nil {
					log.WithError(err).Error("RPC call fail")
					return
				}
			}()
		}
	}()

	//// wait for protocols to finish
	//protoW.Wait()
	//
	//// terminate subscription loops and unsubscribe
	//sub_one.Unsubscribe()
	//sub_two.Unsubscribe()

	// stop the servers

	//srv_one.Stop()
	//srv_two.Stop()

	select  {}

}
