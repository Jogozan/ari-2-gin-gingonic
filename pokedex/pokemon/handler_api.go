package pokemon

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func RegisterAPIRoutes(rg *gin.RouterGroup) {
	rg.GET("/pokemons", getPokemons)
	rg.GET("/pokemons/:id", getPokemonByID)
	rg.POST("/pokemons", createPokemon)
	rg.DELETE("/pokemons/:id", deletePokemon)
}

func getPokemons(c *gin.Context) {
	all := GetAll()
	respondOK(c, all)
}

func getPokemonByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
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
