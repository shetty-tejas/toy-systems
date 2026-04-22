package main

import (
	"bufio"
	"flag"
	"fmt"
	"ledger/types"
	"os"
)

func main() {
	directory := *flag.String("directory", "segments", "directory to store and retrieve segments from.")
	limit := *flag.Uint("limit", 1000, "maximum messages per segment.")
	current := *flag.Uint("current", 0, "position to start reading the messages from.")

	sw, err := types.NewSegmentWriter(directory, limit)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sr, err := types.NewSegmentReader(directory, current, limit)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	messages, _ := sr.Read()

	go func(m <-chan *types.Message) {
		for {
			message := <-m
			fmt.Printf("reader: message received: %s ## size: %d ## offset: %d ## position in segment: %d", message.Entry, message.StoreEntrySize, message.Offset, message.Position)
		}
	}(messages)

	for {
		bc := bufio.NewScanner(os.Stdin)
		fmt.Println("Enter your message: ")

		if bc.Scan() {
			input := bc.Text()
			message, err := sw.Write(input)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Printf("writer: message added: %s ## size: %d ## offset: %d ## position in segment: %d", message.Entry, message.StoreEntrySize, message.Offset, message.Position)
		}
	}

}
