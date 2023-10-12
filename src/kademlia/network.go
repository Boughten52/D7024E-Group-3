package kademlia

import (
	"d7024e/protobuf"
	"d7024e/utils"
	"fmt"
	"net"
	"sync"
	"time"
)

// Mutex
var mComs sync.RWMutex
var mRoutingtable sync.RWMutex

// Defines the different message types sent over the network.
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
	rt      *RoutingTable
	storage *Storage
	coms    map[string]chan map[string]string

	k     int
	alpha int
	ttl   time.Duration
}

// Create a new Network instance.
func NewNetwork(rt *RoutingTable, k int, alpha int, ttl time.Duration) *Network {
	return &Network{rt, NewStorage(ttl), make(map[string]chan map[string]string), k, alpha, ttl}
}

// Listens for incoming messages on a specified port.
func (network *Network) Listen(ip string, port int) {

	// Resolve the UDP address to bind to
	address := fmt.Sprintf("%s:%d", ip, port)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		utils.LogError("Network.Listen net resolve %w", err)
	}

	// Create a UDP connection to listen on the specified address
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		utils.LogError("Network.Listen listen on udp %w", err)
	}
	defer conn.Close()

	for {
		buffer := make([]byte, 4096) // Adjust buffer size as needed
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			utils.LogError("Network.Listen reading from udp %w", err)
			continue
		}

		// Handle incoming message in a separate goroutine
		go func() {
			values, err := protobuf.DeserializeMessage(buffer[:n])
			if err != nil {
				utils.LogError("Listen failed to deserialize message %s", err)
				return
			}
			utils.Log(1, "Recieved %s message from %s", values["type"], values["sender_address"])
			contact := NewContact(NewKademliaID(values["sender_id"]), values["sender_address"])

			switch values["type"] {
			case PING:
				network.SendPongMessage(&contact, NewKademliaID(values["rpc_id"]))

			case FIND_NODE:
				network.sendFindContactResponseMessage(values, &contact)

			case FIND_VALUE:
				// Similar to FIND_NODE, but return the value if found instead of contacts
				data, exist := network.storage.FetchData(values["key"])
				if !exist {
					network.sendFindContactResponseMessage(values, &contact)
					break
				}

				response := make(map[string]string)
				response["rpc_id"] = values["rpc_id"]
				response["sender_id"] = network.rt.me.ID.String()
				response["sender_address"] = network.rt.me.Address
				response["key"] = values["key"]
				response["data"] = string(data)
				response["type"] = FIND_VALUE_RESPONSE

				data, err := protobuf.SerializeMessage(response)
				if err != nil {
					utils.LogError("Listen FIND_VALUE could not serialize data %s", err)
					return
				}

				utils.Log(1, "Sending %s message to %s", response["type"], values["sender_address"])
				network.sendMessage(contact.Address, data)

			case STORE:
				network.storage.StoreData(values["key"], []byte(values["data"]), network.ttl)

			default:
				network.TransmitResponse(NewKademliaID(values["rpc_id"]), values)
			}

			// Update routing table with sender
			mRoutingtable.Lock()
			network.rt.AddContact(contact)
			mRoutingtable.Unlock()
		}()
	}
}

// Sends a ping message to contact.
func (network *Network) SendPingMessage(contact *Contact, rpcID *KademliaID) {
	// Create a map to hold the values for the Ping message
	values := make(map[string]string)
	values["rpc_id"] = rpcID.String()
	values["sender_id"] = network.rt.me.ID.String()
	values["sender_address"] = network.rt.me.Address
	values["type"] = PING

	data, err := protobuf.SerializeMessage(values)
	if err != nil {
		utils.LogError("SendPingPongMessage: could not build message (%s)", err)
		return
	}

	utils.Log(1, "Sending %s message to %s", PING, contact.Address)
	network.sendMessage(contact.Address, data)
}

// Sends a pong message to contact.
func (network *Network) SendPongMessage(contact *Contact, rpcID *KademliaID) {
	// Create a map to hold the values for the Ping message
	values := make(map[string]string)
	values["rpc_id"] = rpcID.String()
	values["sender_id"] = network.rt.me.ID.String()
	values["sender_address"] = network.rt.me.Address
	values["type"] = PONG

	data, err := protobuf.SerializeMessage(values)
	if err != nil {
		utils.LogError("SendPingPongMessage: could not build message (%s)", err)
		return
	}

	utils.Log(1, "Sending %s message to %s", PONG, contact.Address)
	network.sendMessage(contact.Address, data)
}

