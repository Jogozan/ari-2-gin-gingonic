package pokemon

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var list []Pokemon
	if err := json.Unmarshal(data, &list); err != nil {
		return err
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

func LevelUp(id int) (Pokemon, error) {
	mu.Lock()
	defer mu.Unlock()

	for i, p := range pokemons {
		// exemple de règle: niveau max 100, +1 niveau = +10 HP
		if p.ID == id {
			if p.Stats.HP >= 400 { // règle arbitraire de “max”
				return Pokemon{}, ErrMaxLevel
			}
			p.Stats.HP += 10
			p.BaseExperience += 5
			pokemons[i] = p
			return p, nil
		}
	}
	return Pokemon{}, errors.New("not found")
}
