package pokemon

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidations registers project-specific validation rules
// into the provided validator instance. Current example registers the
// `pokemon_type` rule used to validate allowed pokemon types.
func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidation("pokemon_type", func(fl validator.FieldLevel) bool {
		allowed := map[string]bool{
			"normal": true, "fire": true, "water": true,
			"grass": true, "electric": true, "ice": true,
			"fighting": true, "poison": true, "ground": true,
			"flying": true, "psychic": true, "bug": true,
			"rock": true, "ghost": true, "dragon": true,
			"dark": true, "steel": true, "fairy": true,
		}

		if s, ok := fl.Field().Interface().(string); ok {
			s = strings.ToLower(strings.TrimSpace(s))
			return allowed[s]
		}
		return false
	})
}
