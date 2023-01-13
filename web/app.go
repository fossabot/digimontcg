package web

import (
	"embed"
	"net/http"

	"github.com/theandrew168/digimontcg/model"
)

//go:embed static
var staticFS embed.FS

type Application struct {
	sets  []model.Set
	cards []model.Card
}

func NewApplication(sets []model.Set, cards []model.Card) *Application {
	app := Application{
		sets:  sets,
		cards: cards,
	}
	return &app
}

func (app *Application) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.handleIndex)
	mux.HandleFunc("/api/v1/", app.handleRapidoc)
	mux.HandleFunc("/api/v1/sets", app.handleSets)
	mux.HandleFunc("/api/v1/cards", app.handleCards)
	mux.Handle("/static/", http.FileServer(http.FS(staticFS)))
	return mux
}
