package pokemon

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(rg *gin.RouterGroup) {
	rg.GET("/pokemons", getPokemons)
	rg.GET("/pokemons/:id", getPokemonByID)
	rg.POST("/pokemons", createPokemon)
	rg.DELETE("/pokemons/:id", deletePokemon)
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
