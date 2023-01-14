package model

import (
	"encoding/json"
)

type Effect string

type DigivolveCost struct {
	Cost json.Number   `json:"cost,omitempty"`
	From DigivolveFrom `json:"from,omitempty"`
}

type DNADigivolveCost struct {
	Cost json.Number     `json:"cost,omitempty"`
	From []DigivolveFrom `json:"from,omitempty"`
}

type DigivolveFrom struct {
	Color         string      `json:"color,omitempty"`
	Level         json.Number `json:"level,omitempty"`
	Name          string      `json:"name,omitempty"`
	Sources       []string    `json:"sources,omitempty"`
	Traits        []string    `json:"traits,omitempty"`
	TraitsContain []string    `json:"traitsContain,omitempty"`
	Tamer         bool        `json:"tamer,omitempty"`
}

type DigiXrosCost struct {
	Cost       json.Number        `json:"cost,omitempty"`
	FromNames  DigiXrosFromNames  `json:"fromNames,omitempty"`
	FromTraits DigiXrosFromTraits `json:"fromTraits,omitempty"`
}

type DigiXrosFromNames []string

type DigiXrosFromTraits struct {
	Count         json.Number `json:"count,omitempty"`
	Traits        []string    `json:"traits,omitempty"`
	TraitsContain []string    `json:"traitsContain,omitempty"`
	UniqueNames   bool        `json:"uniqueNames,omitempty"`
	UniqueNumbers bool        `json:"uniqueNumbers,omitempty"`
}

type Card struct {
	Set                         string             `json:"set,omitempty"`
	Name                        string             `json:"name,omitempty"`
	NameIncludes                []string           `json:"nameIncludes,omitempty"`
	NameTreatedAs               []string           `json:"nameTreatedAs,omitempty"`
	Number                      string             `json:"number,omitempty"`
	Rarity                      string             `json:"rarity,omitempty"`
	Type                        string             `json:"type,omitempty"`
	Colors                      []string           `json:"colors,omitempty"`
	Images                      []string           `json:"images,omitempty"`
	Form                        string             `json:"form,omitempty"`
	Attributes                  []string           `json:"attributes,omitempty"`
	Types                       []string           `json:"types,omitempty"`
	Effects                     []Effect           `json:"effects,omitempty"`
	InheritedEffects            []Effect           `json:"inheritedEffects,omitempty"`
	SecurityEffects             []Effect           `json:"securityEffects,omitempty"`
	Cost                        json.Number        `json:"cost,omitempty"`
	PlayCost                    json.Number        `json:"playCost,omitempty"`
	DP                          json.Number        `json:"dp,omitempty"`
	Level                       json.Number        `json:"level,omitempty"`
	Abilities                   []string           `json:"abilities,omitempty"`
	DigivolutionRequirements    []DigivolveCost    `json:"digivolutionRequirements,omitempty"`
	DNADigivolutionRequirements []DNADigivolveCost `json:"dnaDigivolutionRequirements,omitempty"`
	DigiXrosRequirements        []DigiXrosCost     `json:"digiXrosRequirements,omitempty"`
}
