package kademlia

import (
	"fmt"
	"net"
	"strings"
)

/*
 * Defines the different message types sent over the network
 */
const (
	PING                string = "ping"
	PONG                string = "pong"
	FIND_NODE           string = "find_node"
	FIND_VALUE          string = "find_value"
	FIND_NODE_RESPONSE  string = "find_node_response"
	FIND_VALUE_RESPONSE string = "find_value_response"
	STORE               string = "store"
)

type Network struct {
	rt    *RoutingTable
	kad   *Kademlia
	k     int
	alpha int

	coms map[string]chan map[string]string
}

func NewNetwork(rt *RoutingTable, kad *Kademlia, k int, alpha int) *Network {
	return &Network{rt, kad, k, alpha, make(map[string]chan map[string]string)}
}

func (network *Network) JoinNetwork(contact *Contact) {
	id := NewRandomKademliaID()
	network.sendPingMessage(contact, id)

	// Wait for response
	<-network.coms[id.String()]

	fmt.Println("My contacts before find node:")
	for _, contact := range network.rt.FindClosestContacts(network.rt.me.ID, network.k) {
		fmt.Printf("Contact in my routing table: %s\n", contact.ID.String())
	}

	network.nodeLookup(network.rt.me.ID)

	fmt.Println("My contacts after find node:")
	for _, contact := range network.rt.FindClosestContacts(network.rt.me.ID, network.k) {
		fmt.Printf("Contact in my routing table: %s\n", contact.ID.String())
	}

}

func (network *Network) Put(data string) string {
	var hash string
	/*target := NewKademliaID(hash)
	// TODO: remake?
	network.sendFindContactMessage(target, network.rt.FindClosestContacts(target, network.k))
	for _, node := range network.rt.FindClosestContacts(target, network.k) {
		network.sendStoreMessage(data, &node)
	}*/

	return hash
}

func (network *Network) Get(hash string) (string, error) {
	// Check if its locally stored
	/*data, isLocal := network.kad.LookUpData(hash)
	if isLocal {
		return data, nil
	}

	network.sendFindDataMessage(hash)*/

	return "dummy data", nil // Remove
}

func (network *Network) Forget(hash string) {
	// TODO: Implement
}

/*
 * Listens for incoming UDP messages on the specified IP and port
 *
 * @param ip: IP address to listen on
 * @param port: Port to listen on
 * @return error: Error if any
 */
func (network *Network) Listen(ip string, port int) {
	fmt.Printf("Listening for incoming UDP messages on %s:%d...\n", ip, port)

	// Resolve the UDP address to bind to
	address := fmt.Sprintf("%s:%d", ip, port)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Listen: net resolve \n%w", err)
	}

	// Create a UDP connection to listen on the specified address
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Listen: listen on udp \n%w", err)
	}
	defer conn.Close()

	buffer := make([]byte, 4096) // Adjust buffer size as needed

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		// Handle incoming message in a separate goroutine
		go network.handleMessage(buffer[:n])
	}

}

/*
 * Sends a ping message to contact
 *
 * @param contact: Contact to ping
 */
func (network *Network) sendPingMessage(contact *Contact, id *KademliaID) {
	// Create a map to hold the values for the Ping message
	values := make(map[string]string)
	values["rpc_id"] = id.String()
	values["sender_id"] = network.rt.me.ID.String()
	values["sender_address"] = network.rt.me.Address
	values["type"] = PING

	// Build message
	data, err := BuildMessage(values)
	if err != nil {
		fmt.Println("sendPingMessage: could not build message \n%w", err)
	}

	network.coms[values["rpc_id"]] = make(chan map[string]string)

	// Send message
	sendMessage(contact.Address, data)
}

/*
 * Sends a pong message to contact
 *
 * @param contact: Receiver of message
 * @param rpcID: Unique RPC ID
 */
func (network *Network) sendPongMessage(contact *Contact, rpcID *KademliaID) {
	// Create a map to hold the values for the Pong message
	values := make(map[string]string)
	values["rpc_id"] = rpcID.String()
	values["sender_id"] = network.rt.me.ID.String()
	values["sender_address"] = network.rt.me.Address
	values["type"] = PONG

	// Build message
	data, err := BuildMessage(values)
	if err != nil {
		fmt.Println("sendPongMessage: could not build message \n%w", err)
	}

	// Send message
	sendMessage(contact.Address, data)
}

/*
 * Sends a find contact message to contact
 *
 * @param id: ID to find
 * @param nodes: Nodes to send message to
 */
func (network *Network) sendFindContactMessage(id *KademliaID, nodes []Contact, rpcID *KademliaID) {
	for _, node := range nodes {
		// Create a map to hold the values for the FindContact message
		values := make(map[string]string)
		values["rpc_id"] = rpcID.String()
		values["sender_id"] = network.rt.me.ID.String()
		values["sender_address"] = network.rt.me.Address
		values["key"] = id.String()
		values["type"] = FIND_NODE

		// Build message
		data, err := BuildMessage(values)
		if err != nil {
			fmt.Println("SendFindContactMessage: could not build message \n%w", err)
		}

		network.coms[values["rpc_id"]] = make(chan map[string]string)

		// Send message
		sendMessage(node.Address, data)
	}
}

