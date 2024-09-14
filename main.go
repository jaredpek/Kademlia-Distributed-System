// TODO: Add package documentation for `main`, like this:
// Package main something something...
package main

import (
	"d7024e/kademlia"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Pretending to run the kademlia app...")
	// Using stuff from the kademlia package here. Something like...
	id := kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	contact := kademlia.NewContact(id, "localhost:8000")
	fmt.Println(contact.String())
	fmt.Printf("%v\n", contact)

	arg := os.Args[1]
	if arg == "listen" {
		fmt.Println("Listening...")
		kademlia.TestListen()
	} else if arg == "send" {
		kademlia.TestSend(os.Args[2])
	} else if arg == "ping" {
		kademlia.TestSendPing(os.Args[2])
	}

	//kademlia.TestListen() //TODO: send and listen on one execution
	//kademlia.TestSend()
}
