package kademlia

import (
	"testing"
)

func TestNewKademliaID(t *testing.T) {
	id := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	var idType interface{} = id

	// test if func return a KademliaID
	_, ok := idType.(*KademliaID)

	if !ok {
		t.Fatalf("NewKademliaID() does not return type 'KademliaID'")
	}

	// test that correct value is saved
	corrVal := "ffffffff00000000000000000000000000000000"
	if id.String() != corrVal {
		t.Fatalf("Value stored in the Kademlia is incorrect. \n%s != %s", corrVal, id.String())
	}
}

// not sure how to test this. function should maybe be rewritten so it takes a seed
func TestNewRandomKademliaID(t *testing.T) {

}

func TestKademliaIdLess(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	id2 := NewKademliaID("FFFFFFFF00000000000000000000000000000001")
	id3 := NewKademliaID("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")

	if id1.Less(id1) { // test equal case
		t.Fatalf("An ID cannot be less than itself!")
	}

	if !id1.Less(id2) { // test more case
		t.Fatalf("Incorrect result from Less!")
	}

	if id1.Less(id3) { // test less case
		t.Fatalf("Incorrect result from Less!")
	}
}

func TestKademliaIdEquals(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	id2 := NewKademliaID("FFFFFFFF00000000000000000000000000000001")

	if !id1.Equals(id1) { // test equal case
		t.Fatalf("Incorrect result from Equal!")
	}

	if id1.Equals(id2) { // test non-equal case
		t.Fatalf("Incorrect result from Equal!")
	}
}

func TestKademliaIdCalcDistance(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	id2 := NewKademliaID("EFFFFFFF00000000000000000000000000000000")

	dist1 := id1.CalcDistance(id2)
	dist2 := id2.CalcDistance(id1)
	dist := NewKademliaID("1000000000000000000000000000000000000000")

	if *dist1 != *dist2 || *dist1 != *dist {
		t.Fatalf("The calculated ID distance is incorrect! \ndist = %d \ndist1 = %d \ndist2 = %d", *dist, *dist1, *dist2)
	}
}

func TestKademliaIdString(t *testing.T) {
	id := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	corr := "ffffffff00000000000000000000000000000000"

	if id.String() != corr {
		t.Fatalf("The returned string is incorrect! \n%s != %s", id, corr)
	}
}
