package host

import (
	"fmt"
)

type HostSettings struct {
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
	Port 	 string `json:"port"`
}

func (host *HostSettings) Server() string {
	return fmt.Sprintf("%s:%s", host.Address, host.Port);
}