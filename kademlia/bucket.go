package kademlia

import (
	"container/list"
)

// bucket definition
// contains a List
type bucket struct {
	list *list.List
}

// newBucket returns a new instance of a bucket
func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *bucket) AddContact(contact Contact, ping func(*Contact, chan Message)) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() { // check if the contact is in the bucket
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}
	if element == nil { // if contact not in bucket
		if bucket.list.Len() < bucketSize { // and bucket not full
			bucket.list.PushFront(contact) // add new contact to head
		} else { // if bucket is full and contact not in bucket
			oldestElement := bucket.list.Back()
			oldest := oldestElement.Value.(Contact)
			responseCh := make(chan Message)
			go ping(&oldest, responseCh)
			response := <-responseCh
			if response.MsgType == "TIMEOUT" { // if oldest node fails to respond
				bucket.list.Remove(oldestElement) // remove it
				bucket.list.PushFront(contact)
			} else {
				bucket.list.MoveToFront(oldestElement)
			}
		}
	} else { // if contact is in bucket
		bucket.list.MoveToFront(element)
	}
}

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

// Len return the size of the bucket
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
