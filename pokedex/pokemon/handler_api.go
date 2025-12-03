package pokemon

// A fournir entierement

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// RegisterAPIRoutes registers the HTTP routes for the JSON API.
// It mounts handlers for listing, reading, creating and deleting pokemons
// and exposes an admin subgroup which demonstrates middleware usage.
func RegisterAPIRoutes(rg *gin.RouterGroup) {
	rg.GET("/pokemons", getPokemons)
	rg.GET("/pokemons/:id", getPokemonByID)
	rg.POST("/pokemons", createPokemon)
	rg.DELETE("/pokemons/:id", deletePokemon)
	rg.POST("/pokemons/:id/level-up", levelUpPokemon)

	// Admin subgroup — demonstrates group middleware (simple auth + optional rate limit)
	admin := rg.Group("/admin")
	// tiny hard-coded secret for the exercise; in real apps use env/config
	admin.Use(SimpleAuth("admin-secret"))

	// POST /api/v1/admin/pokemons/:id/level-up
	// Protect this route with a small rate limiter and optional server fatigue
	admin.POST("/pokemons/:id/level-up", RateLimitMiddleware(5, 10*time.Second), levelUpPokemon)
}

// getPokemons handles GET /pokemons for the API.
// It supports optional query filters for `type` and `minLevel` and
// returns a list of PokemonResponse DTOs (including computed Power).
func getPokemons(c *gin.Context) {
	typeFilter := c.Query("type")
	minLevelStr := c.Query("minLevel")

	all := GetAll()
	var filtered []Pokemon

	for _, p := range all {
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
			// on utilise ici HP comme proxy de “niveau”
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
	RespondOK(c, resp)
}

// getPokemonByID handles GET /pokemons/:id and returns a single pokemon.
// It validates the id parameter, looks up the pokemon and returns 404 when missing.
func getPokemonByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, []string{"ID invalide."})
		return
	}

	p, err := GetByID(id)
	if err != nil {
		RespondError(c, http.StatusNotFound, []string{"Pokemon non trouvé"})
		return
	}
	RespondOK(c, toResponse(*p))
}

// createPokemon handles POST /pokemons for the JSON API.
// It binds and validates the incoming JSON payload, maps validation errors
// to friendly messages and returns the created resource on success.
func createPokemon(c *gin.Context) {
	var input CreatePokemonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		// replace
		/*
			RespondError(c, http.StatusBadRequest, []string{"Données invalides"})
			return
		*/
		// by
		var messages []string
		if verrs, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range verrs {
				switch fe.Field() {
				case "Name":
					messages = append(messages, "Le nom est obligatoire et max 50 caractères.")
				case "Types":
					messages = append(messages, "Types invalides (ex : Fire, Water, Grass).")
				case "BaseExperience":
					messages = append(messages, "L'expérience de base est obligatoire et doit être comprise entre 1 et 1000.")
				case "Weight":
					messages = append(messages, "Le poids est obligatoire et doit être compris entre 1 et 10000.")
				case "Height":
					messages = append(messages, "La taille est obligatoire et doit être comprise entre 1 et 100.")
				case "Stats":
					messages = append(messages, "Les statistiques sont obligatoires et doivent être valides.")
				case "Sprites":
					messages = append(messages, "Les sprites sont obligatoires et doivent être valides.")
				default:
					messages = append(messages, fe.Field()+" invalide.")
				}
			}
		}
		RespondError(c, 400, messages)
		return
	}

	p := Create(input)
	RespondCreated(c, p)
}

// deletePokemon handles DELETE /pokemons/:id.
// It validates the id parameter and deletes the pokemon if it exists.
func deletePokemon(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, []string{"ID invalide"})
		return
	}

	if err := Delete(id); err != nil {
		RespondError(c, http.StatusNotFound, []string{"Pokemon non trouvé"})
		return
	}

	RespondOK(c, "Pokemon supprimé")
}

// levelUpPokemon allows an authenticated admin to "level up" a Pokemon.
// This handler demonstrates: auth middleware (group), route middleware (rate limit),
// context propagation (trainer and target_pokemon) and state change.
// levelUpPokemon handles POST /admin/pokemons/:id/level-up.
// This admin-only endpoint increments a pokemon's level by a given amount
// and demonstrates usage of group + route middleware (auth + rate limit)
// and context propagation from earlier middlewares.
func levelUpPokemon(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, []string{"ID invalide"})
		return
	}

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

	RespondOK(c, toResponse(p))
}
