package server

import (
	"net/http"
	"time"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/pat"
	"github.com/etherlabsio/healthcheck"

	"github.com/Sirupsen/logrus"
	"context"
)

var (
	router                *pat.Router
	log                   *logrus.Entry = logrus.WithField("package", "server")

)

var DEBUG = false

// SetLogger set the logger
func SetLogger(loggers *logrus.Entry) {
	log = loggers.WithFields(log.Data)
}

func init() {
	CockroachClient = SetupCockroachDB()

}

//NewServer return pointer to new created server object
func NewServer(Port string) *http.Server {
	router = InitRouting()
	return &http.Server{
		Addr:    ":" + Port,
		Handler: router,
	}
}

//StartServer start and listen @server
func StartServer(Port string, loggers *logrus.Entry) {
	//init log
	SetLogger(loggers)


	log.Info("Starting server")
	s := NewServer(Port)
	log.Info("Server starting --> " + Port)

	//enable graceful shutdown
	err := gracehttp.Serve(
		s,
	)

	if err != nil {
		log.Error("Error: %v", err)
		//os.Exit(0)
	}

}

func InitRouting() *pat.Router {

	r := pat.New()

	r.Handle("/healthcheck", healthcheck.Handler(
		// WithTimeout allows you to set a max overall timeout.
		healthcheck.WithTimeout(5*time.Second),

		healthcheck.WithChecker(
			"redis", healthcheck.CheckerFunc(
				func(ctx context.Context) error {
					return nil
				},
			),
		),
	))

	return r
}
