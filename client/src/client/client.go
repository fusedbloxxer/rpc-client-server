package client

import (
	conf "aio/client/src/settings"
	logg "aio/common/src/logger"
	mod "aio/common/src/model"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

type IsActive interface {
	Status() bool
	Update(bool)
}

type Alive struct {
	Mutex sync.Mutex
	Flag  bool
}

func (a *Alive) Status() (status bool) {
	a.Mutex.Lock()
	status = a.Flag
	a.Mutex.Unlock()
	return status
}

func (a *Alive) Update(status bool) {
	a.Mutex.Lock()
	a.Flag = status
	a.Mutex.Unlock()
}

type Client struct {
	Settings           *conf.ClientSettings
	Logger             *logg.Logger
	Connection         *net.Conn
	IsConnectionActive IsActive
}

func (c *Client) Name() string {
	return c.Settings.ClientName
}

func (c *Client) Init(configFilePath string) (err error) {
	if c.Settings, err = conf.ReadSettings(configFilePath); err != nil {
		return
	}

	c.IsConnectionActive = new(Alive)
	c.IsConnectionActive.Update(false)

	c.Logger = new(logg.Logger)
	c.Logger.Entity = c
	return
}

func (c *Client) SendLoopSync(callback func(c *Client) (int, error)) (err error) {
	var res int
	for err == nil && res == 0 && c.IsConnectionActive.Status() {
		res, err = callback(c)
	}
	return
}

func (c *Client) RecvLoopAsync() {
	go func() {
		for c.IsConnectionActive.Status() {
			// Receive the response
			if res, err := c.Receive(); err != nil {
				c.IsConnectionActive.Update(false)
				c.Logger.Fatal(err)
				return
			} else if res.Status == mod.OK {
				if err = c.Send(&mod.Ack{Response: *res}); err != nil {
					c.IsConnectionActive.Update(false)
					c.Logger.Fatal(err)
					return
				}
			}
		}
	}()
}

func (c *Client) Disconnect() (err error) {
	if c.Connection == nil {
		return nil
	}

	if err = c.Send(&mod.Bye{}); err != nil {
		return err
	}

	if err = c.Close(); err != nil {
		return err
	}

	c.IsConnectionActive.Update(false)
	c.Logger.Log("connection closed\n")
	return nil
}

func (c *Client) Close() error {
	if c.Connection == nil {
		return nil
	}

	if err := (*c.Connection).Close(); err != nil {
		return err
	}

	return nil
}

func (c *Client) Connect() (err error) {
	// Assure no other connection is on-going
	if c.Connection != nil {
		return fmt.Errorf("connection already established")
	}

	// Try to establish the connection multiple times
	if _, err := c.DialWithRetries(); err != nil {
		return fmt.Errorf(
			"failed to connect to %s %v",
			c.Settings.Host.Server(),
			err,
		)
	}

	// Connection established
	c.IsConnectionActive.Update(true)

	// Send Salute and check for error
	if err = c.Send(&mod.Salute{}); err != nil {
		defer c.Close()
		return err
	}

	// Receive Salute Response
	var res *mod.ResponseModel
	if res, err = c.Receive(); err != nil {
		defer c.Close()
		return err
	}

	// Check Salute Status
	if res.Status != mod.OK {
		defer c.Close()
		return fmt.Errorf("%v", res.Content)
	}

	// No errors occurred
	c.IsConnectionActive.Update(true)
	return nil
}

func (c *Client) DialWithRetries() (*net.Conn, error) {
	var (
		conn net.Conn
		err  error
	)

	for retry := c.Settings.MaxRetries; retry != 0; retry-- {
		c.Logger.Log("attempts left ", retry, " to connect...\n")

		conn, err = net.Dial(
			c.Settings.Host.Protocol,
			c.Settings.Host.Server(),
		)

		if err == nil {
			c.Connection = new(net.Conn)
			(*c.Connection) = conn
			c.Logger.Log("connection established\n")
			return c.Connection, nil
		}

		c.Logger.Log("Retrying in ", int(c.Settings.Timeout), " ms...\n")
		time.Sleep(time.Millisecond * c.Settings.Timeout)
	}

	return nil, err
}

func (c *Client) Send(req mod.Request) (err error) {
	if c.Connection == nil {
		return fmt.Errorf("connection does not exist")
	}

	// Map object to json request
	request := mod.RequestModel{
		Sender:  c.Settings.ClientName,
		Content: req.Content(),
		Type:    req.Type(),
	}

	// Transform to json
	var raw []byte
	if raw, err = json.Marshal(request); err != nil {
		return fmt.Errorf("could not convert obj to json: %v", err)
	}

	// Log the request that is being sent
	json := string(raw)
	c.Logger.Log("sending request " + json + "...\n")

	// Write to server
	if c.IsConnectionActive.Status() {
		fmt.Fprintf(*c.Connection, "%s\n", json)
	}

	// Return nil
	return nil
}

func (c *Client) Receive() (*mod.ResponseModel, error) {
	// Read message sent by the server
	reader := bufio.NewReader(*c.Connection)

	// Split by newline
	c.Logger.Log("waiting for responses...\n")
	input, err := reader.ReadString('\n')
	fmt.Printf("(Server) received response %v", input)

	// Check for errors (ex. if the connection is still on-going)
	if err != nil {
		c.IsConnectionActive.Update(false)
		return nil, err
	}

	// Parse the response
	var res *mod.ResponseModel = new(mod.ResponseModel)
	if err = json.Unmarshal([]byte(input), res); err != nil {
		return nil, err
	}

	// Return the response
	return res, nil
}
