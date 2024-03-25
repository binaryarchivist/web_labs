package utils

import (
	"encoding/json"
	"os"
	"sync"
)

const cacheFilePath = "cache.json"

var (
	cache  map[string]string
	mutex  sync.Mutex
	loaded = false
)

func loadCache() error {
	mutex.Lock()
	defer mutex.Unlock()

	if loaded {
		return nil
	}

	cache = make(map[string]string)

	file, err := os.ReadFile(cacheFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	err = json.Unmarshal(file, &cache)
	if err != nil {
		return err
	}

	loaded = true
	return nil
}

func saveCache() error {
	mutex.Lock()
	defer mutex.Unlock()

	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	err = os.WriteFile(cacheFilePath, data, 0644)
	return err
}

func Set(key, value string) error {
	if err := loadCache(); err != nil {
		return err
	}

	cache[key] = value
	return saveCache()
}

func Get(key string) (string, bool, error) {
	if err := loadCache(); err != nil {
		return "", false, err
	}

	val, found := cache[key]
	return val, found, nil
}
