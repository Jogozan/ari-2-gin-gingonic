package pokemon

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterHTMLRoutes(r *gin.Engine) {
	r.GET("/pokemons", listPokemonsHTML)
	r.GET("/pokemons/new", newPokemonFormHTML)
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
				if t == typeFilter {
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

	// on peut pré-calculer le power pour l’affichage
	type PokemonView struct {
		Pokemon
		Power int
	}
	var views []PokemonView
	for _, p := range filtered {
		power := p.Stats.HP * p.Stats.Attack
		views = append(views, PokemonView{Pokemon: p, Power: power})
	}

	c.HTML(http.StatusOK, "pokemons_index.tmpl", gin.H{
		"title":    "Pokédex",
		"pokemons": views,
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

// Formulaire HTML
func newPokemonFormHTML(c *gin.Context) {
	// TODO
	c.HTML(http.StatusOK, "pokemons_form.tmpl", gin.H{
		"title":  "Nouveau Pokémon",
		"errors": []string{},
		"input":  CreatePokemonInput{},
	})
}
