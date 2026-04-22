package main

import (
	"bufio"
	"fmt"
	"ledger/api"
	"os"
)

func main() {
	writer := api.NewWriter("segments", 10)
	defer writer.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()

		entry := writer.Append(message)

		fmt.Println("\n===Message===")
		fmt.Printf("%s\n", entry.Message)
		fmt.Println("===Message===")
		fmt.Println("===Stats===")
		fmt.Printf("Segment: %d || Position: %d || Offset: %d\n", entry.Segment, entry.Position, entry.Offset)
		fmt.Printf("===Stats===\n\n")
	}
}
