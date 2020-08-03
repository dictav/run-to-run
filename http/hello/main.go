package main

import (
	"io"
	_log "log"
	"net/http"
	"os"

	"cloud.google.com/go/compute/metadata"
)

var (
	log = _log.New(os.Stderr, "", 0) //nolint:gochecknoglobals
)

func main() {
	worldURL := os.Getenv("RUN_WORLD_URL")
	if worldURL == "" {
		log.Fatal("RUN_WORLD_URL is required")
	}

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("headers (%d):", len(r.Header))
		for k, v := range r.Header {
			log.Printf("%s: %v", k, v)
		}

		tokenURL := "/instance/service-accounts/default/identity?audience=" + worldURL
		idToken, err := metadata.Get(tokenURL)
		if err != nil {
			log.Println(err)
			http.Error(rw, "could not create token", http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest(http.MethodGet, worldURL+"?text=hello", nil)
		if err != nil {
			log.Println(err)
			http.Error(rw, "invalid request", http.StatusInternalServerError)
			return
		}

		req.Header.Add("Authorization", "Bearer "+idToken)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			http.Error(rw, "could not run", http.StatusInternalServerError)
			return
		}

		if res.StatusCode != http.StatusOK {
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
