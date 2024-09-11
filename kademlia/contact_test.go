package kademlia

import "testing"

func TestNewContact(t *testing.T) {
	// Check that contact with correct info is returned
	id := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	address := "localhost:8000"

	c := NewContact(id, address)
	var cType interface{} = c

	// test if the type is Contact
	_, ok := cType.(Contact)

	if !ok {
		t.Fatalf("NewContact() does not return contact of type 'Contact'")
	}

	// test if information is correct
	if c.Address != address || c.ID != id || c.distance != nil {
		t.Fatalf("The information returned is not same as saved information. Saved information: %v", c)
	}
}

func TestCalcDistance(t *testing.T) {
	// Check that distance to a target is correct
	c := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")

	c.CalcDistance(NewKademliaID("1111111100000000000000000000000000000000"))

	cRes := NewKademliaID("eeeeeeee00000000000000000000000000000000")

	if *c.distance != *cRes {
		t.Fatalf("The distance between nodes is incorrect. Should be: %v ; result: %v", cRes, c.distance)
	}
}

func TestLess(t *testing.T) {
	// Check that less functions as expected
	c1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	c1.distance = NewKademliaID("1FFFFFFF00000000000000000000000000000000")
	c2 := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8000")
	c2.distance = NewKademliaID("2FFFFFFF00000000000000000000000000000000")

	if !c1.Less(&c2) {
		t.Fatalf("Less function is wrong! Thinks %v > %v", c1.distance, c2.distance)
	}
}

func TestString(t *testing.T) {
	// Check that string functions as expected

}

func TestCandidatesAppend(t *testing.T) {

}

func TestCandidatesGetContacts(t *testing.T) {

}

func TestCandidatesSort(t *testing.T) {

}

func TestCandidatesLen(t *testing.T) {

}

func TestCandidatesSwap(t *testing.T) {

}

func TestCandidatesLess(t *testing.T) {

}
