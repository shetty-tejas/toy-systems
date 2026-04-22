package models

import (
	"os"
	"slices"
	"strconv"
)

type Position struct {
	Cursor uint
	Start  uint
}

func PositionForWrite(directory string) *Position {
	segmentNames := getSegmentNames(directory)
	if len(segmentNames) == 0 {
		return &Position{0, 0}
	}

	s, err := strconv.Atoi(segmentNames[len(segmentNames)-1])
	if err != nil {
		panic(err)
	}

	segment := uint(s)
	return &Position{0, segment}
}

func PositionForRead(directory string, position uint) *Position {
	segmentNames := getSegmentNames(directory)
	size := uint(len(segmentNames))
	if size == 0 {
		panic("segments are not yet initialised")
	}

	index := binarySearchFloor(segmentNames, position, 0, size-1)
	s, err := strconv.Atoi(segmentNames[index])
	if err != nil {
		panic(err)
	}

	segment := uint(s)
	return &Position{position - segment, segment}
}

func getSegmentNames(directory string) []string {
	entries, err := os.ReadDir(directory)
	if err != nil {
		panic(err)
	}

	directories := make([]string, 0, len(entries))

	for _, e := range entries {
		if e.IsDir() {
			directories = append(directories, e.Name())
		}
	}

	return slices.Clip(directories)
}

func binarySearchFloor(directories []string, search, left, right uint) uint {
	if left == right {
		return left
	}

	middle := left + (right-left+1)/2
	m, err := strconv.Atoi(directories[middle])
	if err != nil {
		panic(err)
	}

	value := uint(m)

	if value == search {
		return middle
	} else if value > search {
		return binarySearchFloor(directories, search, left, middle-1)
	} else {
		return binarySearchFloor(directories, search, middle, right)
	}
}
