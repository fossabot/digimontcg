package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
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

//go:embed data/cards/*.json
var cardsFS embed.FS

func main() {
	os.Exit(run())
}

func run() int {
	log.SetOutput(os.Stdout)
	log.SetPrefix("")
	log.SetFlags(0)

	// all set data lives within a single data file
	var sets []model.Set
	err := json.Unmarshal(setsJSON, &sets)
	if err != nil {
		log.Println(err)
		return 1
	}

	// read each set's worth of card data one at a time
	var cards []model.Card
	err = fs.WalkDir(cardsFS, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		data, _ := cardsFS.ReadFile(path)

		var set []model.Card
		err = json.Unmarshal(data, &set)
		if err != nil {
			return err
		}

		cards = append(cards, set...)
		return nil
	})
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
