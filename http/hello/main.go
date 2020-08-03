package main

import (
	"io"
	_log "log"
	"net/http"
	"os"
)

var (
	log = _log.New(os.Stderr, "", 0)

	worldURL = ""
)

func main() {
	worldURL := os.Getenv("RUN_WORLD_URL")
	if worldURL == "" {
		log.Fatal("RUN_WORLD_URL is required")
	}

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		for k, v := range r.URL.Query() {
			log.Printf("%s: %v", k, v)
		}

		res, err := http.Get(worldURL + "?text=hello")
		if err != nil {
			log.Println(err)
			http.Error(rw, "could not run", http.StatusInternalServerError)
			return
		}

		if res.StatusCode != 200 {
			log.Printf("unexpected status: %d", res.StatusCode)
			http.Error(rw, "run error", http.StatusInternalServerError)
			return
		}

		defer res.Body.Close()

		if _, err := io.Copy(rw, res.Body); err != nil {
			log.Println(err)
			http.Error(rw, "could not hello world", http.StatusInternalServerError)
			return
		}
	})

	addr := ":" + os.Getenv("PORT")

	log.Println("http-hello service is listening on ", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println(err)
	}
}
