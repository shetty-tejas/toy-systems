package concerns

import (
	"fmt"
	"os"
	"time"
)

func FileClose(file *os.File) {
	err := file.Close()
	if err != nil {
		panic(err)
	}
}

func FileOpen(path string, mode int) *os.File {
	file, err := os.OpenFile(path, mode, 0666)
	if err != nil {
		panic(err)
	}

	return file
}

func FileRead(file *os.File, offset, size uint) []byte {
	data := make([]byte, size)

	bytes, err := file.ReadAt(data, int64(offset))
	if err != nil {
		panic(err)
	}

	if bytes != int(size) {
		panic(fmt.Sprintf("expected %d bytes to be read, but %d bytes were read", size, bytes))
	}

	return data
}

func FileWrite(file *os.File, data []byte) {
	bytes, err := file.Write(data)
	if err != nil {
		panic(err)
	}

	if s := len(data); bytes != len(data) {
		panic(fmt.Sprintf("expected %d bytes to be written, but %d bytes were written", s, bytes))
	}
}

func PollTillDirExists(path string) {
	for i := 0; !dirExists(path); i++ {
		if i == 5 {
			panic(fmt.Sprintln("could not find expected directory at", path))
		}

		time.Sleep(2 * time.Second)
	}

	time.Sleep(time.Second)
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		panic(err)
	}

	if info.IsDir() {
		return true
	}

	panic("path exists as a file and not as a directory")
}

func fileExists(path string) bool {

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		panic(err)
	}

	if !info.IsDir() {
		return true
	}

	panic("path exists as a directory and not as a file")
}
