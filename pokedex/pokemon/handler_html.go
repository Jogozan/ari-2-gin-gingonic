package pokemon

import (
	"net/http"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
)

func RegisterHTMLRoutes(r *gin.Engine) {
	r.GET("/pokemons", listPokemonsHTML)
	r.GET("/pokemons/:id", pokemonDetailHTML)
}

func listPokemonsHTML(c *gin.Context) {
	typeFilter := c.Query("type")
	minLevelStr := c.Query("minLevel")

	all := GetAll()
	var filtered []Pokemon
	for _, p := range all {
		// même logique que pour l’API
		if typeFilter != "" {
			found := false
			for _, t := range p.Types {
				if t == strings.ToLower(typeFilter) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if minLevelStr != "" {
			minLevel, err := strconv.Atoi(minLevelStr)
			if err == nil && p.Stats.HP < minLevel {
				continue
			}
		}
		filtered = append(filtered, p)
	}

	// map -> DTO avec Power
	var resp []PokemonResponse
	for _, p := range filtered {
		resp = append(resp, toResponse(p))
	}

	c.HTML(http.StatusOK, "pokemons_index.tmpl", gin.H{
		"title":    "Pokédex",
		"pokemons": resp,
		"type":     typeFilter,
		"minLevel": minLevelStr,
	})
}

func pokemonDetailHTML(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "ID invalide")
		return
	}

	p, err := GetByID(id)
	if err != nil {
		c.String(http.StatusNotFound, "Pokemon non trouvé")
		return
	}

	c.HTML(http.StatusOK, "pokemons_detail.tmpl", gin.H{
		"title":    p.Name,
		"pokemon":  p,
		"pokemons": GetAll(), // pour navigation si besoin
	})
}
