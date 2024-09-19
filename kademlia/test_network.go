package kademlia

// send dummy message to ip
/*func TestSend(ip string) { //TODO: add assertions
	udpAddr, err := net.ResolveUDPAddr("udp", ":1234")
	if err != nil {
		panic(err)
	}

	c := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), ip+":1234")
	ch := make(chan Message)
	n := Network{ListenAddr: udpAddr, PacketSize: 512, ExpectedResponses: make(map[KademliaID]chan Message)}

	go n.Listen()
	go n.SendPingMessage(&c, ch)
	response := <-ch
	log.Println("Got response: ", response.MsgType)
	log.Println(response.RPCID)
	go n.SendPingMessage(&c, ch)
	response = <-ch
	log.Println("Got response: ", response.MsgType)
	log.Println(response.RPCID)
}*/

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

/*func TestListen() { //TODO: add assertions
	udpAddr, err := net.ResolveUDPAddr("udp", ":1234")
	if err != nil {
		panic(err)
	}

	n := Network{ListenAddr: udpAddr, PacketSize: 512}

	n.Listen()
}*/
