package kademlia

import (
	"d7024e/utils"
	"strings"
	"sync"
	"time"
)

var shortListMutex = &sync.RWMutex{}
var respondedNodesMutex = &sync.RWMutex{}

type Kademlia struct {
	network       *Network
	DataStore     map[string]string
	ClosestPeers  map[string][]Contact
	RefreshTicker *time.Ticker
}

// Create a new Kademlia instance.
func NewKademlia(network *Network) *Kademlia {
	return &Kademlia{network, make(map[string]string), make(map[string][]Contact), time.NewTicker(network.refreshInterval)}
}

// Join the network by pinging the contact node and then performing a node lookup.
func (kademlia *Kademlia) JoinNetwork(contact *Contact) {
	rpcID := NewRandomKademliaID()
	kademlia.network.SendPingMessage(contact, rpcID)
	<-kademlia.network.CreateChannel(rpcID) // TODO: Listen with timeout

	utils.Log(1, "My routing table before node lookup:")
	for _, contact := range kademlia.network.rt.FindClosestContacts(kademlia.network.rt.me.ID, kademlia.network.k) {
		utils.Log(1, "%v, %v", contact.Address, contact.ID)
	}

	kademlia.LookupContact(kademlia.network.rt.me.ID)

	utils.Log(1, "My routing table after node lookup:")
	for _, contact := range kademlia.network.rt.FindClosestContacts(kademlia.network.rt.me.ID, kademlia.network.k) {
		utils.Log(1, "%v, %v", contact.Address, contact.ID)
	}
}

// Lookup a contact by performing a node lookup.
func (kademlia *Kademlia) LookupContact(target *KademliaID) {
	utils.Log(1, "Looking up contact %v", target)

	closestContacts, _ := kademlia.nodeLookup(target, FIND_NODE)

	utils.Log(1, "Closest contacts found to %v after node lookup:", target)
	for _, contact := range closestContacts {
		utils.Log(1, "%v, %v", contact.Address, contact.ID)
	}
}

// Lookup data on the network by performing a node lookup. Returns the data.
func (kademlia *Kademlia) LookupData(hash string) []byte {
	utils.Log(1, "Looking up data for hash %v", hash)

	closestContactsWithoutValue, dataResult := kademlia.nodeLookup(NewKademliaID(hash), FIND_VALUE)

	if dataResult != nil && len(closestContactsWithoutValue) > 0 {
		utils.Log(1, "Closest contacts found to %v after data lookup that didn't return value:", hash)
		for _, contact := range closestContactsWithoutValue {
			utils.Log(1, "%v, %v", contact.Address, contact.ID)
		}

		// Store data on closest contact that didn't return the value (cache it)
		utils.Log(1, "Storing data %v on closest contact that didn't return the value", dataResult)
		utils.Log(1, "%v, %v", closestContactsWithoutValue[0].Address, closestContactsWithoutValue[0].ID)
		kademlia.network.SendStoreMessage(NewKademliaID(hash), dataResult, &closestContactsWithoutValue[0], NewRandomKademliaID())
	}

	return dataResult
}

// Store data on the network by performing a node lookup and then storing the data on the closest contacts. Returns the hash of the data.
func (kademlia *Kademlia) Store(data []byte) string {
	utils.Log(1, "Storing data %v", data)

	hash := utils.Hash(data)
	key := NewKademliaID(hash)
	closestContacts, _ := kademlia.nodeLookup(key, STORE)

	// Store data on closest contacts
	utils.Log(1, "Closest contacts found to %v to store data at:", key)
	for _, contact := range closestContacts {
		kademlia.network.SendStoreMessage(key, data, &contact, NewRandomKademliaID())
		utils.Log(1, "%v, %v", contact.Address, contact.ID)
	}

	// Save closestContacts for this hash
	kademlia.ClosestPeers[key.String()] = closestContacts

	return hash
}

// Start refresh routine for refreshing the closest peers to stored values
func (kademlia *Kademlia) StartRefreshRoutine() {
	go func() {
		for range kademlia.RefreshTicker.C {
			kademlia.refreshClosestPeers()
		}
	}()
}

// Refresh the closest peers to stored values.
func (kademlia *Kademlia) refreshClosestPeers() {
	for hash, contacts := range kademlia.ClosestPeers {
		for _, contact := range contacts {
			go func(hash string, contact Contact) {
				// Use a goroutine to prevent blocking the loop
				// Implement SendRefreshMessage asynchronously
				kademlia.network.SendRefreshMessage(NewKademliaID(hash), &contact, NewRandomKademliaID())
			}(hash, contact)
		}
	}
}

