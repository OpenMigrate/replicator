package main

import (
	"net/http"
	"replicator/internal"
	"replicator/logger"
)

func main() {
  logger.Init(logger.Options{Verbose:true, File:"", JSON:false});
  log := logger.Get();

  log.Info("Replicate server started")
  r := internal.NewRouter()

  log.Info("Listening on port 4000")
  err := http.ListenAndServe(":4000", r)
  if err != nil {
    log.Error(err.Error())
  }

}
