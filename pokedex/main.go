package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"pokedex/pokemon"
)

func main() {
	// Charger les données depuis pokemons.json
	if err := pokemon.LoadFromFile("pokemons.json"); err != nil {
		log.Fatalf("Impossible de charger pokemons.json: %v", err)
	}

	router := gin.Default()

	// Templates HTML
	router.LoadHTMLGlob("templates/*.tmpl")

	// Route de test
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// API JSON
	api := router.Group("/api/v1")
	pokemon.RegisterAPIRoutes(api)

	// Routes HTML
	pokemon.RegisterHTMLRoutes(router)

	// Lancement du serveur
	router.Run(":8080")
}
