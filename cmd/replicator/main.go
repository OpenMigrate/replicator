package main

import (
	"net/http"
	"replicator/internal/api"
	"replicator/internal/storage"
	"replicator/logger"
)

func main() {
	logger.Init(logger.Options{Verbose: true, File: "", JSON: false})
	log := logger.Get()

	store, err := storage.Init("")
	if err != nil {
		log.Error("Unable to connect to db", err.Error)
	}

	log.Info("Replicate server started")
	r := api.NewRouter(store)

	log.Info("Listening on port 4000")
	err = http.ListenAndServe(":4000", r)
	if err != nil {
		log.Error(err.Error())
	}

}
