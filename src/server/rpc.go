package server

import (
	"fmt"
	"net"
	"github.com/ethereum/go-ethereum/rpc"

	"time"
	"os"
)

var (
	//protoW                 = &sync.WaitGroup{}
	//messageW               = &sync.WaitGroup{}
	msgC                   = make(chan string)
	ipcpath                = ".demo.ipc"
)

type FooAPI struct {
	sent bool
}

func (api *FooAPI) SendMsg(content string) error {
	log.Info("SendMsg, ", content)

	msgC <- content

	return nil
}

func NewRPCServer() (*rpc.Server, error) {
	// set up the RPC server
	rpcsrv := rpc.NewServer()
	err := rpcsrv.RegisterName("foo", &FooAPI{})
	if err != nil {
		return nil, fmt.Errorf("register API method(s) fail: %v", err)
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
			log.WithError(err).Fatal("Mount RPC on IPC fail")
		} else {
			log.Info("RPC server stared")
		}

	}()

	return rpcsrv, nil
}

func StartRPC() *rpc.Server {
	log.Info("creating and starting RPC server")
	rpcsrv, err := NewRPCServer()
	if err != nil {
		log.WithError(err).Error("Could not StartRPC")
	} else {

		// create an IPC client
		rpcclient, err := rpc.Dial(ipcpath)
		if err != nil {
			os.Remove(ipcpath)
			log.WithError(err).Fatal("IPC dial fail")
		} else {

			go func() {
				for {
					<-time.After(5 * time.Second)
					go func() {
						// call the RPC method

						log.WithField("node", 1).WithField("Peers", srv.PeerCount()).Info("PEER INFO ")

						err = rpcclient.Call(nil, "foo_sendMsg", "foobar")

						if err != nil {
							log.WithError(err).Error("RPC call fail")
						}
					}()
				}
			}()
		}
	}

	return rpcsrv
}
