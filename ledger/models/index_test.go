package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func openTestFile(t *testing.T, dir, name string) *os.File {
	t.Helper()

	file, err := os.OpenFile(filepath.Join(dir, name), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o666)
	if err != nil {
		t.Fatal(err)
	}

	return file
}

func expectPanic(t *testing.T, contains string, fn func()) {
	t.Helper()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic containing %q", contains)
		}

		if contains == "" {
			return
		}

		message := fmt.Sprint(r)
		if !strings.Contains(message, contains) {
			t.Fatalf("expected panic containing %q, got %q", contains, message)
		}
	}()

	fn()
}

func TestIndexAppendReadAndSize(t *testing.T) {
	index := NewIndex(openTestFile(t, t.TempDir(), "index.data"))
	defer index.Close()

	if size := index.GetSize(); size != 0 {
		t.Fatalf("expected empty index size 0, got %d", size)
	}

	offsets := []uint{3, 9, 42}
	for position, offset := range offsets {
		gotPosition := index.Append(offset)
		if gotPosition != uint(position) {
			t.Fatalf("expected append position %d, got %d", position, gotPosition)
		}

		expectedSize := uint(position+1) * entrySizeInBytes
		if size := index.GetSize(); size != expectedSize {
			t.Fatalf("expected size %d after append, got %d", expectedSize, size)
		}
	}

	for position, expectedOffset := range offsets {
		if offset := index.ReadAt(uint(position)); offset != expectedOffset {
			t.Fatalf("expected offset %d at position %d, got %d", expectedOffset, position, offset)
		}
	}
}

func TestIndexReadAtPanicsAtEOF(t *testing.T) {
	index := NewIndex(openTestFile(t, t.TempDir(), "index.data"))
	defer index.Close()

	index.Append(17)

	expectPanic(t, "offset is greater than the file size", func() {
		index.ReadAt(1)
	})
}
