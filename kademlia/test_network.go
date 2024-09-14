package kademlia

import (
	"net"
)

// send dummy message to ip
func TestSend(ip string) { //TODO: add assertions
	udpAddr, err := net.ResolveUDPAddr("udp", ":1234")
	if err != nil {
		panic(err)
	}

	c := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), ip+":1234")

	n := Network{ListenAddr: udpAddr, PacketSize: 512}

	n.SendMessage(&c, Message{"Ping", "Test body", "123", []Contact{c}})
}

// send ping message to ip
func TestSendPing(ip string) {
	udpAddr, err := net.ResolveUDPAddr("udp", ":1234")
	if err != nil {
		panic(err)
	}

	c := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), ip+":1234")

	n := Network{ListenAddr: udpAddr, PacketSize: 512}

	n.SendPingMessage(&c)
}

func TestListen() { //TODO: add assertions
	udpAddr, err := net.ResolveUDPAddr("udp", ":1234")
	if err != nil {
		panic(err)
	}

	n := Network{ListenAddr: udpAddr, PacketSize: 512}

	n.Listen()
}
