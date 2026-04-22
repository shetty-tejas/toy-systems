package concerns

import (
	"encoding/binary"
	"fmt"
)

var binaryByteOrder = binary.BigEndian

func BinaryDecode(data []byte) uint {
	var result uint64

	bytes, err := binary.Decode(data, binaryByteOrder, &result)
	if err != nil {
		panic(err)
	}

	if s := len(data); bytes != s {
		panic(fmt.Sprintf("expected %d bytes to be written, but %d bytes were written", s, bytes))
	}

	return uint(result)
}

func BinaryEncode(buffer []byte, data uint64) {
	bytes, err := binary.Encode(buffer, binaryByteOrder, data)
	if err != nil {
		panic(err)
	}
	if s := len(buffer); bytes != s {
		panic(fmt.Sprintf("expected %d bytes to be written, but %d bytes were written", s, bytes))
	}
}
