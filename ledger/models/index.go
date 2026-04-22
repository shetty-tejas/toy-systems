package models

import (
	"fmt"
	"ledger/concerns"
	"os"
)

type Index struct {
	file *os.File
}

const entrySizeInBytes = 8

func NewIndex(file *os.File) *Index {
	return &Index{file}
}

func (i *Index) GetSize() uint {
	fileInfo, err := i.file.Stat()
	if err != nil {
		panic(err)
	}

	return uint(fileInfo.Size())
}

func (i *Index) Append(offset uint) uint {
	position := i.getEntryCount()

	data := make([]byte, entrySizeInBytes)

	concerns.BinaryEncode(data, uint64(offset))
	concerns.FileWrite(i.file, data)

	return position
}

func (i *Index) ReadAt(position uint) uint {
	off := position * entrySizeInBytes
	if s := i.GetSize(); off >= s {
		panic(fmt.Sprintf("offset is greater than the file size (%d > %d)", off, s))
	}

	data := concerns.FileRead(i.file, off, entrySizeInBytes)
	offset := concerns.BinaryDecode(data)

	return uint(offset)
}

func (i *Index) Close() {
	concerns.FileClose(i.file)
}

func (i *Index) getEntryCount() uint {
	return i.GetSize() / entrySizeInBytes
}
