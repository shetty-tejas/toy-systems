package types

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

type SegmentReader struct {
	current *segment

	directory string
	limit     uint
	position  uint
}

func NewSegmentReader(directory string, current, limit uint) (*SegmentReader, error) {
	if limit == 0 {
		return nil, fmt.Errorf("limit cannot be zero")
	}

	var err error

	segmentReader := &SegmentReader{directory: directory, limit: limit, position: current}
	segmentReader.current, err = fetchSegmentByCurrent(directory, current)
	if err != nil {
		return nil, err
	}

	return segmentReader, nil
}

func (sr *SegmentReader) Read() (<-chan *Message, error) {
	messages := make(chan *Message, 1)

	go func(sr *SegmentReader) {
		for {
			time.Sleep(time.Second)

			if sr.position == sr.limit {
				err := sr.current.close()
				if err != nil {
					log.Fatalf("something went wrong: %s", err)
					break
				}

				sr.position = 0
				sr.current, err = fetchSegmentByCurrent(sr.directory, sr.current.start+sr.limit)
				if err != nil {
					log.Fatalf("something went wrong: %s", err)
					break
				}
			}

			currentSize, err := sr.current.entryCount()
			if err != nil {
				log.Fatalf("something went wront: %s", err)
				break
			}

			if sr.position >= currentSize {
				continue
			}

			message, err := sr.current.read(sr.position)
			if err != nil {
				log.Fatalf("something went wrong reading: %s", err)
				break
			}

			messages <- message
			sr.position++
		}
	}(sr)

	return messages, nil
}

func fetchSegmentByCurrent(directory string, current uint) (*segment, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	directories := make([]string, 0, len(entries))

	for _, e := range entries {
		if e.IsDir() {
			directories = append(directories, e.Name())
		}
	}

	idx, err := searchDirectoriesForCurrentSegmentName(directories, 0, len(directories)-1, int(current))
	if err != nil {
		return nil, fmt.Errorf("could not find a valid segment: %s", err)
	}

	segment, err := fetchSegment(path.Join(directory, directories[idx]))
	if err != nil {
		return nil, err
	}

	return segment, nil
}

func searchDirectoriesForCurrentSegmentName(directories []string, left, right, value int) (int, error) {
	if left == right {
		return left, nil
	}

	mid := left + (right-left+1)/2
	midVal, err := strconv.Atoi(directories[mid])
	if err != nil {
		return 0, errors.New("not a valid segment")
	}

	if value == midVal {
		return mid, nil
	} else if midVal > value {
		return searchDirectoriesForCurrentSegmentName(directories, left, mid-1, value)
	} else {
		return searchDirectoriesForCurrentSegmentName(directories, mid, right, value)
	}
}
