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

func NewCLI(kademlia *kademlia.Kademlia, exit *sync.WaitGroup) CLI {
	return CLI{kademlia, exit}
}

// TODO: Implement
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
		args := strings.SplitN(text, " ", 2)

		switch {
		case text == "exit":
			cli.exit()
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

func (cli *CLI) put(content string) {
	// TODO: Implement
}

func (cli *CLI) get(hash string) {
	// TODO: Implement
}

func (cli *CLI) forget(hash string) {
	// TODO: Implement
}

func (cli *CLI) exit() {
	cli.syncExit.Done()
}
