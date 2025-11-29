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
type CreatePokemonInput struct {
	Name           string   `json:"name" binding:"required"`
	BaseExperience int      `json:"baseExperience" binding:"required,min=1"`
	Weight         int      `json:"weight" binding:"required,min=1"`
	Height         int      `json:"height" binding:"required,min=1"`
	Types          []string `json:"types" binding:"required,min=1,dive,required"`
	Stats          Stats    `json:"stats" binding:"required"`
	Sprites        Sprites  `json:"sprites" binding:"required"`
}
