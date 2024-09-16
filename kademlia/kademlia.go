package kademlia

const Alpha = 3

type Kademlia struct {
	Network Network
	Rt      RoutingTable
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) []Contact {
	var closest ContactCandidates
	var contacted map[string]bool = map[string]bool{}

	// For each contact of the initial k-closest contacts to the target
	for _, contact := range kademlia.Rt.FindClosestContacts(target, bucketSize) {
		// Calculate the distance to the target
		contact.CalcDistance(target)

		// Create record for new contact
		contacted[contact.Address] = false

		// Add it to the closest list
		closest.Append([]Contact{contact})
	}

	for {
		var ids []KademliaID

		// Sort the contacts by their distance
		closest.Sort()

		// For each contact of the k-closest
		for _, contact := range closest.GetContacts(bucketSize) {
			// Continue to the next contact if already contacted
			if contacted[contact.Address] {
				continue
			}

			// Send node lookup request to the node and append resulting list of nodes
			ids = append(ids, kademlia.Network.SendFindContactMessage(target, &contact)...)

			// Update contact record status
			contacted[contact.Address] = true

			// If it has reached alpha contacts then finish
			if len(ids) == Alpha {
				break
			}
		}

		// If there are no k closest contacts that are uncontacted, return k closest contacts
		if len(ids) == 0 {
			return closest.GetContacts(bucketSize)
		}
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
