package parties

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/brunojuliao/go-clappie/internal/filestore"
	"github.com/brunojuliao/go-clappie/internal/platform"
)

const partiesDir = "recall/parties"

func partiesPath(root string) string {
	return filepath.Join(root, partiesDir)
}

func gamesPath(root string) string {
	return filepath.Join(partiesPath(root), "games")
}

func simsPath(root string) string {
	return filepath.Join(partiesPath(root), "sims")
}

// ListGames returns all available game definitions.
func ListGames(root string) ([]GameInfo, error) {
	dir := gamesPath(root)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var games []GameInfo
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".txt") {
			continue
		}
		body, blocks, err := filestore.ReadAndParse(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		desc := filestore.GetMetaField(blocks, "meta", "description")
		name := strings.TrimSuffix(e.Name(), ".txt")
		games = append(games, GameInfo{
			Name:        name,
			Description: desc,
			Path:        filepath.Join(dir, e.Name()),
		})
		_ = body
	}
	return games, nil
}

// Rules returns the rules text for a game.
func Rules(root, gameName string) (string, error) {
	path := filepath.Join(gamesPath(root), gameName+".txt")
	body, _, err := filestore.ReadAndParse(path)
	if err != nil {
		return "", fmt.Errorf("game %q not found: %w", gameName, err)
	}
	return body, nil
}

// ParseGame parses a full game definition file.
func ParseGame(path string) (*Game, error) {
	body, blocks, err := filestore.ReadAndParse(path)
	if err != nil {
		return nil, err
	}

	game := &Game{
		Name:           filepath.Base(strings.TrimSuffix(path, ".txt")),
		Path:           path,
		Rules:          body,
		SuggestedState: make(map[string]string),
	}

	meta := filestore.GetMeta(blocks, "meta")
	if meta != nil {
		game.Description = meta.Fields["description"]
	}

	// Parse player cards from meta
	cardsMeta := filestore.GetMeta(blocks, "cards")
	if cardsMeta != nil {
		for name, desc := range cardsMeta.Fields {
			game.PlayerCards = append(game.PlayerCards, PlayerCard{
				Name:        name,
				Description: desc,
			})
		}
	}

	// Parse suggested state
	stateMeta := filestore.GetMeta(blocks, "state")
	if stateMeta != nil {
		game.SuggestedState = stateMeta.Fields
	}

	return game, nil
}

// Init initializes a new simulation from a game definition.
func Init(root, gameName string) (string, error) {
	game, err := findGame(root, gameName)
	if err != nil {
		return "", err
	}

	simDir := simsPath(root)
	if err := platform.EnsureDir(simDir); err != nil {
		return "", err
	}

	simID := fmt.Sprintf("%s-%d", gameName, filestore.Count(simDir)+1)
	simPath := filepath.Join(simDir, simID+".txt")

	blocks := []filestore.MetaBlock{
		{
			Tag: "meta",
			Fields: map[string]string{
				"game":   gameName,
				"status": "initialized",
			},
		},
	}

	// Copy suggested state
	if len(game.SuggestedState) > 0 {
		blocks = append(blocks, filestore.MetaBlock{
			Tag:    "state",
			Fields: game.SuggestedState,
		})
	}

	if err := filestore.WriteWithMeta(simPath, game.Rules, blocks); err != nil {
		return "", err
	}

	return simID, nil
}

// Launch starts a simulation by spawning AI agents.
func Launch(root, simID string) error {
	// TODO: spawn sidekick agents for each participant
	return fmt.Errorf("launch not yet implemented")
}

// Show displays the current simulation state.
func Show(root, simID string) error {
	path := filepath.Join(simsPath(root), simID+".txt")
	content, err := filestore.ReadFile(path)
	if err != nil {
		return fmt.Errorf("simulation %q not found: %w", simID, err)
	}
	fmt.Println(content)
	return nil
}

// End ends a simulation.
func End(root, simID string) error {
	path := filepath.Join(simsPath(root), simID+".txt")
	body, blocks, err := filestore.ReadAndParse(path)
	if err != nil {
		return err
	}
	filestore.SetMetaField(&blocks, "meta", "status", "ended")
	return filestore.WriteWithMeta(path, body, blocks)
}

// Set sets a value in the simulation ledger.
func Set(root, simID, key, value string) error {
	path := filepath.Join(simsPath(root), simID+".txt")
	body, blocks, err := filestore.ReadAndParse(path)
	if err != nil {
		return err
	}
	filestore.SetMetaField(&blocks, "state", key, value)
	return filestore.WriteWithMeta(path, body, blocks)
}

// Get gets a value from the simulation ledger.
func Get(root, simID, key string) (string, error) {
	path := filepath.Join(simsPath(root), simID+".txt")
	_, blocks, err := filestore.ReadAndParse(path)
	if err != nil {
		return "", err
	}
	value := filestore.GetMetaField(blocks, "state", key)
	return value, nil
}

func findGame(root, gameName string) (*Game, error) {
	path := filepath.Join(gamesPath(root), gameName+".txt")
	return ParseGame(path)
}
