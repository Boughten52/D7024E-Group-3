package main

import (
	"d7024e/cli"
	"d7024e/kademlia"
	"d7024e/utils"
	"fmt"
	"sync"
)

func main() {
	var exit sync.WaitGroup
	exit.Add(1)

	fmt.Println("Setting up node")
	ip, err := utils.GetIP()
	if err != nil {
		fmt.Println(err)
		return
	}

	me := kademlia.NewContact(kademlia.NewRandomKademliaID(), ip)
	rt := kademlia.NewRoutingTable(me)
	fmt.Println("Joining kademlia network...")
	network := kademlia.NewNetwork(rt)
	fmt.Println("Kademlia network joined")

	// CLI
	local := cli.NewCLI(&network, &exit)
	go local.Listen()

	// RESTful API
	// TODO: Implement restful api

	exit.Wait()
}
