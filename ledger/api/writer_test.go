package api

import (
	"bytes"
	"fmt"
	"ledger/models"
	"path/filepath"
	"testing"
)

func TestNewWriterPanicsOnZeroLimit(t *testing.T) {
	ledgerDir := filepath.Join(t.TempDir(), "ledger")

	assertPanicContains(t, "segment limit should be greater than zero", func() {
		NewWriter(ledgerDir, 0)
	})
}

func TestWriterAppendReturnsEntriesAndRollsOverSegments(t *testing.T) {
	ledgerDir := filepath.Join(t.TempDir(), "ledger")
	writer := NewWriter(ledgerDir, 2)

	first := writer.Append([]byte("alpha"))
	second := writer.Append([]byte("beta"))
	if !writer.ShouldRollover() {
		t.Fatal("expected writer to rollover after reaching the segment limit")
	}

	third := writer.Append([]byte("gamma"))
	if writer.ShouldRollover() {
		t.Fatal("expected fresh segment to have remaining capacity")
	}
	if writer.segment.Start != 2 {
		t.Fatalf("expected writer to switch to segment 2, got %d", writer.segment.Start)
	}

	offsets := messageOffsets("alpha", "beta")
	assertEntry(t, first, 0, 0, offsets[0], "alpha")
	assertEntry(t, second, 0, 1, offsets[1], "beta")
	assertEntry(t, third, 2, 0, 0, "gamma")

	writer.Close()

	segmentZero := models.SegmentForRead(ledgerDir, 0)
	defer segmentZero.Close()
	segmentTwo := models.SegmentForRead(ledgerDir, 2)
	defer segmentTwo.Close()

	assertEntry(t, segmentZero.ReadAt(0), 0, 0, offsets[0], "alpha")
	assertEntry(t, segmentZero.ReadAt(1), 0, 1, offsets[1], "beta")
	assertEntry(t, segmentTwo.ReadAt(0), 2, 0, 0, "gamma")
}

func TestWriterReopensExistingLedgerDirectory(t *testing.T) {
	ledgerDir := filepath.Join(t.TempDir(), "ledger")
	writer := NewWriter(ledgerDir, 3)
	writer.Append([]byte("first"))
	writer.Append([]byte("second"))
	writer.Close()

	reopened := NewWriter(ledgerDir, 3)
	defer reopened.Close()

	third := reopened.Append([]byte("third"))
	offsets := messageOffsets("first", "second", "third")
	assertEntry(t, third, 0, 2, offsets[2], "third")
	if reopened.segment.Start != 0 {
		t.Fatalf("expected reopened writer to continue segment 0, got %d", reopened.segment.Start)
	}

	segment := models.SegmentForRead(ledgerDir, 0)
	defer segment.Close()

	assertEntry(t, segment.ReadAt(0), 0, 0, offsets[0], "first")
	assertEntry(t, segment.ReadAt(1), 0, 1, offsets[1], "second")
	assertEntry(t, segment.ReadAt(2), 0, 2, offsets[2], "third")
}

func newLedgerWithMessages(t *testing.T, limit uint, messages ...string) string {
	t.Helper()

	ledgerDir := filepath.Join(t.TempDir(), "ledger")
	writer := NewWriter(ledgerDir, limit)
	for _, message := range messages {
		writer.Append([]byte(message))
	}
	writer.Close()

	return ledgerDir
}

func assertEntry(t *testing.T, got *models.Entry, wantSegment, wantPosition, wantOffset uint, wantMessage string) {
	t.Helper()

	if got.Segment != wantSegment {
		t.Fatalf("expected segment %d, got %d", wantSegment, got.Segment)
	}
	if got.Position != wantPosition {
		t.Fatalf("expected position %d, got %d", wantPosition, got.Position)
	}
	if got.Offset != wantOffset {
		t.Fatalf("expected offset %d, got %d", wantOffset, got.Offset)
	}
	if !bytes.Equal(got.Message, []byte(wantMessage)) {
		t.Fatalf("expected message %q, got %q", wantMessage, string(got.Message))
	}
}

func assertPanicContains(t *testing.T, want string, fn func()) {
	t.Helper()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("expected panic containing %q", want)
		}
		if got := fmt.Sprint(recovered); !bytes.Contains([]byte(got), []byte(want)) {
			t.Fatalf("expected panic containing %q, got %q", want, got)
		}
	}()

	fn()
}

func messageOffsets(messages ...string) []uint {
	offsets := make([]uint, len(messages))
	var offset uint
	for i, message := range messages {
		offsets[i] = offset
		offset += uint(len(message) + 8)
	}

	return offsets
}
