package kademlia

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

type Network struct {
	Rt                *RoutingTable
	BootstrapIP       string
	ListenPort        string
	PacketSize        int
	ExpectedResponses map[KademliaID](chan Message) // map of RPCID : message channel used by handler
	lock              sync.Mutex
}

type Message struct {
	MsgType  string
	Sender   Contact
	Body     string
	Key      KademliaID
	RPCID    KademliaID
	Contacts []Contact
}

func (network *Network) Listen() {
	ListenAddr, err := net.ResolveUDPAddr("udp", ":"+network.ListenPort)
	if err != nil {
		log.Fatal(err)
	}

	// start listening
	conn, err := net.ListenUDP("udp", ListenAddr)
	if err != nil {
		log.Fatal(err) // TODO: unsure how to handle the errors should i return them or log.Fatal(err)
	}
	defer conn.Close() // close connection when listening is done

	//spawn a message handler
	messages := make(chan Message)
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

		decoded_message.Sender.Address = addr.IP.String() + ":" + network.ListenPort // ensure the sender has the correct IP

		log.Println("ip:", addr.IP) // for debugging
		log.Println("port:", addr.Port)

		messages <- decoded_message //give received message to the handler
	}
}

// TODO: add testing
func (network *Network) MessageHandler(messages chan Message) {
	for {
		// TODO: perform appropriate routing table operations

		received_message := <-messages

		//add the sender to routing table
		network.Rt.lock.Lock()
		sender := received_message.Sender
		sender.CalcDistance(network.Rt.me.ID) // calc distance to self
		network.Rt.AddContact(sender)

		log.Println(network.Rt.FindClosestContacts(network.Rt.me.ID, 20)) //debug
		network.Rt.lock.Unlock()

		switch received_message.MsgType {
		case "PING":
			go network.SendPongMessage(received_message)
		case "FIND_CONTACT":
			go network.SendFindContactResponse(received_message)
		case "FIND_DATA":
			go network.SendFindDataResponse(received_message)
		case "STORE":
			go network.SendStoreResponse(received_message)
		case "PONG", "FIND_CONTACT_RESPONSE", "FIND_DATA_RESPONSE", "STORE_RESPONSE":
			go network.handleResponse(received_message)
		}
	}
}

func (network *Network) handleResponse(response Message) {
	network.lock.Lock()
	chn := network.ExpectedResponses[response.RPCID] // grab the channel of the waiting sender
	if chn != nil {
		chn <- response // give response to the waiting channel
		close(chn)      // clean up
		delete(network.ExpectedResponses, response.RPCID)
	}
	network.lock.Unlock()
}

/*
Implementation of how send message might work. Since we want a response for every single message that we send each SendMessage should be async and coupled with a listen
(we will be listening with Listen also as an always on server). When sending, we will have our IP and we will send it on some port P. The sender function will reserve P,
and after sending the message it will block until it recieves a response on P. This way we know that we got a response to the actual sent message (we also check
sender IP and RPC ID too).
*/

// send generic message
func (network *Network) SendMessage(contact *Contact, msg Message) {
	// make sure the sender field is always this node
	network.lock.Lock()
	msg.Sender = network.Rt.me
	network.lock.Unlock()

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
// TODO: add testing
func (network *Network) SendPingMessage(contact *Contact, out chan Message) {

	//make the message
	ID := *NewRandomKademliaID()
	m := Message{
		MsgType: "PING",
		RPCID:   ID,
	}
	response := make(chan Message) // channel for receiving a response to the sent message

	network.lock.Lock()
	network.ExpectedResponses[ID] = response // "subscribe" to receive a response
	network.lock.Unlock()

	network.SendMessage(contact, m)
	read := <-response // block here until you get a response

	out <- read // return the response through the out channel

	// debug
	log.Println("Response from sent message:", read.MsgType)
	log.Println("Response ID:", read.RPCID)
}

// send pong response to the subject message
func (network *Network) SendPongMessage(subject Message) {
	log.Print("sending PONG...")

	m := Message{
		MsgType: "PONG",
		RPCID:   subject.RPCID,
	}
	network.SendMessage(&subject.Sender, m)
}

// ask contact about id, receive response in out channel
func (network *Network) SendFindContactMessage(id KademliaID, contact *Contact, out chan Message) {
	// create the message
	ID := *NewRandomKademliaID()
	m := Message{
		MsgType: "FIND_CONTACT",
		RPCID:   ID,
		Key:     id,
	}
	response := make(chan Message)

	network.lock.Lock()
	network.ExpectedResponses[ID] = response
	network.lock.Unlock()

	network.SendMessage(contact, m)
	read := <-response // block here until you get a response

	log.Println("Response from sent message:", read.MsgType)
	out <- read
}

func (network *Network) SendFindContactResponse(subject Message) {
	closest := network.Rt.FindClosestContacts(&subject.Key, bucketSize)

	m := Message{
		MsgType:  "FIND_CONTACT_RESPONSE",
		RPCID:    subject.RPCID,
		Contacts: closest,
	}
	network.SendMessage(&subject.Sender, m)
}

func (network *Network) SendFindDataMessage(hash KademliaID, contact *Contact, out chan Message) {
	ID := *NewRandomKademliaID()
	m := Message{
		MsgType: "FIND_DATA",
		RPCID:   ID,
		Key:     hash,
	}
	response := make(chan Message)

	network.lock.Lock()
	network.ExpectedResponses[ID] = response
	network.lock.Unlock()

	network.SendMessage(contact, m)
	read := <-response // block here until you get a response

	out <- read // return the response
}

func (network *Network) SendFindDataResponse(subject Message) {
	// TODO
}

func (network *Network) SendStoreMessage(key KademliaID, data string, contact *Contact, out chan Message) {
	ID := *NewRandomKademliaID()
	m := Message{
		MsgType: "STORE",
		RPCID:   ID,
		Key:     key,
		Body:    data,
	}
	response := make(chan Message)

	network.lock.Lock()
	network.ExpectedResponses[ID] = response
	network.lock.Unlock()

	network.SendMessage(contact, m)
	read := <-response // block here until you get a response

	out <- read // return the response
}

func (network *Network) SendStoreResponse(subject Message) {
	type StoredData struct {
		ID   KademliaID
		Data string
	}
	// store data
	data := StoredData{subject.RPCID, subject.Body}

	res, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Current JSON:", string(res))

	err = os.WriteFile("stored_values.json", res, 0666)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Values saved!")
}
