package kademlia

import (
	"container/list"
	"testing"
)

func TestAddContact(t *testing.T) {
	// Test case 1: Adding a new contact to an empty bucket
	b := &bucket{list: list.New()}
	newContact := NewContact(NewRandomKademliaID(), "0")
	b.AddContact(newContact)

	if b.list.Len() != 1 {
		t.Errorf("Expected bucket length to be 1, got %d", b.list.Len())
	}

	// Test case 2: Adding a contact that already exists in the bucket
	b.AddContact(newContact)

	if b.list.Len() != 1 {
		t.Errorf("Expected bucket length to still be 1, got %d", b.list.Len())
	}

	// Test case 3: Adding contacts until bucket size is reached
	for i := 0; i < bucketSize-1; i++ {
		b.AddContact(NewContact(NewRandomKademliaID(), "address"))
	}

	if b.list.Len() != bucketSize {
		t.Errorf("Expected bucket length to be %d, got %d", bucketSize, b.list.Len())
	}

	// Test case 4: Adding another contact should not exceed bucket size
	b.AddContact(NewContact(NewRandomKademliaID(), "overflow"))
	if b.list.Len() != bucketSize {
		t.Errorf("Expected bucket length to still be %d, got %d", bucketSize, b.list.Len())
	}
}

func TestLenEmptyBucket(t *testing.T) {
	b := newBucket()
	expectedLength := 0
	result := b.Len()

	if result != expectedLength {
		t.Errorf("Expected bucket length to be %d, but got %d", expectedLength, result)
	}
}

func TestLenNonEmptyBucket(t *testing.T) {
	b := newBucket()
	contacts := []Contact{
		NewContact(NewRandomKademliaID(), "1"),
		NewContact(NewRandomKademliaID(), "2"),
		NewContact(NewRandomKademliaID(), "3"),
	}

	for _, contact := range contacts {
		b.AddContact(contact)
	}

	expectedLength := len(contacts)
	result := b.Len()

	if result != expectedLength {
		t.Errorf("Expected bucket length to be %d, but got %d", expectedLength, result)
	}
}
