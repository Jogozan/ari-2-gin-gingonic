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

// LoadFromFile loads the pokemons list from a JSON file located at `path`.
// It is safe to call multiple times; the file is read only once and stored
// in an in-memory slice protected by a mutex.
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

// GetAll returns the in-memory slice of pokemons.
// The returned slice should be treated as read-only by callers.
func GetAll() []Pokemon {
	mu.RLock()
	defer mu.RUnlock()
	return pokemons
}

// GetByID returns a copy of the Pokemon with the given id or an error when not found.
// A copy is returned to avoid callers accidentally mutating the internal slice state.
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

// Create creates a new Pokemon from the provided input, assigns a new ID
// and appends it to the in-memory store. It returns the created Pokemon.
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

// Delete removes the Pokemon with the given id from the in-memory store.
// Returns an error if the pokemon does not exist.
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

// LevelUp increases the specified Pokemon's level by `levels` and applies
// a simple side-effect on BaseExperience (10 XP per level). It returns a
// copy of the updated Pokemon or an error if the id is not found.
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
