// TODO: Add package documentation for `main`, like this:
// Package main something something...
package main

import (
	"d7024e/kademlia"
	"fmt"
	"log"
	"net"
	"os"
)

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

	fmt.Println(GetLocalIP().String())

	arg := os.Args[1]
	if arg == "listen" {
		fmt.Println("Listening...")
		kademlia.TestListen()
	} else if arg == "send" {
		kademlia.TestSend()
	} else if arg == "store" {
		//kademlia.TestLocalStore(os.Args[2]
		fmt.Println("User input:", os.Args[2])
		kademlia.TestStore(os.Args[2], GetLocalIP().String())
	} else if arg == "find" {
		kademlia.TestFindData(kademlia.NewKademliaID(os.Args[2]))
	} else if arg == "d" {
		kademlia.TestDocker()
	} else if arg == "join" {
		fmt.Println("GET HERE1")
		kademlia.TestJoin(GetLocalIP().String())
		fmt.Println("GET HERE2")
	} else if arg == "rest" {
		kademlia.TestRest()
	} else if arg == "cli" {
		ip := GetLocalIP().String()
		c := kademlia.NewContact(kademlia.NewRandomKademliaID(), ip)
		k := kademlia.NewKademlia(c)
		var cli = kademlia.NewCli(k)

		go k.Network.Listen()
		go k.JoinNetwork()

		for {
			fmt.Println("You are currently using the Kademlia CLI!")
			err := cli.UserInput()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}

	/*else if arg == "ping" {
		kademlia.TestSendPing(os.Args[2])
	}*/

	//kademlia.TestListen() //TODO: send and listen on one execution
	//kademlia.TestSend()
}
