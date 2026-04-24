package api

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewReaderReadsSequentiallyAcrossSegments(t *testing.T) {
	ledgerDir := newLedgerWithMessages(t, 2, "alpha", "beta", "gamma")
	reader := NewReader(ledgerDir, 0)
	defer reader.Close()

	if reader.AtLastSegment() {
		t.Fatal("expected reader to start before the last segment")
	}
	if !reader.CanReadNext() {
		t.Fatal("expected first entry to be readable")
	}

	assertEntry(t, reader.ReadNext(), 0, 0, 0, "alpha")
	if !reader.CanReadNext() {
		t.Fatal("expected second entry to be readable from the current segment")
	}

	assertEntry(t, reader.ReadNext(), 0, 1, messageOffsets("alpha", "beta")[1], "beta")
	if !reader.CanReadNext() {
		t.Fatal("expected reader to continue into the next segment")
	}
	if reader.AtLastSegment() {
		t.Fatal("expected another segment to remain before rollover")
	}

	assertEntry(t, reader.ReadNext(), 2, 0, 0, "gamma")
	if !reader.AtLastSegment() {
		t.Fatal("expected reader to be on the last segment after rollover")
	}
	if reader.CanReadNext() {
		t.Fatal("expected no more entries at EOF on the last segment")
	}

	assertPanicContains(t, "cursor tried to read at EOF. something went wrong.", func() {
		reader.ReadNext()
	})
}

func TestNewReaderStartsFromRequestedPosition(t *testing.T) {
	ledgerDir := newLedgerWithMessages(t, 2, "alpha", "beta", "gamma")
	reader := NewReader(ledgerDir, 1)
	defer reader.Close()

	if reader.AtLastSegment() {
		t.Fatal("expected reader at position 1 to begin before the final segment")
	}
	if !reader.CanReadNext() {
		t.Fatal("expected entry at position 1 to be readable")
	}

	assertEntry(t, reader.ReadNext(), 0, 1, messageOffsets("alpha", "beta")[1], "beta")
	assertEntry(t, reader.ReadNext(), 2, 0, 0, "gamma")
	if reader.CanReadNext() {
		t.Fatal("expected position-based reader to reach EOF after the remaining entries")
	}
}

func TestNewReaderPanicsBeforeFirstSegmentExists(t *testing.T) {
	ledgerDir := filepath.Join(t.TempDir(), "ledger")
	if err := os.Mkdir(ledgerDir, 0o766); err != nil {
		t.Fatalf("creating ledger directory: %v", err)
	}

	assertPanicContains(t, "segments are not yet initialised", func() {
		NewReader(ledgerDir, 0)
	})
}

func TestNewReaderWaitsForLedgerRootAndFirstSegment(t *testing.T) {
	ledgerDir := filepath.Join(t.TempDir(), "ledger")
	result := make(chan readerConstructionResult, 1)

	go func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				result <- readerConstructionResult{panicValue: recovered}
			}
		}()

		result <- readerConstructionResult{reader: NewReader(ledgerDir, 0)}
	}()

	time.Sleep(100 * time.Millisecond)

	writer := NewWriter(ledgerDir, 2)
	writer.Append("ready")
	writer.Close()

	select {
	case outcome := <-result:
		if outcome.panicValue != nil {
			t.Fatalf("expected constructor wait to succeed, got panic %v", outcome.panicValue)
		}
		defer outcome.reader.Close()
		assertEntry(t, outcome.reader.ReadNext(), 0, 0, 0, "ready")
	case <-time.After(8 * time.Second):
		t.Fatal("reader constructor did not return after ledger became available")
	}
}

type readerConstructionResult struct {
	reader     *Reader
	panicValue any
}
