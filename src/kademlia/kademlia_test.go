package kademlia

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewKademlia(t *testing.T) {
	network := NewNetwork(NewRoutingTable(NewContact(NewRandomKademliaID(), "172.20.0.10")), 20, 3, time.Second*60)
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
	net := NewNetwork(rt, 20, 3, time.Second*60)
	kademlia := NewKademlia(net)

	kademlia.LookupContact(NewRandomKademliaID())
}

func TestLookupData(t *testing.T) {
	// Create a Kademlia instance
	me := NewContact(NewRandomKademliaID(), "172.20.0.10")
	rt := NewRoutingTable(me)
	net := NewNetwork(rt, 20, 3, time.Second*60)
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
	net := NewNetwork(rt, 20, 3, time.Second*60)
	kademlia := NewKademlia(net)

	hash := kademlia.Store([]byte("hello world"))
	if hash == "" {
		t.Error("Store() returned an empty hash")
	}
}

func TestUpdateShortList(t *testing.T) {
	target := NewKademliaID("0000000000000000000000000000000000000000")
	id1 := NewKademliaID("1111111111111111111111111111111111111111")
	id2 := NewKademliaID("2222222222222222222222222222222222222222")
	id3 := NewKademliaID("3333333333333333333333333333333333333333")
	id4 := NewKademliaID("4444444444444444444444444444444444444444")

	// Initialize shortList and new contacts
	shortList := &ContactCandidates{contacts: []Contact{
		{ID: id1, Address: "192.168.0.1", distance: id1.CalcDistance(target)},
		{ID: id2, Address: "192.168.0.2", distance: id2.CalcDistance(target)},
	}}
	newContacts := []Contact{
		{ID: id3, Address: "192.168.0.3", distance: id3.CalcDistance(target)},
		{ID: id4, Address: "192.168.0.4", distance: id4.CalcDistance(target)},
	}

	// Call the function
	nodesReplaced := updateShortList(shortList, newContacts, 3)

	// Print distances for debugging
	fmt.Println("Distances after update:")
	for _, contact := range shortList.contacts {
		fmt.Printf("Node ID: %s, Distance: %s\n", contact.ID.String(), contact.distance.String())
	}

	// Assert the results
	expectedContacts := []Contact{
		{ID: id1, Address: "192.168.0.1", distance: id1.CalcDistance(target)},
		{ID: id2, Address: "192.168.0.2", distance: id2.CalcDistance(target)},
		{ID: id3, Address: "192.168.0.3", distance: id3.CalcDistance(target)},
	}

	if !compareContactArrays(shortList.contacts, expectedContacts) {
		t.Errorf("Short list not updated as expected.")
	}

	if nodesReplaced {
		t.Errorf("Nodes should just be added, and not replaced.")
	}
}

func compareContactArrays(a, b []Contact) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].ID.String() != b[i].ID.String() || a[i].Address != b[i].Address {
			return false
		}
	}
	return true
}

func TestWaitForFastest(t *testing.T) {
	// Prepare test data
	testData := []byte("test data")
	channel := make(chan []byte, 1)
	channel <- testData

	// Use a WaitGroup to simulate asynchronous tasks
	var wg sync.WaitGroup
	wg.Add(2)

	// Synchronization channel to ensure proper closure
	done := make(chan struct{}, 2)

	// Test the function
	go func() {
		defer wg.Done()
		receivedData := waitForFastest(&wg, channel)

		// Assert that the received data matches the sent data
		if string(receivedData) != string(testData) {
			t.Errorf("Expected %s, but got %s", string(testData), string(receivedData))
		}

		// Signal that the goroutine is done
		done <- struct{}{}
	}()

	// Simulate another asynchronous task
	go func() {
		defer wg.Done()
		// Simulate some time-consuming task
		time.Sleep(100 * time.Millisecond)
		channel <- []byte("slow data")

		// Signal that the goroutine is done
		done <- struct{}{}
	}()

	// Wait for both goroutines to finish
	wg.Wait()

	// Close the synchronization channel after all goroutines are done
	close(done)

	// Ensure that the channel is closed properly
	for range done {
		// Channel is properly closed
	}

	// Ensure that all goroutines have completed before exiting the test case
}
