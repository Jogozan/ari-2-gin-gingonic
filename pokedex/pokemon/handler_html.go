package pokemon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterHTMLRoutes(r *gin.Engine) {
	r.GET("/pokemons", listPokemonsHTML)
	r.GET("/pokemons/:id", pokemonDetailHTML)
	r.POST("/pokemons/:id/level-up", pokemonLevelUpHTML)
	// A small HTML-only "release" action — implemented as POST for form friendliness
	r.POST("/pokemons/:id/release", pokemonReleaseHTML)
	r.GET("/pokemons/stats", pokemonsStatsHTML)
}

//	func listPokemonsHTML(c *gin.Context) {
//		all := GetAll()
//		c.HTML(http.StatusOK, "pokemons_index.tmpl", gin.H{
//			"title":    "Pokédex",
//			"pokemons": all,
//		})
//	}
func listPokemonsHTML(c *gin.Context) {
	all := GetAll()
	// Support simple sorting by level or power (default keeps original order)
	sortBy := c.Query("sort")
	// Make a copy before sorting so the in-memory global slice isn't re-ordered
	copyList := make([]Pokemon, len(all))
	copy(copyList, all)

	if sortBy == "level" {
		sort.Slice(copyList, func(i, j int) bool { return copyList[i].Level > copyList[j].Level })
	} else if sortBy == "power" {
		power := func(p Pokemon) int {
			return p.Stats.Attack + p.Stats.Defense + p.Stats.Speed + p.Stats.HP + p.BaseExperience
		}
		sort.Slice(copyList, func(i, j int) bool { return power(copyList[i]) > power(copyList[j]) })
	}
	data, _ := json.MarshalIndent(copyList, "", "  ")
	fmt.Println(string(data))
	// precompute a simple power metric for each pokemon (used by template)
	powerMap := map[int]int{}
	for _, p := range copyList {
		powerMap[p.ID] = p.Stats.Attack + p.Stats.Defense + p.Stats.Speed + p.Stats.HP + p.BaseExperience
	}

	c.HTML(http.StatusOK, "pokemons_index.tmpl", gin.H{
		"title":    "Pokédex",
		"pokemons": copyList,
		// current request URI (including query) is passed so forms can return here
		"current_url": c.Request.URL.RequestURI(),
		"message":     c.Query("msg"),
		// small palette used by the HTML templates for badges
		"type_colors": typeColorMap(),
		"power_map":   powerMap,
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
