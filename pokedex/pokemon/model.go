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
	ID   int    `json:"id"`
	Name string `json:"name"`
	// Level is an in-memory convenience field used by the TP to demonstrate "level-up".
	// It's not present in the original JSON but will be initialized to 1 when loading.
	Level          int      `json:"level"`
	BaseExperience int      `json:"baseExperience"`
	Weight         int      `json:"weight"`
	Height         int      `json:"height"`
	Types          []string `json:"types"`
	Stats          Stats    `json:"stats"`
	Sprites        Sprites  `json:"sprites"`
}

// Pour la création/édition via API JSON
// give example for Name and Types, and ask to student to do the other attibut validation (ex: Weight est obligatoire et comprit entre 1 et 1000)
type CreatePokemonInput struct {
	Name           string   `json:"name" form:"name" binding:"required,max=50"`
	Types          []string `json:"types" binding:"required,min=1,max=2,dive,pokemon_type"`
	BaseExperience int      `json:"baseExperience" form:"baseExperience" binding:"required,min=1,max=1000"`
	Weight         int      `json:"weight" form:"weight" binding:"required,min=1,max=10000"`
	Height         int      `json:"height" form:"height" binding:"required,min=1,max=100"`
	Stats          Stats    `json:"stats" form:"stats" binding:"required"`
	Sprites        Sprites  `json:"sprites" form:"sprites" binding:"required"`
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
	// toResponse maps an internal Pokemon model to the API DTO PokemonResponse.
	// It computes presentation-only fields (like Power) and returns a value
	// suitable for JSON encoding to clients.
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
