package kademlia

import (
	"encoding/hex"
	"testing"
)

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
	lBucket.AddContact(contact)

	checkFirstElement(t, contact, lBucket)

	// add element that is already in bucket and check that it is moved to the front
	lBucket.AddContact(NewContact(NewKademliaID("1FFFFFFF00000000000000000000000000000000"), "localhost:8000"))
	lBucket.AddContact(contact)

	checkFirstElement(t, contact, lBucket)

	// check that a new element can not be added to bucket if bucket is already full
	var s, _ = hex.DecodeString("FFFFFFFFFFFFFFFF000000000000000000000000")
	for i := 0; i < bucketSize-2; i++ {
		s[0] -= 1
		contact = NewContact(NewKademliaID(hex.EncodeToString(s)), "localhost:8000")
		lBucket.AddContact(contact)
	}

	var l = lBucket.Len()
	if l != bucketSize {
		t.Fatalf("The bucket is not full! It only contains %d elements.", l)
	}

	lBucket.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8000"))
	checkFirstElement(t, contact, lBucket)
}

func TestGetContactAndCalcDistance(t *testing.T) {

}
