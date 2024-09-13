package kademlia

import (
	"reflect"
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

	if !(res[0] == c1 && reflect.DeepEqual(res2, []Contact{c1, c2}) && reflect.DeepEqual(res3, []Contact{c1, c2, c3})) {
		t.Fatalf("The GetContacts do not return the expected values! \nres: %v \nres2: %v \nres3: %v \nc1: %v \nc2: %v \n c3: %v", res, res2, res3, c1, c2, c3)
	}
}

func TestCandidatesSort(t *testing.T) {
	c1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	c2 := NewContact(NewKademliaID("EFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	c3 := NewContact(NewKademliaID("DFFFFFFF00000000000000000000000000000000"), "localhost:8000")

	// The sorting is done based on distance. Since c2 is further away the final order should be: 1 3 2
	c1.distance = NewKademliaID("1000000000000000000000000000000000000000")
	c2.distance = NewKademliaID("3000000000000000000000000000000000000000")
	c3.distance = NewKademliaID("2000000000000000000000000000000000000000")

	cList := []Contact{c1, c2, c3}

	var cc ContactCandidates

	cc.Append(cList)

	cc.Sort()

	if cc.contacts[1] != c3 || cc.contacts[2] != c2 {
		t.Fatalf("ContactCandidates where not sorted! \n%v", cc.GetContacts(3))
	}
}

func TestCandidatesLen(t *testing.T) {
	c1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	c2 := NewContact(NewKademliaID("EFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	c3 := NewContact(NewKademliaID("DFFFFFFF00000000000000000000000000000000"), "localhost:8000")

	cList := []Contact{c1, c2, c3}

	var cc ContactCandidates

	cc.Append(cList)

	if cc.Len() != 3 {
		t.Fatalf("ContactCandidates is of the wrong size! \ncc: %v", cc.GetContacts(3))
	}
}

func TestCandidatesSwap(t *testing.T) {
	c1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	c2 := NewContact(NewKademliaID("EFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	c3 := NewContact(NewKademliaID("DFFFFFFF00000000000000000000000000000000"), "localhost:8000")

	cList := []Contact{c1, c2, c3}

	var cc ContactCandidates

	cc.Append(cList)

	cc.Swap(0, 2)

	if !(cc.contacts[0] == c3 && cc.contacts[2] == c1) {
		t.Fatalf("The candidates were not swapped! \ncc: %v", cc)
	}
}

func TestCandidatesLess(t *testing.T) {
	c1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	c2 := NewContact(NewKademliaID("EFFFFFFF00000000000000000000000000000000"), "localhost:8000")

	c1.distance = NewKademliaID("1000000000000000000000000000000000000000")
	c2.distance = NewKademliaID("2000000000000000000000000000000000000000")

	cList := []Contact{c1, c2}

	var cc ContactCandidates

	cc.Append(cList)

	if !cc.Less(0, 1) {
		t.Fatalf("CandidatesLess does not work as intended!")
	}
}
