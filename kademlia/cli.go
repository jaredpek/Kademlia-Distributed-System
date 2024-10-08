package kademlia

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type cli struct {
	Kademlia *Kademlia
}

// Creates a new instance of a cli struct. Takes an instance of kademlia as an input
func NewCli(kademlia *Kademlia) *cli {
	cli := &cli{}
	cli.Kademlia = kademlia
	return cli
}

// Takes in new user input
func (cli *cli) UserInput() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a command: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input) // Remove any leading/trailing whitespace

	// Split the input into parts
	parts := strings.Fields(input)

	if len(parts) == 0 {
		return fmt.Errorf("CLI Error: No command entered")
	}

	command := parts[0]

	var data string

	if command == "put" {
		// "put" can accept multiple words after it
		if len(parts) > 1 {
			data = strings.Join(parts[1:], " ")
			fmt.Println("Put command with data:", data)
		} else {
			return fmt.Errorf("CLI Error: Invalid put command. No data provided")
		}
	} else if command == "get" {
		// "get" can only accept a single word after it
		if len(parts) == 2 {
			data = parts[1]
			fmt.Println("Get command with data:", data)
		} else {
			return fmt.Errorf("CLI Error: Invalid get command. Only provide the name of the file after 'get'")
		}
	} else {
		return fmt.Errorf("CLI Error: Invalid command. Must start with 'put' or 'get'")
	}

	return cli.HandleInput(command, input)
}

// Handles the users input. If the user has entered a command that is not recognised by the implementation
// the implementation panics. Should maybe be an error.
func (cli *cli) HandleInput(command, input string) error {
	err := fmt.Errorf("CLI error: disallowed input")

	if input != "" {
		switch command {
		case "put":
			cli.Put(input)
		case "get":
			cli.Get(input)
		default:
			return err
		}
	} else {
		switch command {
		case "show":
			cli.Show()
		case "exit":
			cli.Exit()
		default:
			return err
		}
	}

	return err
}

// Stores the input by calling the "Store" function in kademlia
func (cli *cli) Put(input string) {
	data := []byte(input)
	err, hash := cli.Kademlia.Store(data)

	if err != nil { // print of result should maybe not be here
		fmt.Println("An error occured:", err)
	} else {
		fmt.Println("The file has been uploaded successfully. \nHash:", hash)
	}
}

// Tries to get the data corresponding to the hash.
func (cli *cli) Get(hash string) {
	fmt.Println(cli.Kademlia.LookupData(hash)) // print of result should maybe not be here
}

func (cli *cli) Show() {
	rtInfo := "Routing table:\n"

	currRt := cli.Kademlia.Rt.buckets

	for i, val := range currRt {
		rtInfo += "Content in bucket " + strconv.Itoa(i) + "\n"
		for e := val.list.Front(); e != nil; e = e.Next() {
			rtInfo += "  " + e.Value.(Contact).ID.String() + "\n"
		}
	}

	fmt.Println(rtInfo)
}

// Terminates the node
func (cli *cli) Exit() {
	os.Exit(0)
}
