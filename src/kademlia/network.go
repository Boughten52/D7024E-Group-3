package kademlia

import (
	"d7024e/protobuf"
	"fmt"
	"net"
	"time"
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
	rt   *RoutingTable
	coms map[string]chan map[string]string

	k     int
	alpha int
}

func NewNetwork(contact *Contact) {

}

func (network *Network) Listen(ip string, port int) {

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

	for {
		buffer := make([]byte, 4096) // Adjust buffer size as needed
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		// Handle incoming message in a separate goroutine
		go func() {
			values, err := protobuf.DeserializeMessage(buffer[:n])
			if err != nil {
				fmt.Println(err)
				return
			}

			contact := NewContact(NewKademliaID(values["sender_id"]), values["sender_address"])

			switch values["type"] {
			case PING:
				network.SendPingPongMessage(&contact, NewKademliaID(values["rpc_id"]), PONG)
			case FIND_NODE:
				// TODO: Implement
			case FIND_VALUE:
				// TODO: Implement
			case STORE:
				// TODO: Implement
			default:
				network.SendResponse(NewKademliaID(values["rpc_id"]), values)
			}

			// Update routing table with sender
			network.rt.AddContact(contact)
		}()
	}
}

/*
 * Sends a ping message to contact
 *
 * @param contact: Contact to ping/pong
 */
func (network *Network) SendPingPongMessage(contact *Contact, rpcID *KademliaID, msgType string) {
	// Create a map to hold the values for the Ping message
	values := make(map[string]string)
	values["rpc_id"] = rpcID.String()
	values["sender_id"] = network.rt.me.ID.String()
	values["sender_address"] = network.rt.me.Address
	values["type"] = msgType

	data, err := protobuf.SerializeMessage(values)
	if err != nil {
		fmt.Printf("SendPingPongMessage: could not build message (%s)", err)
		return
	}

	network.sendMessage(contact.Address, data)

}

func (network *Network) SendFindContactMessage(contact *Contact, rpcID *KademliaID) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string, rpcID *KademliaID) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte, rpcID *KademliaID) {
	// TODO
}

/*
 * Sends a message to address
 *
 * @param address: Address to send message to
 * @param port: Port to send message to
 * @param data: Data to send
 */
func (network *Network) sendMessage(address string, data []byte) {
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
 * Listens on a specified channel for set amount of time before timing out
 */
func (network *Network) ListenWithTimeout(rpcID *KademliaID, sec int) (map[string]string, error) {
	network.CreateChannel(rpcID) // makes sure it exist

	select {
	case res := <-network.coms[rpcID.String()]:
		return res, nil

	case <-time.After(time.Duration(sec) * time.Second):
		return nil, fmt.Errorf("timeout occured")

	}
}

/*
 * Sends data on a specified channel
 */
func (network *Network) SendResponse(rpcID *KademliaID, response map[string]string) {
	network.CreateChannel(rpcID) // makes sure it exist
	network.coms[rpcID.String()] <- response
}

/*
 * Creates channel for rpc id if a channel does not already exist
 */
func (network *Network) CreateChannel(rpcID *KademliaID) {
	_, exist := network.coms[rpcID.String()]
	if !exist {
		network.coms[rpcID.String()] = make(chan map[string]string)
	}
}

/*
 * Deletes channel for rpc id if it exists
 */
func (network *Network) RemoveChannel(rpcID *KademliaID) {
	_, exist := network.coms[rpcID.String()]
	if exist {
		delete(network.coms, rpcID.String())
	}
}
