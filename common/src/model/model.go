package model

// Request Types
const (
	COMMAND = "command"
	SALUTE  = "salute"
	BYE     = "bye"
	ACK 	= "ack"
)

// Command Verbs
const (
	SOLVE = "solve"	// Answer problems using the server as the solver
	LIST  = "list" 	// List various information on the server
)

// Status Codes
const (
	ALRDREGISTERED = "alreadydregistered"
	UNREGISTERED   = "unregistered"
	BADREQUEST     = "badrequest"
	BADNAME        = "badname"
	ERROR          = "error"
	LOG			   = "log"
	OK             = "ok"
)

type ResponseModel struct {
	Content interface{} `json:"content"`
	Status  string      `json:"status"`
}

type Request interface {
	Content() interface{}
	Type() string
}

type RequestModel struct {
	Type    string      `json:"requestType"`
	Content interface{} `json:"content"`
	Sender  string      `json:"sender"`
}

type Ack struct {
	Response ResponseModel `json:"response"`
}

func (a *Ack) Content() interface{} {
	return a.Response
}

func (a *Ack) Type() string {
	return ACK
}

type Command struct {
	Verb string      `json:"verb"`
	Args interface{} `json:"args"`
}

type SolveCommand struct {
	Problem string        `json:"problem"`
	Array   []interface{} `json:"array"`
}

type ListCommand struct {
	Entity string `json:"entity"`
}

func (c *Command) Type() string {
	return COMMAND
}

func (c *Command) Content() interface{} {
	return *c
}

type Salute struct{}

func (s *Salute) Type() string {
	return SALUTE
}

func (s *Salute) Content() interface{} {
	return ""
}

type Bye struct{}

func (b *Bye) Type() string {
	return BYE
}

func (b *Bye) Content() interface{} {
	return ""
}
