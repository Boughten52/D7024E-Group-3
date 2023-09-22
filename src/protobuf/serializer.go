package protobuf

import (
	"fmt"

	proto "google.golang.org/protobuf/proto"
)

/*
 * SerializeMessage takes a map of values and returns the serialized data
 *
 * @param values: map containing message values
 * @return tuple containing serialized message and error if any
 */
func SerializeMessage(values map[string]string) ([]byte, error) {

	// Define message
	msg := &KademliaMessage{
		Sender: &Node{
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
		return nil, fmt.Errorf("BuildMessage: failed to marshal data %w", err)
	}

	return data, nil
}

/*
 * DeserializeMessage takes a serialized message and returns a map of values
 *
 * @param data: serialized message
 * @return tuple containing map of values and error if any
 */
func DeserializeMessage(data []byte) (map[string]string, error) {
	msg := &KademliaMessage{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, fmt.Errorf("DeconstructMessage: failed to unmarshal data %w", err)
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
