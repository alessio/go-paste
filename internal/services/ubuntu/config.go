package ubuntu

import (
	"os/user"
)

type Config struct {
	Poster     string `ini:"poster"`
	Syntax     string `ini:"syntax"`
	Expiration string `ini:"expiration"`
}

func DefaultConfig() *Config {
	return &Config{Poster: defaultPoster(), Syntax: "text"}
}

func defaultPoster() string {
	u, err := user.Current()
	if err != nil {
		return "anonymous"
	}

	return u.Username
}
