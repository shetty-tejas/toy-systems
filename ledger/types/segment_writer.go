package types

import (
	"errors"
	"fmt"
	"os"
	"path"
)

var ErrSegmentUninitialized = errors.New("segment is not initialised")

type SegmentWriter struct {
	current *segment

	directory string
	limit     uint
}

func NewSegmentWriter(directory string, limit uint) (*SegmentWriter, error) {
	if limit == 0 {
		return nil, fmt.Errorf("limit cannot be zero")
	}

	segmentWriter := &SegmentWriter{directory: directory, limit: limit}

	err := os.MkdirAll(directory, 0766)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return nil, err
	}

	segmentPath, err := fetchLatestSegmentPath(directory)
	switch err {
	case nil:
		segmentWriter.current, err = fetchSegment(path.Join(directory, segmentPath))
		if err != nil {
			return nil, err
		}
	case ErrSegmentUninitialized:
		segmentWriter.current, err = newSegment(directory, 0)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("could not build segment writer: %s", err)
	}

	return segmentWriter, nil
}

func (sw *SegmentWriter) Close() error {
	return sw.current.close()
}

func (sw *SegmentWriter) Write(payload string) (*Message, error) {
	currentSize, err := sw.current.entryCount()
	if err != nil {
		return nil, fmt.Errorf("could not fetch current size: %s", err)
	}

	if currentSize != 0 && currentSize%sw.limit == 0 {
		err := sw.current.close()
		if err != nil {
			return nil, err
		}

		sw.current, err = newSegment(sw.directory, sw.current.start+sw.limit)
		if err != nil {
			return nil, fmt.Errorf("could not create newer segment: %s", err)
		}
	}

	return sw.current.write([]byte(payload))
}

func fetchLatestSegmentPath(directory string) (string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return "", err
	}
	if len(entries) == 0 {
		return "", ErrSegmentUninitialized
	}

	return entries[len(entries)-1].Name(), nil // We depend on the sorting done by the OS.ReadDir. If it fails, then we need to do a custom sort.
}
