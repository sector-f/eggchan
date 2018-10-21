package main

import (
	"os"
)

func main() {
	app := App{}
	app.Initialize(
		os.Getenv("EGGCHAN_DB_USERNAME"),
		os.Getenv("EGGCHAN_DB_PASSWORD"),
		os.Getenv("EGGCHAN_DB_NAME"),
	)
	app.Run(os.Args[1])
}
