package kademlia

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewCli(t *testing.T) {
	lCli := NewCli(&Kademlia{})
	var lCliType interface{} = lCli

	// test if func return a cli
	_, ok := lCliType.(*cli)

	if !ok {
		t.Fatalf("NewCli() does not return cli of type 'Cli'")
	}
}

func TestHandleInput(t *testing.T) {
	cli := NewCli(&Kademlia{})
	err := cli.HandleInput("not an input\n", "")

	errStr := err.Error()

	if errStr != "CLI Error: Disallowed input" {
		t.Fatalf("This input should not be allowed!")
	}

	err = cli.HandleInput("not an input", "not an input")

	errStr = err.Error()

	if errStr != "CLI Error: Disallowed input" {
		t.Fatalf("This input should not be allowed!")
	}
}

func TestProcessInput(t *testing.T) {
	cli := NewCli(&Kademlia{})
	err := cli.processInput("")

	errStr := err.Error()

	if errStr != "CLI Error: No command entered" {
		t.Fatalf("No error returned when empty string was processed!")
	}

	err = cli.processInput("put asdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasdasda")

	errStr = err.Error()

	if errStr != "CLI Error: Invalid put command. Data longer than 255 characters" {
		t.Fatalf("No error returned for 'put' with string over 255 characters!")
	}

	err = cli.processInput("put")

	errStr = err.Error()

	if errStr != "CLI Error: Invalid put command. No data provided" {
		t.Fatalf("No error returned for 'put' with no object provided!")
	}

	err = cli.processInput("get")

	errStr = err.Error()

	if errStr != "CLI Error: Invalid get command. Only provide the hash of the file after 'get'" {
		t.Fatalf("No error returned for 'get' when something other than a hash was provided")
	}

	err = cli.processInput("show asdasd")

	errStr = err.Error()

	err2 := cli.processInput("exit asdasd")

	errStr2 := err2.Error()

	if errStr != errStr2 || errStr != "CLI Error: Invalid 'show' or 'exit' command. There should be no characters after the 'show' or 'exit' command" {
		t.Fatalf("No error returned for 'show' or 'exit' when extra data was provided!")
	}

	err = cli.processInput("nonsense")

	errStr = err.Error()

	if errStr != "CLI Error: Invalid command. Must start with 'put', 'get', 'show' or 'exit'" {
		t.Fatalf("No error was returned for an CLI-input that does not exist!")
	}
}

func TestShow(t *testing.T) {
	var pingTest = func(_ *Contact, out chan Message) {
		m := Message{
			MsgType: "PONG",
		}
		out <- m
	}

	var details []Detail = []Detail{
		GetContactDetails("0000000000000000000000000000000000000001", "localhost:8000"),
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

	k.Rt.AddContact(contacts[1], pingTest)
	k.Rt.AddContact(contacts[2], pingTest)
	k.Rt.AddContact(contacts[3], pingTest)
	k.Rt.AddContact(contacts[4], pingTest)
	k.Rt.AddContact(contacts[5], pingTest)

	cli := NewCli(k)

	fmt.Println(cli.Show())

	var expectedStrings = []string{"Content in bucket 0\n  ffffffffffffffffffffffffffffffffffffffff\n  ff11111300000000000000000000000000000000", "Content in bucket 3\n  1111111300000000000000000000000000000000\n  1111111200000000000000000000000000000000\n  1111111100000000000000000000000000000000"}

	for _, s := range expectedStrings {
		if !strings.Contains(cli.Show(), s) {
			t.Fatalf("Error in Show! The returned string does not containe the expected values.")
		}
	}
}

func TestPut(t *testing.T) {

}

func TestGet(t *testing.T) {

}

func TestExit(t *testing.T) {

}
