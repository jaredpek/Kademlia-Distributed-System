package kademlia

import (
	"testing"
)

var me = NewContact(NewRandomKademliaID(), "127.0.0.1:1234")
var rt = NewRoutingTable(me)
var n = Network{
	ListenPort:        "1234",
	PacketSize:        1024,
	ExpectedResponses: make(map[KademliaID]chan Message, 10),
	Rt:                rt,
}

func TestPingLocal(t *testing.T) {
	ch := make(chan Message, 5)

	go n.Listen()
	go n.SendPingMessage(&me, ch)
	response := <-ch
	if response.MsgType != "PONG" {
		t.Fatalf("Received message was not of correct type")
	}
}
