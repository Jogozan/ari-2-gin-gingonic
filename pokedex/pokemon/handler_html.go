package pokemon

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterHTMLRoutes(r *gin.Engine) {
	r.GET("/pokemons", listPokemonsHTML)
	r.GET("/pokemons/:id", pokemonDetailHTML)

	// routes de création
	r.GET("/pokemons/new", newPokemonFormHTML)
	r.POST("/pokemons", createPokemonHTML)
}

func newPokemonFormHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "pokemons_form.tmpl", gin.H{
		"title":  "Nouveau Pokémon",
		"errors": []string{},
		"input":  CreatePokemonInput{}, // valeurs vides
	})
}

func createPokemonHTML(c *gin.Context) {
	var input CreatePokemonInput

	// Bind depuis un formulaire HTML (x-www-form-urlencoded)
	if err := c.ShouldBind(&input); err != nil {
		// Erreur générique de parsing
		c.HTML(http.StatusBadRequest, "pokemons_form.tmpl", gin.H{
			"title":  "Nouveau Pokémon",
			"errors": []string{"Données invalides."},
			"input":  input,
		})
		return
	}

	// Validation avancée déjà définie par les tags `binding`
	// Si tu veux des messages plus précis :
	// import "github.com/go-playground/validator/v10"
	// et traite err.(validator.ValidationErrors) comme pour l’API JSON.

	p := Create(input)

	// Redirection vers la page détail du Pokémon
	c.Redirect(http.StatusSeeOther, "/pokemons/"+strconv.Itoa(p.ID))
}

func listPokemonsHTML(c *gin.Context) {
	all := GetAll()
	c.HTML(http.StatusOK, "pokemons_index.tmpl", gin.H{
		"title":    "Pokédex",
		"pokemons": all,
	})
}

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
