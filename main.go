// TODO: Add package documentation for `main`, like this:
// Package main something something...
package main

import (
	"d7024e/kademlia"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func printer(k kademlia.Kademlia) {
	for {
		log.Println(k.Rt.FindClosestContacts(kademlia.NewRandomKademliaID(), 30))
		time.Sleep(5 * time.Second)
	}
}

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}

func main() {
	fmt.Println("Pretending to run the kademlia app...")
	// Using stuff from the kademlia package here. Something like...
	id := kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	contact := kademlia.NewContact(id, "localhost:8000")
	fmt.Println(contact.String())
	fmt.Printf("%v\n", contact)

	fmt.Println(GetLocalIP())

	arg := os.Args[1]
	if arg == "listen" {
		fmt.Println("Listening...")
		kademlia.TestListen()
	} else if arg == "send" {
		kademlia.TestSend()
	} else if arg == "store" {
		kademlia.TestStore(os.Args[2])
	} else if arg == "find" {
		kademlia.TestFindData(kademlia.NewKademliaID(os.Args[2]))
	} else if arg == "t" {
		kademlia.TestSendListenLocal1()
	}

	/*else if arg == "ping" {
		kademlia.TestSendPing(os.Args[2])
	}*/

	//kademlia.TestListen() //TODO: send and listen on one execution
	//kademlia.TestSend()
}
