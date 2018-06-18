package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"

	"pscoin/src/server"
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

	server.StartServer(PORT, log)
}

func handlePanic() {
	if r := recover(); r != nil {
		log.WithError(fmt.Errorf("%+v", r)).Error(fmt.Sprintf("Application %s panic", APP_NAME))
	}
	time.Sleep(time.Second * 5)
}
