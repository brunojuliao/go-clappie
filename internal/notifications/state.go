package notifications

import (
	"time"

	"github.com/brunojuliao/go-clappie/internal/filestore"
	"github.com/brunojuliao/go-clappie/internal/platform"
)

// AddDirty adds a raw notification to the dirty queue.
func AddDirty(root string, item DirtyItem) error {
	dir := platform.NotificationsDirtyDir(root)
	if err := platform.EnsureDir(dir); err != nil {
		return err
	}

	name := item.Name
	if name == "" {
		name = filestore.TimestampedName("notification")
	}
	path := filestore.FilePath(dir, name)

	blocks := []filestore.MetaBlock{{
		Tag: "meta",
		Fields: map[string]string{
			"source":    item.Source,
			"source_id": item.SourceID,
			"created":   time.Now().Format("2006-01-02 15:04:05"),
		},
	}}

	return filestore.WriteWithMeta(path, item.Body, blocks)
}

// ListDirty returns all dirty notifications.
func ListDirty(root string) ([]DirtyItem, error) {
	dir := platform.NotificationsDirtyDir(root)
	entries, err := filestore.List(dir)
	if err != nil {
		return nil, err
	}

	var items []DirtyItem
	for _, entry := range entries {
		body, blocks, err := filestore.ReadAndParse(entry.Path)
		if err != nil {
			continue
		}
		item := DirtyItem{
			Name:     entry.Name,
			Path:     entry.Path,
			Body:     body,
			Source:   filestore.GetMetaField(blocks, "meta", "source"),
			SourceID: filestore.GetMetaField(blocks, "meta", "source_id"),
		}
		if created := filestore.GetMetaField(blocks, "meta", "created"); created != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", created); err == nil {
				item.Created = t
			}
		}
		items = append(items, item)
	}
	return items, nil
}

// AddClean adds a curated notification to the clean queue.
func AddClean(root string, item CleanItem) error {
	dir := platform.NotificationsCleanDir(root)
	if err := platform.EnsureDir(dir); err != nil {
		return err
	}

	name := item.Name
	if name == "" {
		name = filestore.TimestampedName("notification")
	}
	path := filestore.FilePath(dir, name)

	blocks := []filestore.MetaBlock{{
		Tag: "meta",
		Fields: map[string]string{
			"source_id": item.SourceID,
			"context":   item.Context,
			"created":   time.Now().Format("2006-01-02 15:04:05"),
		},
	}}

	return filestore.WriteWithMeta(path, item.Body, blocks)
}

// ListClean returns all clean notifications.
func ListClean(root string) ([]CleanItem, error) {
	dir := platform.NotificationsCleanDir(root)
	entries, err := filestore.List(dir)
	if err != nil {
		return nil, err
	}

	var items []CleanItem
	for _, entry := range entries {
		body, blocks, err := filestore.ReadAndParse(entry.Path)
		if err != nil {
			continue
		}
		item := CleanItem{
			Name:     entry.Name,
			Path:     entry.Path,
			Body:     body,
			SourceID: filestore.GetMetaField(blocks, "meta", "source_id"),
			Context:  filestore.GetMetaField(blocks, "meta", "context"),
		}
		if created := filestore.GetMetaField(blocks, "meta", "created"); created != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", created); err == nil {
				item.Created = t
			}
		}
		items = append(items, item)
	}
	return items, nil
}

// CountClean returns the number of clean notifications.
func CountClean(root string) int {
	return filestore.Count(platform.NotificationsCleanDir(root))
}

// DismissBySourceID removes notifications matching a source_id.
func DismissBySourceID(root, sourceID string) error {
	items, err := ListClean(root)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.SourceID == sourceID {
			filestore.DeleteFile(item.Path)
		}
	}
	return nil
}

// ClearDirty removes all dirty notifications.
func ClearDirty(root string) error {
	items, err := ListDirty(root)
	if err != nil {
		return err
	}
	for _, item := range items {
		filestore.DeleteFile(item.Path)
	}
	return nil
}
