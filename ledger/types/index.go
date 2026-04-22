package types

import (
	"encoding/binary"
	"fmt"
	"os"
	"path"
)

type index struct {
	file *os.File
}

const recordBytes = 8

func newIndex(segmentPath string) (*index, error) {
	file, err := os.Create(indexPath(segmentPath))
	if err != nil {
		return nil, err
	}

	return &index{file: file}, nil
}

func fetchIndex(segmentPath string) (*index, error) {
	file, err := os.OpenFile(indexPath(segmentPath), os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &index{file: file}, nil
}

func (index *index) close() error {
	return index.file.Close()
}

func (index *index) read(position uint) (uint, error) {
	maxPos, err := index.position()
	if err != nil {
		return 0, err
	}

	if position >= maxPos {
		return 0, fmt.Errorf("position bigger than expected")
	}

	readBytes := make([]byte, recordBytes)

	n, err := index.file.ReadAt(readBytes, int64(position)*recordBytes)
	if err != nil {
		return 0, err
	}

	if n != recordBytes {
		return 0, fmt.Errorf("something went wrong. bytes read '%d' are less than expected '%d'", n, recordBytes)
	}

	var offset uint64

	n, err = binary.Decode(readBytes, binary.LittleEndian, &offset)
	if err != nil {
		return 0, err
	}

	return uint(offset), nil
}

func (index *index) size() (uint, error) {
	info, err := index.file.Stat()
	if err != nil {
		return 0, err
	}

	return uint(info.Size()), nil
}

func (index *index) write(offset uint) (uint, int, error) {
	position, err := index.position()
	if err != nil {
		return 0, 0, err
	}

	payload := make([]byte, recordBytes)

	_, err = binary.Encode(payload, binary.LittleEndian, uint64(offset))
	if err != nil {
		return 0, 0, err
	}

	size, err := index.file.Write(payload)
	if err != nil {
		return 0, 0, err
	}

	return position, size, nil
}

func (index *index) position() (uint, error) {
	size, err := index.size()
	if err != nil {
		return 0, err
	}

	return size / recordBytes, nil
}

func indexPath(segmentPath string) string {
	return path.Join(segmentPath, "index.idx")
}
