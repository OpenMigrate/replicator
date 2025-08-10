package storage

import (
	"maps"
	"replicator/internal/models"
	"slices"
)

// TODO: to implement actual tinydb & mutex lock for the map
var (
  store = make(map[string]models.Metadata)
)

func SaveServer(md models.Metadata){
  store[md.ID] = md
}

func ListServers() []models.Metadata{
  return  slices.Collect(maps.Values(store))
}

func GetServer(id string) (models.Metadata, bool){
  v, ok := store[id]
	return v, ok
}
