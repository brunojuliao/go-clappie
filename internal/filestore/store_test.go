package filestore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStoreRoundtrip(t *testing.T) {
	dir := t.TempDir()

	// Write
	path := filepath.Join(dir, "test.txt")
	err := WriteFile(path, "Hello, world!")
	if err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Read
	content, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if content != "Hello, world!" {
		t.Errorf("content = %q", content)
	}

	// Exists
	if !Exists(path) {
		t.Error("file should exist")
	}

	// Delete
	if err := DeleteFile(path); err != nil {
		t.Fatalf("DeleteFile: %v", err)
	}
	if Exists(path) {
		t.Error("file should not exist after delete")
	}
}

func TestList(t *testing.T) {
	dir := t.TempDir()

	// Create some files
	WriteFile(filepath.Join(dir, "a.txt"), "A")
	WriteFile(filepath.Join(dir, "b.txt"), "B")
	WriteFile(filepath.Join(dir, "c.json"), "C") // not .txt
	os.Mkdir(filepath.Join(dir, "subdir"), 0755)

	entries, err := List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("entries = %d, want 2", len(entries))
	}
}

func TestListNonexistent(t *testing.T) {
	entries, err := List("/nonexistent/path")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries for nonexistent dir")
	}
}

func TestCount(t *testing.T) {
	dir := t.TempDir()
	WriteFile(filepath.Join(dir, "a.txt"), "A")
	WriteFile(filepath.Join(dir, "b.txt"), "B")

	if count := Count(dir); count != 2 {
		t.Errorf("Count = %d, want 2", count)
	}
}

func TestReadAndParse(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	content := "Body text\n\n---\n[meta]\nkey: value\n"
	WriteFile(path, content)

	body, blocks, err := ReadAndParse(path)
	if err != nil {
		t.Fatalf("ReadAndParse: %v", err)
	}
	if body != "Body text" {
		t.Errorf("body = %q", body)
	}
	if len(blocks) != 1 {
		t.Fatalf("blocks = %d", len(blocks))
	}
	if blocks[0].Fields["key"] != "value" {
		t.Errorf("field = %q", blocks[0].Fields["key"])
	}
}

func TestWriteWithMeta(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	blocks := []MetaBlock{{
		Tag:    "meta",
		Fields: map[string]string{"status": "ok"},
	}}
	if err := WriteWithMeta(path, "Hello", blocks); err != nil {
		t.Fatalf("WriteWithMeta: %v", err)
	}

	body, parsedBlocks, err := ReadAndParse(path)
	if err != nil {
		t.Fatalf("ReadAndParse: %v", err)
	}
	if body != "Hello" {
		t.Errorf("body = %q", body)
	}
	if GetMetaField(parsedBlocks, "meta", "status") != "ok" {
		t.Error("expected status=ok")
	}
}
