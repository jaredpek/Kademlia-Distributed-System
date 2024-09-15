package kademlia

import (
	"testing"
)

type Detail struct {
	id string;
	addr string;
}

func GetContactDetails(id string, addr string) Detail {
	return Detail{id, addr}
}

func GetContact(detail Detail) Contact {
	return NewContact(NewKademliaID(detail.id), detail.addr)
}

func TestRoutingTable(t *testing.T) {
	// Sample contacts
	var details []Detail = []Detail{
		GetContactDetails("ffffffff00000000000000000000000000000000", "localhost:8000"), // 1101 1110 1110 1110 1110 1110 1110 1011 -> 6
		GetContactDetails("1111111100000000000000000000000000000000", "localhost:8001"), // 0011 0000 0000 0000 0000 0000 0000 0101 -> 3
		GetContactDetails("1111111200000000000000000000000000000000", "localhost:8002"), // 0011 0000 0000 0000 0000 0000 0000 0110 -> 4
		GetContactDetails("1111111300000000000000000000000000000000", "localhost:8003"), // 0011 0000 0000 0000 0000 0000 0000 0111 -> 5
		GetContactDetails("1111111400000000000000000000000000000000", "localhost:8004"), // 0011 0000 0000 0000 0000 0000 0000 0000 -> 2
		GetContactDetails("2111111400000000000000000000000000000000", "localhost:8005"), // 0000 0000 0000 0000 0000 0000 0000 0000 -> 1 (target)
	}

	var quantity = len(details)
	var contacts []Contact = []Contact{}
	for _, detail := range details {
		contacts = append(contacts, GetContact(detail))
	}

	// Add dummy contacts that should not be added to table
	contacts = append(contacts, contacts...)

	// Test routing table creation
	table := NewRoutingTable(contacts[0])

	// Test routing table population
	for _, contact := range contacts {
		table.AddContact(contact)
	}
	added := table.FindClosestContacts(contacts[0].ID, bucketSize)
	if len(added) < quantity {
		t.Error("[FAIL] Incorrect number of contacts added")
	}

	// Test routing table closest contacts search that should be in order of [closest -> ... -> furthest]
	closest := table.FindClosestContacts(contacts[5].ID, bucketSize)
	if 
		closest[0].ID != contacts[5].ID ||
		closest[1].ID != contacts[4].ID ||
		closest[2].ID != contacts[1].ID ||
		closest[3].ID != contacts[2].ID ||
		closest[4].ID != contacts[3].ID ||
		closest[5].ID != contacts[0].ID {
			t.Error("[FAIL] Incorrect closest contacts found")
	}
}
