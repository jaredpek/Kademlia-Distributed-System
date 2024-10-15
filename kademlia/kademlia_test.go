package kademlia

import (
	"fmt"
	"testing"
)

func TestLookupContact(t *testing.T) {
	// Sample contacts
	/*var details []Detail = []Detail{
		GetContactDetails("ffffffff00000000000000000000000000000000", "localhost:8000"),
		GetContactDetails("1111111100000000000000000000000000000000", "localhost:8001"),
		GetContactDetails("1111111200000000000000000000000000000000", "localhost:8002"),
		GetContactDetails("1111111300000000000000000000000000000000", "localhost:8003"),
	}

	var quantity int = len(details)
	var contacts []Contact = []Contact{}
	for _, detail := range details {
		contacts = append(contacts, GetContact(detail))
	}

	var k Kademlia = *NewKademlia(contacts[0])
	for _, contact := range contacts {
		k.Rt.AddContact(contact)
	}
	var response []Contact = k.LookupContact(*contacts[1].ID)
	if len(response) != quantity {
		t.Error("[FAIL] Incorrect closest contacts returned")
	}*/
}

func TestNewKademlia(t *testing.T) {
	lKademlia := NewKademlia(Contact{})
	var lKademliaType interface{} = lKademlia

	// test if func return a Kademlia
	_, ok := lKademliaType.(*Kademlia)

	if !ok {
		t.Fatalf("NewKademlia() does not return Kademlia of type 'Kademlia'")
	}
}

func TestUpdateContact(t *testing.T) {
	var target KademliaID = *NewKademliaID("6FFFFFFF00000000000000000000000000000000")
	var closest ContactCandidates
	var contacted map[string]bool = map[string]bool{}
	var contacts []Contact
	var localContacts []Contact = []Contact{
		NewContact(NewKademliaID("1FFFFFFF00000000000000000000000000000000"), "127.0.0.2:1234"),
		NewContact(NewKademliaID("2FFFFFF000000000000000000000000000000000"), "127.0.0.3:1234"),
		NewContact(NewKademliaID("3FFFFFFF00000000000000000000000000000000"), "127.0.0.4:1234"),
		NewContact(NewKademliaID("4FFFFFFF00000000000000000000000000000000"), "127.0.0.5:1234"),
		NewContact(NewKademliaID("5FFFFFFF00000000000000000000000000000000"), "127.0.0.6:1234"),
		NewContact(NewKademliaID("6FFFFFFF00000000000000000000000000000000"), "127.0.0.7:1234"),
		NewContact(NewKademliaID("7FFFFFFF00000000000000000000000000000000"), "127.0.0.8:1234"),
		NewContact(NewKademliaID("8FFFFFFF00000000000000000000000000000000"), "127.0.0.9:1234"),
		NewContact(NewKademliaID("9FFFFFFF00000000000000000000000000000000"), "127.0.0.10:1234"),
		NewContact(NewKademliaID("AFFFFFFF00000000000000000000000000000000"), "127.0.0.11:1234"),
		NewContact(NewKademliaID("BFFFFFFF00000000000000000000000000000000"), "127.0.0.9:1234"),
		NewContact(NewKademliaID("CFFFFFFF00000000000000000000000000000000"), "127.0.0.10:1234"),
		NewContact(NewKademliaID("DFFFFFFF00000000000000000000000000000000"), "127.0.0.11:1234"),
	}
	var otherKademlias []Kademlia = []Kademlia{
		*NewKademlia(localContacts[0]),
		*NewKademlia(localContacts[1]),
		*NewKademlia(localContacts[2]),
		*NewKademlia(localContacts[3]),
	}
	var pingTest = func(_ *Contact, out chan Message) {
		m := Message{
			MsgType: "PONG",
		}
		out <- m
	}
	otherKademlias[0].Rt.AddContact(localContacts[1], pingTest)
	otherKademlias[0].Rt.AddContact(localContacts[2], pingTest)
	otherKademlias[0].Rt.AddContact(localContacts[3], pingTest)
	otherKademlias[1].Rt.AddContact(localContacts[4], pingTest)
	otherKademlias[1].Rt.AddContact(localContacts[5], pingTest)
	otherKademlias[1].Rt.AddContact(localContacts[6], pingTest)
	otherKademlias[2].Rt.AddContact(localContacts[7], pingTest)
	otherKademlias[2].Rt.AddContact(localContacts[8], pingTest)
	otherKademlias[2].Rt.AddContact(localContacts[9], pingTest)
	otherKademlias[3].Rt.AddContact(localContacts[10], pingTest)
	otherKademlias[3].Rt.AddContact(localContacts[11], pingTest)
	otherKademlias[3].Rt.AddContact(localContacts[12], pingTest)
	responses := make(chan Message, 5)

	var me = NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:1234")

	var k Kademlia = *NewKademlia(me)

	k.Rt.AddContact(localContacts[0], pingTest)
	k.Rt.AddContact(localContacts[1], pingTest)
	k.Rt.AddContact(localContacts[2], pingTest)
	k.Rt.AddContact(localContacts[3], pingTest)

	// finds the kademlia with the contact ID and gets the closest contacts to kId
	findFunc := func(kId KademliaID, c *Contact, ch chan Message) {
		for _, i := range otherKademlias {
			if i.Rt.me.ID == c.ID {
				foundContacts := i.Network.Rt.FindClosestContacts(&kId, bucketSize)
				ch <- Message{
					Contacts: foundContacts,
				}
			}
		}
	}

	for _, contact := range k.Rt.FindClosestContacts(&target, bucketSize) {
		// Calculate the distance to the target
		contact.CalcDistance(&target)

		// Create record for new contact
		contacted[contact.Address] = false

		// Add it to the closest list
		closest.Append([]Contact{contact})
	}

	k.updateContacts(&contacted, &closest, &contacts, responses, target, findFunc)

	res := <-responses

	fmt.Println(len(contacts), contacts, res)

	res = <-responses

	fmt.Println(res)

	res = <-responses

	fmt.Println(res)
}
