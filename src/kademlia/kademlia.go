package kademlia

import (
	"d7024e/utils"
)

type Kademlia struct {
	network   *Network
	DataStore map[string]string
}

// Create a new Kademlia instance.
func NewKademlia(network *Network) *Kademlia {
	return &Kademlia{network, make(map[string]string)}
}

// Join the network by pinging the contact node and then performing a node lookup.
func (kademlia *Kademlia) JoinNetwork(contact *Contact) {
	rpcID := NewRandomKademliaID()
	kademlia.network.SendPingPongMessage(contact, rpcID, PING)
	<-kademlia.network.CreateChannel(rpcID)

	utils.Log(1, "My contacts before node lookup: %v", kademlia.network.rt.FindClosestContacts(kademlia.network.rt.me.ID, kademlia.network.k))

	kademlia.LookupContact(kademlia.network.rt.me.ID)
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) {
	utils.Log(1, "Looking up contact %v", target)

	closestContacts, _ := kademlia.nodeLookup(target, FIND_NODE)

	// TODO: Print closest contacts
}

func (kademlia *Kademlia) LookupData(hash string) []byte {
	utils.Log(1, "Looking up data %v", hash)

	closestContacts, data := kademlia.nodeLookup(NewKademliaID(hash), FIND_VALUE)

	// TODO: Store data on the closest contact that didn't return the data, and return data
}

func (kademlia *Kademlia) Store(data []byte) {
	utils.Log(1, "Storing data %v", data)

	closestContacts, _ := kademlia.nodeLookup(NewRandomKademliaID(), STORE)

	// TODO: Store data on closest contacts
}

func (kademlia *Kademlia) nodeLookup(target *KademliaID, opType string) ([]Contact, []byte) {
	// Pick the alpha closest nodes to the target ID from the buckets and add to shortList.
	shortList := ContactCandidates{kademlia.network.rt.FindClosestContacts(target, kademlia.network.alpha)}

	// Create a list of nodes that have already been contacted.
	contactedNodes := ContactCandidates{make([]Contact, 0)}

	// Create a channel to receive the results from the RPCs.
	rpcID := NewRandomKademliaID()
	channel := kademlia.network.CreateChannel(rpcID)

	var data []byte

	for {
		// Filter out nodes that are already in contactedNodes
		alphaNodes := ContactCandidates{make([]Contact, 0)}
		allNodesContacted := true
		for _, node := range shortList.contacts {
			if !Contains(contactedNodes.contacts, node) && !node.ID.Equals(kademlia.network.rt.me.ID) {
				alphaNodes.Append([]Contact{node})
				allNodesContacted = false
			}
		}

		// For FIND_NODE and STORE, terminate when all nodes have been contacted.
		if allNodesContacted && (opType == FIND_NODE || opType == STORE) {
			break
		}

		// Limit the number of nodes in alphaNodes to alpha
		if alphaNodes.Len() > kademlia.network.alpha {
			// Sort inorder to prioritize messaging closer nodes
			alphaNodes.Sort()
			alphaNodes.contacts = alphaNodes.contacts[:kademlia.network.alpha]
		}

		// Send RPCs to alphaNodes
		if opType == FIND_NODE {
			kademlia.network.SendFindContactMessage(target, alphaNodes.contacts, rpcID)
		} else if opType == FIND_VALUE {
			kademlia.network.SendFindDataMessage(alphaNodes.contacts, rpcID, target)
		} else if opType == STORE {
			kademlia.network.SendStoreMessage(alphaNodes.contacts, rpcID, target)
		}
	}

	return shortList.contacts, data
}
