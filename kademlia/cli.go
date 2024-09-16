package kademlia

import (
	"fmt"
	"os"
)

type cli struct {
	Kademlia *Kademlia
}

func newCli(kademlia *Kademlia) *cli {
	cli := &cli{}
	cli.Kademlia = kademlia
	return cli
}

func (cli *cli) handleInput() {
	var command, input string

	fmt.Scanf("%s %s", &command, &input)
	fmt.Println("Your input: ", command, " ", input)
}

func (cli *cli) put(data []byte) {
	cli.Kademlia.Store(data)
}

func (cli *cli) get(data string) {
	cli.Kademlia.LookupData(data)
}

func exit() {
	os.Exit(0)
}