/*
 * Sends a find contact response message to contact
 *
 * @param contacts: Contacts to send
 * @param id: ID to find
 * @param contact: Receiver of message
 * @param rpcID: Unique RPC ID
 */
func (network *Network) sendFindContactResponseMessage(contacts string, id *KademliaID, contact *Contact, rpcID *KademliaID) {
	// Create a map to hold the values for the FindContactResponse message
	values := make(map[string]string)
	values["rpc_id"] = rpcID.String()
	values["sender_id"] = network.rt.me.ID.String()
	values["sender_address"] = network.rt.me.Address
	values["key"] = id.String()
	values["data"] = contacts
	values["type"] = FIND_NODE_RESPONSE

	// Build message
	data, err := BuildMessage(values)
	if err != nil {
		fmt.Println("SendFindContactResponseMessage: could not build message \n%w", err)
	}

	// Send message
	sendMessage(contact.Address, data)
}

/*
 * Sends a find data message to contact
 *
 * @param hash: Hash to find
 *
 * @return string: Data retrieved
 */
func (network *Network) sendFindDataMessage(hash string) (string, error) {
	fmt.Println("Retrieving content for hash:", hash)

	// Create a map to hold the values for the FindData message
	values := make(map[string]string)
	values["rpc_id"] = NewRandomKademliaID().String()
	values["sender_id"] = network.rt.me.ID.String()
	values["sender_address"] = network.rt.me.Address
	values["key"] = hash
	values["type"] = FIND_VALUE

	//network.coms[values["rpc_id"]] = make(chan map[string]string)

	// TODO: Create a way to find alpha closest contact to hash

	// TODO: wait on channel to return data

	// Build message
	/*data, err := BuildMessage(values)
	if err != nil {
		fmt.Println("SendFindDataMessage: could not build message \n%w", err)
	}*/

	// Send message
	//SendMessage(contact.Address, data)
	return "dummy data", nil
}

/*
 * Sends a store message to contact
 *
 * @param data: Data to store
 */
func (network *Network) sendStoreMessage(data string, contact *Contact) {
	fmt.Println("Uploading content: ", data)

	// Create a map to hold the values for the Store message
	values := make(map[string]string)
	values["rpc_id"] = NewRandomKademliaID().String()
	values["sender_id"] = network.rt.me.ID.String()
	values["sender_address"] = network.rt.me.Address
	values["key"] = hashKey(data).String()
	values["data"] = data
	values["type"] = STORE

	// TODO: Create a way to find alpha closest contact to hash

	// Build message
	/*data, err := BuildMessage(values)
	if err != nil {
		fmt.Println("SendStoreMessage: could not build message \n%w", err)
	}

	// Send message
	SendMessage(contact.Address, data)*/
}

/*
 * Sends a message to address
 *
 * @param address: Address to send message to
 * @param data: Data to send
 */
func sendMessage(address string, data []byte) {
	fmt.Printf("Sending UDP message to %s\n", address+":80")

	// Create UDP connection
	conn, err := net.Dial("udp", address+":80")
	if err != nil {
		fmt.Println("SendMessage: ", err)
	}

	// Write data to address
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("SendMessage: ", err)
	}

	// Close connection
	conn.Close()
}

/*
 * Handles incoming messages
 *
 * @param data: Data received
 */
