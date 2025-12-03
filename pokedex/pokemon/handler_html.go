package pokemon

import (
	"net/http"
	"strconv"
	"strings"
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
}

// listPokemonsHTML handles GET /pokemons and renders the HTML index page.
// It supports the same `type` and `minLevel` filters as the API and computes
// a simple Power value for display purposes.
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
		"title":    p.Name,
		"pokemon":  p,
		"pokemons": GetAll(), // pour navigation si besoin
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
				var msg string
				switch fe.Field() {
				case "Name":
					msg = "Le nom est obligatoire et max 50 caractères."
				case "Types":
					msg = "Merci de fournir 2 types valides maximum."
				case "Types[0]":
					msg = "Erreur sur le premier type."
				case "Types[1]":
					msg = "Erreur sur le deuxième type."
				case "BaseExperience":
					msg = "L'expérience de base est obligatoire et doit être comprise entre 1 et 1000."
				case "Weight":
					msg = "Le poids est obligatoire et doit être compris entre 1 et 10000."
				case "Height":
					msg = "La taille est obligatoire et doit être comprise entre 1 et 100."
				case "Stats":
					msg = "Les statistiques sont obligatoires et doivent être valides."
				case "Sprites":
					msg = "Les sprites sont obligatoires et doivent être valides."
				default:
					msg = fe.Field() + " invalide."
				}

				field := fe.Field() // ex: Name, Types, BaseExperience...
				fieldErrors[field] = append(fieldErrors[field], msg)
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
