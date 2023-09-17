package main

import (
	"fmt"
	"net"
)

type Network struct {
	rt *RoutingTable
}

/*
 * Defines the different message types sent over the network
 */
const (
	PING string = "ping"
	PONG string = "pong"
	FIND_NODE string = "find_node"
	FIND_VALUE string = "find_value"
	STORE string = "store"
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
func (network *Network) SendPingMessage(contact *Contact, rpcID *KademliaID) {
	// Create a map to hold the values for the Ping message
    values := make(map[string]string)
	values["rpc_id"] = rpcID.String()
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
 * Sends a find contact message to contact
 *
 * @param id: ID to find
 * @param contact: Receiver of message
 * @param rpcID: Unique RPC ID
 */
func (network *Network) SendFindContactMessage(id *KademliaID, contact *Contact, rpcID *KademliaID) {
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
func (network *Network) SendFindDataMessage(hash string, contact *Contact, rpcID *KademliaID) {
	// Create a map to hold the values for the FindData message
	values := make(map[string]string)
	values["rpc_id"] = rpcID.String()
	values["sender_id"] = network.rt.me.ID.String()
	values["sender_address"] = network.rt.me.Address
	values["key"] = hash
	values["type"] = FIND_VALUE

	// Build message
	data, err := BuildMessage(values)
	if err != nil {
		fmt.Println("SendFindDataMessage: could not build message \n%w", err)
	}
	
	// Send message
	SendMessage(contact.Address, data)
}

/*
 * Sends a store message to contact
 *
 * @param key: Key to store
 * @param data: Data to store
 * @param contact: Receiver of message
 * @param rpcID: Unique RPC ID
 */
func (network *Network) SendStoreMessage(key *KademliaID, data []byte, contact *Contact, rpcID *KademliaID) {
	// Create a map to hold the values for the Store message
	values := make(map[string]string)
	values["rpc_id"] = rpcID.String()
	values["sender_id"] = network.rt.me.ID.String()
	values["sender_address"] = network.rt.me.Address
	values["key"] = key.String()
	values["data"] = string(data)
	values["type"] = STORE

	// Build message
	data, err := BuildMessage(values)
	if err != nil {
		fmt.Println("SendStoreMessage: could not build message \n%w", err)
	}

	// Send message
	SendMessage(contact.Address, data)
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

/// Private Functions ///

func (network *Network) handleMessage(data []byte) {
	values, err := DeconstructMessage(data)
	if err != nil {
		fmt.Println("handleMessage: could not deconstruct message \n%w", err)
	}
	
	switch values["type"] {
	case PING:
		// Create a new contact
		// contact := NewContact(NewKademliaID(values["sender_id"]), values["sender_address"])
		
		// TODO: Send a pong message back to the sender
	case PONG:

	case FIND_NODE:

	case FIND_VALUE:
		
	case STORE:

	default:
		fmt.Println("handleMessage: message type not recognized \n%w", err)
		// TODO: Don't update routing table for default case
	}

	// Update the routing table with the sender
	contact := NewContact(NewKademliaID(values["sender_id"]), values["sender_address"])
	network.rt.Update(contact)
}
