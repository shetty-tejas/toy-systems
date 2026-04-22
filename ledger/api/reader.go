package api

import (
	"fmt"
	"ledger/concerns"
	"ledger/models"
)

type Reader struct {
	segment *models.Segment

	Cursor    uint
	Directory string
}

func NewReader(directory string, position uint) *Reader {
	concerns.PollTillDirExists(directory)

	pos := models.PositionForRead(directory, position)

	reader := &Reader{
		segment:   models.SegmentForRead(directory, pos.Start),
		Cursor:    pos.Cursor,
		Directory: directory,
	}

	if reader.isCursorOverflowing() {
		panic(fmt.Sprintf("cursor is expected to be less than or equal to current segment size (%d <= %d)", pos.Cursor, reader.segment.GetSize()))
	}

	return reader
}

func (r *Reader) CanReadNext() bool {
	if r.isCursorReadPending() {
		return true
	} else if r.isCursorAtEOF() {
		return !r.AtLastSegment()
	}

	panic("cursor has overflown checking before reading next. something went wrong.")
}

func (r *Reader) ReadNext() *models.Entry {
	if r.isCursorOverflowing() {
		panic("cursor has overflown before reading next. something went wrong.")
	}

	if r.isCursorAtEOF() {
		if r.AtLastSegment() {
			panic("cursor tried to read at EOF. something went wrong.")
		}

		r.Rollover()
	}

	entry := r.ReadAt(r.Cursor)
	r.Cursor++

	return entry

}

func (r *Reader) ReadAt(cursor uint) *models.Entry {
	if r.isReadPending(cursor) {
		return r.segment.ReadAt(cursor)
	}

	panic(fmt.Sprintf("cursor is expected to be less than current segment size (%d <= %d)", cursor, r.segment.GetSize()))
}

func (r *Reader) Close() {
	r.segment.Close()
}

func (r *Reader) AtLastSegment() bool {
	pos := models.PositionForRead(r.Directory, r.probableNextSegment())

	return pos.Start == r.segment.Start
}

func (r *Reader) Rollover() {
	start := r.probableNextSegment()
	r.segment.Close()

	r.segment = models.SegmentForRead(r.Directory, start)
	r.Cursor = 0
}

func (r *Reader) probableNextSegment() uint {
	return r.segment.Start + r.segment.GetSize()
}

func (r *Reader) isCursorReadPending() bool {
	return r.isReadPending(r.Cursor)
}

func (r *Reader) isReadPending(cursor uint) bool {
	return cursor < r.segment.GetSize()
}

func (r *Reader) isCursorAtEOF() bool {
	return r.Cursor == r.segment.GetSize()
}

func (r *Reader) isCursorOverflowing() bool {
	return r.Cursor > r.segment.GetSize()
}
