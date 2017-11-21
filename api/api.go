package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/powerman/go-service-narada4d-example/dal"
	"github.com/powerman/structlog"
)

var log = structlog.New()

var key, ver string

// Init must be called before using this package.
func Init(apiKey, appVer string) error {
	key, ver = apiKey, appVer
	http.HandleFunc("/", handler)
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("key") != key {
		log.Warn("wrong API key")
		http.Error(w, "Wrong API key", http.StatusForbidden)
		return
	}

	fmt.Fprintf(w, "%s\n\n", ver)

	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "Path = %q\n", r.URL.Path)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	fmt.Fprintf(w, "X-Forwarded-For: %q\n", r.Header.Get("X-Forwarded-For"))
	fmt.Fprintf(w, "X-Forwarded-Proto: %q\n", r.Header.Get("X-Forwarded-Proto"))
	fmt.Fprintf(w, "X-Forwarded-HTTPS: %q\n\n", r.Header.Get("X-Forwarded-HTTPS"))

	counter, err := dal.Count()
	fmt.Fprintf(w, "Counter: %d (err=%v)\n\n", counter, err)

	t := time.Now()
	resp, err := http.Get("https://google.com/")
	var body []byte
	if err == nil {
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
	}
	fmt.Fprintf(w, "Fetched https://google.com/: %d bytes in %s (err=%v)\n\n", len(body), time.Since(t), err)

	log.Info("request handled")
}
