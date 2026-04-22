package main

import (
	"fmt"
	"ledger/api"
	"time"
)

func main() {
	reader := api.NewReader("segments", 0)
	defer reader.Close()

	for {
		if reader.CanReadNext() {
			entry := reader.ReadNext()

			fmt.Printf("===Message===\n")
			fmt.Printf("%s\n", entry.Message)
			fmt.Printf("===Message===\n")
			fmt.Printf("===Stats===\n")
			fmt.Printf("Segment: %d || Position: %d || Offset: %d\n", entry.Segment, entry.Position, entry.Offset)
			fmt.Printf("===Stats===\n")
		}

		time.Sleep(2 * time.Second)
	}
}
