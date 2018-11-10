package eggchan

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type EggchanService interface {
	ListCategories() ([]Category, error)
	ShowCategory(name string) ([]Board, error)
	ListBoards() ([]Board, error)
	ShowBoard(board string) (BoardReply, error)
	ShowThread(board string, id int) (ThreadReply, error)
	MakeThread(board, comment, author subject string) (int, error)
	MakeComment(board string, thread int, comment string, author string)
}

type AdminService interface {
	AddUser(user string) error
	DeleteUser(user string) error
	ListUsers() ([]User, error)
	GrantPermissions(user string, perms []Permission) error
	RevokePermissions(user string, perms []Permission) error
	ListPermissions() ([]Permission, error)
	AddBoard(board, description, category string) error
	DeleteThread(board string, thread int) (int64, error)
	DeleteComment(board string, thread int) (int64, error)
}

type Category struct {
	Name string `json:"name"`
}

type BoardReply struct {
	Board   Board    `json:"board"`
	Threads []Thread `json:"threads"`
}

type ThreadReply struct {
	Thread Thread `json:"op"`
	Posts  []Post `json:"posts"`
}

type Board struct {
	Name        string      `json:"name"`
	Description null.String `json:"description"`
	Category    null.String `json:"category"`
}

type Thread struct {
	PostNum         int         `json:"post_num"`
	Subject         null.String `json:"subject"`
	Author          string      `json:"author"`
	Time            time.Time   `json:"post_time"`
	NumReplies      int         `json:"num_replies"`
	SortLatestReply time.Time   `json:"latest_reply_time"`
	Comment         string      `json:"comment"`
}

type Post struct {
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
