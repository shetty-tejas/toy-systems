package concerns

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFileOpenWriteReadCloseRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "segment.bin")
	written := []byte{1, 2, 3, 4, 5}

	file := FileOpen(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC)
	FileWrite(file, written)

	if got := FileRead(file, 1, 3); !reflect.DeepEqual(got, []byte{2, 3, 4}) {
		t.Fatalf("read bytes = %v, want %v", got, []byte{2, 3, 4})
	}

	FileClose(file)

	if _, err := file.Write([]byte{9}); err == nil {
		t.Fatalf("expected closed file to reject writes")
	}

	reopened := FileOpen(path, os.O_RDONLY)
	defer FileClose(reopened)

	if got := FileRead(reopened, 0, uint(len(written))); !reflect.DeepEqual(got, written) {
		t.Fatalf("persisted bytes = %v, want %v", got, written)
	}
}

func TestFileReadPanicsPastEOF(t *testing.T) {
	path := filepath.Join(t.TempDir(), "segment.bin")
	file := FileOpen(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC)
	defer FileClose(file)

	FileWrite(file, []byte{1, 2, 3})

	capturePanic(t, func() {
		FileRead(file, 2, 2)
	})
}

func TestDirExists(t *testing.T) {
	root := t.TempDir()
	dirPath := filepath.Join(root, "segments")
	filePath := filepath.Join(root, "segment.bin")

	if err := os.Mkdir(dirPath, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filePath, []byte("data"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if dirExists(filepath.Join(root, "missing")) {
		t.Fatalf("expected missing path to report false")
	}
	if !dirExists(dirPath) {
		t.Fatalf("expected directory path to report true")
	}

	capturePanic(t, func() {
		dirExists(filePath)
	})
}

func TestFileExists(t *testing.T) {
	root := t.TempDir()
	dirPath := filepath.Join(root, "segments")
	filePath := filepath.Join(root, "segment.bin")

	if err := os.Mkdir(dirPath, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filePath, []byte("data"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if fileExists(filepath.Join(root, "missing")) {
		t.Fatalf("expected missing path to report false")
	}
	if !fileExists(filePath) {
		t.Fatalf("expected file path to report true")
	}

	capturePanic(t, func() {
		fileExists(dirPath)
	})
}

func TestPollTillDirExistsReturnsForExistingDirectory(t *testing.T) {
	PollTillDirExists(t.TempDir())
}
