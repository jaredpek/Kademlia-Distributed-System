package kademlia

import (
	"testing"
)

func TestHandleInput(t *testing.T) {
	cli := newCli(&Kademlia{})
	err := cli.HandleInput("not an input\n", "")

	errStr := err.Error()

	if errStr != "CLI error: disallowed input" {
		t.Fatalf("This input should not be allowed!")
	}

	// test when implemented
	/*pass := cli.handleInput("get", "asdadfpok")

	pass2 := cli.handleInput("put", "asdasdas")*/
}

func TestShow(t *testing.T) {
	var details []Detail = []Detail{
		GetContactDetails("0000000000000000000000000000000000000001", "localhost:8000"),
		GetContactDetails("0000000000000000000000000000000000000002", "localhost:8000"),
		GetContactDetails("1111111100000000000000000000000000000000", "localhost:8000"),
		GetContactDetails("1111111200000000000000000000000000000000", "localhost:8000"),
		GetContactDetails("1111111300000000000000000000000000000000", "localhost:8000"),
		GetContactDetails("ff11111300000000000000000000000000000000", "localhost:8000"),
		GetContactDetails("ffffffffffffffffffffffffffffffffffffffff", "localhost:8000"),
	}

	var contacts []Contact = []Contact{}
	for _, detail := range details {
		contacts = append(contacts, GetContact(detail))
	}

	k := NewKademlia(contacts[0])

	k.Rt.AddContact(contacts[1])
	k.Rt.AddContact(contacts[2])
	k.Rt.AddContact(contacts[3])
	k.Rt.AddContact(contacts[4])
	k.Rt.AddContact(contacts[5])
	k.Rt.AddContact(contacts[6])

	cli := newCli(k)

	cli.Show()
	// need to add assertion
}

func TestPut(t *testing.T) {

}

func TestGet(t *testing.T) {

}

func TestExit(t *testing.T) {

}
