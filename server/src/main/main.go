// Run with this command:
// go build -o ./server/bin ./server/src/server.go && ./server/bin/server.exe
// or
// go run ./server/src/server.go

package main

import (
	mod "aio/common/src/model"
	"aio/server/src/server"
	"strings"
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	// Allocate memory for the server instance
	var err error
	s := new(server.Server)

	// Initialize the server from the settings file
	if err = s.Init("./server/assets/appsettings.json"); err != nil {
		log.Fatal("Failed to initialize server: ", err)
	}

	// Start handling requests from clients
	if err = s.Start(handler, errorHandler); err != nil {
		s.Logger.Fatal("server encountered an err: ", err)
	}
}

func handler(s *server.Server, conn net.Conn, e chan error) {
	for {
		// Read request from the client
		message, err := bufio.NewReader(conn).ReadString('\n')

		// Assure the connection is stil ongoing
		if err != nil {
			e <- err
			break
		}

		// Remove delimiter
		message = strings.TrimSuffix(message, "\n");

		// Parse the received request json
		var req *mod.RequestModel = new(mod.RequestModel)
		if err = s.Parse(message, req); err != nil {
			e <- err
			continue
		}

		if req.Type == mod.ACK {
			// Log the client confirmation of OK messages
			s.Logger.Log(
				fmt.Sprintf(
					"client " + req.Sender + " has received: %v\n",
					req.Content,
				),
			)
			continue
		}

		if req.Type != mod.SALUTE {
			// Inform the user that the server has received the request
			s.Logger.Log("received client request\n")
			if err = s.Send(conn, &mod.ResponseModel{
				Content: "server has received the request",
				Status: mod.LOG,
			}); err != nil {
				e <- err
				break
			}

			// Inform the user that the server is processing the data
			s.Logger.Log("processing client request\n")
			if err = s.Send(conn, &mod.ResponseModel{
				Content: "server is processing the request",
				Status: mod.LOG,
			}); err != nil {
				e <- err
				break
			}
		}

		// Process the request
		var res *mod.ResponseModel
		if res, err = s.ProcessRequest(req); err != nil {
			e <- err
			continue
		}

		if req.Type == mod.BYE && res.Status == mod.OK {
			return
		}

		s.Send(conn, res)
	}
}

func errorHandler(s *server.Server, err error) {
	s.Logger.Log("Error encountered: ", err, "\n");
}