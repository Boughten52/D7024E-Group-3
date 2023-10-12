package cli

import (
	"bufio"
	"d7024e/kademlia"
	"fmt"
	"os"
	"strings"
	"sync"
)

type CLI struct {
	kademlia *kademlia.Kademlia
	syncExit *sync.WaitGroup
}

// Create a new CLI instance.
func NewCLI(kademlia *kademlia.Kademlia, exit *sync.WaitGroup) CLI {
	return CLI{kademlia, exit}
}

// Listen for user input and execute commands.
func (cli *CLI) Listen() {

	stdin := bufio.NewReader(os.Stdin)

	for {
		// Print commands
		fmt.Println("")
		fmt.Println("DEFINED COMMANDS:")
		fmt.Println("put [content]")
		fmt.Println("get [hash]")
		fmt.Println("forget [hash]")
		fmt.Println("exit")
		fmt.Println("")

		text, _ := stdin.ReadString('\n')
		text = strings.TrimRight(text, "\n")
		args := strings.SplitN(text, " ", 2)

		switch {
		case text == "exit":
			cli.exit()
			return
		case args[0] == "put" && len(args) > 1:
			cli.put(args[1])
		case args[0] == "get" && len(args) > 1:
			cli.get(args[1])
		case args[0] == "forget" && len(args) > 1:
			cli.forget(args[1])
		default:
			fmt.Println("Invalid command.")
		}
	}
}

// Handle put command by storing content on the network.
func (cli *CLI) put(content string) {
	hash := cli.kademlia.Store([]byte(content))
	fmt.Println("Stored content with hash", hash)
}

// Handle get command by retrieving data from the network.
func (cli *CLI) get(hash string) {
	if len([]byte(hash)) != 40 {
		fmt.Printf("Expected hash length of 20 but got %d", len([]byte(hash)))
		return
	}

	data := cli.kademlia.LookupData(hash)
	if data == nil {
		fmt.Println("Data not found")
	} else {
		fmt.Println("Data:", string(data))
	}
}

func (cli *CLI) forget(hash string) {
	if len([]byte(hash)) != 40 {
		fmt.Printf("Expected hash length of 20 but got %d", len([]byte(hash)))
		return
	}

	cli.kademlia.Forget(hash)
}

// Handle exit command by exiting the program.
func (cli *CLI) exit() {
	cli.syncExit.Done()
}
