package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
)

func TestSendListenLocal1() {
	me := NewContact(NewRandomKademliaID(), "127.0.0.1:1234")
	rt := NewRoutingTable(me)
	n := Network{
		ListenPort:        "1234",
		PacketSize:        1024,
		ExpectedResponses: make(map[KademliaID]chan Message, 10),
		Rt:                rt,
	}
	ch := make(chan Message, 5)

	//k := Kademlia{&n, rt}

	go n.Listen()
	go n.SendPingMessage(&me, ch) // ping bootstrap
	response := <-ch
	log.Println("Got response: ", response.MsgType)
	log.Println(response.RPCID)
}

func TestDocker() {
	rt := NewRoutingTable(NewContact(NewRandomKademliaID(), "127.0.0.1"))
	n := Network{
		ListenPort:        "1234",
		PacketSize:        1024,
		ExpectedResponses: make(map[KademliaID]chan Message, 10),
		Rt:                rt,
	}

	k := Kademlia{&n, rt}

	go n.Listen()
	k.JoinNetwork()
	for {
	}
}

func TestJoin() {
	rt := NewRoutingTable(NewContact(NewRandomKademliaID(), "127.0.0.1"))
	n := Network{
		ListenPort:        "1234",
		PacketSize:        1024,
		ExpectedResponses: make(map[KademliaID]chan Message, 10),
		Rt:                rt,
		BootstrapIP:       "172.26.0.2:1234",
	}

	k := Kademlia{&n, rt}

	go n.Listen()
	k.JoinNetwork()
}

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

func TestStore(val string) {
	rt := NewRoutingTable(NewContact(NewRandomKademliaID(), "127.0.0.1"))
	n := Network{
		ListenPort:        "1234",
		PacketSize:        1024,
		ExpectedResponses: make(map[KademliaID]chan Message, 10),
		Rt:                rt,
	}

	h := sha1.New()

	io.WriteString(h, val)

	res := h.Sum(nil)

	id := NewKademliaID(hex.EncodeToString(res))

	m := Message{
		MsgType:  "",
		Sender:   Contact{},
		Body:     val,
		Key:      *id,
		RPCID:    KademliaID{},
		Contacts: nil,
	}

	fmt.Println("KademliaID:", id)

	n.SendStoreResponse(m)
}

func TestFindData(id *KademliaID) {
	rt := NewRoutingTable(NewContact(NewRandomKademliaID(), "127.0.0.1"))
	n := Network{
		ListenPort:        "1234",
		PacketSize:        1024,
		ExpectedResponses: make(map[KademliaID]chan Message, 10),
		Rt:                rt,
	}

	m := Message{
		MsgType:  "",
		Sender:   Contact{},
		Body:     "",
		Key:      *id,
		RPCID:    KademliaID{},
		Contacts: nil,
	}

	n.FindData(m.Key.String())
}

func TestRest() {
	c := NewContact(NewRandomKademliaID(), "127.0.0.1")
	k := NewKademlia(c)
	r := newRest(k)

	r.StartServer("127.0.0.1:8080")
}
