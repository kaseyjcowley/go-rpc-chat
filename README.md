# Gotta love Golang
This is a simple Golang implementation of an RPC chat client/server. By running the server, any number of users can join the chat room and send messages to each other.

**Starting the Server**

Make sure to build the server by running `go build server.go`. Once it's built, run the server with usage `./server <host>:<port>`. If you omit the host and port, it will default to port 3410 on localhost.

Once the server runs, it will accept connections at the specified host and port. Throughout the life of the chat, the server will log/report on the activity. 

Ex. `2014-10-12 09:00:00 fred has joined the chat`

**Starting the client and connecting to the server**

As with the server, make sure to build the client by running `go build client.go`. Once it's build, you run the client in like manner to the server with usage `./client -user <username> -host <host>:<port>`. If the user and host are ommitted, it defaults to fred and localhost:3410.

**Client API**

* `list` - Lists the users in the chat
* `tell <user> <msg>` - Send a private message to a user
* `say <msg>` - Tell all users a message
* `logout` - Logout of the chat
* `shutdown` - Kills the chat for everyone (weird...I know)