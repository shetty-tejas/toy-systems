package types

type Message struct {
	Entry []byte

	StoreEntrySize uint
	Offset         uint
	Position       uint
}
