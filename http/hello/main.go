package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		url := os.Getenv("RUN_WORLD_ADDR")
		if url == "" {
			http.Error(rw, "can not run", http.StatusInternalServerError)
			return
		}

		res, err := http.Get(url + "?text=hello")
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

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Println(err)
	}
}
