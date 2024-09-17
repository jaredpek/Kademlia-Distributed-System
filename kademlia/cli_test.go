package kademlia

import "testing"

func TestHandleInput(t *testing.T) {
	input := "not an input\n"

	cli := newCli(&Kademlia{})
	cli.handleInput(input, "")
}
