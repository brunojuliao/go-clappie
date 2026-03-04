package parties

import (
	"path/filepath"
	"strings"

	"github.com/brunojuliao/go-clappie/internal/filestore"
	"github.com/brunojuliao/go-clappie/internal/platform"
)

func identitiesPath(root string) string {
	return filepath.Join(partiesPath(root), "identities")
}

// ListIdentities returns all identity templates.
func ListIdentities(root string) ([]Identity, error) {
	dir := identitiesPath(root)
	entries, err := filestore.List(dir)
	if err != nil {
		return nil, err
	}

	var identities []Identity
	for _, entry := range entries {
		body, blocks, err := filestore.ReadAndParse(entry.Path)
		if err != nil {
			continue
		}

		memories := filestore.GetMetaField(blocks, "meta", "memories")
		var memList []string
		if memories != "" {
			memList = strings.Split(memories, "|")
		}

		identities = append(identities, Identity{
			Name:     entry.Name,
			Template: body,
			Memories: memList,
		})
	}
	return identities, nil
}

// GetIdentity returns a specific identity by name.
func GetIdentity(root, name string) (*Identity, error) {
	path := filestore.FilePath(identitiesPath(root), name)
	body, blocks, err := filestore.ReadAndParse(path)
	if err != nil {
		return nil, err
	}

	memories := filestore.GetMetaField(blocks, "meta", "memories")
	var memList []string
	if memories != "" {
		memList = strings.Split(memories, "|")
	}

	return &Identity{
		Name:     name,
		Template: body,
		Memories: memList,
	}, nil
}

// SaveIdentity saves an identity template.
func SaveIdentity(root string, id Identity) error {
	dir := identitiesPath(root)
	if err := platform.EnsureDir(dir); err != nil {
		return err
	}

	path := filestore.FilePath(dir, id.Name)
	blocks := []filestore.MetaBlock{{
		Tag: "meta",
		Fields: map[string]string{
			"memories": strings.Join(id.Memories, "|"),
		},
	}}
	return filestore.WriteWithMeta(path, id.Template, blocks)
}

// AddMemory adds a memory to an identity.
func AddMemory(root, name, memory string) error {
	id, err := GetIdentity(root, name)
	if err != nil {
		return err
	}
	id.Memories = append(id.Memories, memory)
	return SaveIdentity(root, *id)
}
