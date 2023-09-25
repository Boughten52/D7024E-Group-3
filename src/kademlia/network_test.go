package kademlia

import (
	"testing"
)

func TestSendMessage(t *testing.T) {
	me := NewContact(NewRandomKademliaID(), "172.20.0.10:80")
	rt := NewRoutingTable(me)
	net := NewNetwork(rt, 20, 3)
	contact := NewContact(NewRandomKademliaID(), "172.20.0.10:80")
	values := make(map[string]string)
	values["rpc_id"] = NewRandomKademliaID().String()
	values["sender_id"] = me.ID.String()
	values["sender_address"] = me.Address
	values["key"] = NewRandomKademliaID().String()
	values["data"] = "hello world!"
	values["type"] = STORE

	net.SendPingMessage(&contact, NewRandomKademliaID())
	net.SendPongMessage(&contact, NewRandomKademliaID())
	net.SendFindContactMessage(NewRandomKademliaID(), &contact, NewRandomKademliaID())
	net.SendFindDataMessage(NewRandomKademliaID(), &contact, NewRandomKademliaID())
	net.SendStoreMessage(NewRandomKademliaID(), []byte("hello world"), &contact, NewRandomKademliaID())
	net.sendFindContactResponseMessage(values, &contact)
}

func TestComs(t *testing.T) {
	me := NewContact(NewRandomKademliaID(), "172.20.0.10:80")
	rt := NewRoutingTable(me)
	net := NewNetwork(rt, 20, 3)
	rpc := NewRandomKademliaID()

	values := make(map[string]string)
	values["rpc_id"] = NewRandomKademliaID().String()
	values["sender_id"] = me.ID.String()
	values["sender_address"] = me.Address
	values["key"] = NewRandomKademliaID().String()
	values["data"] = "hello world!"
	values["type"] = STORE

	net.ListenWithTimeout(rpc, 1)
	net.CreateChannel(rpc) <- values
	net.ListenWithTimeout(rpc, 1)
	net.RemoveChannel(rpc)
}
