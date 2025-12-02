package pokemon

import "github.com/go-playground/validator/v10"

func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidation("pokemon_type", func(fl validator.FieldLevel) bool {
		allowed := map[string]bool{
			"Normal": true, "Fire": true, "Water": true,
			"Grass": true, "Electric": true, "Ice": true,
			"Fighting": true, "Poison": true, "Ground": true,
			"Flying": true, "Psychic": true, "Bug": true,
			"Rock": true, "Ghost": true, "Dragon": true,
			"Dark": true, "Steel": true,
			"Fairy": true,
		}
		if s, ok := fl.Field().Interface().(string); ok {
			return allowed[s]
		}
		return false
	})
}
