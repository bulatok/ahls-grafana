package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	HTTP_ADDR = "0.0.0.0:9001"
)

func ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"message": "pong"}`))
	}
}

func work() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// query parsing...
		statusCode := r.URL.Query().Get("code")
		code, err := strconv.Atoi(statusCode)
		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
			return
		}

		sleepTime := r.URL.Query().Get("sleep")
		sleepTimeSec, err := strconv.Atoi(sleepTime)
		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
			return
		}

		// business logic
		log.Println(fmt.Sprintf("starting working %ds...", sleepTimeSec))
		time.Sleep(time.Duration(sleepTimeSec) * time.Second)

		// response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write([]byte(`{"result":"ok"}`))
	}
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func applyMiddleware(next http.Handler) http.Handler {
	return logMiddleware(next)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	log.Println("starting v2")

	// metrics
	metrics := NewMetrics()
	if err := metrics.Register(); err != nil {
		log.Fatal(err)
	}

	// handlers
	http.Handle("/metrics", logMiddleware(promhttp.Handler()))
	http.Handle("/ping", metrics.ApplyMiddleware(applyMiddleware(ping())))
	http.Handle("/work", metrics.ApplyMiddleware(applyMiddleware(work())))

	// starting
	log.Printf("starting listening on %s http addr", HTTP_ADDR)
	if err := http.ListenAndServe(HTTP_ADDR, nil); err != nil {
		log.Fatal(err)
	}
}
