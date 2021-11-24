# rpc-client-server
Client-Server App Using RPC to Communicate

# How to run the application
## Start the server
From the root execute the following command in a terminal:

```bash
go run ./server/src/main
```

## Start the client
In another therminal, fom the root, execute the following command to start a client:

```bash
go run ./client/src/main
```

# How to use the client
As the architecture is based on [Remote Procedure Call (RPC)](https://en.wikipedia.org/wiki/Remote_procedure_call) the following request types are supported:

`
// send a registration request to the server
// if the sender name is already registered the client will receive an error message
salute 

// close the connection with the server and remove the client from the registed user list
bye    

// send an acknowledgement
ack

// represents a command supported by the server (see below)
{command} args params
`

The following commands can be used by a client to communicate with the server:

`
// send a registration request to the server
// if the sender name is already registered the client will receive an error message
salute 

// close the connection with the server and remove the client from the registed user list
bye

// list the currently registed clients in the server
list clients

// solve the first problem
solve 1 ["casa","masa","trei","tanc","4321"]

// solve the second problem
solve 2 ["abd4g5","1sdf6fd","fd2fdsf5"]

// solve the third problem
solve 3 [12,13,14] 

// solve the eighth problem
solve 8 [23,17,15,3,18]
`