// Sends a find node message to contact.
func (network *Network) SendFindContactMessage(id *KademliaID, contact *Contact, rpcID *KademliaID) {
	// Create a map to hold the values for the FindContact message
	values := make(map[string]string)
	values["rpc_id"] = rpcID.String()
	values["sender_id"] = network.rt.me.ID.String()
	values["sender_address"] = network.rt.me.Address
	values["key"] = id.String()
	values["type"] = FIND_NODE

	// Build message
	data, err := protobuf.SerializeMessage(values)
	if err != nil {
		utils.LogError("SendFindContactMessage: could not build message (%s)", err)
		return
	}

	// Send message
	utils.Log(1, "Sending %s message to %s", FIND_NODE, contact.Address)
	network.sendMessage(contact.Address, data)
}

// Sends a find data message to contact.
func (network *Network) SendFindDataMessage(key *KademliaID, contact *Contact, rpcID *KademliaID) {
	// Create a map to hold the values for the FindData message
	values := make(map[string]string)
	values["rpc_id"] = rpcID.String()
	values["sender_id"] = network.rt.me.ID.String()
	values["sender_address"] = network.rt.me.Address
	values["key"] = key.String()
	values["type"] = FIND_VALUE

	// Build message
	data, err := protobuf.SerializeMessage(values)
	if err != nil {
		utils.LogError("SendFindDataMessage: could not build message (%s)", err)
		return
	}

	// Send message
	utils.Log(1, "Sending %s message to %s", FIND_VALUE, contact.Address)
	network.sendMessage(contact.Address, data)
}

// Sends a store message to contact.
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
	data, err := protobuf.SerializeMessage(values)
	if err != nil {
		utils.LogError("SendStoreMessage: could not build message (%s)", err)
		return
	}

	// Send message
	utils.Log(1, "Sending %s message to %s", STORE, contact.Address)
	network.sendMessage(contact.Address, data)
}

// Sends a find node response message to contact.
func (network *Network) sendFindContactResponseMessage(values map[string]string, contact *Contact) {
	contacts := ""
	for _, node := range network.rt.FindClosestContacts(NewKademliaID(values["key"]), network.k) {
		contacts += node.String() + "\n"
	}

	response := make(map[string]string)
	response["rpc_id"] = values["rpc_id"]
	response["sender_id"] = network.rt.me.ID.String()
	response["sender_address"] = network.rt.me.Address
	response["key"] = values["key"]
	response["data"] = contacts
	response["type"] = FIND_NODE_RESPONSE

	data, err := protobuf.SerializeMessage(response)
	if err != nil {
		utils.LogError("Listen %s could not serialize data %s", values["type"], err)
		return
	}

	utils.Log(1, "Sending %s message to %s", response["type"], values["sender_address"])
	network.sendMessage(contact.Address, data)
}

// Sends a message to address.
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

// Listens on a specified channel for set amount of time before timing out.
func (network *Network) ListenWithTimeout(rpcID *KademliaID, sec int) (map[string]string, error) {
	network.CreateChannel(rpcID) // makes sure it exist

	select {
	case res := <-network.coms[rpcID.String()]:
		return res, nil

	case <-time.After(time.Duration(sec) * time.Second):
		return nil, fmt.Errorf("timeout occured")

	}
}

// Sends data on a specified channel.
func (network *Network) TransmitResponse(rpcID *KademliaID, response map[string]string) {
	network.CreateChannel(rpcID) <- response
}

// Creates channel for rpc id if a channel does not already exist.
func (network *Network) CreateChannel(rpcID *KademliaID) chan map[string]string {
	mComs.Lock()
	_, exist := network.coms[rpcID.String()]
	if !exist {
		network.coms[rpcID.String()] = make(chan map[string]string, 50)
	}
	mComs.Unlock()

	return network.coms[rpcID.String()]
}

// Deletes channel for rpc id if it exists.
func (network *Network) RemoveChannel(rpcID *KademliaID) {
	mComs.Lock()
	_, exist := network.coms[rpcID.String()]
	if exist {
		delete(network.coms, rpcID.String())
	}
	mComs.Unlock()
}
