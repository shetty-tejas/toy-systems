package types

import (
	"encoding/binary"
	"fmt"
	"os"
	"path"
)

type store struct {
	file *os.File
}

const sizeBytes = 8

func newStore(segmentPath string) (*store, error) {
	file, err := os.Create(storePath(segmentPath))
	if err != nil {
		return nil, err
	}

	return &store{file: file}, nil
}

func fetchStore(segmentPath string) (*store, error) {
	file, err := os.OpenFile(storePath(segmentPath), os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &store{file: file}, nil
}

func (store *store) close() error {
	return store.file.Close()
}

func (store *store) read(offset uint) ([]byte, error) {
	size := make([]byte, sizeBytes)

	n, err := store.file.ReadAt(size, int64(offset))
	if err != nil {
		return nil, err
	}

	if s := len(size); n != s {
		return nil, fmt.Errorf("something went wrong. bytes read '%d' are less than expected '%d'", n, s)
	}

	var entrySize uint64

	n, err = binary.Decode(size, binary.LittleEndian, &entrySize)
	if err != nil {
		return nil, err
	}

	entry := make([]byte, entrySize)

	n, err = store.file.ReadAt(entry, int64(offset+recordBytes))
	if err != nil {
		return nil, err
	}

	if n != int(entrySize) {
		return nil, fmt.Errorf("something went wrong. bytes read '%d' are less than expected '%d'", n, entrySize)
	}

	return entry, nil
}

func (store *store) size() (uint, error) {
	info, err := store.file.Stat()
	if err != nil {
		return 0, err
	}

	return uint(info.Size()), nil
}

func (store *store) write(entry []byte) (uint, int, error) {
	offset, err := store.size()
	if err != nil {
		return 0, 0, err
	}

	entrySize := uint64(len(entry))
	payload := make([]byte, sizeBytes+entrySize)

	_, err = binary.Encode(payload[:sizeBytes], binary.LittleEndian, entrySize)
	if err != nil {
		return 0, 0, nil
	}

	copy(payload[sizeBytes:], entry)

	size, err := store.file.Write(payload)
	if err != nil {
		return 0, 0, err
	}

	return offset, size, nil
}

func storePath(segmentPath string) string {
	return path.Join(segmentPath, "store.log")
}
