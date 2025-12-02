package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"pokedex/pokemon"
)

func main() {
	// Charger les données depuis pokemons.json
	if err := pokemon.LoadFromFile("pokemons.json"); err != nil {
		log.Fatalf("Impossible de charger pokemons.json: %v", err)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		pokemon.RegisterCustomValidations(v)
	}

	router := gin.Default()

	// Global middlewares (applied to all routes):
	// - EnrichedLogger: copies X-Trainer header into context and logs it
	// - FatigueMiddleware: if a request includes header X-Server-Fatigue=true we add delay
	router.Use(pokemon.EnrichedLogger())
	router.Use(pokemon.FatigueMiddleware(500 * time.Millisecond))

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
