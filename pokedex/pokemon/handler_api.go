package pokemon

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func RegisterAPIRoutes(rg *gin.RouterGroup) {
	rg.GET("/pokemons", getPokemons)
	rg.GET("/pokemons/:id", getPokemonByID)
	rg.POST("/pokemons", createPokemon)
	rg.DELETE("/pokemons/:id", deletePokemon)

	// Admin subgroup — demonstrates group middleware (simple auth + optional rate limit)
	admin := rg.Group("/admin")
	// tiny hard-coded secret for the exercise; in real apps use env/config
	admin.Use(SimpleAuth("admin-secret"))

	// POST /api/v1/admin/pokemons/:id/level-up
	// Protect this route with a small rate limiter and optional server fatigue
	admin.POST("/pokemons/:id/level-up", RateLimitMiddleware(5, 10*time.Second), levelUpPokemon)
}

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

	respondOK(c, all)
}

func getPokemonByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, []string{"ID invalide."})
		return
	}

	p, err := GetByID(id)
	if err != nil {
		respondError(c, http.StatusNotFound, []string{"Pokemon non trouvé."})
		return
	}

	respondOK(c, p)
}

func createPokemon(c *gin.Context) {
	var input CreatePokemonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		var messages []string

		// Si c’est une erreur de validation, on la détaille
		if verrs, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range verrs {
				field := fe.Field()
				tag := fe.Tag()

				switch field {
				case "Name":
					if tag == "required" {
						messages = append(messages, "Le nom est obligatoire.")
					} else if tag == "max" {
						messages = append(messages, "Le nom ne doit pas dépasser 50 caractères.")
					}
				case "Type":
					if tag == "required" {
						messages = append(messages, "Le type est obligatoire.")
					} else if tag == "oneof" {
						messages = append(messages, "Le type doit être un type Pokémon valide.")
					}
				case "BaseExperience":
					messages = append(messages, "L'expérience de base doit être entre 1 et 1000.")
				case "Weight":
					messages = append(messages, "Le poids doit être positif et raisonnable.")
				case "Height":
					messages = append(messages, "La taille doit être positive et raisonnable.")
				default:
					messages = append(messages, "Champ "+field+" invalide.")
				}
			}
		} else {
			// Erreur de parsing JSON ou autre
			messages = append(messages, "Données invalides.")
		}

		respondError(c, http.StatusBadRequest, messages)
		return
	}

	p := Create(input)
	respondCreated(c, p)
}

func deletePokemon(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	if err := Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pokemon non trouvé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pokemon supprimé"})
}

// levelUpPokemon allows an authenticated admin to "level up" a Pokemon.
// This handler demonstrates: auth middleware (group), route middleware (rate limit),
// context propagation (trainer and target_pokemon) and state change.
func levelUpPokemon(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	// For pedagogical purposes we read number of levels to add from query or default 1
	addLevels := 1
	if lvlStr := c.Query("levels"); lvlStr != "" {
		if lv, err := strconv.Atoi(lvlStr); err == nil && lv > 0 {
			addLevels = lv
		}
	}

	// retrieve trainer information propagated by earlier middleware
	trainerVal, _ := c.Get("trainer")
	trainer, _ := trainerVal.(string)

	// Try to perform level-up
	p, err := LevelUp(id, addLevels)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pokemon non trouvé"})
		return
	}

	// Demonstrate reading target_pokemon attached earlier by logging middleware
	if tp, ok := c.Get("target_pokemon"); ok {
		if tpok, is := tp.(*Pokemon); is {
			// Use the attached target pokemon just to show how middleware propagation works.
			// (p may be different copy, we just illustrate available context values)
			_ = tpok
		}
	}

	// Return the updated pokemon and a small message that includes the trainer if present
	msg := "Pokemon level-up effectué"
	if trainer != "" {
		msg = "Pokemon level-up effectué par " + trainer
	}

	c.JSON(http.StatusOK, gin.H{"message": msg, "data": p})
}

// Standard response helpers

type APIResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error []string    `json:"error,omitempty"`
}

func respondOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{Data: data})
}

func respondCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{Data: data})
}

func respondError(c *gin.Context, status int, messages []string) {
	c.JSON(status, APIResponse{Error: messages})
}
