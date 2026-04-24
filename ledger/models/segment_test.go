package models

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSegmentForWriteAppendReadAndClose(t *testing.T) {
	root := t.TempDir()
	segment := SegmentForWrite(root, 0)

	if segment.Start != 0 {
		t.Fatalf("expected segment start 0, got %d", segment.Start)
	}
	if size := segment.GetSize(); size != 0 {
		t.Fatalf("expected empty segment size 0, got %d", size)
	}

	first := segment.Append("alpha")
	second := segment.Append("beta")

	expectedFirst := &Entry{Segment: 0, Position: 0, Offset: 0, Message: []byte("alpha")}
	if !reflect.DeepEqual(first, expectedFirst) {
		t.Fatalf("expected first entry %+v, got %+v", expectedFirst, first)
	}

	expectedSecondOffset := headerSizeInBytes + uint(len("alpha"))
	expectedSecond := &Entry{Segment: 0, Position: 1, Offset: expectedSecondOffset, Message: []byte("beta")}
	if !reflect.DeepEqual(second, expectedSecond) {
		t.Fatalf("expected second entry %+v, got %+v", expectedSecond, second)
	}

	if size := segment.GetSize(); size != 2 {
		t.Fatalf("expected segment size 2, got %d", size)
	}

	for position, expected := range []*Entry{expectedFirst, expectedSecond} {
		entry := segment.ReadAt(uint(position))
		if entry.Segment != expected.Segment || entry.Position != expected.Position || entry.Offset != expected.Offset {
			t.Fatalf("position %d: expected entry metadata %+v, got %+v", position, expected, entry)
		}
		if !bytes.Equal(entry.Message, expected.Message) {
			t.Fatalf("position %d: expected message %q, got %q", position, expected.Message, entry.Message)
		}
	}

	segment.Close()

	directories, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(directories) != 1 || directories[0].Name() != "0000000000000000" {
		t.Fatalf("expected zero-padded segment directory, got %+v", directories)
	}
}

func TestSegmentForReadReopensExistingData(t *testing.T) {
	root := t.TempDir()

	writer := SegmentForWrite(root, 7)
	writer.Append("first")
	writer.Append("second")
	writer.Close()

	reader := SegmentForRead(root, 7)
	defer reader.Close()

	if size := reader.GetSize(); size != 2 {
		t.Fatalf("expected reopened segment size 2, got %d", size)
	}

	first := reader.ReadAt(0)
	if first.Segment != 7 || first.Position != 0 || first.Offset != 0 || !bytes.Equal(first.Message, []byte("first")) {
		t.Fatalf("unexpected first entry after reopen: %+v", first)
	}

	second := reader.ReadAt(1)
	expectedSecondOffset := headerSizeInBytes + uint(len("first"))
	if second.Segment != 7 || second.Position != 1 || second.Offset != expectedSecondOffset || !bytes.Equal(second.Message, []byte("second")) {
		t.Fatalf("unexpected second entry after reopen: %+v", second)
	}
}

func TestSegmentDirectoriesSupportRolloverStarts(t *testing.T) {
	root := t.TempDir()

	first := SegmentForWrite(root, 0)
	first.Append("alpha")
	first.Append("beta")
	first.Close()

	second := SegmentForWrite(root, 2)
	thirdEntry := second.Append("gamma")
	second.Close()

	directories, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}

	names := make([]string, 0, len(directories))
	for _, directory := range directories {
		names = append(names, directory.Name())
	}

	expectedNames := []string{"0000000000000000", "0000000000000002"}
	if !reflect.DeepEqual(names, expectedNames) {
		t.Fatalf("expected segment directories %v, got %v", expectedNames, names)
	}

	position := PositionForWrite(root)
	if position.Start != 2 || position.Cursor != 0 {
		t.Fatalf("expected write position at latest segment, got %+v", position)
	}

	reopened := SegmentForRead(root, 2)
	defer reopened.Close()

	entry := reopened.ReadAt(0)
	if entry.Segment != thirdEntry.Segment || entry.Position != thirdEntry.Position || entry.Offset != thirdEntry.Offset {
		t.Fatalf("expected reopened entry metadata %+v, got %+v", thirdEntry, entry)
	}
	if !bytes.Equal(entry.Message, thirdEntry.Message) {
		t.Fatalf("expected reopened message %q, got %q", thirdEntry.Message, entry.Message)
	}

	segmentPath := filepath.Join(root, "0000000000000002")
	for _, fileName := range []string{"index.data", "store.data"} {
		if _, err := os.Stat(filepath.Join(segmentPath, fileName)); err != nil {
			t.Fatalf("expected %s to exist in rollover segment: %v", fileName, err)
		}
	}
}
