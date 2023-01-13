package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/daemon"
)

//go:embed index.html
var index []byte

//go:embed api.html
var api []byte

//go:embed static
var static embed.FS

func main() {
	os.Exit(run())
}

func run() int {
	logger := log.New(os.Stdout, "", 0)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/v1/", handleAPI)
	mux.Handle("/static/", http.FileServer(http.FS(static)))

	port := "5000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	addr := fmt.Sprintf("127.0.0.1:%s", port)

	srv := http.Server{
		Handler: mux,

		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// open up the socket listener
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Println(err)
		return 1
	}

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)
	logger.Printf("started server on %s\n", addr)

	// kick off a goroutine to listen for SIGINT and SIGTERM
	shutdownError := make(chan error)
	go func() {
		// idle until a signal is caught
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		// give the web server 5 seconds to shutdown gracefully
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// shutdown the web server and track any errors
		logger.Println("stopping server")
		srv.SetKeepAlivesEnabled(false)
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		shutdownError <- nil
	}()

	// serve the app, check for ErrServerClosed (expected after shutdown)
	err = srv.Serve(l)
	if !errors.Is(err, http.ErrServerClosed) {
		logger.Println(err)
		return 1
	}

	// check for shutdown errors
	err = <-shutdownError
	if err != nil {
		logger.Println(err)
		return 1
	}

	logger.Println("stopped server")
	return 0
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write(index)
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Write(api)
}
