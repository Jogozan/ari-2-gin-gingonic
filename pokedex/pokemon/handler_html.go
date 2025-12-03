package pokemon

import (
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// RegisterHTMLRoutes registers routes that render HTML pages for the pokedex.
// These handlers return rendered templates (index, detail, and form) instead of JSON.
func RegisterHTMLRoutes(r *gin.Engine) {
	r.GET("/pokemons", listPokemonsHTML)
	r.GET("/pokemons/new", newPokemonFormHTML)
	r.POST("/pokemons", createPokemonHTML)
	r.GET("/pokemons/:id", pokemonDetailHTML)
	r.POST("/pokemons/:id/level-up", pokemonLevelUpHTML)
	r.POST("/pokemons/:id/release", pokemonReleaseHTML)
	r.GET("/pokemons/stats", pokemonsStatsHTML)
}

// ---- Helpers métier ----

func matchesType(p Pokemon, typeFilter string) bool {
	if typeFilter == "" {
		return true
	}
	for _, t := range p.Types {
		if strings.EqualFold(t, typeFilter) {
			return true
		}
	}
	return false
}

func matchesMinLevel(p Pokemon, minLevelStr string) bool {
	if minLevelStr == "" {
		return true
	}
	minLevel, err := strconv.Atoi(minLevelStr)
	if err != nil {
		return true // si minLevel est invalide, on ignore le filtre
	}
	return p.Stats.HP >= minLevel
}

func filterPokemons(list []Pokemon, typeFilter, minLevelStr string) []Pokemon {
	var filtered []Pokemon
	for _, p := range list {
		if !matchesType(p, typeFilter) {
			continue
		}
		if !matchesMinLevel(p, minLevelStr) {
			continue
		}
		filtered = append(filtered, p)
	}
	return filtered
}

func pokemonPower(p Pokemon) int {
	return p.Stats.Attack + p.Stats.Defense + p.Stats.Speed + p.Stats.HP + p.BaseExperience
}

func sortPokemons(list []Pokemon, sortBy string) []Pokemon {
	copyList := make([]Pokemon, len(list))
	copy(copyList, list)

	switch sortBy {
	case "level":
		sort.Slice(copyList, func(i, j int) bool {
			return copyList[i].Level > copyList[j].Level
		})
	case "power":
		sort.Slice(copyList, func(i, j int) bool {
			return pokemonPower(copyList[i]) > pokemonPower(copyList[j])
		})
		// default: pas de tri, on garde l'ordre
	}

	return copyList
}

func buildPowerMap(list []Pokemon) map[int]int {
	powerMap := make(map[int]int, len(list))
	for _, p := range list {
		powerMap[p.ID] = pokemonPower(p)
	}
	return powerMap
}

func toResponses(list []Pokemon) []PokemonResponse {
	resp := make([]PokemonResponse, 0, len(list))
	for _, p := range list {
		resp = append(resp, toResponse(p))
	}
	return resp
}

// ---- Handler HTML ----

func listPokemonsHTML(c *gin.Context) {
	typeFilter := c.Query("type")
	minLevelStr := c.Query("minLevel")
	sortBy := c.Query("sort")

	all := GetAll()

	// 1) Filtrage
	filtered := filterPokemons(all, typeFilter, minLevelStr)

	// 2) Tri
	sorted := sortPokemons(filtered, sortBy)

	// 3) Power pré-calculé + DTO
	powerMap := buildPowerMap(sorted)
	resp := toResponses(sorted)

	c.HTML(http.StatusOK, "pokemons_index.tmpl", gin.H{
		"title":       "Pokédex",
		"pokemons":    resp, // ou `sorted` si le template attend le modèle brut
		"type":        typeFilter,
		"minLevel":    minLevelStr,
		"sort":        sortBy,
		"current_url": c.Request.URL.RequestURI(),
		"message":     c.Query("msg"),
		"type_colors": typeColorMap(),
		"power_map":   powerMap,
	})
}

// pokemonDetailHTML handles GET /pokemons/:id and renders the detail page
// for a single pokemon. It returns HTTP errors as plain text when the id is invalid.
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
		"title":       p.Name,
		"pokemon":     p,
		"pokemons":    GetAll(), // pour navigation si besoin
		"current_url": c.Request.URL.RequestURI(),
		"type_colors": typeColorMap(),
	})
}

