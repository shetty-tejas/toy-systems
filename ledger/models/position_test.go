package models

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func makeSegmentDirs(t *testing.T, root string, starts ...uint) []string {
	t.Helper()

	names := make([]string, 0, len(starts))
	for _, start := range starts {
		name := fmt.Sprintf("%016d", start)
		if err := os.MkdirAll(filepath.Join(root, name), 0o766); err != nil {
			t.Fatal(err)
		}
		names = append(names, name)
	}

	return names
}

func TestPositionForWriteStartsAtZeroWhenNoSegmentsExist(t *testing.T) {
	position := PositionForWrite(t.TempDir())

	expected := &Position{Cursor: 0, Start: 0}
	if !reflect.DeepEqual(position, expected) {
		t.Fatalf("expected %+v, got %+v", expected, position)
	}
}

func TestPositionForWriteUsesLatestSegmentDirectory(t *testing.T) {
	root := t.TempDir()
	makeSegmentDirs(t, root, 0, 3, 9)

	ignoredFile := filepath.Join(root, "not-a-segment.txt")
	if err := os.WriteFile(ignoredFile, []byte("ignored"), 0o666); err != nil {
		t.Fatal(err)
	}

	position := PositionForWrite(root)
	if position.Start != 9 || position.Cursor != 0 {
		t.Fatalf("expected start 9 cursor 0, got %+v", position)
	}
}

func TestPositionForReadAcrossSegmentBoundaries(t *testing.T) {
	root := t.TempDir()
	makeSegmentDirs(t, root, 0, 3, 7, 12)

	tests := []struct {
		position uint
		expected *Position
	}{
		{position: 0, expected: &Position{Cursor: 0, Start: 0}},
		{position: 2, expected: &Position{Cursor: 2, Start: 0}},
		{position: 3, expected: &Position{Cursor: 0, Start: 3}},
		{position: 6, expected: &Position{Cursor: 3, Start: 3}},
		{position: 7, expected: &Position{Cursor: 0, Start: 7}},
		{position: 11, expected: &Position{Cursor: 4, Start: 7}},
		{position: 12, expected: &Position{Cursor: 0, Start: 12}},
		{position: 20, expected: &Position{Cursor: 8, Start: 12}},
	}

	for _, tt := range tests {
		position := PositionForRead(root, tt.position)
		if !reflect.DeepEqual(position, tt.expected) {
			t.Fatalf("position %d: expected %+v, got %+v", tt.position, tt.expected, position)
		}
	}
}

func TestPositionForReadPanicsWhenSegmentsAreNotInitialised(t *testing.T) {
	expectPanic(t, "segments are not yet initialised", func() {
		PositionForRead(t.TempDir(), 0)
	})
}

func TestBinarySearchFloorEdgeCases(t *testing.T) {
	directories := []string{
		"0000000000000005",
		"0000000000000010",
		"0000000000000020",
	}

	tests := []struct {
		name     string
		search   uint
		expected uint
	}{
		{name: "below first", search: 0, expected: 0},
		{name: "exact first", search: 5, expected: 0},
		{name: "between first and second", search: 9, expected: 0},
		{name: "exact middle", search: 10, expected: 1},
		{name: "between second and third", search: 19, expected: 1},
		{name: "exact last", search: 20, expected: 2},
		{name: "above last", search: 999, expected: 2},
	}

	for _, tt := range tests {
		index := binarySearchFloor(directories, tt.search, 0, uint(len(directories)-1))
		if index != tt.expected {
			t.Fatalf("%s: expected index %d, got %d", tt.name, tt.expected, index)
		}
	}

	if index := binarySearchFloor(directories[:1], 42, 0, 0); index != 0 {
		t.Fatalf("single element search: expected index 0, got %d", index)
	}
}
