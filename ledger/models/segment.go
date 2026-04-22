package models

import (
	"fmt"
	"ledger/concerns"
	"os"
	"path"
)

type Segment struct {
	index *Index
	store *Store

	Start uint
}

func SegmentForRead(directory string, start uint) *Segment {
	directory = path.Join(directory, fmt.Sprintf("%016d", start))

	segment := &Segment{Start: start}

	file := concerns.FileOpen(path.Join(directory, "index.data"), os.O_RDONLY)
	segment.index = NewIndex(file)

	file = concerns.FileOpen(path.Join(directory, "store.data"), os.O_RDONLY)
	segment.store = NewStore(file)

	return segment
}

func SegmentForWrite(directory string, start uint) *Segment {
	segmentDirectory := path.Join(directory, fmt.Sprintf("%016d", start))

	err := os.MkdirAll(segmentDirectory, 0766)
	if err != nil {
		panic(err)
	}

	segment := &Segment{Start: start}

	file := concerns.FileOpen(path.Join(segmentDirectory, "index.data"), os.O_CREATE|os.O_RDWR|os.O_APPEND)
	segment.index = NewIndex(file)

	file = concerns.FileOpen(path.Join(segmentDirectory, "store.data"), os.O_CREATE|os.O_RDWR|os.O_APPEND)
	segment.store = NewStore(file)

	return segment
}

func (s *Segment) GetSize() uint {
	return s.index.getEntryCount()
}

func (s *Segment) Append(message string) *Entry {
	entry := &Entry{Segment: s.Start}

	entry.Message = []byte(message)
	entry.Offset = s.store.Append(entry.Message)
	entry.Position = s.index.Append(entry.Offset)

	return entry
}

func (s *Segment) ReadAt(position uint) *Entry {
	entry := &Entry{Segment: s.Start, Position: position}

	entry.Offset = s.index.ReadAt(position)
	entry.Message = s.store.ReadAt(entry.Offset)

	return entry
}

func (s *Segment) Close() {
	s.index.Close()
	s.store.Close()
}
