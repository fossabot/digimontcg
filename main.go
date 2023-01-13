package main

import (
	"context"
	_ "embed"
	"encoding/json"
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

	"github.com/theandrew168/digimontcg/model"
	"github.com/theandrew168/digimontcg/web"
)

//go:embed data/sets.json
var setsJSON []byte

//go:embed data/cards.json
var cardsJSON []byte

func main() {
	os.Exit(run())
}

func run() int {
	log.SetOutput(os.Stdout)
	log.SetPrefix("")
	log.SetFlags(0)

	var sets []model.Set
	err := json.Unmarshal(setsJSON, &sets)
	if err != nil {
		log.Println(err)
		return 1
	}

	var cards []model.Card
	err = json.Unmarshal(cardsJSON, &cards)
	if err != nil {
		log.Println(err)
		return 1
	}

	app := web.NewApplication(sets, cards)

	port := "5000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	addr := fmt.Sprintf("127.0.0.1:%s", port)

	srv := http.Server{
		Handler: app.Handler(),

		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// open up the socket listener
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err)
		return 1
	}

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)
	log.Printf("started server on %s\n", addr)

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
		log.Println("stopping server")
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
		log.Println(err)
		return 1
	}

	// check for shutdown errors
	err = <-shutdownError
	if err != nil {
		log.Println(err)
		return 1
	}

	log.Println("stopped server")
	return 0
}
