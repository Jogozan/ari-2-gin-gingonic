package pokemon

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func RegisterHTMLRoutes(r *gin.Engine) {
	r.GET("/pokemons", listPokemonsHTML)
	r.GET("/pokemons/new", newPokemonFormHTML)
	r.POST("/pokemons", createPokemonHTML)
	r.GET("/pokemons/:id", pokemonDetailHTML)
}

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
			minLevel, err := strconv.Atoi(minLevelStr)
			if err == nil && p.Stats.HP < minLevel {
				continue
			}
		}
		filtered = append(filtered, p)
	}

	// on peut pré-calculer le power pour l’affichage
	type PokemonView struct {
		Pokemon
		Power int
	}
	var views []PokemonView
	for _, p := range filtered {
		power := p.Stats.HP * p.Stats.Attack
		views = append(views, PokemonView{Pokemon: p, Power: power})
	}

	c.HTML(http.StatusOK, "pokemons_index.tmpl", gin.H{
		"title":    "Pokédex",
		"pokemons": views,
		"type":     typeFilter,
		"minLevel": minLevelStr,
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

// Formulaire HTML
func newPokemonFormHTML(c *gin.Context) {
	// TODO
	c.HTML(http.StatusOK, "pokemons_form.tmpl", gin.H{
		"title":  "Nouveau Pokémon",
		"errors": []string{},
		"input":  CreatePokemonInput{},
	})
}

func createPokemonHTML(c *gin.Context) {
	var input CreatePokemonInput
	if err := c.ShouldBindJSON(&input); err != nil {
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
