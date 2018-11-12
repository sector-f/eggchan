package server

import (
	"log"
	"net/http"
	"time"
)

func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		var ip_addr string
		addrs, present := r.Header["X-Real-Ip"]
		if present {
			ip_addr = addrs[0]
		} else {
			ip_addr = r.RemoteAddr
		}

		log.Printf(
			"%s\t%s\t%s\t%s",
			ip_addr,
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

