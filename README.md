# Eggchan
A headless textboard that exposes a JSON API which allows frontends to be written in any language.

## Routes
```
Verb   Route                              Returns     Description

GET    /                                              List routes
GET    /categories                        []Category  List categories
GET    /categories/{category}             Category    List boards in a specific category
GET    /boards                            []Board     List boards
GET    /boards/{board}                    BoardReply  List threads in a specific board
POST   /boards/{board}                                Post to a specific board
GET    /boards/{board}/{thread}           ThreadReply Show a specific thread
POST   /boards/{board}/{thread}                       Post to a specific thread
DELETE /boards/{board}/threads/{thread}               Delete a specific thread
DELETE /boards/{board}/comments/{comment}             Delete a specific comment
POST   /new/boards/{board}                            Create a new board
```

## Objects
```go
type Board struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type Category struct {
	Name   string  `json:"name"`
	Boards []Board `json:"boards"`
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
```
