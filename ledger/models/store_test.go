package models

import (
	"bytes"
	"testing"
)

func TestStoreAppendReadAndSize(t *testing.T) {
	store := NewStore(openTestFile(t, t.TempDir(), "store.data"))
	defer store.Close()

	if size := store.GetSize(); size != 0 {
		t.Fatalf("expected empty store size 0, got %d", size)
	}

	messages := [][]byte{
		[]byte("alpha"),
		{},
		[]byte("beta"),
	}

	expectedOffset := uint(0)
	for _, message := range messages {
		offset := store.Append(message)
		if offset != expectedOffset {
			t.Fatalf("expected append offset %d, got %d", expectedOffset, offset)
		}

		expectedOffset += headerSizeInBytes + uint(len(message))
		if size := store.GetSize(); size != expectedOffset {
			t.Fatalf("expected size %d after append, got %d", expectedOffset, size)
		}
	}

	offset := uint(0)
	for index, expectedMessage := range messages {
		message := store.ReadAt(offset)
		if !bytes.Equal(message, expectedMessage) {
			t.Fatalf("expected message %q at offset %d, got %q", expectedMessage, offset, message)
		}

		offset += headerSizeInBytes + uint(len(expectedMessage))
		if offset > store.GetSize() && index != len(messages)-1 {
			t.Fatalf("offset advanced past store size: %d > %d", offset, store.GetSize())
		}
	}
}

func TestStoreReadAtPanicsAtEOF(t *testing.T) {
	store := NewStore(openTestFile(t, t.TempDir(), "store.data"))
	defer store.Close()

	store.Append([]byte("alpha"))

	expectPanic(t, "offset is greater than the file size", func() {
		store.ReadAt(store.GetSize())
	})
}
