package main

import (
  "fmt"
  "net/http"
  "os"
  "replicator/config"
  "replicator/internal/api"
  "replicator/internal/storage"
  "replicator/logger"
)

func main() {
  cfg := config.LoadConfig()

  _, _, err := logger.Init(logger.Options{Verbose: cfg.Verbose, File: cfg.LogPath, JSON: cfg.JSON})
  if err != nil {
    fmt.Fprintln(os.Stderr, "logger init:", err)
    os.Exit(1)
  }
  log := logger.Get()

  store, err := storage.Init(cfg.DBURL)
  if err != nil {
    log.Error("Unable to connect to db", "msg", err.Error)
  }
  log.Info("db", "data", store)

  log.Info("Replicate server started")
  r := api.NewRouter(store)

  log.Info("Listening on port 4000")
  err = http.ListenAndServe(":4000", r)
  if err != nil {
    log.Error(err.Error())
  }

}
