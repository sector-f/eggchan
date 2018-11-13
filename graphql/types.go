package graphql

import "time"

type Category struct {
	name   string
	boards []string
}

type Board struct {
	name        string
	description string
	category    string
	threads     []int
}

type Thread struct {
	board           string
	postNum         int
	subject         string
	author          string
	time            time.Time
	numReplies      int
	latestReplyTime time.Time
	comment         string
	posts           []int
}

type Post struct {
	postNum int
	author  string
	time    time.Time
	comment string
}
