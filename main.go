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

var thisIP string = GetLocalIP().String()
var rt *kademlia.RoutingTable = kademlia.NewRoutingTable(kademlia.NewContact(kademlia.NewRandomKademliaID(), thisIP))
var network kademlia.Network = kademlia.Network{
	ListenPort:        "1234",
	PacketSize:        1024 * 4,
	ExpectedResponses: make(map[kademlia.KademliaID]chan kademlia.Message, 10),
	Rt:                rt,
	BootstrapIP:       "172.26.0.2:1234",
	Messenger:         &kademlia.UDPMessenger{Rt: rt},
}
var k kademlia.Kademlia = kademlia.Kademlia{&network, rt}

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
	fmt.Println("This nodes IP: " + GetLocalIP().String())

	arg := os.Args[1]
	if arg == "listen" {
		fmt.Println("Listening...")
		network.Listen()
	} else if arg == "join" {
		go k.JoinNetwork()
		network.Listen()
	} else if arg == "cli" {
		var cli = kademlia.NewCli(&k)

		go network.Listen()
		go k.JoinNetwork()

		for {
			fmt.Println("You are currently using the Kademlia CLI!")
			err := cli.UserInput()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
