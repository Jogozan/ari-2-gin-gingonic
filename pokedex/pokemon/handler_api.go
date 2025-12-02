package pokemon

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(rg *gin.RouterGroup) {
	rg.GET("/pokemons", getPokemons)
	rg.GET("/pokemons/:id", getPokemonByID)
	rg.POST("/pokemons", createPokemon)
	rg.DELETE("/pokemons/:id", deletePokemon)
	rg.POST("/pokemons/:id/level-up", levelUpPokemon)
}

func levelUpPokemon(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Niveau maximum atteint"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Pokemon non trouvé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toResponse(p)})
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

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func getPokemonByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	p, err := GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pokemon non trouvé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toResponse(*p)})
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
