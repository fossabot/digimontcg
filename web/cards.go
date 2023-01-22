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
		Name   string
		Number string
		Set    string
		Rarity string

		Sort string

		Page int
		Size int
	}

	v := NewValidator()
	qs := r.URL.Query()

	input.Name = readString(qs, "name", "")
	input.Number = readString(qs, "number", "")
	input.Set = readString(qs, "set", "")
	input.Rarity = readString(qs, "rarity", "")

	input.Sort = readString(qs, "sort", "number")

	input.Page = readInt(qs, "page", 1, v)
	input.Size = readInt(qs, "size", 20, v)

	validSorts := []string{"name", "-name", "number", "-number"}
	v.Check(v.In(input.Sort, validSorts...), "sort", "invalid sort value")

	v.Check(input.Page > 0, "page", "must be greater than zero")
	v.Check(input.Size > 0, "size", "must be greater than zero")
	v.Check(input.Size <= 100, "size", "must be a maximum of 100")

	if !v.Valid() {
		failedValidationResponse(w, r, v.Errors)
		return
	}

	var cards []model.Card

	// apply filtering
	for _, card := range app.cards {
		if input.Name != "" {
			if !strings.Contains(strings.ToLower(card.Name), strings.ToLower(input.Name)) {
				continue
			}
		}

		if input.Number != "" {
			if strings.ToLower(card.Number) != strings.ToLower(input.Number) {
				continue
			}
		}

		if input.Set != "" {
			if strings.ToLower(card.Set) != strings.ToLower(input.Set) {
				continue
			}
		}

		if input.Rarity != "" {
			if strings.ToLower(card.Rarity) != strings.ToLower(input.Rarity) {
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
	start := (input.Page - 1) * input.Size
	if start >= len(cards) {
		cards = []model.Card{}
		start = 0
	}

	end := start + input.Size
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
