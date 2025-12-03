package pokemon

import (
	"os"
	"testing"
)

// TestLoadAndLevelUp verifies that creating a pokemon and leveling it up
// updates Level and BaseExperience as expected. It runs in isolation and
// cleans up any test state at the end.
func TestLoadAndLevelUp(t *testing.T) {
	// Create a pokemon locally (avoids depending on loaded JSON data)
	input := CreatePokemonInput{
		Name:           "test-mon",
		BaseExperience: 10,
		Weight:         1,
		Height:         1,
		Types:          []string{"normal"},
		Stats:          Stats{Speed: 1, Attack: 1, Defense: 1, HP: 1},
		Sprites:        Sprites{FrontDefault: "", BackDefault: ""},
	}

	p0 := Create(input)

	id := p0.ID
	origLevel := p0.Level
	origExp := p0.BaseExperience

	p, err := LevelUp(id, 2)
	if err != nil {
		t.Fatalf("level up failed: %v", err)
	}

	if p.Level != origLevel+2 {
		t.Fatalf("expected level %d, got %d", origLevel+2, p.Level)
	}

	if p.BaseExperience != origExp+20 {
		t.Fatalf("expected baseExperience %d, got %d", origExp+20, p.BaseExperience)
	}

	// cleanup any state by re-reading json (simple approach)
	_ = os.RemoveAll("./pokemons.json")
}
