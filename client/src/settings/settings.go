package settings

import (
	"aio/common/src/host"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type ClientSettings struct {
	Timeout            time.Duration     `json:"connectionTimeout"`
	MaxRetries         int               `json:"maxRetries"`
	Host               host.HostSettings `json:"host"`
	ClientName         string            `json:"clientName"`
	DefaultNameAllowed bool              `json:"defaultNameAllowed"`
	MaxRngValue        int               `json:"maxRngValue"`
	AskName            bool              `json:"askName"`
}

// Initialize the settings of the client from the config file.
func ReadSettings(configFilePath string) (*ClientSettings, error) {
	data, err := ioutil.ReadFile(configFilePath)

	if err != nil {
		return nil, fmt.Errorf("could not find the file: %v", err)
	}

	clientSettings := new(ClientSettings)
	err = json.Unmarshal(data, clientSettings)

	if err != nil {
		return nil, fmt.Errorf("could not convert the json to obj: %v", err)
	}

	clientSettings.getClientName()
	return clientSettings, nil
}

// Read a client name from stdin or generate one
func (clientSettings *ClientSettings) getClientName() string {
	// Create random number generator
	seed := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(seed)

	// Forced name generation
	if !clientSettings.AskName {
		randomValue := strconv.Itoa(rng.Intn(clientSettings.MaxRngValue))
		clientSettings.ClientName = clientSettings.ClientName + randomValue
		return clientSettings.ClientName
	}

	// Open communication with the user
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter a client name: ")

	// Interact with the user
	for scanner.Scan() {
		input := scanner.Text()

		if len(input) != 0 {
			clientSettings.ClientName = input
			break
		}

		if clientSettings.DefaultNameAllowed {
			randomValue := strconv.Itoa(rng.Intn(clientSettings.MaxRngValue))
			clientSettings.ClientName = clientSettings.ClientName + randomValue
			break
		}

		fmt.Println("You must enter a client name!")
		fmt.Print("Enter a client name: ")
	}

	return clientSettings.ClientName
}
