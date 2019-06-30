package main

type channel struct {
	name             string
	topic            string
	subscribed_users []*user
}
