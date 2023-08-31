package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Started node")

	// Block the main goroutine by waiting for input (just so the node won't exit)
	fmt.Println("Press Enter to exit...")
	_, err := os.Stdin.Read([]byte{0})
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
}
