package kademlia

import (
	"crypto/sha1"
	"math/big"
)

type Kademlia struct {
	ID        *big.Int
	DataStore map[string]string
}

/*
 * Looks up stored data
 *
 * Takes:   the key for the object
 * Returns: the object and a bool if the retreval was successful or not
 */
func (kademlia *Kademlia) LookUpData(key string) (string, bool) {
	keyHash := hashKey(key)
	closestNode := kademlia.lookUpContact(keyHash)
	value, ok := closestNode.DataStore[keyHash.String()]
	return value, ok
}

/*
 * Stores given data
 *
 * Takes:   the key for the object and the value
 */
func (kademlia *Kademlia) StoreData(key, value string) {
	keyHash := hashKey(key)
	closestNode := kademlia.lookUpContact(keyHash)
	closestNode.DataStore[keyHash.String()] = value
}

/////////////////////////
/// Private Functions ///
/////////////////////////

func hashKey(key string) *big.Int {
	hash := sha1.Sum([]byte(key))
	return new(big.Int).SetBytes(hash[:])
}

func (kademlia *Kademlia) lookUpContact(key *big.Int) *Kademlia {
	// TODO: Network operation to find k closest to given key.
	// ATM just return current node.
	return kademlia
}
