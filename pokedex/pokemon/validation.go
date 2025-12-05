package pokemon

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

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

func validationMessage(field string) string {
	switch field {
	case "Name":
		return "Le nom est obligatoire et max 50 caractères."
	case "Types":
		return "Merci de fournir 2 types valides maximum."
	case "Types[0]":
		return "Erreur sur le premier type."
	case "Types[1]":
		return "Erreur sur le deuxième type."
	case "BaseExperience":
		return "L'expérience de base est obligatoire et doit être comprise entre 1 et 1000."
	case "Weight":
		return "Le poids est obligatoire et doit être compris entre 1 et 10000."
	case "Height":
		return "La taille est obligatoire et doit être comprise entre 1 et 100."
	case "Stats":
		return "Les statistiques sont obligatoires et doivent être valides."
	case "Sprites":
		return "Les sprites sont obligatoires et doivent être valides."
	default:
		return field + " invalide."
	}
}
