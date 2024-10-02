package kademlia

import "testing"

func TestNewRest(t *testing.T) {
	var lKademlia = NewKademlia(NewContact(NewRandomKademliaID(), ""))

	var lRest = newRest(lKademlia)
	var lRestType interface{} = lRest

	// test if func return an instance of Rest
	_, ok := lRestType.(*Rest)

	if !ok {
		t.Fatalf("newRest() does not return value of type 'Rest'")
	}
}

func TestGetObject(t *testing.T) {
	// create kademlia environment for test

	// save data and see if it can be get through REST

}

func TestCreateObject(t *testing.T) {
	// create kademlia environment for test
	var lKademlia1 = NewKademlia(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:8000"))
	lKademlia1.Rt.AddContact(NewContact(NewKademliaID("1FFFFFFF00000000000000000000000000000000"), "127.0.0.1:8001"))
	lKademlia1.Rt.AddContact(NewContact(NewKademliaID("2FFFFFFF00000000000000000000000000000000"), "127.0.0.1:8002"))
	// var lKademlia2 = NewKademlia(NewContact(NewKademliaID("1FFFFFFF00000000000000000000000000000000"), "127.0.0.1:8001"))
	// var lKademlia3 = NewKademlia(NewContact(NewKademliaID("2FFFFFFF00000000000000000000000000000000"), "127.0.0.1:8002"))

	// var lRest = newRest(lKademlia1)

	// save some data on node and see if it can be found

	// save some data on another node and see if it can be found

	// try and find some data that does not exist
}
