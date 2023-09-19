package cli

import (
	"d7024e/kademlia"
	"fmt"
	"regexp"
	"sync"
)

type CLI struct {
	network  *kademlia.Network
	syncExit *sync.WaitGroup
}

func NewCLI(network *kademlia.Network, exit *sync.WaitGroup) CLI {
	return CLI{network, exit}
}

// TODO: Implement
func (cli *CLI) Listen() {

	re := regexp.MustCompile(`(\w+) (\w+)`)
	var input string

	for {
		// Print commands
		fmt.Println("Defined Commands:")
		fmt.Println("put [content]")
		fmt.Println("get [hash]")
		fmt.Println("forget [hash]")
		fmt.Println("exit")
		fmt.Println("")

		// Read input
		fmt.Print("Enter command: ")
		_, err := fmt.Scan(&input)
		if err != nil {
			fmt.Println("CLI: could not read input \n%w", err)
			continue
		}

		// Find input variant
		args := re.FindStringSubmatch(input)
		switch args[0] {
		case "put":
			cli.put(args[1])
		case "get":
			cli.get(args[1])
		case "forget":
			cli.forget(args[1])
		case "exit":
			cli.exit()
		default:
			fmt.Println("Invalid command.")
		}

	}
}

func (cli *CLI) put(content string) {
	fmt.Println("Uploading content: ", content)
	// TODO: Calculate the hash and output it
}

func (cli *CLI) get(hash string) {
	fmt.Println("Retrieving content for hash:", hash)
	// TODO: Retrieve and output the content and the node it was retrieved from
}

func (cli *CLI) forget(hash string) {
	fmt.Println("Stopping refresh for hash:", hash)
	// TODO: Implement the 'forget' logic here
}

func (cli *CLI) exit() {
	fmt.Println("Terminating the node")
	cli.syncExit.Done()
}
