package kademlia

import (
	"log"
)

// send dummy message to ip
func TestSend() { //TODO: add assertions

	rt := NewRoutingTable(NewContact(NewRandomKademliaID(), "127.0.0.1"))
	n := Network{
		ListenPort:        "1234",
		PacketSize:        1024,
		ExpectedResponses: make(map[KademliaID]chan Message, 10),
		Rt:                rt,
	}
	c := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "172.26.0.2:1234") // bootstrap node
	ch := make(chan Message, 5)

	//k := Kademlia{&n, rt}

	go n.Listen()
	go n.SendPingMessage(&c, ch) // ping bootstrap
	response := <-ch
	log.Println("Got response: ", response.MsgType)
	log.Println(response.RPCID)
	go n.SendPingMessage(&c, ch)
	response = <-ch
	log.Println("Got response: ", response.MsgType)
	log.Println(response.RPCID)
}

/*
// send ping message to ip
func TestSendPing(ip string) {
	udpAddr, err := net.ResolveUDPAddr("udp", ":1234")
	if err != nil {
		panic(err)
	}

	c := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), ip+":1234")

	n := Network{ListenAddr: udpAddr, PacketSize: 512}

	n.SendPingMessage(&c)
}*/

func TestListen() { //TODO: add assertions

	rt := NewRoutingTable(NewContact(NewRandomKademliaID(), "127.0.0.1"))
	n := Network{
		ListenPort:        "1234",
		PacketSize:        1024,
		ExpectedResponses: make(map[KademliaID]chan Message, 10),
		Rt:                rt,
	}

	n.Listen()
}
