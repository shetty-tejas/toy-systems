package models

type Entry struct {
	Message  []byte
	Offset   uint
	Position uint
	Segment  uint
}
