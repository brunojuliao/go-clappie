package filestore

import (
	"testing"
)

func TestParseFile(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantBody   string
		wantBlocks int
	}{
		{
			name:       "body only",
			input:      "Hello world",
			wantBody:   "Hello world",
			wantBlocks: 0,
		},
		{
			name: "body with meta",
			input: `This is the body

---
[chore-meta]
title: My Chore
status: pending
icon: 📧`,
			wantBody:   "This is the body",
			wantBlocks: 1,
		},
		{
			name: "multiple meta blocks",
			input: `Body content

---
[meta]
key1: value1

[chore-meta]
title: Test
status: approved`,
			wantBody:   "Body content",
			wantBlocks: 2,
		},
		{
			name: "empty body with meta",
			input: `---
[heartbeat-meta]
interval: 5m
last_run: 2025-01-01 12:00:00`,
			wantBody:   "",
			wantBlocks: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, blocks := ParseFile(tt.input)
			if body != tt.wantBody {
				t.Errorf("body = %q, want %q", body, tt.wantBody)
			}
			if len(blocks) != tt.wantBlocks {
				t.Errorf("blocks = %d, want %d", len(blocks), tt.wantBlocks)
			}
		})
	}
}

func TestParseFileMetaFields(t *testing.T) {
	input := `Body text

---
[chore-meta]
title: Test Chore
summary: A test
icon: 📧
status: pending
created: 2025-03-03 14:30`

	body, blocks := ParseFile(input)
	if body != "Body text" {
		t.Errorf("body = %q", body)
	}
	if len(blocks) != 1 {
		t.Fatalf("blocks = %d, want 1", len(blocks))
	}

	block := blocks[0]
	if block.Tag != "chore-meta" {
		t.Errorf("tag = %q, want chore-meta", block.Tag)
	}
	if block.Fields["title"] != "Test Chore" {
		t.Errorf("title = %q", block.Fields["title"])
	}
	if block.Fields["status"] != "pending" {
		t.Errorf("status = %q", block.Fields["status"])
	}
	if block.Fields["icon"] != "📧" {
		t.Errorf("icon = %q", block.Fields["icon"])
	}
}

func TestFormatFile(t *testing.T) {
	blocks := []MetaBlock{{
		Tag: "meta",
		Fields: map[string]string{
			"key": "value",
		},
	}}

	output := FormatFile("Hello", blocks)
	body, parsed := ParseFile(output)

	if body != "Hello" {
		t.Errorf("roundtrip body = %q", body)
	}
	if len(parsed) != 1 {
		t.Fatalf("roundtrip blocks = %d", len(parsed))
	}
	if parsed[0].Fields["key"] != "value" {
		t.Errorf("roundtrip key = %q", parsed[0].Fields["key"])
	}
}

func TestGetMeta(t *testing.T) {
	blocks := []MetaBlock{
		{Tag: "meta", Fields: map[string]string{"a": "1"}},
		{Tag: "chore-meta", Fields: map[string]string{"b": "2"}},
	}

	m := GetMeta(blocks, "chore-meta")
	if m == nil {
		t.Fatal("expected chore-meta block")
	}
	if m.Fields["b"] != "2" {
		t.Errorf("field b = %q", m.Fields["b"])
	}

	m = GetMeta(blocks, "nonexistent")
	if m != nil {
		t.Error("expected nil for nonexistent block")
	}
}

func TestSetMetaField(t *testing.T) {
	var blocks []MetaBlock

	// Set field on new block
	SetMetaField(&blocks, "meta", "key", "value")
	if len(blocks) != 1 {
		t.Fatalf("blocks = %d", len(blocks))
	}
	if blocks[0].Fields["key"] != "value" {
		t.Errorf("field = %q", blocks[0].Fields["key"])
	}

	// Update existing field
	SetMetaField(&blocks, "meta", "key", "updated")
	if len(blocks) != 1 {
		t.Fatalf("blocks = %d after update", len(blocks))
	}
	if blocks[0].Fields["key"] != "updated" {
		t.Errorf("updated field = %q", blocks[0].Fields["key"])
	}
}
