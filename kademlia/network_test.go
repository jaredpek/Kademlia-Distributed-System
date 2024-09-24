package kademlia

import (
	"testing"
)

func TestPingLocal(t *testing.T) {
	me := NewContact(NewRandomKademliaID(), "127.0.0.1:1234")
	rt := NewRoutingTable(me)
	n := Network{
		ListenPort:        "1234",
		PacketSize:        1024,
		ExpectedResponses: make(map[KademliaID]chan Message, 10),
		Rt:                rt,
	}
	ch := make(chan Message, 5)

	go n.Listen()
	go n.SendPingMessage(&me, ch) // ping bootstrap
	response := <-ch
	if response.MsgType != "PONG" {
		t.Fatalf("Received message was not of correct type")
	}
}