// pokemonLevelUpHTML handles simple form POSTs originating from HTML views
// It performs a LevelUp(+1) and redirects back to the index with a message
func pokemonLevelUpHTML(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "ID invalide")
		return
	}

	// By default +1 level; optionally caller can provide 'levels' form value
	levels := 1
	if lvl := c.PostForm("levels"); lvl != "" {
		if v, e := strconv.Atoi(lvl); e == nil && v > 0 {
			levels = v
		}
	}

	p, err := LevelUp(id, levels)
	if err != nil {
		c.String(http.StatusNotFound, "Pokemon non trouvé")
		return
	}

	// choose redirect target in order: explicit form 'redirect', Referer header, default index
	target := c.PostForm("redirect")
	if target == "" {
		target = c.Request.Referer()
	}
	if target == "" {
		target = "/pokemons"
	}

	// Append message to target URL preserving its existing query parameters
	msg := p.Name + " a gagné " + strconv.Itoa(levels) + " niveau(s) !"
	u, err := url.Parse(target)
	if err != nil {
		// fallback
		u = &url.URL{Path: "/pokemons"}
	}
	q := u.Query()
	q.Set("msg", msg)
	u.RawQuery = q.Encode()
	c.Redirect(http.StatusSeeOther, u.String())
}

// pokemonReleaseHTML handles releasing (deleting) a pokemon from HTML
// and performs a redirect back to the index page.
func pokemonReleaseHTML(c *gin.Context) {
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

	if err := Delete(id); err != nil {
		c.String(http.StatusInternalServerError, "Impossible de supprimer")
		return
	}

	target := c.PostForm("redirect")
	if target == "" {
		target = c.Request.Referer()
	}
	if target == "" {
		target = "/pokemons"
	}

	msg := p.Name + " a été relâché 😢"
	u, err := url.Parse(target)
	if err != nil {
		u = &url.URL{Path: "/pokemons"}
	}
	q := u.Query()
	q.Set("msg", msg)
	u.RawQuery = q.Encode()
	c.Redirect(http.StatusSeeOther, u.String())
}

// pokemonsStatsHTML builds a small aggregation of pokemon counts per type
func pokemonsStatsHTML(c *gin.Context) {
	all := GetAll()
	counts := map[string]int{}
	for _, p := range all {
		for _, t := range p.Types {
			counts[t]++
		}
	}

	c.HTML(http.StatusOK, "pokemons_stats.tmpl", gin.H{
		"title":       "Statistiques du dresseur",
		"pokemons":    all,
		"type_counts": counts,
		"type_colors": typeColorMap(),
	})
}

// typeColorMap returns a map of type->css color used for badges.
// Keep the palette small and readable; it's purely presentational for the TP.
func typeColorMap() map[string]string {
	return map[string]string{
		"grass":    "#78C850",
		"fire":     "#F08030",
		"water":    "#6890F0",
		"electric": "#F8D030",
		"ice":      "#98D8D8",
		"psychic":  "#F85888",
		"ghost":    "#705898",
		"dark":     "#705848",
		"rock":     "#B8A038",
		"steel":    "#B8B8D0",
		"ground":   "#E0C068",
		"flying":   "#A890F0",
		"bug":      "#A8B820",
		"poison":   "#A040A0",
		"normal":   "#A8A878",
	}
}

// Formulaire HTML
// newPokemonFormHTML renders the empty form used to create a new pokemon.
// The template expects `title`, `errors` and `input` values in its context.
func newPokemonFormHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "pokemons_form.tmpl", gin.H{
		"title":  "Nouveau Pokémon",
		"errors": []string{},
		"input":  CreatePokemonInput{},
	})
}

// createPokemonHTML accepts JSON payloads and returns HTML-friendly error maps.
// It demonstrates converting validator errors into a `map[string][]string`
// structure so templates / clients can display per-field messages.
func createPokemonHTML(c *gin.Context) {
	var input CreatePokemonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		// Map[string][]string : nom du champ -> liste de messages
		fieldErrors := map[string][]string{}

		if verrs, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range verrs {
				var msg string
				switch fe.Field() {
				case "Name":
					msg = "Le nom est obligatoire et max 50 caractères."
				case "Types":
					msg = "Merci de fournir 2 types valides maximum."
				case "Types[0]":
					msg = "Erreur sur le premier type."
				case "Types[1]":
					msg = "Erreur sur le deuxième type."
				case "BaseExperience":
					msg = "L'expérience de base est obligatoire et doit être comprise entre 1 et 1000."
				case "Weight":
					msg = "Le poids est obligatoire et doit être compris entre 1 et 10000."
				case "Height":
					msg = "La taille est obligatoire et doit être comprise entre 1 et 100."
				case "Stats":
					msg = "Les statistiques sont obligatoires et doivent être valides."
				case "Sprites":
					msg = "Les sprites sont obligatoires et doivent être valides."
				default:
					msg = fe.Field() + " invalide."
				}

				field := fe.Field() // ex: Name, Types, BaseExperience...
				fieldErrors[field] = append(fieldErrors[field], msg)
			}
		}

		// Si tu veux garder une fonction utilitaire :
		// RespondError(c, http.StatusBadRequest, fieldErrors)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": fieldErrors,
		})
		return
	}

	p := Create(input)
	// Réponse succès JSON
	c.JSON(http.StatusCreated, gin.H{
		"message": "Pokémon créé",
		"pokemon": p,
	})
}
