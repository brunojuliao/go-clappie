package sidekicks

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/brunojuliao/go-clappie/internal/filestore"
	"github.com/brunojuliao/go-clappie/internal/platform"
)

const sidekicksDir = "recall/sidekicks"

func sidekicksPath(root string) string {
	return filepath.Join(root, sidekicksDir)
}

// ListActive returns all active sidekick sessions.
func ListActive(root string) ([]SidekickInfo, error) {
	dir := sidekicksPath(root)
	entries, err := filestore.List(dir)
	if err != nil {
		return nil, err
	}

	var sidekicksList []SidekickInfo
	for _, entry := range entries {
		body, blocks, err := filestore.ReadAndParse(entry.Path)
		if err != nil {
			continue
		}
		_ = body

		status := filestore.GetMetaField(blocks, "meta", "status")
		if status != "active" {
			continue
		}

		sk := SidekickInfo{
			ID:     entry.Name,
			Prompt: filestore.GetMetaField(blocks, "meta", "prompt"),
			Model:  filestore.GetMetaField(blocks, "meta", "model"),
			Squad:  filestore.GetMetaField(blocks, "meta", "squad"),
			Status: status,
			PaneID: filestore.GetMetaField(blocks, "meta", "pane_id"),
		}
		if created := filestore.GetMetaField(blocks, "meta", "created"); created != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", created); err == nil {
				sk.CreatedAt = t
			}
		}
		sidekicksList = append(sidekicksList, sk)
	}
	return sidekicksList, nil
}

// Get returns a specific sidekick by ID.
func Get(root, id string) (*SidekickInfo, error) {
	path := filestore.FilePath(sidekicksPath(root), id)
	_, blocks, err := filestore.ReadAndParse(path)
	if err != nil {
		return nil, err
	}

	return &SidekickInfo{
		ID:     id,
		Prompt: filestore.GetMetaField(blocks, "meta", "prompt"),
		Model:  filestore.GetMetaField(blocks, "meta", "model"),
		Squad:  filestore.GetMetaField(blocks, "meta", "squad"),
		Status: filestore.GetMetaField(blocks, "meta", "status"),
		PaneID: filestore.GetMetaField(blocks, "meta", "pane_id"),
	}, nil
}

// createSession creates a new sidekick session file.
func createSession(root string, config SpawnConfig, paneID string) (string, error) {
	dir := sidekicksPath(root)
	if err := platform.EnsureDir(dir); err != nil {
		return "", err
	}

	id := fmt.Sprintf("sk-%d", time.Now().UnixMilli())
	path := filestore.FilePath(dir, id)

	blocks := []filestore.MetaBlock{{
		Tag: "meta",
		Fields: map[string]string{
			"prompt":  config.Prompt,
			"model":   config.Model,
			"squad":   config.Squad,
			"status":  "active",
			"pane_id": paneID,
			"created": time.Now().Format("2006-01-02 15:04:05"),
		},
	}}

	if err := filestore.WriteWithMeta(path, config.Prompt, blocks); err != nil {
		return "", err
	}
	return id, nil
}

// updateStatus updates the status of a sidekick session.
func updateStatus(root, id, status string) error {
	path := filestore.FilePath(sidekicksPath(root), id)
	body, blocks, err := filestore.ReadAndParse(path)
	if err != nil {
		return err
	}
	filestore.SetMetaField(&blocks, "meta", "status", status)
	return filestore.WriteWithMeta(path, body, blocks)
}
