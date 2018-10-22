package main

import (
	"flag"
	"os"
)

func main() {
	addrPtr := flag.String("bind", "127.0.0.1:8000", "Address/port to bind to (Default: 127.0.0.1:8000)")
	flag.Parse()

	server := Server{}
	server.Initialize(
		os.Getenv("EGGCHAN_DB_USERNAME"),
		os.Getenv("EGGCHAN_DB_PASSWORD"),
		os.Getenv("EGGCHAN_DB_NAME"),
	)
	server.Run(*addrPtr)
}
