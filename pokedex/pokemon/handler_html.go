package pokemon

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// RegisterHTMLRoutes registers routes that render HTML pages for the pokedex.
// These handlers return rendered templates (index, detail, and form) instead of JSON.
func RegisterHTMLRoutes(r *gin.Engine) {
	r.GET("/pokemons", listPokemonsHTML)
	r.GET("/pokemons/new", newPokemonFormHTML)
	r.POST("/pokemons", createPokemonHTML)
	r.GET("/pokemons/:id", pokemonDetailHTML)
	r.POST("/pokemons/:id/level-up", levelUpPokemonHTML)
	r.POST("/pokemons/:id/release", releasePokemonHTML)
	r.GET("/pokemons/stats", pokemonsStatsHTML)
}

func listPokemonsHTML(c *gin.Context) {
	typeFilter := c.Query("type")
	minLevelStr := c.Query("minLevel")
	sortBy := c.Query("sort")

	all := GetAll()

	// 1) Filtrage
	filtered := filterPokemons(all, typeFilter, minLevelStr)

	// 3) Wrapper dans le DTO
	resp := toResponses(filtered)

	// 2) Tri
	pokemons := sortPokemons(resp, sortBy)

	c.HTML(http.StatusOK, "pokemons_index.tmpl", gin.H{
		"title":       "Pokédex",
		"pokemons":    pokemons, // ou `sorted` si le template attend le modèle brut
		"type":        typeFilter,
		"minLevel":    minLevelStr,
		"sort":        sortBy,
		"current_url": c.Request.URL.RequestURI(),
		"message":     c.Query("msg"),
		"type_colors": typeColorMap(),
	})
}

// pokemonDetailHTML handles GET /pokemons/:id and renders the detail page
// for a single pokemon. It returns HTTP errors as plain text when the id is invalid.
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
		"title":       p.Name,
		"pokemon":     p,
		"pokemons":    GetAll(), // pour navigation si besoin
		"current_url": c.Request.URL.RequestURI(),
		"type_colors": typeColorMap(),
	})
}

// levelUpPokemonHTML handles simple form POSTs originating from HTML views
// It performs a LevelUp(+1) and redirects back to the index with a message
func levelUpPokemonHTML(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, []string{"ID invalide"})
		return
	}

	// By default +1 level; optionally caller can provide 'levels' form value
	addLevels := 1
	if lvlStr := c.Query("levels"); lvlStr != "" {
		if lv, err := strconv.Atoi(lvlStr); err == nil && lv > 0 {
			addLevels = lv
		}
	}

	p, err := LevelUp(id, addLevels)
	if err != nil {
		if errors.Is(err, ErrMaxLevel) {
			RespondError(c, http.StatusBadRequest, []string{"Niveau maximum atteint"})
			return
		}
		RespondError(c, http.StatusNotFound, []string{"Pokemon non trouvé"})
		return
	}

	msg := p.Name + " a gagné " + strconv.Itoa(addLevels) + " niveau(x) !"
	redirectWithMessage(c, msg)
}

// releasePokemonHTML handles releasing (deleting) a pokemon from HTML
// and performs a redirect back to the index page.
func releasePokemonHTML(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, []string{"ID invalide"})
		return
	}

	p, err := GetByID(id)
	if err != nil {
		RespondError(c, http.StatusNotFound, []string{"Pokemon non trouvé"})
		return
	}

	if err := Delete(id); err != nil {
		RespondError(c, http.StatusInternalServerError, []string{"Impossible de supprimer"})
		return
	}

	msg := p.Name + " a été relâché 😢"
	redirectWithMessage(c, msg)
}

// pokemonsStatsHTML builds a small aggregation of pokemon counts per type
func pokemonsStatsHTML(c *gin.Context) {
	all := GetAll()
	counts := map[string]int{}
	for _, p := range all {
		for _, t := range p.Types {
			counts[t]++
		}
	}

	c.HTML(http.StatusOK, "pokemons_stats.tmpl", gin.H{
		"title":       "Statistiques du dresseur",
		"type_counts": counts,
		"type_colors": typeColorMap(),
	})
}

// Formulaire HTML
// newPokemonFormHTML renders the empty form used to create a new pokemon.
// The template expects `title`, `errors` and `input` values in its context.
func newPokemonFormHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "pokemons_form.tmpl", gin.H{
		"title":  "Nouveau Pokémon",
		"errors": []string{},
		"input":  CreatePokemonInput{},
	})
}

// createPokemonHTML accepts JSON payloads and returns HTML-friendly error maps.
// It demonstrates converting validator errors into a `map[string][]string`
// structure so templates / clients can display per-field messages.
func createPokemonHTML(c *gin.Context) {
	var input CreatePokemonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		// Map[string][]string : nom du champ -> liste de messages
		fieldErrors := map[string][]string{}

		if verrs, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range verrs {
				field := fe.Field()
				fieldErrors[field] = append(fieldErrors[field], validationMessage(field))
			}
		}

		// Si tu veux garder une fonction utilitaire :
		// RespondError(c, http.StatusBadRequest, fieldErrors)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": fieldErrors,
		})
		return
	}

	p := Create(input)
	// Réponse succès JSON
	c.JSON(http.StatusCreated, gin.H{
		"message": "Pokémon créé",
		"pokemon": p,
	})
}
