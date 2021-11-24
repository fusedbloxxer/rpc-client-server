package server

import (
	"net"
	"fmt"
	"strings"
	"encoding/json"
	"aio/server/src/pmap"
	"aio/server/src/settings"
	mod "aio/common/src/model"
	logg "aio/common/src/logger"
	prob "aio/server/src/problems"
)

type Server struct {
	Clients			*pmap.PMap
	Logger 	 		*logg.Logger
	Listener		*net.Listener
	Settings 		*settings.ServerSettings
	ProblemMapper	map[string]func([]interface{}) (string, error)
}

func (s *Server) Init(configFilePath string) (err error) {
	// Read the config file settings
	if s.Settings, err = settings.ReadSettings(configFilePath); err != nil {
		return
	}

	// Map the requests to the solution functions
	s.ProblemMapper = map[string]func([]interface{}) (string, error) {
		"1": prob.Problem1,
		"2": prob.Problem2,
		"3": prob.Problem3,
		"8": prob.Problem8,
	}

	// Construct and initialize objs
	s.Clients = new(pmap.PMap)
	s.Clients.Map = make(map[string]bool)
	s.Listener = new(net.Listener)
	s.Logger = new(logg.Logger)
	s.Logger.Entity = s
	return
}

func (s *Server) Start(
	callback func(*Server, net.Conn, chan error),
	errorHandler func(*Server, error),
) (err error) {
	// Start listening to requests
	*s.Listener, err = net.Listen(
		s.Settings.Host.Protocol,
		s.Settings.Host.Server(),
	)

	// Assure no error occurred
	if err != nil {
		return
	}

	// Create error channel
	var e chan error = make(chan error)

	// Launch error handler watcher
	go func(e chan error) {
		for {
			// Wait for an error to occurr
			err := <-e
			go errorHandler(s, err)
		}
	}(e)

	// Close the server connection
	defer func() {
		(*s.Listener).Close()
	}()

	// Start accepting connections
	for {
		var c net.Conn
		c, err = (*s.Listener).Accept()

		if err != nil {
			return
		}

		go callback(s, c, e)
	}
}

func (s *Server) Parse(request string, req *mod.RequestModel) (err error) {
	if err = json.Unmarshal([]byte(request), req); err != nil {
		return err
	}
	return nil
}

func (s *Server) Solve(problem string, arr []interface{}) (string, error) {
	if solver, ok := s.ProblemMapper[problem]; !ok {
		return "", fmt.Errorf("cannot handle request %v", problem)
	} else {
		return solver(arr)
	}
}

func (s *Server) Send(conn net.Conn, response *mod.ResponseModel) (err error) {
	if response == nil {
		return
	}

	// Transform response to json
	var raw []byte
	if raw, err = json.Marshal(*response); err != nil {
		return err
	}

	// Send the message to the client
	fmt.Fprintf(conn, "%v\n", string(raw))
	return nil
}

func (s *Server) ProcessRequest(
	req *mod.RequestModel,
) (res *mod.ResponseModel, err error) {
	// Log client request
	s.Logger.Log(
		fmt.Sprintf("client %s made a request: %v\n", req.Sender, *req),
	)

	// Create a response
	res = new(mod.ResponseModel)

	// Compute the response
	switch req.Type {
	case mod.COMMAND:
		// Validate the client
		if !s.Clients.Exists(req.Sender) {
			*res = mod.ResponseModel{
				Content: "client not registered",
				Status: mod.UNREGISTERED,
			}
			// Log server response
			s.Logger.Log(
				fmt.Sprintf("send response %v to client %s\n", *res, req.Sender),
			)
			return
		}

		// Extract the command
		com := mod.Command{
			Verb: req.Content.(map[string]interface{})["verb"].(string),
			Args: req.Content.(map[string]interface{})["args"],
		}

		// Solve each specific verb
		switch com.Verb {
		case mod.SOLVE:
			err = s.ResolveSolveCommand(com, res)
		case mod.LIST:
			err = s.ResolveListCommand(com, res)
		default:
			*res = mod.ResponseModel{
				Content: "invalid verb",
				Status: mod.BADREQUEST,
			}
		}
	case mod.SALUTE:
		if s.Clients.Exists(req.Sender) {
			*res = mod.ResponseModel{
				Content: "client already registered",
				Status: mod.ALRDREGISTERED,
			}
				// Log server response
			s.Logger.Log(
				fmt.Sprintf("send response %v to client %s\n", *res, req.Sender),
			)
			return
		}
		s.Clients.Add(req.Sender)

		// Log client connection
		defer s.Logger.Log("client " + req.Sender + " has connected\n")

		*res = mod.ResponseModel{
			Content: "registration successful",
			Status: mod.OK,
		}
	case mod.BYE:
		s.Clients.Delete(req.Sender)

		// Log client connection
		defer s.Logger.Log("client " + req.Sender + " has disconnected\n")

		*res = mod.ResponseModel{
			Content: "connection stopped",
			Status: mod.OK,
		}
	case mod.ACK:
		// Log the client confirmation of OK messages
		s.Logger.Log(
			fmt.Sprintf(
				"client " + req.Sender + " has received: %v\n",
				req.Content,
			),
		)
	default:
		*res = mod.ResponseModel{
			Content: fmt.Sprintf("invalid request type %v", req.Type),
			Status: mod.BADREQUEST,
		}
	}

	// Log server response
	s.Logger.Log(
		fmt.Sprintf("send response %v to client %s\n", *res, req.Sender),
	)

	// Return the response
	return
}

func (s *Server) ResolveListCommand(com mod.Command, res *mod.ResponseModel) (err error) {
	entity := com.Args.(map[string]interface{})["entity"].(string)

	switch entity {
	case "clients":
		clients := strings.Join(s.Clients.Keys(), ",")
		*res = mod.ResponseModel{
			Content: fmt.Sprintf("registered clients: %v", clients),
			Status: mod.OK,
		}
	default:
		*res = mod.ResponseModel{
			Content: "invalid list entity",
			Status: mod.BADREQUEST,
		}
	}

	return
}

func (s *Server) ResolveSolveCommand(com mod.Command, res *mod.ResponseModel) (err error) {
	// Extract the solve command
	solve := mod.SolveCommand{
		Array: com.Args.(map[string]interface{})["array"].([]interface{}),
		Problem: com.Args.(map[string]interface{})["problem"].(string),
	}

	// Validate the data
	if len(solve.Array) > s.Settings.MaxArrLen {
		*res = mod.ResponseModel{
			Content: fmt.Sprintf("array length should be leq %v", s.Settings.MaxArrLen),
			Status: mod.BADREQUEST,
		}
		return
	}

	// Solve the problem
	var e error
	var sol string
	if sol, e = s.Solve(solve.Problem, solve.Array); e != nil {
		*res = mod.ResponseModel{
			Content: e.Error(),
			Status: mod.ERROR,
		}
		return
	}

	// Send the result
	*res = mod.ResponseModel{
		Content: sol,
		Status: mod.OK,
	}

	// Return the result
	return
}

func (s *Server) Name() string {
	return s.Settings.Name
}