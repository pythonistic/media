package main

import (
	"os"
	"time"
	"crypto/tls"
	"net/http"
	"os/signal"
	"syscall"
	"github.com/pythonistic/media"
	"math/rand"
)

const EXIT_ERROR_HTTP = 1
const EXIT_ERROR_TLS = 2
const IS_TLS = false

var running bool
var server *http.Server

func main() {
	rand.Seed(time.Now().UnixNano())

	// a flag that can be used for graceful shutdowns
	running = true

	addr := ":8088"

	HandleSignals()

	media.LoadMedia()
	router := media.ConstructRoutes()

	if err := ListenAndServe(addr, router); err != nil {
		os.Stderr.WriteString("FATAL: listenning for web connections: " + err.Error() + "\n")
		os.Exit(EXIT_ERROR_HTTP)
	}
}

func ListenAndServe(addr string, handler http.Handler) error {
	readTimeout := 5 * time.Second
	writeTimeout := 25 * time.Second

	var tlsConfig *tls.Config
	var certFile, keyFile string

	if IS_TLS {
		certFile = "server.crt"
		keyFile = "server.key"
		insecureSkipVerify := false

		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			os.Stderr.WriteString("FATAL: error building TLS configuration: " + err.Error() + "\n")
			os.Exit(EXIT_ERROR_TLS)
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			InsecureSkipVerify: insecureSkipVerify,
		}
	}

	server = &http.Server{
		ReadTimeout: readTimeout,
		WriteTimeout: writeTimeout,
		Addr: addr,
		TLSConfig: tlsConfig,
		Handler: handler,
	}

	if IS_TLS {
		return server.ListenAndServeTLS(certFile, keyFile)
	} else {
		return server.ListenAndServe()
	}
}

func HandleSignals() {
	sigHandler := make(chan os.Signal, 1)

	// register to handle signals
	signal.Notify(sigHandler, syscall.SIGHUP)
	signal.Notify(sigHandler, syscall.SIGQUIT)

	// block until we receive SIGHUP
	go func() {
		for running {
			switch <-sigHandler {
			case syscall.SIGHUP:
				os.Stderr.WriteString("INFO: SIGHUP received\n")
				// TODO reload configuration
			case syscall.SIGQUIT:
				running = false
				server.Shutdown(nil)
			}
		}
	}()
}
