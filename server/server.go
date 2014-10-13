package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Nothing bool

type Message struct {
	User   string
	Target string
	Msg    string
}

type ChatServer struct {
	port         string
	messageQueue map[string][]string
	users        []string
	shutdown     chan bool
}

// Register registers a client username with the chat server.
// It sends a message to all users notifying a user has joined
func (c *ChatServer) Register(username string, reply *string) error {
	*reply = "Welcome to GoChat v1.0!\n"
	*reply += "List of users online:\n"

	c.users = append(c.users, username)
	c.messageQueue[username] = nil

	for _, value := range c.users {
		*reply += value + "\n"
	}

	for k, _ := range c.messageQueue {
		c.messageQueue[k] = append(c.messageQueue[k], username+" has joined.")
	}

	log.Printf("%s has joined the chat.\n", username)

	return nil
}

func (c *ChatServer) CheckMessages(username string, reply *[]string) error {
	*reply = c.messageQueue[username]
	c.messageQueue[username] = nil
	return nil
}

func (c *ChatServer) List(none Nothing, reply *[]string) error {
	*reply = append(*reply, "Current online users:")

	for i := range c.users {
		*reply = append(*reply, c.users[i])
	}

	log.Println("Dumped list of users to client output")

	return nil
}

func (c *ChatServer) Tell(msg Message, reply *Nothing) error {

	if queue, ok := c.messageQueue[msg.Target]; ok {
		m := msg.User + " tells you " + msg.Msg
		c.messageQueue[msg.Target] = append(queue, m)
	} else {
		m := msg.Target + " does not exist"
		c.messageQueue[msg.User] = append(queue, m)
	}

	*reply = false

	return nil
}

func (c *ChatServer) Say(msg Message, reply *Nothing) error {

	for k, v := range c.messageQueue {
		m := msg.User + " says " + msg.Msg
		c.messageQueue[k] = append(v, m)
	}

	*reply = true

	return nil
}

func (c *ChatServer) Logout(username string, reply *Nothing) error {

	delete(c.messageQueue, username)

	for i := range c.users {
		if c.users[i] == username {
			c.users = append(c.users[:i], c.users[i+1:]...)
		}
	}

	for k, v := range c.messageQueue {
		c.messageQueue[k] = append(v, username+" has logged out.")
	}

	fmt.Println("User " + username + " has logged out.")

	*reply = false

	return nil
}

func (elt *ChatServer) Shutdown(nothing Nothing, reply *Nothing) error {

	log.Println("Server shutdown...Goodbye.")
	*reply = false
	elt.shutdown <- true

	return nil
}

func parseFlags(cs *ChatServer) {
	flag.StringVar(&cs.port, "port", "3410", "port for chat server to listen on")
	flag.Parse()

	cs.port = ":" + cs.port
}

func RunServer(cs *ChatServer) {
	rpc.Register(cs)
	rpc.HandleHTTP()

	log.Printf("Listening on port %s...\n", cs.port)

	l, err := net.Listen("tcp", cs.port)
	if err != nil {
		log.Panicf("Can't bind port to listen. %q", err)
	}

	go http.Serve(l, nil)
}

func main() {
	cs := new(ChatServer)
	cs.messageQueue = make(map[string][]string)
	cs.shutdown = make(chan bool, 1)

	parseFlags(cs)
	RunServer(cs)

	<-cs.shutdown
}
