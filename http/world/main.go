package main

import (
	_log "log"
	"net/http"
	"os"
)

var log = _log.New(os.Stderr, "", 0) //nolint:gochecknoglobals

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		for k, v := range r.URL.Query() {
			log.Printf("%s: %v", k, v)
		}

		txt := r.URL.Query().Get("text")
		_, err := rw.Write([]byte(txt + ", world!"))
		if err != nil {
			log.Println(err)
			http.Error(rw, "could not write", http.StatusInternalServerError)
		}
	})

	addr := ":" + os.Getenv("PORT")

	log.Println("http-world is listening on ", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println(err)
	}
}
