package main

import (
	"d7024e/api"
	"d7024e/cli"
	"d7024e/kademlia"
	"d7024e/utils"
	"fmt"
	"sync"
	"time"
)

var bootstrap = kademlia.NewContact(kademlia.NewKademliaID("0000000000000000000000000000000000000000"), "172.20.0.10:80")
var k = 20
var alpha = 3
var ttl = time.Second * 86430             // 24 hours and 30 seconds
var refreshInterval = time.Second * 86400 // 24 hours
var port = 80

func main() {

	// Prevent main from closing before user wants to terminate node
	var exit sync.WaitGroup
	exit.Add(1)

	// Fetch host ip
	ip, err := utils.GetIP()
	if err != nil {
		fmt.Println(err)
		return
	}

	utils.Log(1, "Hello I exist and my ip is %s", ip)

	address := fmt.Sprintf("%s:%d", ip, port)
	me := kademlia.NewContact(kademlia.NewRandomKademliaID(), address)
	rt := kademlia.NewRoutingTable(me)
	net := kademlia.NewNetwork(rt, k, alpha, ttl, refreshInterval)
	kad := kademlia.NewKademlia(net)
	kad.StartRefreshRoutine()

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
	go api.StartServer(kad, port)

	exit.Wait()
	utils.Log(3, "NODE TERMINATED")
}
