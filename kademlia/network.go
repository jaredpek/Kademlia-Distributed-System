package kademlia

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

const timeout = 10 * time.Second

// interfaces and structs for Messenger
type Messenger interface {
	SendMessage(contact *Contact, msg Message)
}

type UDPMessenger struct {
	Rt *RoutingTable
}

type MockMessenger struct {
	Rt       *RoutingTable
	Messages []Message
}

type Network struct {
	Rt                *RoutingTable
	BootstrapIP       string
	ListenPort        string
	PacketSize        int
	ExpectedResponses map[KademliaID](chan Message) // map of RPCID : message channel used by handler
	lock              sync.Mutex
	Messenger         Messenger
}

type Message struct {
	MsgType  string
	Sender   Contact
	Body     string
	Key      KademliaID
	RPCID    KademliaID
	Contacts []Contact
}

// send generic message
func (m *UDPMessenger) SendMessage(contact *Contact, msg Message) {
	log.Println("Sending message: ", msg)
	// make sure the sender field is always this node
	// network.lock.Lock()
	msg.Sender = m.Rt.me
	// network.lock.Unlock()

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

// Mock version of send message. Used for testing
func (m *MockMessenger) SendMessage(_ *Contact, msg Message) {
	msg.Sender = m.Rt.me
	m.Messages = append(m.Messages, msg)
}

// Get latest message from mock version of send message
func (m *MockMessenger) GetLatestMessage() (Message, error) {
	if len(m.Messages) == 0 {
		return Message{}, fmt.Errorf("MOCK MESSAGE ERROR: There are no more messages! Returning empty message")
	}
	mes := m.Messages[0]

	m.Messages = append(m.Messages[:0], m.Messages[1:]...)

	return mes, nil
}

func (network *Network) AddContact(contact Contact) {
	bucketIndex := network.Rt.getBucketIndex(contact.ID)
	correctBucket := network.Rt.buckets[bucketIndex]
	if correctBucket.Len() == bucketSize { // if bucket is full
		network.Rt.lock.Lock()
		oldestElement := correctBucket.list.Back()
		oldestContact := oldestElement.Value.(Contact) // get oldest contact
		network.Rt.lock.Unlock()

		responseChan := make(chan Message)
		log.Println("Network add contact ping.")
		go network.SendPingMessage(&oldestContact, responseChan) // ping oldest contact
		response := <-responseChan
		if response.MsgType == "TIMEOUT" { // if oldest timed out
			network.Rt.lock.Lock()
			network.Rt.AddContact(contact) // add new contact
			network.Rt.lock.Unlock()
		} else {
			network.Rt.lock.Lock()
			correctBucket.list.MoveToFront(oldestElement) // else move oldest to front
			network.Rt.lock.Unlock()
		}
	} else { // if bucket is not full
		network.Rt.lock.Lock()
		network.Rt.AddContact(contact)
		network.Rt.lock.Unlock()
	}
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
	messages := make(chan Message, 50) // large buffer because of network startup, might receive 50 at once
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

		log.Println("received message: ", decoded_message) // for debugging

		messages <- decoded_message //give received message to the handler
	}
}

// TODO: add testing
func (network *Network) MessageHandler(messages chan Message) {
	for {
		// TODO: perform appropriate routing table operations

		received_message := <-messages

		//add the sender to routing table
		sender := received_message.Sender
		sender.CalcDistance(network.Rt.me.ID) // calc distance to self
		go network.AddContact(sender)

		switch received_message.MsgType {
		case "PING":
			go network.SendPongMessage(received_message)
		case "FIND_CONTACT":
			go network.SendFindContactResponse(received_message)
		case "FIND_DATA":
			go network.SendFindDataResponse(received_message)
		case "STORE":
			fmt.Println("GOT STORE MESSAGE")
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

// Send message to contact and await a response. Time out after 10 seconds
func (network *Network) SendAndAwaitResponse(contact *Contact, message Message) Message {
	response := make(chan Message) // channel for receiving a response to the sent message

	network.lock.Lock()
	network.ExpectedResponses[message.RPCID] = response // "subscribe" to receive a response
	network.lock.Unlock()

	network.Messenger.SendMessage(contact, message)
	fmt.Println(network.ExpectedResponses[message.RPCID])

	select {
	case read := <-response: // got a response
		return read
	case <-time.After(timeout): // no response
		log.Println("Time out while waiting for message: ", message.RPCID)
		return Message{MsgType: "TIMEOUT", RPCID: message.RPCID}
	}
}

// Send ping message to contact and wait for a response
// TODO: add timeout
// TODO: add testing
func (network *Network) SendPingMessage(contact *Contact, out chan Message) {
	// make the message
	ID := *NewRandomKademliaID()
	m := Message{
		MsgType: "PING",
		RPCID:   ID,
	}

	response := network.SendAndAwaitResponse(contact, m) // send message, get a response or a timeout
	out <- response                                      // return the response through the out channel
}

// send pong response to the subject message
func (network *Network) SendPongMessage(subject Message) {
	m := Message{
		MsgType: "PONG",
		RPCID:   subject.RPCID,
	}
	network.Messenger.SendMessage(&subject.Sender, m)
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

	response := network.SendAndAwaitResponse(contact, m) // send message, get a response or a timeout
	out <- response                                      // return the response through the out channel
}

func (network *Network) SendFindContactResponse(subject Message) {
	closest := network.Rt.FindClosestContacts(&subject.Key, bucketSize)

	m := Message{
		MsgType:  "FIND_CONTACT_RESPONSE",
		RPCID:    subject.RPCID,
		Contacts: closest,
	}
	network.Messenger.SendMessage(&subject.Sender, m)
}

func (network *Network) SendFindDataMessage(hash KademliaID, contact *Contact, out chan Message) {
	ID := *NewRandomKademliaID()
	m := Message{
		MsgType: "FIND_DATA",
		RPCID:   ID,
		Key:     hash,
	}

	response := network.SendAndAwaitResponse(contact, m) // send message, get a response or a timeout
	out <- response                                      // return the response through the out channel
}

func (network *Network) SendFindDataResponse(subject Message) {
	// TODO

	// find data
	// res, err := network.FindData(subject.Key.String())

	// send response
}

func (network *Network) FindData(key string) (string, error) {
	fmt.Println("filename:", key)
	path := "kademlia/values/" + key
	res, err := os.ReadFile(path)

	if err != nil {
		fmt.Println("Could not find file...")
		return "", err
	}

	fmt.Println("This value was found:", string(res))

	return string(res), nil
}

func (network *Network) SendStoreMessage(key KademliaID, data []byte, contact *Contact /*, out chan Message*/) {
	ID := *NewRandomKademliaID()
	m := Message{
		MsgType: "STORE",
		RPCID:   ID,
		Key:     key,
		Body:    string(data),
	}

	fmt.Println("Value in body:", hex.EncodeToString(data))

	// since Store does not return a message should not be needed
	/*response, err := network.SendAndAwaitResponse(contact, m) // send message, get a response or a timeout
	if err != nil {                                           // there is a timeout, no response
		return
	}
	out <- response // return the response through the out channel*/

	network.Messenger.SendMessage(contact, m)
}

func (network *Network) SendStoreResponse(subject Message) {
	// store data
	path := "kademlia/values/" + subject.Key.String()
	err := os.WriteFile(path, []byte(subject.Body), 0666)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Values saved!")
}
