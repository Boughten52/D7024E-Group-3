package cli

import (
	"flag"
	"fmt"
	"os"
)

// startCLI initializes the command-line interface
func StartCLI() {
	// Define flags for command-line arguments
	putCmd := flag.NewFlagSet("put", flag.ExitOnError)
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	forgetCmd := flag.NewFlagSet("forget", flag.ExitOnError)
	exitCmd := flag.NewFlagSet("exit", flag.ExitOnError)

	// Define command-line flags for put command
	putContent := putCmd.String("content", "", "Contents of the file to upload")

	// Define command-line flags for get command
	getHash := getCmd.String("hash", "", "Hash of the object to retrieve")

	// Define command-line flags for forget command
	forgetHash := forgetCmd.String("hash", "", "Hash of the object to stop refreshing")

	// Parse command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./node-cli [command] [options]")
		fmt.Println("Commands:")
		fmt.Println("  put -content [content]: Upload a file and output its hash")
		fmt.Println("  get -hash [hash]: Retrieve an object by its hash")
		fmt.Println("  forget -hash [hash]: Stop refreshing an object by its hash")
		fmt.Println("  exit: Terminate the node")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "put":
		putCmd.Parse(os.Args[2:])
		if *putContent == "" {
			fmt.Println("Please provide content with the -content flag")
			os.Exit(1)
		}
		put(*putContent)
	case "get":
		getCmd.Parse(os.Args[2:])
		if *getHash == "" {
			fmt.Println("Please provide a hash with the -hash flag")
			os.Exit(1)
		}
		get(*getHash)
	case "forget":
		forgetCmd.Parse(os.Args[2:])
		if *forgetHash == "" {
			fmt.Println("Please provide a hash with the -hash flag")
			os.Exit(1)
		}
		forget(*forgetHash)
	case "exit":
		exitCmd.Parse(os.Args[2:])
		exit()
	default:
		fmt.Println("Invalid command. Use 'put', 'get', or 'exit'.")
		os.Exit(1)
	}
}

func put(content string) {
	fmt.Println("Uploading content:", content)
	// TODO: Calculate the hash and output it
}

func get(hash string) {
	fmt.Println("Retrieving content for hash:", hash)
	// TODO: Retrieve and output the content and the node it was retrieved from
}

func forget(hash string) {
	fmt.Println("Stopping refresh for hash:", hash)
	// TODO: Implement the 'forget' logic here
}

func exit() {
	fmt.Println("Terminating the node")
	// TODO: Your Kademlia node termination logic goes here
}
