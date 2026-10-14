package store

import (
	"errors"
)

type Store struct {
	dictionary map[string]string
}

func (s Store) Exists(key string) bool {
	_, doesExist := s.dictionary[key]
	return doesExist
}

func (s Store) Set(key, value string) bool {
	wasUpdated := s.Exists(key)
	s.dictionary[key] = value
	return wasUpdated
}

func (s Store) Get(key string) (string, error) {
	if !s.Exists(key) {
		return "", errors.New("key not found")
	}
	return s.dictionary[key], nil
}

func (s Store) Del(key string) bool {
	if !s.Exists(key) {
		return false
	}
	delete(s.dictionary, key)
	return true
}