func (kademlia *Kademlia) Forget(hash string) {
	key := NewKademliaID(hash)

	if _, ok := kademlia.ClosestPeers[key.String()]; !ok {
		utils.Log(1, "No closest peers found for hash %s", hash)
		return
	}

	utils.Log(1, "Forgetting hash %s", hash)
	delete(kademlia.ClosestPeers, key.String())
}

// Perform a node lookup on the network.
func (kademlia *Kademlia) nodeLookup(target *KademliaID, opType string) ([]Contact, []byte) {

	var data []byte
	var closerFound chan bool
	dataFound := make(chan []byte, kademlia.network.alpha)

	// Pick the alpha closest nodes to the target ID from the buckets and add to shortList.
	shortList := ContactCandidates{kademlia.network.rt.FindClosestContacts(target, kademlia.network.alpha)}

	utils.Log(1, "Shortlist at start of nodeLookup:")
	for _, contact := range shortList.contacts {
		utils.Log(1, "%v, %v", contact.Address, contact.ID)
	}

	// Create a list of nodes that have already been contacted.
	contactedNodes := ContactCandidates{make([]Contact, 0)}

	// Create a list of nodes that have responded but not returned a value.
	respondedNodesWithoutValue := ContactCandidates{make([]Contact, 0)}

OUTER_LOOP:
	for {
		var iterativeSync sync.WaitGroup

		closerFound = make(chan bool, kademlia.network.alpha)

		// Filter out nodes that are already in contactedNodes
		alphaNodes := ContactCandidates{make([]Contact, 0)}
		allNodesContacted := true
		shortListMutex.RLock()
		for _, node := range shortList.contacts {
			if !Contains(contactedNodes.contacts, node) && !node.ID.Equals(kademlia.network.rt.me.ID) {
				alphaNodes.Append([]Contact{node})
				allNodesContacted = false
			}
		}
		shortListMutex.RUnlock()

		// Terminate when all nodes have been contacted.
		if allNodesContacted {
			break
		}

		// Limit the number of nodes in alphaNodes to alpha
		if alphaNodes.Len() > kademlia.network.alpha {
			// Sort inorder to prioritize messaging closer nodes
			alphaNodes.Sort()
			alphaNodes.contacts = alphaNodes.contacts[:kademlia.network.alpha]
		}

		// Send RPCs to alphaNodes
		for _, node := range alphaNodes.contacts {
			rpcID := NewRandomKademliaID()

			if opType == FIND_NODE || opType == STORE {
				kademlia.network.SendFindContactMessage(target, &node, rpcID)
			} else {
				kademlia.network.SendFindDataMessage(target, &node, rpcID)
			}

			go kademlia.waitForResponse(&iterativeSync, closerFound, dataFound, rpcID, &shortList, &respondedNodesWithoutValue, &node)
		}
		contactedNodes.Append(alphaNodes.contacts)

		// Loose parallelism
		time.Sleep(5 * time.Second)

		// Wait for a response from the alpha closest nodes in state.
		data = waitForFastest(&iterativeSync, dataFound)
		if data != nil {
			utils.Log(1, "Recieved FIND_VALUE_RESPONSE from node, Im done searching")
			respondedNodesMutex.Lock()
			respondedNodesWithoutValue.Sort()
			shortListMutex.Lock()
			shortList = ContactCandidates{respondedNodesWithoutValue.contacts}
			shortListMutex.Unlock()
			respondedNodesMutex.Unlock()
			break OUTER_LOOP
		}

		// If a closer node was found, set a flag that tells us to keep iterating
		closerFoundFlag := false
		for value := range closerFound {
			if value {
				closerFoundFlag = true
				break
			}
		}

		// Otherwise, send a find node message to the k closest nodes in state that have not been contacted yet.
		if !closerFoundFlag {
			var iterativeSync sync.WaitGroup
			closerFound = make(chan bool, kademlia.network.k)
			for _, node := range shortList.contacts {
				rpcID := NewRandomKademliaID()

				if !Contains(contactedNodes.contacts, node) {
					if opType == FIND_NODE || opType == STORE {
						kademlia.network.SendFindContactMessage(target, &node, rpcID)
					} else {
						kademlia.network.SendFindDataMessage(target, &node, rpcID)
					}
					contactedNodes.Append([]Contact{node})
					go kademlia.waitForResponse(&iterativeSync, closerFound, dataFound, rpcID, &shortList, &respondedNodesWithoutValue, &node)
				}

			}

			// Wait for a response from the k closest nodes in state.
			data = waitForFastest(&iterativeSync, dataFound)
			if data != nil {
				utils.Log(1, "Recieved FIND_VALUE_RESPONSE from node, Im done searching")
				respondedNodesMutex.Lock()
				respondedNodesWithoutValue.Sort()
				shortListMutex.Lock()
				shortList = ContactCandidates{respondedNodesWithoutValue.contacts}
				shortListMutex.Unlock()
				respondedNodesMutex.Unlock()
				break OUTER_LOOP
			}
		}
	}

	return shortList.contacts, data
}

