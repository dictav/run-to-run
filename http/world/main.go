package main

import (
	_log "log"
	"net/http"
	"os"
)

var log = _log.New(os.Stderr, "", 0)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		txt := r.URL.Query().Get("text")
		rw.Write([]byte(txt + ", world!"))
	})

	addr := ":" + os.Getenv("PORT")

	log.Println("http-world is listening on ", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println(err)
	}
}
