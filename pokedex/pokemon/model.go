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

// Validation avancée pour la création / édition via API JSON.
type CreatePokemonInput struct {
	Name           string   `json:"name" binding:"required,max=50"`
	Types          []string `json:"types" binding:"required"` // plus de dive ici
	BaseExperience int      `json:"baseExperience" binding:"required,min=1,max=1000"`
	Weight         int      `json:"weight" binding:"required,min=1,max=10000"`
	Height         int      `json:"height" binding:"required,min=1,max=100"`
	Stats          Stats    `json:"stats" binding:"required"` // plus de dive
	Sprites        Sprites  `json:"sprites" binding:"required"`
}
