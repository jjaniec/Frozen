package main

import (
	"fmt"
	"strings"
)

type user struct {
	username string
	nickname string
	password string
	client	*connection
}

func (u *user) send_message(msg string) {
	fmt.Println("User: ", u.nickname, " sent: ", msg)
}

func (u *user) update_nickname(conn_words []string) (response string) {
	fmt.Println(conn_words)
	if len(conn_words) == 1 {
		return "No nickname specified, usage: NICK new_nickname"
	}
	new_nickname := conn_words[1]
	if len(new_nickname) == 0 || len(strings.Fields(new_nickname)) == 0 {
		return "Invalid arguments"
	}
	// Check if nickname already taken
	for _, e := range current_users {
		if (e.nickname == new_nickname) {
			return "Nickname already in use."
		}
	}

	fmt.Println("User nickname: ", u.nickname, " updated as ", new_nickname)
	u.nickname = new_nickname
	current_users = append(current_users, u)
	return
}
