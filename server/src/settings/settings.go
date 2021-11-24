package settings

import (
	"fmt"
	"time"
	"io/ioutil"
	"encoding/json"
	"aio/common/src/host"
)

type ServerSettings struct {
	ErrorPolling	time.Duration 	   `json:"errorPolling"`
	MaxClients		int				   `json:"maxClients"`
	Name			string 			   `json:"serverName"`
	MaxArrLen		int				   `json:"maxArrLen"`
	Host			*host.HostSettings `json:"host"`
}

func ReadSettings(configFilePath string) (*ServerSettings, error) {
	data, err := ioutil.ReadFile(configFilePath)

	if err != nil {
		return nil, fmt.Errorf(
			"cannot read the input file %s: %v",
			configFilePath,
			err,
		)
	}

	serverSettings := new(ServerSettings)
	if err = json.Unmarshal(data, serverSettings); err != nil {
		return nil, fmt.Errorf("cannot parse the server settings file: %v", err)
	}

	return serverSettings, nil
}
