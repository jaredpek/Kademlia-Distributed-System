package kademlia

import (
	"encoding/hex"
	"testing"
)

var pingTest = func(_ *Contact, out chan Message) {
	m := Message{
		MsgType: "PONG",
	}
	out <- m
}

var pingTestTimeout = func(_ *Contact, out chan Message) {
	m := Message{
		MsgType: "TIMEOUT",
	}
	out <- m
}

func TestNewBucket(t *testing.T) {
	var lBucket = newBucket()
	var lBucketType interface{} = lBucket

	// test if func return a bucket
	_, ok := lBucketType.(*bucket)

	if !ok {
		t.Fatalf("newBucket() does not return bucket of type 'bucket'")
	}

	// test if bucket is empty
	var l = lBucket.list.Len()
	if l != 0 {
		t.Fatalf("The bucket is not empty, it contains %d enteries", l)
	}
}

func TestAddContact(t *testing.T) {
	// function for checking if a contact can be found in front
	checkFirstElement := func(t *testing.T, contact Contact, bucket *bucket) {
		var fElement = bucket.list.Front().Value
		if fElement != contact {
			t.Fatalf("The contact (%s) != from first element (%s)", &contact, fElement)
		}
	}

	var lBucket = newBucket()

	// add element to bucket and check that it can be found in the front
	id := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	var contact = NewContact(id, "localhost:8000")
	lBucket.AddContact(contact, pingTest)

	checkFirstElement(t, contact, lBucket)

	// add element that is already in bucket and check that it is moved to the front
	lBucket.AddContact(NewContact(NewKademliaID("1FFFFFFF00000000000000000000000000000000"), "localhost:8000"), pingTest)
	lBucket.AddContact(contact, pingTest)

	checkFirstElement(t, contact, lBucket)

	// check that a new element can not be added to bucket if bucket is already full
	var s, _ = hex.DecodeString("FFFFFFFFFFFFFFFF000000000000000000000000")
	for i := 0; i < bucketSize-2; i++ {
		s[0] -= 1
		contact = NewContact(NewKademliaID(hex.EncodeToString(s)), "localhost:8000")
		lBucket.AddContact(contact, pingTest)
	}

	var l = lBucket.Len()
	if l != bucketSize {
		t.Fatalf("The bucket is not full! It only contains %d elements.", l)
	}

	// check that element is removed from bucket if timeout happens
	contact2 := NewContact(NewKademliaID("11FFFFFF00000000000000000000000000000000"), "localhost:8000")
	lBucket.AddContact(contact2, pingTestTimeout)

	checkFirstElement(t, contact2, lBucket)
}

func TestGetContactAndCalcDistance(t *testing.T) {
	var lBucket = newBucket()

	lBucket.AddContact(NewContact(NewKademliaID("1FFFFFFF00000000000000000000000000000000"), "localhost:8000"), pingTest)
	lBucket.AddContact(NewContact(NewKademliaID("2FFFFFFF00000000000000000000000000000000"), "localhost:8000"), pingTest)
	lBucket.AddContact(NewContact(NewKademliaID("3FFFFFFF00000000000000000000000000000000"), "localhost:8000"), pingTest)

	var contacts = lBucket.GetContactAndCalcDistance(NewKademliaID("FFFFFFFF00000000000000000000000000000000"))

	// test that the calculated distance if correct
	if !(contacts[2].distance.String() == "e000000000000000000000000000000000000000" &&
		contacts[1].distance.String() == "d000000000000000000000000000000000000000" &&
		contacts[0].distance.String() == "c000000000000000000000000000000000000000") {
		t.Fatalf("The calculated distances are incorrect! \n%s \n%s \n%s", contacts[0].distance.String(), contacts[1].distance.String(), contacts[2].distance.String())
	}
}

func TestLen(t *testing.T) {
	var lBucket = newBucket()

	lBucket.AddContact(NewContact(NewKademliaID("1FFFFFFF00000000000000000000000000000000"), "localhost:8000"), pingTest)
	lBucket.AddContact(NewContact(NewKademliaID("2FFFFFFF00000000000000000000000000000000"), "localhost:8000"), pingTest)
	lBucket.AddContact(NewContact(NewKademliaID("3FFFFFFF00000000000000000000000000000000"), "localhost:8000"), pingTest)

	bucketLen := lBucket.Len()

	if bucketLen != 3 {
		t.Fatalf("The returned length is incorrect. lBucket.len() = %d != 3", bucketLen)
	}
}
