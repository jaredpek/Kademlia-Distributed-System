package kademlia

import (
	"fmt"
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

func TestHandleResponse(t *testing.T) {
	fmt.Println("hello")
	id := *NewRandomKademliaID()
	state := Network{
		ListenPort:        "1234",
		PacketSize:        1024,
		ExpectedResponses: make(map[KademliaID]chan Message, 10),
		Rt:                rt,
	}
	responseCh := make(chan Message)
	state.ExpectedResponses[id] = responseCh
	msg := Message{RPCID: id, MsgType: "test756756756", Body: "This is in the state channel"}
	go state.handleResponse(msg)

	response := <-responseCh
	if (response.MsgType != msg.MsgType) || (response.RPCID != msg.RPCID) || (response.Body != msg.Body) {
		t.Fatalf("Message was not successfuly retrieved")
	}
	if _, ok := state.ExpectedResponses[id]; ok {
		t.Fatalf("Map entry for the message was not removed")
	}
}
