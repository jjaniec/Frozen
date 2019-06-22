package main

import (
	"fmt"
)

func (c *connection) handle_cmd_pass(password string) {
	// https://tools.ietf.org/html/rfc1459#section-4.1.1
	c.session.password = password
}

func (c *connection) handle_cmd_user(username string, hostname string, servername string, realname string) (string){
	// https://tools.ietf.org/html/rfc1459#section-4.1.3
	fmt.Println("TODO user")
	return ":Welcome to the Internet Relay Network"
}
