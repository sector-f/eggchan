package eggchan

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type BoardService interface {
	ShowBoardReply(board string) (BoardReply, error)
	ShowThreadReply(board string, id int) (ThreadReply, error)

	ListCategories() ([]Category, error)
	ShowCategory(name string) ([]Board, error)
	ListBoards() ([]Board, error)
	ShowBoard(board string) ([]Thread, error)
	ShowThread(board string, id int) ([]Post, error)
	MakeThread(board, comment, author, subject string) (int, error)
	MakeComment(board string, thread int, comment string, author string) (int, error)
	ShowBoardDesc(board string) (Board, error)
	ShowThreadOP(board string, id int) (Thread, error)

	AddBoard(board, description, category string) error
	AddCategory(category string) error
	DeleteThread(board string, thread int) (int64, error)
	DeleteComment(board string, thread int) (int64, error)
}

type AuthService interface {
	AddUser(user, password string) error
	DeleteUser(user string) error
	ListUsers() ([]User, error)
	ValidatePassword(user string, password []byte) (bool, error)
	CheckPermission(user, permission string) (bool, error)
	GrantPermissions(user string, perms []Permission) error
	RevokePermissions(user string, perms []Permission) error
	ListPermissions() ([]Permission, error)
}

type Category struct {
	Name   string  `json:"name"`
	Boards []Board `json:"boards"`
}

type Board struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type Thread struct {
	Board           string      `json:"board"`
	PostNum         int         `json:"post_num"`
	Subject         null.String `json:"subject"`
	Author          string      `json:"author"`
	Time            time.Time   `json:"post_time"`
	NumReplies      int         `json:"num_replies"`
	SortLatestReply time.Time   `json:"latest_reply_time"`
	Comment         string      `json:"comment"`
}

type BoardReply struct {
	Board   Board    `json:"board"`
	Threads []Thread `json:"threads"`
}

type ThreadReply struct {
	Board  Board  `json:"board"`
	Thread Thread `json:"op"`
	Posts  []Post `json:"posts"`
}

type Post struct {
	ReplyTo int       `json:"reply_to"`
	PostNum int       `json:"post_num"`
	Author  string    `json:"author"`
	Time    time.Time `json:"time"`
	Comment string    `json:"comment"`
}

type User struct {
	Name  string
	Perms []string
}

type Permission struct {
	Name string
}

type PostThreadResponse struct {
	PostNum int `json:"post_num"`
}

type PostCommentResponse struct {
	ReplyTo int `json:"reply_to"`
	PostNum int `json:"post_num"`
}
