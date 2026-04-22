package types

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
)

type segment struct {
	store *store
	index *index

	path  string
	start uint
}

func newSegment(directory string, start uint) (*segment, error) {
	segment := &segment{
		path:  path.Join(directory, fmt.Sprintf("%016d", start)),
		start: start,
	}

	err := os.MkdirAll(segment.path, 0766)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return nil, err
	}

	segment.store, err = newStore(segment.path)
	if err != nil {
		err = errors.Join(err, os.RemoveAll(segment.path))
		return nil, fmt.Errorf("cannot build segment store: %s", err)
	}

	segment.index, err = newIndex(segment.path)
	if err != nil {
		err = errors.Join(err, os.RemoveAll(segment.path))
		return nil, fmt.Errorf("cannot build segment index: %s", err)
	}

	return segment, nil
}

func fetchSegment(segmentPath string) (*segment, error) {
	var err error
	segment := &segment{path: segmentPath}

	segment.store, err = fetchStore(segment.path)
	if err != nil {
		return nil, fmt.Errorf("cannot open segment store: %s", err)
	}

	segment.index, err = fetchIndex(segment.path)
	if err != nil {
		return nil, fmt.Errorf("cannot open segment index: %s", err)
	}

	start, err := strconv.Atoi(path.Base(segmentPath))
	if err != nil {
		return nil, fmt.Errorf("cannot decode start value from segment: %d", start)
	}

	segment.start = uint(start)

	return segment, nil
}

func (segment *segment) close() error {
	return errors.Join(segment.store.close(), segment.index.close())
}

func (segment *segment) entryCount() (uint, error) {
	size, err := segment.index.position()
	if err != nil {
		return 0, fmt.Errorf("could not fetch size: %s", err)
	}

	return size, nil
}

func (segment *segment) write(entry []byte) (*Message, error) {
	message := &Message{Entry: entry}

	offset, size, err := segment.store.write(entry)
	if err != nil {
		return nil, err
	}

	message.Offset = offset
	message.StoreEntrySize = uint(size)
	message.Position, _, err = segment.index.write(offset)

	if err != nil {
		return nil, err
	}

	return message, nil
}

func (segment *segment) read(position uint) (*Message, error) {
	var err error

	message := &Message{Position: position}
	message.Offset, err = segment.index.read(position)
	if err != nil {
		return nil, err
	}

	message.Entry, err = segment.store.read(message.Offset)
	if err != nil {
		return nil, err
	}

	return message, nil
}
