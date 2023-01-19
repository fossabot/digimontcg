package web

import (
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/theandrew168/digimontcg/model"
)

func (app *Application) handleCards(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string
		Page     int
		PageSize int
		Sort     string
	}

	v := NewValidator()
	qs := r.URL.Query()

	input.Name = readString(qs, "name", "")
	input.Page = readInt(qs, "page", 1, v)
	input.PageSize = readInt(qs, "page_size", 20, v)
	input.Sort = readString(qs, "sort", "number")

	v.Check(input.Page > 0, "page", "must be greater than zero")
	v.Check(input.Page <= 1_000_000, "page", "must be a maximum of 1 million")
	v.Check(input.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(input.PageSize <= 100, "page_size", "must be a maximum of 100")

	validSorts := []string{"name", "-name", "number", "-number"}
	v.Check(v.In(input.Sort, validSorts...), "sort", "invalid sort value")

	if !v.Valid() {
		failedValidationResponse(w, r, v.Errors)
		return
	}

	var cards []model.Card

	// apply filtering
	for _, card := range app.cards {
		if input.Name != "" {
			name := strings.ToLower(card.Name)
			if !strings.Contains(name, strings.ToLower(input.Name)) {
				continue
			}
		}

		cards = append(cards, card)
	}

	// apply sorting
	switch input.Sort {
	case "name":
		sort.Slice(cards, func(i, j int) bool {
			return cards[i].Name < cards[j].Name
		})
	case "-name":
		sort.Slice(cards, func(i, j int) bool {
			return cards[i].Name > cards[j].Name
		})
	case "number":
		sort.Slice(cards, func(i, j int) bool {
			return cards[i].Number < cards[j].Number
		})
	case "-number":
		sort.Slice(cards, func(i, j int) bool {
			return cards[i].Number > cards[j].Number
		})
	}

	// apply pagination
	start := (input.Page - 1) * input.PageSize
	if start >= len(cards) {
		cards = []model.Card{}
		start = 0
	}

	end := start + input.PageSize
	if end > len(cards) {
		end = len(cards)
	}

	cards = cards[start:end]

	err := writeJSON(w, 200, envelope{"cards": cards})
	if err != nil {
		log.Println(err)

		code := 500
		http.Error(w, http.StatusText(code), code)
		return
	}
}
