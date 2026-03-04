package filestore

import (
	"fmt"
	"strings"
)

// MetaBlock represents a parsed metadata block from a text file.
// Files use the format:
//
//	[body content]
//	---
//	[meta-tag]
//	key: value
//	key: value
type MetaBlock struct {
	Tag    string            // e.g. "meta", "chore-meta", "heartbeat-meta"
	Fields map[string]string // key-value pairs
}

// ParseFile parses a text file into body and metadata blocks.
func ParseFile(content string) (body string, blocks []MetaBlock) {
	// Split on --- separator
	parts := strings.SplitN(content, "\n---\n", 2)
	if len(parts) == 1 {
		// Try with just ---
		parts = strings.SplitN(content, "---\n", 2)
		if len(parts) == 1 {
			return strings.TrimSpace(content), nil
		}
		body = ""
	} else {
		body = strings.TrimSpace(parts[0])
	}

	metaSection := parts[len(parts)-1]
	blocks = parseMetaBlocks(metaSection)
	return body, blocks
}

func parseMetaBlocks(section string) []MetaBlock {
	var blocks []MetaBlock
	var current *MetaBlock

	for _, line := range strings.Split(section, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check for block header like [meta] or [chore-meta]
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			tag := line[1 : len(line)-1]
			blocks = append(blocks, MetaBlock{
				Tag:    tag,
				Fields: make(map[string]string),
			})
			current = &blocks[len(blocks)-1]
			continue
		}

		// Parse key: value pairs
		if current != nil {
			if idx := strings.Index(line, ": "); idx > 0 {
				key := strings.TrimSpace(line[:idx])
				value := strings.TrimSpace(line[idx+2:])
				current.Fields[key] = value
			} else if idx := strings.Index(line, ":"); idx > 0 {
				key := strings.TrimSpace(line[:idx])
				value := strings.TrimSpace(line[idx+1:])
				current.Fields[key] = value
			}
		}
	}

	return blocks
}

// FormatFile creates a file string from body and metadata blocks.
func FormatFile(body string, blocks []MetaBlock) string {
	var sb strings.Builder

	if body != "" {
		sb.WriteString(body)
	}

	if len(blocks) > 0 {
		if body != "" {
			sb.WriteString("\n\n")
		}
		sb.WriteString("---\n")
		for i, block := range blocks {
			if i > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(fmt.Sprintf("[%s]\n", block.Tag))
			for key, value := range block.Fields {
				sb.WriteString(fmt.Sprintf("%s: %s\n", key, value))
			}
		}
	}

	return sb.String()
}

// GetMeta finds a metadata block by tag name.
func GetMeta(blocks []MetaBlock, tag string) *MetaBlock {
	for i := range blocks {
		if blocks[i].Tag == tag {
			return &blocks[i]
		}
	}
	return nil
}

// GetMetaField gets a field from a specific meta block.
func GetMetaField(blocks []MetaBlock, tag, field string) string {
	block := GetMeta(blocks, tag)
	if block == nil {
		return ""
	}
	return block.Fields[field]
}

// SetMetaField sets a field in a specific meta block, creating the block if needed.
func SetMetaField(blocks *[]MetaBlock, tag, field, value string) {
	for i := range *blocks {
		if (*blocks)[i].Tag == tag {
			(*blocks)[i].Fields[field] = value
			return
		}
	}
	// Block doesn't exist, create it
	*blocks = append(*blocks, MetaBlock{
		Tag:    tag,
		Fields: map[string]string{field: value},
	})
}
