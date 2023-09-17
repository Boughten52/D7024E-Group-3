package main

import (
	"errors"
	"fmt"
	"net"
	pb "d7024e/protobuf"
	proto "github.com/golang/protobuf/proto"
)

type Network struct {
	rt     *RoutingTable
}

/*
 * Defines the different message types sent over the network
 */
type MessageType int32
const (
	Ping MessageType = 0
	Pong MessageType = 1
	FindNode MessageType = 2
	FindValue MessageType = 3
)

func Listen(ip string, port int) {
	// TODO
}

/*
 * Sends a ping message to contact
 *
 * Takes: reciever, unique rpc id
 */
func (network *Network) SendPingMessage(contact *Contact, rpcID *KademliaID) {

	msg := &pb.Ping {
		Sender: &pb.Node {
			Id: network.rt.me.ID.String(),
			Address: network.rt.me.Address,
		},
		RpcID: rpcID.String(),
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("SendPingMessage: failed to marshal data %w", err)
	}

	sendMessage(contact.Address, Ping, data)
}

func (network *Network) SendFindContactMessage(id *KademliaID, contact *Contact, rpcID *KademliaID) {
	msg := &pb.Find {
		Sender: &pb.Node {
			Id: network.rt.me.ID.String(),
			Address: network.rt.me.Address,
		},
		Key: id.String(),
		RpcID: rpcID.String(),
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("SendFindContactMessage: failed to marshal data %w", err)
	}

	sendMessage(contact.Address, FindNode, data)
}

func (network *Network) SendFindDataMessage(hash string, contact *Contact, rpcID *KademliaID) {
	msg := &pb.Find {
		Sender: &pb.Node {
			Id: network.rt.me.ID.String(),
			Address: network.rt.me.Address,
		},
		Key: hash,
		RpcID: rpcID.String(),
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("SendFindDataMessage: failed to marshal data %w", err)
	}

	sendMessage(contact.Address, FindValue, data)
}

func (network *Network) SendStoreMessage(data []byte, contact *Contact, rpcID *KademliaID) {
	// TODO
}

/////////////////////////
/// Private Functions ///
/////////////////////////

func handleMessage(channel chan []byte, self *Contact, network *Network) (error) {
	data := <- channel
	mType, message, err := unwrapMessage(data)
	if err != nil {
		return err
	}

	switch mType {
	case Ping:
		ping := &pb.Ping{}
		err := proto.Unmarshal(message, ping)
		if err != nil {
			return err
		}

	default:
		return errors.New("Handle Message: unknown message type")

	}


	return nil

}

/*
 * Sends a protobuf message over UDP connection
 *
 * Takes: address of reciever, type of message, and marshalled message
 */
func sendMessage(address string, mType MessageType, data []byte) {
	// Wrap message for informed deserialization by reciever
	msg, err := wrapMessage(mType, data)
	if (err != nil) {
		fmt.Println("SendMessage: ", err)
	}

	// Create UDP connection
	conn, err := net.Dial("udp", address)
	if err != nil {
		fmt.Println("SendMessage: ", err)
	}

	// Write serialized data to address
	_, err = conn.Write(msg)
	if (err != nil) {
		fmt.Println("SendMessage: ", err)
	}

	// Close connection
	conn.Close()
}

/*
 * Wraps protobuf message with associated message type for informed deserialization later
 *
 * Takes: type of message, and marshalled message
 * Returns: marshalled data ready to be sent
 */
func wrapMessage(mType MessageType, data []byte) ([]byte, error) {

	wrapper := &pb.KademliaMessage {
		MessageType: int32(mType),
		Data: data,
	}

	return proto.Marshal(wrapper)
}

/*
 * Unwraps protobuf message with associated message type for informed deserialization later
 *
 * Takes: marshalled message
 * Returns: message type of marshalled data, Marshalled data
 */
func unwrapMessage(data []byte) (MessageType, []byte, error) {

	message := &pb.KademliaMessage{}
	err := proto.Unmarshal(data, message)
	if err != nil {
		return 0, []byte{0}, err
	}
	mType := MessageType(message.MessageType)
	unwrapped := message.Data

	return mType, unwrapped, nil
}


