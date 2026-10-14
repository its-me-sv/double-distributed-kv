package main

import (
	"github.com/its-me-sv/double-distributed-kv/internal/config"
	"github.com/its-me-sv/double-distributed-kv/internal/repl"
	"github.com/its-me-sv/double-distributed-kv/internal/store"
)

func main() {
	str := store.CreateNewStore()
	cfg := config.RetrieveConfig()

	repl.StartRepl(cfg, str)
}

type Config struct {
	appName string
}
