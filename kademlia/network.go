package kademlia

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"sync"
)

type Network struct {
	ListenAddr        *net.UDPAddr
	PacketSize        int
	ContactChan       chan ReceivedMessage
	DataChan          chan ReceivedMessage
	PingChan          chan ReceivedMessage
	StoreChan         chan ReceivedMessage
	ExpectedResponses map[KademliaID](chan ReceivedMessage) // map of RPCID : message channel used by handler
	lock              sync.Mutex
}

type Message struct {
	MsgType  string
	Body     string
	Key      string
	RPCID    KademliaID
	Contacts []Contact
}

type ReceivedMessage struct { // you need the senders ip when passing from listener to message handler
	Msg    Message
	Sender net.UDPAddr
}

func (network *Network) Listen() {
	conn, err := net.ListenUDP("udp", network.ListenAddr) // start listening
	if err != nil {
		log.Fatal(err) // TODO: unsure how to handle the errors should i return them or log.Fatal(err)
	}
	defer conn.Close() // close connection when listening is done

	//spawn a message handler
	messages := make(chan ReceivedMessage)
	go network.MessageHandler(messages)

	// read messages in a loop
	for {
		buf := make([]byte, network.PacketSize)
		n, addr, err := conn.ReadFromUDP(buf[0:]) // place read message in buf
		if err != nil {
			log.Fatal(err)
		}

		dec := gob.NewDecoder(bytes.NewBuffer(buf[:n])) // give message as input to decoder
		var decoded_message Message
		if err := dec.Decode(&decoded_message); err != nil { //place the decoded message in decoded_message
			log.Fatal(err)
		}

		fmt.Println("ip:", addr.IP)
		fmt.Println("port:", addr.Port)
		fmt.Println(addr.Zone)
		messages <- ReceivedMessage{decoded_message, *addr} //give received message to the handler
	}
}

func (network *Network) MessageHandler(messages chan ReceivedMessage) {
	for {
		// TODO: perform appropriate routing table operations
		received_message := <-messages
		switch received_message.Msg.MsgType {
		case "PING":
			fmt.Println("Got ping in handler")
			network.SendPongMessage(&Contact{Address: received_message.Sender.IP.String() + ":1234"}, received_message.Msg.RPCID)
		case "FIND_CONTACT":
			network.ContactChan <- received_message
		case "FIND_DATA":
			network.DataChan <- received_message
		case "STORE":
			network.StoreChan <- received_message
		case "PONG":
			fmt.Println("Got pong in handler")
			network.lock.Lock()
			chn := network.ExpectedResponses[received_message.Msg.RPCID]
			fmt.Println("handler putting pong in channel")
			chn <- received_message
			network.lock.Unlock()
		}
	}
}

/*
Implementation of how send message might work. Since we want a response for every single message that we send each SendMessage should be async and coupled with a listen
(we will be listening with Listen also as an always on server). When sending, we will have our IP and we will send it on some port P. The sender function will reserve P,
and after sending the message it will block until it recieves a response on P. This way we know that we got a response to the actual sent message (we also check
sender IP and RPC ID too).
*/

// send generic message
func (network *Network) SendMessage(contact *Contact, msg Message) {
	// set up the connection
	udpAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf) // encoded bytes go to buf

	if err := enc.Encode(msg); err != nil { // encode
		log.Fatal(err)
	}

	_, err = conn.Write(buf.Bytes()) // send encoded message
	if err != nil {
		log.Fatal(err)
	}
}

// Send ping message to contact and wait for a response
// TODO: add timeout
func (network *Network) SendPingMessage(contact *Contact) {
	/*network.SendMessage(
		contact,
		Message{
			MsgType:  "ping",
			Body:     "Ping!",
			Key:      "",
			Contacts: nil,
		},
	)
	fmt.Printf("Sending ping to %s", contact.Address)
	// TODO*/

	fmt.Println("Sending PING...")

	ID := *NewKademliaID("FFFFFFFF10000000000000000000000000000000")
	m := Message{MsgType: "PING", RPCID: ID}
	response := make(chan ReceivedMessage)

	network.lock.Lock()
	network.ExpectedResponses[ID] = response
	network.lock.Unlock()

	network.SendMessage(contact, m)
	read := <-response // block here until you get a response

	network.lock.Lock()
	close(response)
	delete(network.ExpectedResponses, m.RPCID)
	network.lock.Unlock()

	fmt.Println("Response from sent message:", read.Msg.MsgType)
}

func (network *Network) SendPongMessage(contact *Contact, ID KademliaID) {
	fmt.Print("sending PONG...")
	m := Message{MsgType: "PONG", RPCID: ID}
	network.SendMessage(contact, m)
}

func (network *Network) SendFindContactMessage(contact *Contact) {

}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
