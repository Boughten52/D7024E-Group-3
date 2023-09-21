package cli

import (
	"d7024e/kademlia"
	"fmt"
	"strings"
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
		args := strings.Fields(input)
		if !(len(args) > 1 || (len(args) > 0 && args[0] == "exit")) {
			fmt.Println("Too few arguments")
			continue
		}

		switch args[0] {
		case "exit":
			cli.exit()
		case "put":
			cli.put(args[1])
		case "get":
			cli.get(args[1])
		case "forget":
			cli.forget(args[1])
		default:
			fmt.Println("Invalid command.")
		}
	}
}

func (cli *CLI) put(content string) {
	fmt.Println("Succesfully stored '", content, "' with the key: ", cli.network.Put(content))
}

func (cli *CLI) get(hash string) {
	res, err := cli.network.Get(hash)
	if err != nil {
		fmt.Println("Could not retrieve content...")
	}
	fmt.Println("Content: ", res)
}

func (cli *CLI) forget(hash string) {
	fmt.Println("Stopping refresh for hash:", hash)
	// TODO: Implement the 'forget' logic here
}

func (cli *CLI) exit() {
	cli.syncExit.Done()
}
