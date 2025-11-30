package pokemon

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterHTMLRoutes(r *gin.Engine) {
	r.GET("/pokemons", listPokemonsHTML)

	r.GET("/pokemons/new", newPokemonFormHTML)
	r.POST("/pokemons", createPokemonUnified)

	r.GET("/pokemons/:id", pokemonDetailHTML)
}

// Formulaire HTML
func newPokemonFormHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "pokemons_form.tmpl", gin.H{
		"title":  "Nouveau Pokémon",
		"errors": []string{},
		"input":  CreatePokemonInput{},
	})
}

// Handler unifié JSON + HTML
func createPokemonUnified(c *gin.Context) {
	var input CreatePokemonInput
	var isJSON bool

	// Détecter si Content-Type = application/json
	if c.ContentType() == "application/json" {
		isJSON = true
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Sinon, binder formulaire HTML
		if err := c.ShouldBind(&input); err != nil {
			c.HTML(http.StatusBadRequest, "pokemons_form.tmpl", gin.H{
				"title":  "Nouveau Pokémon",
				"errors": []string{"Données invalides."},
				"input":  input,
			})
			return
		}

		// Binder manuellement les stats
		input.Stats.HP = parseInt(c.PostForm("stats.hp"))
		input.Stats.Attack = parseInt(c.PostForm("stats.attack"))
		input.Stats.Defense = parseInt(c.PostForm("stats.defense"))
		input.Stats.Speed = parseInt(c.PostForm("stats.speed"))

		// Binder manuellement les sprites
		input.Sprites.FrontDefault = c.PostForm("sprites.front_default")
		input.Sprites.BackDefault = c.PostForm("sprites.back_default")

		// Binder manuellement les types
		rawTypes := c.PostForm("types")
		input.Types = strings.Split(rawTypes, ",")
		for i := range input.Types {
			input.Types[i] = strings.TrimSpace(input.Types[i])
		}

		if len(input.Types) == 0 || len(input.Types) > 2 {
			c.HTML(http.StatusBadRequest, "pokemons_form.tmpl", gin.H{
				"title":  "Nouveau Pokémon",
				"errors": []string{"Veuillez renseigner 1 ou 2 types valides."},
				"input":  input,
			})
			return
		}
	}

	// Création du Pokémon
	p := Create(input)

	// Réponse selon le format demandé
	if isJSON {
		c.JSON(http.StatusCreated, p)
	} else {
		c.Redirect(http.StatusSeeOther, "/pokemons/"+strconv.Itoa(p.ID))
	}
}

// Liste des Pokémons
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

// Détail d’un Pokémon
func pokemonDetailHTML(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "ID invalide")
		return
	}

	p, err := GetByID(id)
	if err != nil {
		c.String(http.StatusNotFound, "Pokémon non trouvé")
		return
	}

	c.HTML(http.StatusOK, "pokemons_detail.tmpl", gin.H{
		"title":   p.Name,
		"pokemon": p,
	})
}

// parseInt helper
func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
