package pokemon

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var (
	pokemons []Pokemon
	mu       sync.RWMutex
	loaded   bool
)

func LoadFromFile(path string) error {
	mu.Lock()
	defer mu.Unlock()

	if loaded {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var list []Pokemon
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}

	// Ensure every Pokemon has a default Level value so level-up can be demonstrated
	for i := range list {
		if list[i].Level == 0 {
			list[i].Level = 1
		}
	}

	pokemons = list
	loaded = true
	return nil
}

func GetAll() []Pokemon {
	mu.RLock()
	defer mu.RUnlock()
	return pokemons
}

func GetByID(id int) (*Pokemon, error) {
	mu.RLock()
	defer mu.RUnlock()
	for _, p := range pokemons {
		if p.ID == id {
			copy := p
			return &copy, nil
		}
	}
	return nil, errors.New("not found")
}

func Create(input CreatePokemonInput) Pokemon {
	mu.Lock()
	defer mu.Unlock()

	maxID := 0
	for _, p := range pokemons {
		if p.ID > maxID {
			maxID = p.ID
		}
	}
	newID := maxID + 1

	p := Pokemon{
		ID:             newID,
		Name:           input.Name,
		Level:          1,
		BaseExperience: input.BaseExperience,
		Weight:         input.Weight,
		Height:         input.Height,
		Types:          input.Types,
		Stats:          input.Stats,
		Sprites:        input.Sprites,
	}

	pokemons = append(pokemons, p)
	return p
}

func Delete(id int) error {
	mu.Lock()
	defer mu.Unlock()

	for i, p := range pokemons {
		if p.ID == id {
			pokemons = append(pokemons[:i], pokemons[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

var ErrMaxLevel = errors.New("max level reached")

func LevelUp(id int, levels int) (Pokemon, error) {
	mu.Lock()
	defer mu.Unlock()

	for i := range pokemons {
		if pokemons[i].ID == id {
			if pokemons[i].Level >= 100 {
				return Pokemon{}, ErrMaxLevel
			}
			pokemons[i].Level += levels
			pokemons[i].BaseExperience += 10 * levels
			pokemons[i].Stats.HP += 3 * levels
			pokemons[i].Stats.Speed += 2 * levels
			pokemons[i].Stats.Attack += 1 * levels
			pokemons[i].Stats.Defense += 2 * levels
			copy := pokemons[i]
			return copy, nil
		}
	}

	return Pokemon{}, errors.New("not found")
}
