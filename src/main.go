package main

import (
	"d7024e/cli"
	"d7024e/kademlia"
	"d7024e/utils"
	"fmt"
	"sync"
)

var bootstrap = kademlia.NewContact(kademlia.NewKademliaID("0000000000000000000000000000000000000000"), "127.0.0.1:8000")
var k = 20
var alpha = 3
var port = 80

func main() {
	utils.Log(1, "Initializing node")

	// Prevent main from closing before user wants to terminate node
	var exit sync.WaitGroup
	exit.Add(1)

	// Fetch host ip
	ip, err := utils.GetIP()
	if err != nil {
		fmt.Println(err)
		return
	}

	me := kademlia.NewContact(kademlia.NewRandomKademliaID(), ip)
	rt := kademlia.NewRoutingTable(me)
	net := kademlia.NewNetwork(rt, k, alpha)
	kad := kademlia.NewKademlia(net)

	// Start listening on network
	utils.Log(1, "Listening on %s:%d", ip, port)
	go net.Listen(ip, port)

	// if this is bootsrap node
	if me.Address == bootstrap.Address {
		utils.Log(1, "Im the bootstrap node")
		me.ID = bootstrap.ID
	} else {
		utils.Log(1, "Joining kademlia network...")
		kad.JoinNetwork(&bootstrap)
		utils.Log(1, "Kademlia network joined")
	}

	// CLI
	local := cli.NewCLI(kad, &exit)
	go local.Listen()

	// RESTful API
	// TODO: Implement restful api

	exit.Wait()
}
