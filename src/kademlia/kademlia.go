package kademlia

type Kademlia struct {
	network *Network
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	kademlia.network.SendFindContactMessage(target, NewRandomKademliaID())
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
