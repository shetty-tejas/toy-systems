package models

import (
	"fmt"
	"ledger/concerns"
	"os"
)

type Store struct {
	file *os.File
}

const headerSizeInBytes = 8

func NewStore(file *os.File) *Store {
	return &Store{file}
}

func (s *Store) GetSize() uint {
	fileInfo, err := s.file.Stat()
	if err != nil {
		panic(err)
	}

	return uint(fileInfo.Size())
}

func (s *Store) Append(message []byte) uint {
	offset := s.GetSize()

	size := len(message)
	data := make([]byte, headerSizeInBytes+size)

	concerns.BinaryEncode(data[:headerSizeInBytes], uint64(size))

	bytes := copy(data[headerSizeInBytes:], message)
	if bytes != size {
		panic(fmt.Sprintf("expected %d bytes to be written, but %d bytes were written", headerSizeInBytes, bytes))
	}

	concerns.FileWrite(s.file, data)

	return offset
}

func (s *Store) ReadAt(offset uint) []byte {
	if s := s.GetSize(); offset >= s {
		panic(fmt.Sprintf("offset is greater than the file size (%d > %d)", offset, s))
	}

	header := concerns.FileRead(s.file, offset, headerSizeInBytes)
	size := concerns.BinaryDecode(header)
	data := concerns.FileRead(s.file, offset+headerSizeInBytes, size)

	return data
}

func (s *Store) Close() {
	concerns.FileClose(s.file)
}
