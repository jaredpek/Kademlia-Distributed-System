package kademlia

import (
	"net"
)

func TestSend() { //TODO: add assertions
	udpAddr, err := net.ResolveUDPAddr("udp", ":1234")
	if err != nil {
		panic(err)
	}

	c := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:1234")

	n := Network{ListenAddr: udpAddr, PacketSize: 512}

	n.SendMessage(&c, Message{"Ping", "Test body"})
}

func TestListen() { //TODO: add assertions
	udpAddr, err := net.ResolveUDPAddr("udp", ":1234")
	if err != nil {
		panic(err)
	}

	n := Network{ListenAddr: udpAddr, PacketSize: 512}

	n.Listen()
}
