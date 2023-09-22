package main

import (
	"d7024e/kademlia"
	"d7024e/utils"
	"fmt"
	"sync"
	"time"
	//"net/http"
)

func main() {
	/*resp, err := http.Get("http://d7024e-group-3_kademlia_network-6")
		if err != nil {
	        panic(err)
	    }*/

	// Known contact to join network
	friend := kademlia.NewContact(kademlia.NewKademliaID("0000000000000000000000000000000000000000"), "172.20.0.10")
	k := 20
	alpha := 3
	tExpire := 86400    // 24 hours
	tRefresh := 3600    // 1 hour
	tReplicate := 3600  // 1 hour
	tRepublish := 86400 // 24 hours

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
	rt := kademlia.NewRoutingTable(&me)
	kad := kademlia.NewKademlia(0, make(map[string]string)) // TODO: what does id do?
	net := kademlia.NewNetwork(rt, kad, k, alpha, tExpire, tRefresh, tReplicate, tRepublish)

	// Network
	go net.Listen(ip, 80)

	// Check if entry node or not
	if friend.Address == me.Address {
		fmt.Println("Im the entry node")
		me.ID = friend.ID
	} else {
		time.Sleep(10 * time.Second) // TODO: remove

		fmt.Println("Joining kademlia network...")
		net.JoinNetwork(&friend)
		fmt.Println("Kademlia network joined")
	}

	// CLI
	//local := cli.NewCLI(net, &exit)
	//go local.Listen()

	// RESTful API
	// TODO: Implement restful api

	exit.Wait()
	fmt.Println("NODE TERMINATED")
}
