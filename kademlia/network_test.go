package kademlia

import (
	"testing"
	"time"
)

/*func TestPingLocal(t *testing.T) {
	ch := make(chan Message, 5)

	go n.Listen()
	go n.SendPingMessage(&me, ch)
	response := <-ch
	if response.MsgType != "PONG" {
		t.Fatalf("Received message was not of correct type")
	}
}*/

func TestHandleResponse(t *testing.T) {
	// environment for test, set locally so tests don't affect eachother
	/*-----------------------------------------------------------------------------------------------*/
	var me = NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:1234")
	var rt = NewRoutingTable(me)
	/*-----------------------------------------------------------------------------------------------*/

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

// tests that the mock function for sending messages works as expected
func TestMockSendMessage(t *testing.T) {
	// environment for test, set locally so tests don't affect eachother
	/*-----------------------------------------------------------------------------------------------*/
	var me = NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:1234")
	var rt = NewRoutingTable(me)
	var n = Network{
		ListenPort:        "1234",
		PacketSize:        1024,
		ExpectedResponses: make(map[KademliaID]chan Message, 10),
		Rt:                rt,
		Messenger:         &MockMessenger{Rt: rt},
	}
	/*-----------------------------------------------------------------------------------------------*/

	mb1 := "test message"
	mb2 := "test message 2"

	m1 := Message{
		Body:   mb1,
		Sender: me,
	}
	m2 := Message{
		Body:   mb2,
		Sender: me,
	}

	n.Messenger.SendMessage(&me, m1)
	n.Messenger.SendMessage(&me, m2)

	res1, _ := n.Messenger.(*MockMessenger).GetLatestMessage()
	res2, _ := n.Messenger.(*MockMessenger).GetLatestMessage()
	_, err := n.Messenger.(*MockMessenger).GetLatestMessage()

	if !(res1.Body == mb1 && res2.Body == mb2) {
		t.Fatalf("GetLatestMessage is returning messages in the wrong order!")
	}

	if !(err.Error() == "MOCK MESSAGE ERROR: There are no more messages! Returning empty message") {
		t.Fatalf("GetLatestMessage is not returning an error on an empty queue!")
	}
}

func TestSendPongMessage(t *testing.T) {
	// environment for test, set locally so tests don't affect eachother
	/*-----------------------------------------------------------------------------------------------*/
	var me = NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:1234")
	var other = NewContact(NewKademliaID("1FFFFFFF00000000000000000000000000000000"), "127.0.0.1:1235")
	var rt = NewRoutingTable(me)
	var n = Network{
		ListenPort:        "1234",
		PacketSize:        1024,
		ExpectedResponses: make(map[KademliaID]chan Message, 10),
		Rt:                rt,
		Messenger:         &MockMessenger{Rt: rt},
	}
	/*-----------------------------------------------------------------------------------------------*/

	id := *NewRandomKademliaID()

	m := Message{
		RPCID:  id,
		Sender: other,
	}
	n.SendPongMessage(m)

	res := n.Messenger.(*MockMessenger).Messages[0]

	if !(res.RPCID == id && *res.Sender.ID == *me.ID) {
		t.Fatalf("Pong message could not be found!")
	}
}

func TestSendFindContactResponse(t *testing.T) {
	// environment for test, set locally so tests don't affect eachother
	/*-----------------------------------------------------------------------------------------------*/
	var me = NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:1234")
	var rt = NewRoutingTable(me)
	var n = Network{
		ListenPort:        "1234",
		PacketSize:        1024,
		ExpectedResponses: make(map[KademliaID]chan Message, 10),
		Rt:                rt,
		Messenger:         &MockMessenger{Rt: rt},
	}
	/*-----------------------------------------------------------------------------------------------*/

	// other nodes
	contacts := []Contact{
		NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000"),
		NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8000"),
		NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8000"),
		NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8000"),
		NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8000"),
		NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8000"),
	}

	var pingTest = func(_ *Contact, out chan Message) {
		m := Message{
			MsgType: "PONG",
		}
		out <- m
	}

	for _, c := range contacts {
		rt.AddContact(c, pingTest)
	}

	key := NewKademliaID("FFF1111100000000000000000000000000000000")

	closest := rt.FindClosestContacts(key, bucketSize)

	id := *NewRandomKademliaID()

	m := Message{
		RPCID:  id,
		Key:    *key,
		Sender: me,
	}

	n.SendFindContactResponse(m)

	res, _ := n.Messenger.(*MockMessenger).GetLatestMessage()

	for i := 0; i < bucketSize; i++ {
		if res.Contacts[i].ID != closest[i].ID {
			t.Fatalf("The Contacts returned from 'SendFindContactResponse' are incorrect! \n%s != %s", res.Contacts[i].String(), closest[i].String())
		}
	}
}

func TestSendAndAwaitResponse(t *testing.T) {
	// environment for test, set locally so tests don't affect eachother
	/*-----------------------------------------------------------------------------------------------*/
	var me = NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:1234")
	var rt = NewRoutingTable(me)
	var n = Network{
		ListenPort:        "1234",
		PacketSize:        1024,
		ExpectedResponses: make(map[KademliaID]chan Message, 10),
		Rt:                rt,
		Messenger:         &MockMessenger{Rt: rt},
	}
	/*-----------------------------------------------------------------------------------------------*/

	id := *NewRandomKademliaID()

	m := Message{
		RPCID: id,
	}

	res := n.SendAndAwaitResponse(&me, m)

	// test timeout, takes 10 seconds
	if res.MsgType != "TIMEOUT" {
		t.Fatalf("The 'SendAndAwaitResponse' should timeout after a given period if there is no response!")
	}

	m.MsgType = "NOTTIMEOUT"

	// Ugly solution, but works
	go func() {
		n.lock.Lock()
		ch, ok := n.ExpectedResponses[id]
		n.lock.Unlock()
		for !ok {
			time.Sleep(10 * time.Millisecond)
			n.lock.Lock()
			ch, ok = n.ExpectedResponses[id]
			n.lock.Unlock()
		}
		n.lock.Lock()
		ch <- m
		n.lock.Unlock()
	}()

	// test SendAndAwait where it does not timeout
	res = n.SendAndAwaitResponse(&me, m)

	if res.MsgType != "NOTTIMEOUT" {
		t.Fatalf("The 'SendAndAwaitResponse' does not read response from array! %s", res.MsgType)
	}
}

func TestSendPingMessage(t *testing.T) {
	// environment for test, set locally so tests don't affect eachother
	/*-----------------------------------------------------------------------------------------------*/
	var me = NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:1234")
	var other = NewContact(NewKademliaID("1FFFFFFF00000000000000000000000000000000"), "127.0.0.1:1235")
	var rt = NewRoutingTable(me)
	var n = Network{
		ListenPort:        "1234",
		PacketSize:        1024,
		ExpectedResponses: make(map[KademliaID]chan Message, 10),
		Rt:                rt,
		Messenger:         &MockMessenger{Rt: rt},
	}
	/*-----------------------------------------------------------------------------------------------*/

	res := Message{
		MsgType: "PONG",
		Sender:  me,
	}

	resCh := make(chan Message)

	go n.SendPingMessage(&other, resCh)

	// another ugly solution, that technically does not cross process bounderies
	go func() {
		m, err := n.Messenger.(*MockMessenger).GetLatestMessage()

		// wait for message to be readable
		for err != nil {
			m, err = n.Messenger.(*MockMessenger).GetLatestMessage()
			time.Sleep(10 * time.Millisecond)
		}

		id := m.RPCID

		res.RPCID = id

		// wait for channel to be active and send PONG message
		n.lock.Lock()
		ch, ok := n.ExpectedResponses[id]
		n.lock.Unlock()
		for !ok {
			time.Sleep(10 * time.Millisecond)
			n.lock.Lock()
			ch, ok = n.ExpectedResponses[id]
			n.lock.Unlock()
		}

		n.lock.Lock()
		ch <- res
		n.lock.Unlock()
	}()

	// response sent back from SendPingMessage
	sentRes := <-resCh

	if sentRes.MsgType != "PONG" {
		t.Fatalf("The 'SendPingMessage' does not read the response from array!")
	}
}

func TestSendStoreMessage(t *testing.T) {
	// environment for test, set locally so tests don't affect eachother
	/*-----------------------------------------------------------------------------------------------*/
	var me = NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:1234")
	var rt = NewRoutingTable(me)
	var n = Network{
		ListenPort:        "1234",
		PacketSize:        1024,
		ExpectedResponses: make(map[KademliaID]chan Message, 10),
		Rt:                rt,
		Messenger:         &MockMessenger{Rt: rt},
	}
	/*-----------------------------------------------------------------------------------------------*/

	key := *NewKademliaID("FFF1111100000000000000000000000000000000")

	data := "this is a string"

	n.SendStoreMessage(key, []byte(data), &me)

	res, _ := n.Messenger.(*MockMessenger).GetLatestMessage()

	if !(res.Key.String() == key.String() && res.Body == data && res.Sender.ID.String() == me.ID.String()) {
		t.Fatalf("The 'SendStoreMessage' does not send the correct message!")
	}
}

func TestSendFindContactMessage(t *testing.T) {
	// environment for test, set locally so tests don't affect eachother
	/*-----------------------------------------------------------------------------------------------*/
	var me = NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:1234")
	var other = NewContact(NewKademliaID("1FFFFFFF00000000000000000000000000000000"), "127.0.0.1:1235")
	var rt = NewRoutingTable(me)
	var n = Network{
		ListenPort:        "1234",
		PacketSize:        1024,
		ExpectedResponses: make(map[KademliaID]chan Message, 10),
		Rt:                rt,
		Messenger:         &MockMessenger{Rt: rt},
	}
	/*-----------------------------------------------------------------------------------------------*/

	res := Message{
		MsgType: "SendFindContactResponse",
		Sender:  other,
	}

	id := NewRandomKademliaID()
	resCh := make(chan Message)

	go n.SendFindContactMessage(*id, &me, resCh)

	// another ugly solution, that technically does not cross process bounderies
	go func() {
		m, err := n.Messenger.(*MockMessenger).GetLatestMessage()

		// wait for message to be readable
		for err != nil {
			m, err = n.Messenger.(*MockMessenger).GetLatestMessage()
			time.Sleep(10 * time.Millisecond)
		}

		id := m.RPCID

		res.RPCID = id

		// wait for channel to be active and send PONG message
		n.lock.Lock()
		ch, ok := n.ExpectedResponses[id]
		n.lock.Unlock()
		for !ok {
			time.Sleep(10 * time.Millisecond)
			n.lock.Lock()
			ch, ok = n.ExpectedResponses[id]
			n.lock.Unlock()
		}

		n.lock.Lock()
		ch <- res
		n.lock.Unlock()
	}()

	sentRes := <-resCh

	if !(sentRes.Sender.ID.String() == other.ID.String() && sentRes.MsgType == "SendFindContactResponse") {
		t.Fatalf("The 'SendFindContactMessage' does not read the response from array!")
	}
}

func TestSendFindDataMessage(t *testing.T) {
	// environment for test, set locally so tests don't affect eachother
	/*-----------------------------------------------------------------------------------------------*/
	var me = NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:1234")
	var other = NewContact(NewKademliaID("1FFFFFFF00000000000000000000000000000000"), "127.0.0.1:1235")
	var rt = NewRoutingTable(me)
	var n = Network{
		ListenPort:        "1234",
		PacketSize:        1024,
		ExpectedResponses: make(map[KademliaID]chan Message, 10),
		Rt:                rt,
		Messenger:         &MockMessenger{Rt: rt},
	}
	/*-----------------------------------------------------------------------------------------------*/

	res := Message{
		MsgType: "SendFindDataResponse",
		Sender:  other,
	}

	resCh := make(chan Message)
	hash := NewRandomKademliaID()

	go n.SendFindDataMessage(*hash, &me, resCh)

	// another ugly solution, that technically does not cross process bounderies
	go func() {
		m, err := n.Messenger.(*MockMessenger).GetLatestMessage()

		// wait for message to be readable
		for err != nil {
			m, err = n.Messenger.(*MockMessenger).GetLatestMessage()
			time.Sleep(10 * time.Millisecond)
		}

		id := m.RPCID

		res.RPCID = id

		// wait for channel to be active and send PONG message
		n.lock.Lock()
		ch, ok := n.ExpectedResponses[id]
		n.lock.Unlock()
		for !ok {
			time.Sleep(10 * time.Millisecond)
			n.lock.Lock()
			ch, ok = n.ExpectedResponses[id]
			n.lock.Unlock()
		}

		n.lock.Lock()
		ch <- res
		n.lock.Unlock()
	}()

	sentRes := <-resCh

	if !(sentRes.Sender.ID.String() == other.ID.String() && sentRes.MsgType == "SendFindDataResponse") {
		t.Fatalf("The 'SendFindDataMessage' does not read the response from array!")
	}
}
