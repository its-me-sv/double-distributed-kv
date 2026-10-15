package store

import (
	"errors"
	"fmt"
	"os"
)

type Store struct {
	dictionary map[string]string
}

func CreateNewStore() Store {
	return Store{
		dictionary: map[string]string{},
	}
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

func (s Store) Debug() {
	fmt.Fprintln(os.Stdout, s)
}
