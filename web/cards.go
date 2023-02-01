package web

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/theandrew168/digimontcg/model"
)

func (app *Application) handleCards(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name   *string
		Number *string
		Set    *string
		Rarity *string
		Type   *string
		Color  *string
		Trait  *string
		Effect *string
		Level  *int
		Cost   *int
		DP     *int

		Sort *string

		Page *int
		Size *int
	}

	v := NewValidator()
	qs := r.URL.Query()

	// read all various filter options
	input.Name = readString(qs, "name")
	input.Number = readString(qs, "number")
	input.Set = readString(qs, "set")
	input.Rarity = readString(qs, "rarity")
	input.Type = readString(qs, "type")
	input.Color = readString(qs, "color")
	input.Trait = readString(qs, "trait")
	input.Effect = readString(qs, "effect")
	input.Level = readInt(qs, "level", v)
	input.Cost = readInt(qs, "cost", v)
	input.DP = readInt(qs, "dp", v)

	// read sorting order
	input.Sort = readString(qs, "sort")
	if input.Sort == nil {
		s := "number"
		input.Sort = &s
	}

	// read pagination params
	input.Page = readInt(qs, "page", v)
	if input.Page == nil {
		i := 1
		input.Page = &i
	}
	input.Size = readInt(qs, "size", v)
	if input.Size == nil {
		i := 20
		input.Size = &i
	}

	validSorts := []string{"name", "-name", "number", "-number"}
	v.Check(v.In(*input.Sort, validSorts...), "sort", "invalid sort value")

	v.Check(*input.Page > 0, "page", "must be greater than zero")
	v.Check(*input.Size > 0, "size", "must be greater than zero")
	v.Check(*input.Size <= 100, "size", "must be a maximum of 100")

	if !v.Valid() {
		failedValidationResponse(w, r, v.Errors)
		return
	}

	var cards []model.Card

	// apply filtering
	for _, card := range app.cards {
		if input.Name != nil {
			if !contains(card.Name, *input.Name) {
				continue
			}
		}

		if input.Number != nil {
			if !equals(card.Number, *input.Number) {
				continue
			}
		}

		if input.Set != nil {
			if !equals(card.Set, *input.Set) {
				continue
			}
		}

		if input.Rarity != nil {
			if !equals(card.Rarity, *input.Rarity) {
				continue
			}
		}

		if input.Type != nil {
			if !equals(card.Type, *input.Type) {
				continue
			}
		}

		if input.Color != nil {
			if !anyEquals(card.Colors, *input.Color) {
				continue
			}
		}

		if input.Trait != nil {
			if !contains(card.Form, *input.Trait) &&
				!anyContains(card.Attributes, *input.Trait) &&
				!anyContains(card.Types, *input.Trait) {
				continue
			}
		}

		if input.Effect != nil {
			if !anyContains(card.Effects, *input.Effect) &&
				!anyContains(card.InheritedEffects, *input.Effect) &&
				!anyContains(card.SecurityEffects, *input.Effect) {
				continue
			}
		}

		if input.Level != nil {
			if !equalsNumber(card.Level, *input.Level) {
				continue
			}
		}

		if input.Cost != nil {
			if !equalsNumber(card.Cost, *input.Cost) &&
				!equalsNumber(card.PlayCost, *input.Cost) {
				continue
			}
		}

		if input.DP != nil {
			if !equalsNumber(card.DP, *input.DP) {
				continue
			}
		}

		cards = append(cards, card)
	}

	// apply sorting
	switch *input.Sort {
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
	start := (*input.Page - 1) * *input.Size
	if start >= len(cards) {
		cards = []model.Card{}
		start = 0
	}

	end := start + *input.Size
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

func equals(a, b string) bool {
	return strings.ToLower(a) == strings.ToLower(b)
}

func anyEquals(a []string, b string) bool {
	for _, aa := range a {
		if equals(aa, b) {
			return true
		}
	}

	return false
}

func contains(a, b string) bool {
	return strings.Contains(strings.ToLower(a), strings.ToLower(b))
}

func anyContains(a []string, b string) bool {
	for _, aa := range a {
		if contains(aa, b) {
			return true
		}
	}

	return false
}

func equalsNumber(a json.Number, b int) bool {
	i, _ := a.Int64()
	return int(i) == b
}
