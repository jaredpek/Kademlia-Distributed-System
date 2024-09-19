package kademlia

import (
	"encoding/hex"
	"fmt"
	"os"
)

type cli struct {
	Kademlia *Kademlia
}

// Creates a new instance of a cli struct. Takes an instance of kademlia as an input
func newCli(kademlia *Kademlia) *cli {
	cli := &cli{}
	cli.Kademlia = kademlia
	return cli
}

// Takes in new user input
func (cli *cli) userInput() error {
	var command, input string
	fmt.Scanf("%s %s", &command, &input)
	return cli.handleInput(command, input)
}

// Handles the users input. If the user has entered a command that is not recognised by the implementation
// the implementation panics. Should maybe be an error.
func (cli *cli) handleInput(command, input string) error {
	if input != "" {
		switch command {
		case "put":
			cli.put(input)
		case "get":
			cli.get(input)
		default:
			return fmt.Errorf("CLI error: disallowed input")
		}
	} else if command == "exit" {
		cli.exit()
	} else {
		return fmt.Errorf("CLI error: disallowed input")
	}

	return nil
}

// Stores the input by calling the "Store" function in kademlia
func (cli *cli) put(input string) {
	data, _ := hex.DecodeString(input)
	err, hash := cli.Kademlia.Store(data)

	if err != nil { // print of result should maybe not be here
		fmt.Println("An error occured:", err)
	} else {
		fmt.Println("The file has been uploaded successfully. \nHash:", hash)
	}
}

// Tries to get the data corresponding to the hash.
func (cli *cli) get(hash string) {
	fmt.Println(cli.Kademlia.LookupData(hash)) // print of result should maybe not be here
}

// Terminates the node
func (cli *cli) exit() {
	os.Exit(0)
}