func (network *Network) handleMessage(data []byte) {

	values, err := DeconstructMessage(data)
	if err != nil {
		fmt.Println("handleMessage: could not deconstruct message \n%w", err)
	}

	// Find the sender in the routing table
	contact := NewContact(NewKademliaID(values["sender_id"]), values["sender_address"])

	switch values["type"] {
	case PING:
		fmt.Printf("Received ping message from %s\n", values["sender_address"])

		_, exist := network.coms[values["rpc_id"]]
		if exist {
			network.coms[values["rpc_id"]] <- values
			delete(network.coms, values["rpc_id"])
		}

		// Send a PONG message back to the sender
		network.sendPongMessage(&contact, NewKademliaID(values["rpc_id"]))

	case PONG:
		fmt.Printf("Received pong message from %s\n", values["sender_address"])

		_, exist := network.coms[values["rpc_id"]]
		if exist {
			network.coms[values["rpc_id"]] <- values
			delete(network.coms, values["rpc_id"])
		}

		// TODO: Handle PONG?

	case FIND_NODE:
		fmt.Printf("Received find node message from %s\n", values["sender_address"])

		// Find my k closest contacts to the target ID (can come from multiple buckets, if one is not enough).
		// If all contacts this node knows about < k, return all contacts. This is handled in FindClosestContacts.
		closestContacts := network.rt.FindClosestContacts(NewKademliaID(values["key"]), bucketSize)

		// Create a response message containing the closest nodes' information (ID and address + port)
		response := ""
		for _, node := range closestContacts {
			response += node.String() + "\n"
		}

		// Send response message back to the sender
		network.sendFindContactResponseMessage(response, NewKademliaID(values["key"]), &contact, NewKademliaID(values["rpc_id"]))

	case FIND_NODE_RESPONSE:
		fmt.Printf("Received find node response message from %s and RPC ID %s\n", values["sender_address"], values["rpc_id"])

		_, exist := network.coms[values["rpc_id"]]
		if exist {
			fmt.Println("Coms are still up and im sending the node look up response over")
			network.coms[values["rpc_id"]] <- values
			delete(network.coms, values["rpc_id"])
		}

		// TODO: Continue node lookup session (write to channel?)

		/*for _, str := range strings.Split(values["data"], "\n") {
			contact, err := NewContactFromString(str)
			if err != nil {
				fmt.Println("handleMessage: could not create contact from string \n%w", err)
			}

			// Send a find contact message to the contact
			network.SendFindContactMessage(NewKademliaID(values["key"]))
		}*/

	case FIND_VALUE:
		fmt.Printf("Received find value message from %s\n", values["sender_address"])

		// Return the value if it is stored locally

		// Otherwise, return my k closest contacts to the target ID (can come from multiple buckets, if one is not enough).

	case FIND_VALUE_RESPONSE:
		fmt.Printf("Received find value response message from %s\n", values["sender_address"])

		// If the value was returned, terminate the lookup session.

		// Otherwise, ask non-contacted nodes the value.

	case STORE:
		fmt.Printf("Received store message from %s\n", values["sender_address"])

		// Store the value locally

	default:
		fmt.Println("handleMessage: message type not recognized \n%w", err)
		return
	}

	// Update routing table with sender
	contact.CalcDistance(network.rt.me.ID)
	network.rt.AddContact(contact)
}

func (network *Network) nodeLookup(target *KademliaID) {
	fmt.Println("Starting node lookup for target:", target.String())

	rpcID := NewRandomKademliaID()

	// Init state by copying k closest nodes from own routing table.
	shortList := ContactCandidates{network.rt.FindClosestContacts(target, network.k)}
	contactedNodes := ContactCandidates{make([]Contact, 0)}

	for {
		// Filter out nodes that are already in contactedNodes
		alphaNodes := ContactCandidates{make([]Contact, 0)}
		allNodesContacted := false
		for _, node := range shortList.contacts {
			if !Contains(contactedNodes.contacts, node) {
				alphaNodes.Append([]Contact{node})
				allNodesContacted = false
			}
		}

		// Terminate when all k nodes in the state have been contacted.
		if allNodesContacted {
			fmt.Println("All nodes contacted exiting node lookup")
			break
		}

		// Limit alphaNodes to network.alpha nodes
		if alphaNodes.Len() > network.alpha {
			// Sort inorder to prioritize messaging closer nodes
			alphaNodes.Sort()
			alphaNodes.contacts = alphaNodes.contacts[:network.alpha]
		}

		network.sendFindContactMessage(target, alphaNodes.contacts, rpcID)
		contactedNodes.Append(alphaNodes.contacts)

		// TODO: Wait for response (containing max k nodes)
		fmt.Printf("Waiting on coms for RPC ID %s\n", rpcID.String())
		response := <-network.coms[rpcID.String()]
		fmt.Println("Recieved coms message continue node lookup")

		contacts := []Contact{}
		for _, str := range strings.Split(response["data"], "\n") {
			fmt.Println(str)
			contact, err := NewContactFromString(str)
			if err != nil {
				fmt.Printf("nodeLookup: could not translate string to contact %s", err)
				continue
			}

			fmt.Printf("New contact added to list: %s\n", str)
			contacts = append(contacts, contact)
		}

		// Update shortList
		nodesReplaced := false
		for _, newNode := range contacts {
			// Filter out nodes that are already in shortList
			if Contains(shortList.contacts, newNode) {
				continue
			}

			// If shortList contains less than k nodes, nodes are freely inserted into the state until it contains k nodes.
			if shortList.Len() < network.k {
				shortList.Append([]Contact{newNode})
			} else {
				// shortList already contains k elements, replace the node furthest away with the new node if it is closer than the closest node in shortList
				shortList.Sort()
				if newNode.Less(&shortList.contacts[0]) {
					shortList.contacts[len(shortList.contacts)-1] = newNode
					nodesReplaced = true
				}
			}
		}

		// If at least one node was replaced, send a new find node message to the alpha closest nodes in state.
		if nodesReplaced {
			fmt.Println("Didnt replace any node, looping again")
			continue
		}

		// Otherwise, send a find node message to the k closest nodes in state that have not been contacted yet.
		for _, node := range shortList.contacts {
			if !Contains(contactedNodes.contacts, node) {
				network.sendFindContactMessage(target, []Contact{node}, rpcID)
				contactedNodes.Append([]Contact{node})
			}
		}
	}

	fmt.Println("Node lookup terminated for target:", target.String())
}
