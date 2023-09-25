package kademlia

import (
	"testing"
)

func TestNewContactFromString(t *testing.T) {
	// Test valid input
	str := "contact(0123456789abcdef0123456789abcdef01234567, 192.168.1.1, 0123456789abcdef0123456789abcdef01234568)"
	contact, err := NewContactFromString(str)
	if err != nil {
		t.Errorf("NewContactFromString() returned an error for valid input: %v", err)
	}
	expectedID := "0123456789abcdef0123456789abcdef01234567"
	expectedAddress := "192.168.1.1"
	expectedDistance := "0123456789abcdef0123456789abcdef01234568"

	if contact.ID.String() != expectedID || contact.Address != expectedAddress || contact.distance.String() != expectedDistance {
		t.Errorf("NewContactFromString() returned incorrect Contact data.\nGot: %s\nExpected: %s", contact.String(), str)
	}

	// Test invalid input
	str = "invalid_contact_format"
	_, err = NewContactFromString(str)
	if err == nil {
		t.Error("NewContactFromString() did not return an error for invalid input")
	}
}

func TestRemoveContact(t *testing.T) {
	// Create a ContactCandidates with some contacts
	contact1 := NewContact(NewRandomKademliaID(), "192.168.1.1")
	contact2 := NewContact(NewRandomKademliaID(), "192.168.1.2")
	contact3 := NewContact(NewRandomKademliaID(), "192.168.1.3")

	candidates := &ContactCandidates{
		contacts: []Contact{contact1, contact2, contact3},
	}

	// Test removing a contact
	candidates.RemoveContact(&contact2)
	if len(candidates.contacts) != 2 {
		t.Error("RemoveContact() did not remove the contact from ContactCandidates")
	}
	if Contains(candidates.contacts, contact2) {
		t.Error("RemoveContact() did not remove the contact from ContactCandidates")
	}

	// Test removing a non-existent contact
	nonExistentContact := NewContact(NewRandomKademliaID(), "192.168.1.4")
	candidates.RemoveContact(&nonExistentContact)
	if len(candidates.contacts) != 2 {
		t.Error("RemoveContact() should not modify ContactCandidates if the contact does not exist")
	}
}

func TestContains(t *testing.T) {
	// Create a list of contacts
	contact1 := NewContact(NewRandomKademliaID(), "192.168.1.1")
	contact2 := NewContact(NewRandomKademliaID(), "192.168.1.2")
	contact3 := NewContact(NewRandomKademliaID(), "192.168.1.3")

	list := []Contact{contact1, contact2, contact3}

	// Test contains
	if !Contains(list, contact1) {
		t.Error("Contains() returned false for an existing contact")
	}

	// Test does not contain
	nonExistentContact := NewContact(NewRandomKademliaID(), "192.168.1.4")
	if Contains(list, nonExistentContact) {
		t.Error("Contains() returned true for a non-existent contact")
	}
}
