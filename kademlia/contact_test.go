package kademlia

import (
	"fmt"
	"testing"
)

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
	// Check that String functions as expected
	c := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	s := c.String()
	var sType interface{} = s

	// test if the type is String
	_, ok := sType.(string)

	if !ok {
		t.Fatalf("String() does not return a string")
	}

	// test that string contains expected result
	s2 := `contact("ffffffff00000000000000000000000000000000", "localhost:8000")`

	if s != s2 {
		t.Fatalf("String do not match. %s != %s", s, s2)
	}
}

func TestCandidatesAppend(t *testing.T) {
	// Check that Append functions as expected
	c1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	c2 := NewContact(NewKademliaID("EFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	c3 := NewContact(NewKademliaID("DFFFFFFF00000000000000000000000000000000"), "localhost:8000")

	cList := []Contact{c1, c2, c3}

	var cc ContactCandidates

	cc.Append(cList)

	if cc.Len() != 3 {
		t.Fatalf("ContactCandidates is of the wrong size!")
	}

	if !(cc.contacts[0] == c1 && cc.contacts[1] == c2 && cc.contacts[2] == c3) {
		t.Fatalf("ContactCandidates does not contain the correct values! \nOriginal values = %v,\nSaved values = %v", &cList, cc.GetContacts(3))
	}
}

func TestCandidatesGetContacts(t *testing.T) {
	c1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	c2 := NewContact(NewKademliaID("EFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	c3 := NewContact(NewKademliaID("DFFFFFFF00000000000000000000000000000000"), "localhost:8000")

	cList := []Contact{c1, c2, c3}

	var cc ContactCandidates

	cc.Append(cList)

	res := cc.GetContacts(1)
	res2 := cc.GetContacts(2)
	res3 := cc.GetContacts(3)

	fmt.Printf("%v", res)
	fmt.Printf("%v", res2)
	fmt.Printf("%v", res3)
}

func TestCandidatesSort(t *testing.T) {

}

func TestCandidatesLen(t *testing.T) {

}

func TestCandidatesSwap(t *testing.T) {

}

func TestCandidatesLess(t *testing.T) {

}
