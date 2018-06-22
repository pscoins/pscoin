package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"

	"flag"
	"pscoin/src/server"
)

var (
	PORT     = os.Getenv("PORT")
	APP_NAME = os.Getenv("APP_NAME")
	log      *logrus.Entry
	ipcpath  = ".demo.ipc"
)

func initLog() {
	log = logrus.WithFields(logrus.Fields{
		"app": APP_NAME,
		"env": server.Env(),
	})

}

func init() {
	if PORT == "" {
		PORT = "3010"
	}

	initLog()

}

func main() {
	leader := flag.String("leader", "false", "Leader: true or false")
	port := flag.Int("port", 30300, "Port")
	host := flag.String("host", "0.0.0.0", "Host")

	peerPort := flag.Int("peerport", 0, "Peer Port")
	peerHost := flag.String("peerhost", "0.0.0.0", "Peer Host")
	peerID := flag.String("peerid", "1", "Peer Host")

	help := flag.Bool("help", false, "Display Help")

	flag.Parse()

	if *help {
		fmt.Printf("This program demonstrates a simple p2p chat application using libp2p\n\n")
		fmt.Printf("Usage: Run './pscoin -port <SOURCE_PORT>' where <SOURCE_PORT> can be any port number. Now run './pscoin -host <MULTIADDR>' where <MULTIADDR> is multiaddress of previous listener host.\n")

		os.Exit(0)
	}


	log.Info(leader)


	// start cockroachdb


	if *leader != "false" {
		log.Info("Starting RPC")
		os.Remove(ipcpath)
		<-time.After(1 * time.Second)
		rpcsrv := server.StartRPC()
		defer func() {
			rpcsrv.Stop()
			os.Remove(ipcpath)
		}()
	}

	server.InitP2p(*host, *port, *peerID, *peerHost, *peerPort)

	defer func() {
		handlePanic()
	}()

	server.StartServer(PORT, log)

}


func handlePanic() {
	if r := recover(); r != nil {
		log.WithError(fmt.Errorf("%+v", r)).Error(fmt.Sprintf("Application %s panic", APP_NAME))
	}
	time.Sleep(time.Second * 5)
}
