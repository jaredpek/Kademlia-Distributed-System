package kademlia

import (
	"testing"
)

func TestHandleInput(t *testing.T) {
	cli := newCli(&Kademlia{})
	err := cli.handleInput("not an input\n", "")

	errStr := err.Error()

	if errStr != "CLI error: disallowed input" {
		t.Fatalf("This input should not be allowed!")
	}

	// test when implemented
	/*pass := cli.handleInput("get", "asdadfpok")

	pass2 := cli.handleInput("put", "asdasdas")*/
}

func TestPut(t *testing.T) {

}

func TestGet(t *testing.T) {

}

func TestExit(t *testing.T) {

}
