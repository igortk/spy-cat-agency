package utils

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type CatBreed struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
	cachedBreedsMap map[string]string
	lastFetchTime   time.Time

	catAPIURL string
)

func InitializeConfig(url string) {
	catAPIURL = url
}

func IsValidBreed(breed string) bool {
	if time.Since(lastFetchTime) > 24*time.Hour || cachedBreedsMap == nil {
		fetchBreeds()
	}

	_, ok := cachedBreedsMap[breed]
	return ok
}

func fetchBreeds() {
	resp, err := http.Get(catAPIURL)
	if err != nil {
		log.Errorf("Failed to fetch cat breeds: %v", err)
		return
	}
	defer resp.Body.Close()

	var breeds []CatBreed
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		log.Errorf("Failed to decode cat breeds response: %v", err)
		return
	}

	newBreedsMap := make(map[string]string, len(breeds))
	for _, b := range breeds {
		newBreedsMap[b.Name] = b.ID
	}

	cachedBreedsMap = newBreedsMap
	lastFetchTime = time.Now()
	log.Infof("Successfully fetched and cached %d cat breeds from TheCatAPI.", len(cachedBreedsMap))
}
