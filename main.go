package main

import (
	"context"
	"embed"
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
)

//go:embed data/sets.json
var setsJSON []byte

//go:embed data/cards.json
var cardsJSON []byte

//go:embed template/index.html
var indexHTML []byte

//go:embed template/api.html
var apiHTML []byte

//go:embed static
var staticFS embed.FS

func main() {
	os.Exit(run())
}

func run() int {
	log.SetOutput(os.Stdout)
	log.SetPrefix("")
	log.SetFlags(0)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/v1/", handleAPI)
	mux.HandleFunc("/api/v1/sets", handleSets)
	mux.HandleFunc("/api/v1/cards", handleCards)
	mux.Handle("/static/", http.FileServer(http.FS(staticFS)))

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

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write(indexHTML)
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Write(apiHTML)
}

func handleSets(w http.ResponseWriter, r *http.Request) {
	var sets []model.Set
	err := json.Unmarshal(setsJSON, &sets)
	if err != nil {
		log.Println(err)

		code := 500
		http.Error(w, http.StatusText(code), code)
		return
	}

	err = writeJSON(w, 200, envelope{"sets": sets})
	if err != nil {
		log.Println(err)

		code := 500
		http.Error(w, http.StatusText(code), code)
		return
	}
}

func handleCards(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("cards"))
}

type envelope map[string]interface{}

func writeJSON(w http.ResponseWriter, status int, data envelope) error {
	// attempt to encode data into JSON
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// append a newline for nicer terminal output
	js = append(js, '\n')

	// set content type, set status, and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
