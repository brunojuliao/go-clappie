package chores

import (
	"fmt"
	"time"

	"github.com/brunojuliao/go-clappie/internal/filestore"
	"github.com/brunojuliao/go-clappie/internal/platform"
)

// Create creates a new chore.
func Create(root string, chore Chore) error {
	dir := platform.ChoresHumansDir(root)
	if err := platform.EnsureDir(dir); err != nil {
		return err
	}

	path := filestore.FilePath(dir, chore.Name)

	blocks := []filestore.MetaBlock{{
		Tag: "chore-meta",
		Fields: map[string]string{
			"title":   chore.Title,
			"summary": chore.Summary,
			"icon":    chore.Icon,
			"context": chore.Context,
			"status":  StatusPending,
			"created": time.Now().Format("2006-01-02 15:04"),
		},
	}}

	return filestore.WriteWithMeta(path, chore.Body, blocks)
}

// GetPending returns all pending chores.
func GetPending(root string) ([]Chore, error) {
	return listByStatus(root, StatusPending)
}

// GetAll returns all chores.
func GetAll(root string) ([]Chore, error) {
	return listByStatus(root, "")
}

func listByStatus(root, status string) ([]Chore, error) {
	dir := platform.ChoresHumansDir(root)
	entries, err := filestore.List(dir)
	if err != nil {
		return nil, err
	}

	var chores []Chore
	for _, entry := range entries {
		chore, err := readChore(entry)
		if err != nil {
			continue
		}
		if status == "" || chore.Status == status {
			chores = append(chores, chore)
		}
	}
	return chores, nil
}

func readChore(entry filestore.Entry) (Chore, error) {
	body, blocks, err := filestore.ReadAndParse(entry.Path)
	if err != nil {
		return Chore{}, err
	}

	meta := filestore.GetMeta(blocks, "chore-meta")

	chore := Chore{
		Name: entry.Name,
		Path: entry.Path,
		Body: body,
	}

	if meta != nil {
		chore.Title = meta.Fields["title"]
		chore.Summary = meta.Fields["summary"]
		chore.Icon = meta.Fields["icon"]
		chore.Context = meta.Fields["context"]
		chore.Status = meta.Fields["status"]
		if created, ok := meta.Fields["created"]; ok {
			if t, err := time.Parse("2006-01-02 15:04", created); err == nil {
				chore.Created = t
			}
		}
	}

	if chore.Title == "" {
		chore.Title = entry.Name
	}

	return chore, nil
}

// Approve approves a chore.
func Approve(root, name string) error {
	return setStatus(root, name, StatusApproved)
}

// Complete marks a chore as completed and moves it to logs.
func Complete(root, name string) error {
	if err := setStatus(root, name, StatusCompleted); err != nil {
		return err
	}

	// Move to logs
	srcDir := platform.ChoresHumansDir(root)
	srcPath := filestore.FilePath(srcDir, name)

	logDir := filestore.ChoreLogPath(root)
	if err := platform.EnsureDir(logDir); err != nil {
		return err
	}

	logName := filestore.TimestampedName(name)
	dstPath := filestore.FilePath(logDir, logName)

	content, err := filestore.ReadFile(srcPath)
	if err != nil {
		return err
	}
	if err := filestore.WriteFile(dstPath, content); err != nil {
		return err
	}
	return filestore.DeleteFile(srcPath)
}

// Reject rejects a chore.
func Reject(root, name string) error {
	return setStatus(root, name, StatusRejected)
}

// Shelve shelves a chore for later.
func Shelve(root, name string) error {
	return setStatus(root, name, StatusShelved)
}

func setStatus(root, name, status string) error {
	dir := platform.ChoresHumansDir(root)
	path := filestore.FilePath(dir, name)

	body, blocks, err := filestore.ReadAndParse(path)
	if err != nil {
		return fmt.Errorf("read chore %s: %w", name, err)
	}

	filestore.SetMetaField(&blocks, "chore-meta", "status", status)
	return filestore.WriteWithMeta(path, body, blocks)
}
