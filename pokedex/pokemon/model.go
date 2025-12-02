package pokemon

type Stats struct {
	Speed   int `json:"speed"`
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	HP      int `json:"hp"`
}

type Sprites struct {
	FrontDefault string `json:"front_default"`
	BackDefault  string `json:"back_default"`
}

type Pokemon struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	BaseExperience int      `json:"baseExperience"`
	Weight         int      `json:"weight"`
	Height         int      `json:"height"`
	Types          []string `json:"types"`
	Stats          Stats    `json:"stats"`
	Sprites        Sprites  `json:"sprites"`
}

// Pour la création/édition via API JSON
type CreatePokemonInput struct {
	Name           string   `json:"name"`
	BaseExperience int      `json:"baseExperience"`
	Weight         int      `json:"weight"`
	Height         int      `json:"height"`
	Types          []string `json:"types"`
	Stats          Stats    `json:"stats"`
	Sprites        Sprites  `json:"sprites"`
}

// DTO de réponse pour l’API (avec Power)
type PokemonResponse struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	BaseExperience int      `json:"baseExperience"`
	Weight         int      `json:"weight"`
	Height         int      `json:"height"`
	Types          []string `json:"types"`
	Stats          Stats    `json:"stats"`
	Sprites        Sprites  `json:"sprites"`
	Power          int      `json:"power"`
}

func toResponse(p Pokemon) PokemonResponse {
	power := p.Stats.HP * p.Stats.Attack
	return PokemonResponse{
		ID:             p.ID,
		Name:           p.Name,
		BaseExperience: p.BaseExperience,
		Weight:         p.Weight,
		Height:         p.Height,
		Types:          p.Types,
		Stats:          p.Stats,
		Sprites:        p.Sprites,
		Power:          power,
	}
}