// Wait for a response from a node. If no response is received within 10 seconds, remove the node from the short list.
func (kademlia *Kademlia) waitForResponse(iterWait *sync.WaitGroup, status chan bool, data chan []byte, rpcID *KademliaID, shortList *ContactCandidates, respondedNodesWithoutValue *ContactCandidates, node *Contact) {

	iterWait.Add(1)
	// Wait for 10 sec if no response remove node from short list
	response, err := kademlia.network.ListenWithTimeout(rpcID, 10)
	if err != nil {
		shortListMutex.Lock()
		shortList.RemoveContact(node)
		shortListMutex.Unlock()

		status <- false
		iterWait.Done()
		return
	}

	// If response contains stored data, terminate and return it to the caller
	if response["type"] == FIND_VALUE_RESPONSE {
		utils.Log(1, "waitForResponse: got value from node: %s", node.Address)
		data <- []byte(response["data"])
		iterWait.Done()
		return
	}

	// Node responded and it wasn't a FIND_VALUE_RESPONSE, add it to respondedNodesWithoutValue
	respondedNodesMutex.Lock()
	if !Contains(respondedNodesWithoutValue.contacts, *node) && response["type"] != FIND_VALUE_RESPONSE {
		respondedNodesWithoutValue.Append([]Contact{*node})
	}
	respondedNodesMutex.Unlock()

	// Extract nodes from message
	responeContacts := []Contact{}
	strs := strings.Split(response["data"], "\n")
	conStrs := strs[:len(strs)-1] // remove last element (not a contact)

	for _, str := range conStrs {
		contact, err := NewContactFromString(str)
		if err != nil {
			utils.LogError("nodeLookup: could not translate string to contact %s", err)
			continue
		}
		if contact.ID.Equals(kademlia.network.rt.me.ID) {
			utils.Log(1, "nodeLookup: contact %s is me, discard", contact.Address)
			continue
		}
		responeContacts = append(responeContacts, contact)
	}

	// Update shortList
	shortListMutex.Lock()
	nodesReplaced := updateShortList(shortList, responeContacts, kademlia.network.k)
	shortListMutex.Unlock()

	// If at least one node was replaced, send a new find node message to the alpha closest nodes in state.
	if nodesReplaced {
		utils.Log(1, "At least one node was replaced, looping again")

		status <- false
		iterWait.Done()
		return
	}

	status <- true
	iterWait.Done()
}

// Update the short list with new contacts.
func updateShortList(shortList *ContactCandidates, newContacts []Contact, k int) bool {
	// Filter out nodes that are already in shortList
	nodesReplaced := false
	for _, newNode := range newContacts {
		if Contains(shortList.contacts, newNode) {
			continue
		}

		// If shortList contains less than k nodes, nodes are freely inserted into the state until it contains k nodes.
		if shortList.Len() < k {
			shortList.Append([]Contact{newNode})
		} else {
			// shortList already contains k elements, replace the node furthest away with the new node if it is closer than the closest node in shortList
			shortList.Sort()
			if newNode.Less(&shortList.contacts[0]) {
				shortList.contacts[len(shortList.contacts)-1] = newNode
				nodesReplaced = true
			}
		}
	}
	return nodesReplaced
}

// Wait for the fastest response from a node.
func waitForFastest(wg *sync.WaitGroup, ch chan []byte) (data []byte) {
	done := make(chan struct{})
	var lock = &sync.Mutex{}
	closed := false

	go func(group *sync.WaitGroup) {
		group.Wait()
		lock.Lock()
		if closed {
			lock.Unlock()
			return
		}
		close(done)
		closed = true
		lock.Unlock()
	}(wg)

	go func(channel chan []byte) {
		data = <-channel
		lock.Lock()
		if closed {
			lock.Unlock()
			return
		}
		close(done)
		closed = true
		lock.Unlock()
	}(ch)

	<-done
	return data
}
