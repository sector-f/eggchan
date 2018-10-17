package main

import (
	_ "fmt"
	_ "github.com/gorilla/mux"
	"os"
	// "github.com/lib/pq"
	_ "log"
	"net/http"
)

// [ ] List boards: `GET /boards`
// [ ] List threads on a board: `GET /boards/<name>`
// [ ] Display a thread: `GET /boards/<board>/<id>`
// [ ] Post a thread: `POST /boards/<board>` (multipart/form-data using "comment" field)
// [ ] Post a comment: `POST /boards/<board>/<id>` (multipart/form-data using "comment" field)
// [x] List categories: `GET /categories`
// [ ] List boards in a category: `GET /categories/<category>`

func Unimplemented(_ http.ResponseWriter, _ *http.Request) {
}

func main() {
	app := App{}
	app.Initialize(
		os.Getenv("EGGCHAN_DB_USERNAME"),
		os.Getenv("EGGCHAN_DB_PASSWORD"),
		os.Getenv("EGGCHAN_DB_NAME"),
	)
	app.Run("127.0.0.1:8005")
}
