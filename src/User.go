package main

import (
	"fmt"
	// "strings"
)

type user struct {
	username string
	nickname string
	password string
	realname string
	client   *connection
}

func (u *user) send_message(msg string) {
	fmt.Println("User: ", u.nickname, " sent: ", msg)
}
