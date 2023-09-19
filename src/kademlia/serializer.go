package main

import (
	pb "d7024e/protobuf"
	"fmt"

	proto "google.golang.org/protobuf/proto"
)

/*
 * BuildMessage takes a map of values and returns a serialized message
 *
 * @param values: map containing message values
 * @return tuple containing serialized message and error if any
 */
func BuildMessage(values map[string]string) ([]byte, error) {

	// Define message
	msg := &pb.KademliaMessage{
		Sender: &pb.Node{
			Id:      values["sender_id"],
			Address: values["sender_address"],
		},
		RpcId: values["rpc_id"],
		Type:  values["type"],
		Key:   values["key"],
		Data:  values["data"],
	}

	// Serialize message
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("BuildMessage: failed to marshal data \n%w", err)
	}

	return data, nil
}

/*
 * DeconstructMessage takes a serialized message and returns a map of values
 *
 * @param data: serialized message
 * @return tuple containing map of values and error if any
 */
func DeconstructMessage(data []byte) (map[string]string, error) {
	msg := &pb.KademliaMessage{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, fmt.Errorf("DeconstructMessage: failed to unmarshal data \n%w", err)
	}

	values := make(map[string]string)
	values["sender_id"] = msg.Sender.Id
	values["sender_address"] = msg.Sender.Address
	values["rpc_id"] = msg.RpcId
	values["type"] = msg.Type
	values["key"] = msg.Key
	values["data"] = msg.Data

	return values, nil
}
