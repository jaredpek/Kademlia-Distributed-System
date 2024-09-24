package kademlia

import "log"

const Alpha = 3

// Default network values
const BootstrapIP = "172.26.0.2:1234"
const ListenPort = "1234"
const PacketSize = 1024

type Kademlia struct {
	Network *Network
	Rt      *RoutingTable
}

func NewKademlia(me Contact) *Kademlia {
	Rt := NewRoutingTable(me)
	return &Kademlia{
		Network: &Network{
			Rt:                Rt,
			BootstrapIP:       BootstrapIP,
			ListenPort:        ListenPort,
			PacketSize:        PacketSize,
			ExpectedResponses: make(map[KademliaID]chan Message, 10),
		},
		Rt: Rt,
	}
}

func (kademlia *Kademlia) LookupContact(target KademliaID) []Contact {
	log.Println("Performing lookup contact")
	var closest ContactCandidates
	var contacted map[string]bool = map[string]bool{}
	responses := make(chan Message, 5)

	// For each contact of the initial k-closest contacts to the target
	for _, contact := range kademlia.Rt.FindClosestContacts(&target, bucketSize) {
		// Calculate the distance to the target
		contact.CalcDistance(&target)

		// Create record for new contact
		contacted[contact.Address] = false

		// Add it to the closest list
		closest.Append([]Contact{contact})
	}

	for {
		var contacts []Contact

		// Sort the contacts by their distance
		closest.Sort()

		// For each contact of the k-closest
		for _, closestContact := range closest.GetContacts(bucketSize) {
			// Continue to the next contact if already contacted
			if contacted[closestContact.Address] {
				continue
			}

			// Send node lookup request to the node and append resulting list of nodes
			go kademlia.Network.SendFindContactMessage(target, &closestContact, responses)
			message := <-responses
			for i := 0; i < len(message.Contacts); i++ {
				message.Contacts[i].CalcDistance(&target)
			}
			log.Println("Got contact response: ", message.Contacts)
			message.Contacts = append(message.Contacts, closestContact)
			closest.Append(message.Contacts)

			contacts = append(contacts, closestContact)

			// Update contact record status
			contacted[closestContact.Address] = true

			// If it has reached alpha contacts then finish
			if len(contacts) == Alpha {
				break
			}
		}

		// If there are no k closest contacts that are uncontacted, return k closest contacts
		if len(contacts) == 0 {
			return closest.GetContacts(bucketSize)
		}
	}
}

func (kademlia *Kademlia) JoinNetwork() {
	log.Println("Joining network")
	response := make(chan Message)
	go kademlia.Network.SendPingMessage(&Contact{Address: kademlia.Network.BootstrapIP}, response) // ping bootstrap node so that it is added to routing table
	r := <-response
	log.Println(r.MsgType)                     // wait for response
	kademlia.LookupContact(*kademlia.Rt.me.ID) // lookup on this node to add close nodes to routing table
}

// should return a string with the result. if the data could be found a string with the data and node it
// was retrived from should be returned. otherwise just return that the file could not be found
func (kademlia *Kademlia) LookupData(hash string) string {
	// TODO
	// return "The requested file could not be downloaded"
	panic("LookupData not implemented")
}

// should return the hash of the data if it was successfully uploaded.
// an error should be returned if the data could not be uploaded
func (kademlia *Kademlia) Store(data []byte) (error, string) {
	// TODO
	panic("Store not implemented")
}
