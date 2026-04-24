package concerns

import (
	"encoding/binary"
	"reflect"
	"strconv"
	"testing"
)

func capturePanic(t *testing.T, fn func()) any {
	t.Helper()

	var recovered any
	func() {
		defer func() {
			recovered = recover()
		}()

		fn()
	}()

	if recovered == nil {
		t.Fatalf("expected panic")
	}

	return recovered
}

func TestBinaryEncodeDecodeRoundTrip(t *testing.T) {
	cases := []uint64{
		0,
		1,
		255,
		1<<32 + 7,
		uint64(^uint(0)),
	}

	for _, value := range cases {
		value := value
		t.Run(strconv.FormatUint(value, 10), func(t *testing.T) {
			buffer := make([]byte, 8)
			BinaryEncode(buffer, value)

			var expected [8]byte
			binary.BigEndian.PutUint64(expected[:], value)
			if !reflect.DeepEqual(buffer, expected[:]) {
				t.Fatalf("encoded bytes = %v, want %v", buffer, expected[:])
			}

			decoded := BinaryDecode(buffer)
			if decoded != uint(value) {
				t.Fatalf("decoded value = %d, want %d", decoded, value)
			}
		})
	}
}

func TestBinaryEncodePanicsOnShortOrOversizedBuffer(t *testing.T) {
	t.Run("short buffer", func(t *testing.T) {
		capturePanic(t, func() {
			BinaryEncode(make([]byte, 7), 42)
		})
	})

	t.Run("oversized buffer", func(t *testing.T) {
		capturePanic(t, func() {
			BinaryEncode(make([]byte, 9), 42)
		})
	})
}

func TestBinaryDecodePanicsOnShortOrOversizedBuffer(t *testing.T) {
	t.Run("short buffer", func(t *testing.T) {
		capturePanic(t, func() {
			BinaryDecode(make([]byte, 7))
		})
	})

	t.Run("oversized buffer", func(t *testing.T) {
		buffer := make([]byte, 9)
		binary.BigEndian.PutUint64(buffer[:8], 42)

		capturePanic(t, func() {
			BinaryDecode(buffer)
		})
	})
}
