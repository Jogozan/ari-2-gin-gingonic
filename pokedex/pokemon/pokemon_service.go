package pokemon

import (
	"sort"
	"strconv"
	"strings"
)

// ---- Helpers métier ----

func matchesType(p Pokemon, typeFilter string) bool {
	if typeFilter == "" {
		return true
	}
	for _, t := range p.Types {
		if strings.EqualFold(t, typeFilter) {
			return true
		}
	}
	return false
}

func matchesMinLevel(p Pokemon, minLevelStr string) bool {
	if minLevelStr == "" {
		return true
	}
	minLevel, err := strconv.Atoi(minLevelStr)
	if err != nil {
		return true // si minLevel est invalide, on ignore le filtre
	}
	return p.Stats.HP >= minLevel
}

func filterPokemons(list []Pokemon, typeFilter, minLevelStr string) []Pokemon {
	var filtered []Pokemon
	for _, p := range list {
		if !matchesType(p, typeFilter) {
			continue
		}
		if !matchesMinLevel(p, minLevelStr) {
			continue
		}
		filtered = append(filtered, p)
	}
	return filtered
}

func sortPokemons(list []PokemonResponse, sortBy string) []PokemonResponse {
	copyList := make([]PokemonResponse, len(list))
	copy(copyList, list)

	switch sortBy {
	case "level":
		sort.Slice(copyList, func(i, j int) bool {
			return copyList[i].Level > copyList[j].Level
		})
	case "power":
		sort.Slice(copyList, func(i, j int) bool {
			return copyList[i].Power > copyList[j].Power
		})
		// default: pas de tri, on garde l'ordre
	}

	return copyList
}

func toResponses(list []Pokemon) []PokemonResponse {
	resp := make([]PokemonResponse, 0, len(list))
	for _, p := range list {
		resp = append(resp, toResponse(p))
	}
	return resp
}

// typeColorMap returns a map of type->css color used for badges.
// Keep the palette small and readable; it's purely presentational for the TP.
func typeColorMap() map[string]string {
	return map[string]string{
		"grass":    "#78C850",
		"fire":     "#F08030",
		"water":    "#6890F0",
		"electric": "#F8D030",
		"ice":      "#98D8D8",
		"psychic":  "#F85888",
		"ghost":    "#705898",
		"dark":     "#705848",
		"rock":     "#B8A038",
		"steel":    "#B8B8D0",
		"ground":   "#E0C068",
		"flying":   "#A890F0",
		"bug":      "#A8B820",
		"poison":   "#A040A0",
		"normal":   "#A8A878",
	}
}
