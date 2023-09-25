package kademlia

import (
	"testing"
)

func TestNewKademlia(t *testing.T) {
	network := NewNetwork(NewRoutingTable(NewContact(NewRandomKademliaID(), "172.20.0.10")), 20, 3)
	kademlia := NewKademlia(network)
	if kademlia == nil {
		t.Fatal("NewKademlia returned nil")
	}
	if kademlia.network != nil && kademlia.network != network {
		t.Error("Kademlia network not initialized correctly")
	}
	if kademlia.DataStore == nil {
		t.Error("Kademlia DataStore not initialized correctly")
	}
}

func TestLookupContact(t *testing.T) {
	// Create a Kademlia instance
	me := NewContact(NewRandomKademliaID(), "172.20.0.10")
	rt := NewRoutingTable(me)
	net := NewNetwork(rt, 20, 3)
	kademlia := NewKademlia(net)

	kademlia.LookupContact(NewRandomKademliaID())
}

func TestLookupData(t *testing.T) {
	// Create a Kademlia instance
	me := NewContact(NewRandomKademliaID(), "172.20.0.10")
	rt := NewRoutingTable(me)
	net := NewNetwork(rt, 20, 3)
	kademlia := NewKademlia(net)

	data := kademlia.LookupData("0123456789abcdef0123456789abcdef01234561")
	if data != nil {
		t.Error("LookupData() should return nil if the data does not exist")
	}
}

func TestStore(t *testing.T) {
	// Create a Kademlia instance
	me := NewContact(NewRandomKademliaID(), "172.20.0.10")
	rt := NewRoutingTable(me)
	net := NewNetwork(rt, 20, 3)
	kademlia := NewKademlia(net)

	hash := kademlia.Store([]byte("hello world"))
	if hash == "" {
		t.Error("Store() returned an empty hash")
	}
}
