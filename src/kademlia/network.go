package kademlia

import (
	"fmt"
	"net"
)

type Network struct {
	rt *RoutingTable
}

// TODO: join network when called
func NewNetwork(rt *RoutingTable) Network {
	return Network{rt}
}

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

/*
 * Listens for incoming UDP messages on the specified IP and port
 *
 * @param ip: IP address to listen on
 * @param port: Port to listen on
 * @return error: Error if any
 */
func (network *Network) Listen(ip string, port int) error {
	// Resolve the UDP address to bind to
	address := fmt.Sprintf("%s:%d", ip, port)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	// Create a UDP connection to listen on the specified address
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Printf("Listening for incoming UDP messages on %s...\n", address)

	buffer := make([]byte, 1024) // Adjust buffer size as needed

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
 * @param rpcID: Unique RPC ID
 */
func (network *Network) SendPingMessage(contact *Contact) {
	// Create a map to hold the values for the Ping message
	values := make(map[string]string)
	values["rpc_id"] = NewRandomKademliaID().String()
	values["sender_id"] = network.rt.me.ID.String()
	values["sender_address"] = network.rt.me.Address
	values["type"] = PING

	// Build message
	data, err := BuildMessage(values)
	if err != nil {
		fmt.Println("SendPingMessage: could not build message \n%w", err)
	}

	// Send message
	SendMessage(contact.Address, data)
}

/*
 * Sends a pong message to contact
 *
 * @param contact: Receiver of message
 * @param rpcID: Unique RPC ID
 */
func (network *Network) SendPongMessage(contact *Contact, rpcID *KademliaID) {
	// Create a map to hold the values for the Pong message
	values := make(map[string]string)
	values["rpc_id"] = rpcID.String()
	values["sender_id"] = network.rt.me.ID.String()
	values["sender_address"] = network.rt.me.Address
	values["type"] = PONG

	// Build message
	data, err := BuildMessage(values)
	if err != nil {
		fmt.Println("SendPongMessage: could not build message \n%w", err)
	}

	// Send message
	SendMessage(contact.Address, data)
}

/*
 * Sends a find contact message to contact
 *
 * @param id: ID to find
 * @param contact: Receiver of message
 * @param rpcID: Unique RPC ID
 */
func (network *Network) SendFindContactMessage(id *KademliaID) {
	// Create a map to hold the values for the FindContact message
	values := make(map[string]string)
	values["rpc_id"] = NewRandomKademliaID().String()
	values["sender_id"] = network.rt.me.ID.String()
	values["sender_address"] = network.rt.me.Address
	values["key"] = id.String()
	values["type"] = FIND_NODE

	// TODO: Create a way to find alpha closest contact to hash

	// Build message
	/*data, err := BuildMessage(values)
	if err != nil {
		fmt.Println("SendFindContactMessage: could not build message \n%w", err)
	}*/

	// Send message
	//SendMessage(contact.Address, data)
}

/*
 * Sends a find contact response message to contact
 *
 * @param contacts: Contacts to send
 * @param contact: Receiver of message
 * @param rpcID: Unique RPC ID
 */
func (network *Network) SendFindContactResponseMessage(contacts string, id *KademliaID, contact *Contact, rpcID *KademliaID) {
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
	SendMessage(contact.Address, data)
}

/*
 * Sends a find data message to contact
 *
 * @param hash: Hash to find
 * @param contact: Receiver of message
 * @param rpcID: Unique RPC ID
 */
func (network *Network) SendFindDataMessage(hash string) {
	// Create a map to hold the values for the FindData message
	values := make(map[string]string)
	values["rpc_id"] = NewRandomKademliaID().String()
	values["sender_id"] = network.rt.me.ID.String()
	values["sender_address"] = network.rt.me.Address
	values["key"] = hash
	values["type"] = FIND_VALUE

	// TODO: Create a way to find alpha closest contact to hash

	// Build message
	/*data, err := BuildMessage(values)
	if err != nil {
		fmt.Println("SendFindDataMessage: could not build message \n%w", err)
	}*/

	// Send message
	//SendMessage(contact.Address, data)
}

/*
 * Sends a store message to contact
 *
 * @param key: Key to store
 * @param data: Data to store
 * @param contact: Receiver of message
 * @param rpcID: Unique RPC ID
 */
func (network *Network) SendStoreMessage(data []byte) {
	// Create a map to hold the values for the Store message
	values := make(map[string]string)
	values["rpc_id"] = NewRandomKademliaID().String()
	values["sender_id"] = network.rt.me.ID.String()
	values["sender_address"] = network.rt.me.Address
	values["key"] = hashKey(string(data)).String()
	values["data"] = string(data)
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
func SendMessage(address string, data []byte) {
	// Create UDP connection
	conn, err := net.Dial("udp", address)
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

	switch values["type"] {
	case PING:
		fmt.Printf("Received ping message from %s", values["sender_address"])

		// Create a new contact and add it to the routing table
		contact := NewContact(NewKademliaID(values["sender_id"]), values["sender_address"])
		network.rt.AddContact(contact)

		// Send a pong message back to the sender
		network.SendPongMessage(&contact, NewKademliaID(values["rpc_id"]))

	case PONG:
		fmt.Printf("Received pong message from %s", values["sender_address"])

		// Create a new contact and add it to the routing table
		contact := NewContact(NewKademliaID(values["sender_id"]), values["sender_address"])
		network.rt.AddContact(contact)

	case FIND_NODE:
		fmt.Printf("Received find node message from %s", values["sender_address"])

		// Create a new contact and add it to the routing table
		contact := NewContact(NewKademliaID(values["sender_id"]), values["sender_address"])
		network.rt.AddContact(contact)

		// TODO: Check if we have found the correct node

		// Find closest contacts
		closestContacts := network.rt.FindClosestContacts(NewKademliaID(values["key"]), bucketSize)

		// Create a response message containing the closest nodes' information
		response := ""
		for _, node := range closestContacts {
			response += node.String() + "\n"
		}

		// Send response message back to the sender
		network.SendFindContactResponseMessage(response, NewKademliaID(values["key"]), &contact, NewKademliaID(values["rpc_id"]))

	case FIND_NODE_RESPONSE:
		fmt.Printf("Received find node response message from %s", values["sender_address"])

		// TODO: Is this code still relevant? REMAKE
		// TODO: Check if the correct node is included in data

		/*for _, str := range strings.Split(values["data"], "\n") {
			contact, err := NewContactFromString(str)
			if err != nil {
				fmt.Println("handleMessage: could not create contact from string \n%w", err)
			}

			// Send a find contact message to the contact
			network.SendFindContactMessage(NewKademliaID(values["key"]))
		}*/

	case FIND_VALUE:
		fmt.Printf("Received find value message from %s", values["sender_address"])

		// Create a new contact and add it to the routing table
		contact := NewContact(NewKademliaID(values["sender_id"]), values["sender_address"])
		network.rt.AddContact(contact)

	case FIND_VALUE_RESPONSE:

	case STORE:

	default:
		fmt.Println("handleMessage: message type not recognized \n%w", err)
	}
}
