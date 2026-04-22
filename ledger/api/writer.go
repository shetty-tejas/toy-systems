package api

import (
	"ledger/models"
	"os"
)

type Writer struct {
	segment *models.Segment

	Directory string
	Limit     uint
}

func NewWriter(directory string, limit uint) *Writer {
	if limit == 0 {
		panic("segment limit should be greater than zero")
	}

	err := os.MkdirAll(directory, 0766)
	if err != nil {
		panic(err)
	}

	pos := models.PositionForWrite(directory)

	return &Writer{
		segment:   models.SegmentForWrite(directory, pos.Start),
		Directory: directory,
		Limit:     limit,
	}
}

func (w *Writer) Append(message string) *models.Entry {
	if w.ShouldRollover() {
		w.Rollover()
	}

	return w.segment.Append(message)
}

func (w *Writer) Close() {
	w.segment.Close()
}

func (w *Writer) ShouldRollover() bool {
	return w.segment.GetSize() >= w.Limit
}

func (w *Writer) Rollover() {
	start := w.segment.Start + w.segment.GetSize()
	w.segment.Close()

	w.segment = models.SegmentForWrite(w.Directory, start)
}
