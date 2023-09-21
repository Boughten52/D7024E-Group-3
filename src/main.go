package main

import (
	"d7024e/cli"
	"d7024e/kademlia"
	"d7024e/utils"
	"fmt"
	"sync"
)

func main() {
	// Known contact to join network
	friend := kademlia.NewContact(kademlia.NewKademliaID("0"), "127.0.0.1:8000")
	k := 20
	alpha := 3

	fmt.Println("Initializing node...")
	// Prevent main from closing before user wants to terminate node
	var exit sync.WaitGroup
	exit.Add(1)

	// Fetch host ip
	ip, err := utils.GetIP()
	if err != nil {
		fmt.Println(err)
		return
	}

	println(ip)

	me := kademlia.NewContact(kademlia.NewRandomKademliaID(), ip)
	rt := kademlia.NewRoutingTable(&me)
	kad := kademlia.NewKademlia(0, make(map[string]string)) // TODO: what does id do?
	net := kademlia.NewNetwork(rt, kad, k, alpha)

	// if this node is the entry node
	if friend.Address == me.Address {
		me.ID = friend.ID
	} else {
		fmt.Println("Joining kademlia network...")
		net.JoinNetwork(&friend)
		fmt.Println("Kademlia network joined")
	}

	// Network
	//go net.Listen(ip, 8000)

	// CLI
	local := cli.NewCLI(net, &exit)
	go local.Listen()

	// RESTful API
	// TODO: Implement restful api

	exit.Wait()
	fmt.Println("NODE TERMINATED")
}
