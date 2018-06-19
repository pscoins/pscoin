package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"

	"pscoin/src/server"
	"pscoin/src/p2p"
	"flag"
)

var (
	PORT     = os.Getenv("PORT")
	APP_NAME = os.Getenv("APP_NAME")
	log      *logrus.Entry
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

}

func main() {
	initLog()

	defer handlePanic()

	initBlockchain()

	server.StartServer(PORT, log)
}

func initBlockchain() {
	port := flag.Int("port", 30300, "Port")
	host := flag.String("host", "0.0.0.0", "Host")
	help := flag.Bool("help", false, "Display Help")

	flag.Parse()

	if *help {
		fmt.Printf("This program demonstrates a simple p2p chat application using libp2p\n\n")
		fmt.Printf("Usage: Run './pscoin -port <SOURCE_PORT>' where <SOURCE_PORT> can be any port number. Now run './pscoin -host <MULTIADDR>' where <MULTIADDR> is multiaddress of previous listener host.\n")

		os.Exit(0)
	}

	p2p.InitP2p(*host, *port)
}

func handlePanic() {
	if r := recover(); r != nil {
		log.WithError(fmt.Errorf("%+v", r)).Error(fmt.Sprintf("Application %s panic", APP_NAME))
	}
	time.Sleep(time.Second * 5)
}
