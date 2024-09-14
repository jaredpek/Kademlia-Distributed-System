package kademlia

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type Network struct {
	ListenAddr  *net.UDPAddr
	PacketSize  int
	ContactChan chan ReceivedMessage
	DataChan    chan ReceivedMessage
	PingChan    chan ReceivedMessage
	StoreChan   chan ReceivedMessage
}

type Message struct {
	MsgType  string
	Body     string
	Key      string
	Contacts []Contact
}

type ReceivedMessage struct { // you need the senders ip when passing from listener to message handler
	Msg    Message
	Sender *net.UDPAddr
}

func (network *Network) Listen() error {
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

		messages <- ReceivedMessage{decoded_message, addr} //give received message to the handler
	}
}

func (network *Network) MessageHandler(messages chan ReceivedMessage) {
	for {
		received_message := <-messages
		switch received_message.Msg.MsgType {
		case "ping":
			panic("MessageHandler for ping is not implemented!")
		case "findContact":
			panic("MessageHandler for findContact is not implemented!")
		case "findData":
			panic("MessageHandler for findData is not implemented!")
		case "store":
			panic("MessageHandler for store is not implemented!")
		}
	}
}

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

func (network *Network) SendPingMessage(contact *Contact) {
	network.SendMessage(
		contact,
		Message{
			MsgType:  "ping",
			Body:     "Ping!",
			Key:      "",
			Contacts: nil,
		},
	)
	fmt.Printf("Sending ping to %s", contact.Address)
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
