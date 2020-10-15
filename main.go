package main

import (
	"log"
	"net/http"
	"time"

	"nicheanal.com/config"
	"nicheanal.com/pkg/logger"
	r "nicheanal.com/router"
)

func main() {
	config.LoadConfig()
	logger.Init()

	router := r.GetRouter()
	s := &http.Server{
		Addr:           r.GetPort(),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.SetKeepAlivesEnabled(false)
	log.Printf("Listening on port %s", r.GetPort())
	log.Fatal(s.ListenAndServe())
}
