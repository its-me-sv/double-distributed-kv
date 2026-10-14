package config

type Config struct {
	AppName string
}

func RetrieveConfig() Config {
	return Config{
		AppName: "double-distributed-kv",
	}
}
