package interpreter

import (
	"fmt"
	"regexp"
	"strings"
	"encoding/json"
	mod "aio/common/src/model"
)

func ParseRequest(request string) (req mod.Request, err error) {
	// Remove terminators
	request = strings.TrimSuffix(request, "\r\n")
	request = strings.TrimSuffix(request, "\n")

	// Get the type of the request
	var reqType string
	if reqType, err = GetRequestType(request); err != nil {
		return nil, err
	}

	// Create a request based on the type
	switch reqType {
	case mod.COMMAND:
		if req, err = ParseCommand(request); err != nil {
			return
		}
	case mod.SALUTE:
		req = new(mod.Salute)
		return
	case mod.BYE:
		req = new(mod.Bye)
		return
	default:
		return nil, fmt.Errorf("")
	}

	// Return the result
	return
}

func ParseCommand(request string) (req *mod.Command, err error) {
	// Split the cli args
	values := strings.Split(request, " ")
	params := values[1:]
	verb := values[0]

	// Parse the command by using the verb
	switch verb {
	case mod.SOLVE:
		return ParseSolveCommand(params)
	case mod.LIST:
		return ParseListCommand(params)
	default:
		return nil, fmt.Errorf("invalid command verb")
	}
}

func ParseListCommand(params []string) (c *mod.Command, err error) {
	// Validate the params
	if len(params) != 1 {
		return nil, fmt.Errorf("only one argument should be provided")
	}

	// Extract the params
	entity := params[0]

	// Validate content of the params
	validEntities := regexp.MustCompile(`^(clients)$`)
	if !validEntities.MatchString(entity) {
		return nil, fmt.Errorf("%v cannot be listed", entity)
	}

	// Bundle the command
	c = new(mod.Command)
	c.Verb = mod.LIST
	c.Args = mod.ListCommand{
		Entity: entity,
	}

	// Return the result
	return c, nil
}

func ParseSolveCommand(params []string) (c *mod.Command, err error) {
	// Validate the params
	if len(params) != 2 {
		return nil, fmt.Errorf("two arguments should be provided")
	}

	// Extract the params
	rawArray := params[1]
	problem := params[0]

	// Validate the format of the array
	arrayPattern := regexp.MustCompile(`^\[((([^\s,]+,)+)?[^\s,]+)?\]$`)
	if !arrayPattern.MatchString(rawArray) {
		return nil, fmt.Errorf("invalid array format")
	}

	// Parse the array
	var arr []interface{}
	if err = json.Unmarshal([]byte(rawArray), &arr); err != nil {
		return nil, err
	}

	// Bundle the command
	c = new(mod.Command)
	c.Verb = mod.SOLVE
	c.Args = mod.SolveCommand{
		Problem: problem,
		Array: arr,
	}

	// Return the result
	return c, nil
}

func GetRequestType(request string) (string, error) {
	// Validate and return the type literal as string
	switch {
	case IsByeRequest(request):
		return mod.BYE, nil
	case IsSaluteRequest(request):
		return mod.SALUTE, nil
	case IsCommandRequest(request):
		return mod.COMMAND, nil
	default:
		return "", fmt.Errorf("invalid request type")
	}
}

func IsCommandRequest(request string) bool {
	// Compile regex
	match := regexp.MustCompile(`^\w+(\s+[^\s]+)*$`)

	// Validate
	return match.FindString(request) != ""
}

func IsByeRequest(request string) bool {
	// Compile regex
	match := regexp.MustCompile(`^bye$`)

	// Validate
	return match.FindString(request) != ""
}

func IsSaluteRequest(request string) bool {
	// Compile regex
	match := regexp.MustCompile(`^salute$`)

	// Validate
	return match.FindString(request) != ""
}
