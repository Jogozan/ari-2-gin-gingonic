package pokemon

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
	all := GetAll()
	c.JSON(http.StatusOK, gin.H{
		"data": all,
	})
}

func getPokemonByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	p, err := GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pokemon non trouvé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": p})
}

func createPokemon(c *gin.Context) {
	var input CreatePokemonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	p := Create(input)
	c.JSON(http.StatusCreated, gin.H{"data": p})
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
