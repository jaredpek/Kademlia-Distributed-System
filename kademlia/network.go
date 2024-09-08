package kademlia

import (
	"fmt"
	"net"
)

type Network struct {
	UDPAddr    *net.UDPAddr
	PacketSize int
}

func (network *Network) Listen() {
	conn, err := net.ListenUDP("udp", network.UDPAddr)

	if err != nil {
		panic(err)
	}

	defer conn.Close() // close connection when Listen returns

	// read messages in a loop
	for {
		buf := make([]byte, network.PacketSize)
		_, addr, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			panic(err)
		}

		fmt.Print("> ", string(buf[0:])) // TODO: handle the read message

		// TODO: respond in an appropriate way
		conn.WriteToUDP([]byte("response from listener"), addr)
	}
}

// send generic message
func (network *Network) SendMessage(contact *Contact, msg []byte) {
	udpAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		panic(err)
	}

	_, err = conn.Write(msg)
	if err != nil {
		panic(err)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
