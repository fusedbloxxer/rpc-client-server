// Run with this command:
// go build -o ./client/bin ./client/src/client.go && ./client/bin/client.exe
// or
// go run ./client/src/client.go

package main

import (
	"aio/client/src/client"
	"aio/client/src/interpreter"
	mod "aio/common/src/model"
	"bufio"
	"log"
	"os"
)

func main() {
	// Create a new client
	c := new(client.Client)

	// Initialize the client
	if err := c.Init("./client/assets/appsettings.json"); err != nil {
		log.Fatal(
			"Failed to initialize: ",
			err,
		)
	}

	// Connect to the server
	if err := c.Connect(); err != nil {
		c.Logger.Fatal(err, "\n")
	}

	// Listen asynchronously to messages from the server
	c.RecvLoopAsync()

	// Send messages synchronously to the server
	err := c.SendLoopSync(func(cl *client.Client) (int, error) {
		// Read buffered input from the user
		reader := bufio.NewReader(os.Stdin)

		// Read an entire line
		input, _ := reader.ReadString('\n')

		// Interpret the client input
		var err error
		var req mod.Request
		if req, err = interpreter.ParseRequest(input); err != nil {
			cl.Logger.Log(err, "\n")
			return 0, nil
		}

		// Check if the user terminated the session
		if req.Type() == mod.BYE {
			c.Disconnect()
			return 1, nil
		}

		// Send it to the server
		if err := cl.Send(req); err != nil {
			return -1, err
		}

		// Continue the Loop
		return 0, nil
	})

	// Handle Errors
	if err != nil {
		c.Logger.Fatal(
			"connection ended abruptly: ",
			err,
		)
	}
}
