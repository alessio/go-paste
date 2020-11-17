package ubuntu

import (
	"net/url"
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

func (config *Config) FormValues() url.Values {
	data := url.Values{}

	data.Set("poster", config.Poster)
	data.Set("syntax", config.Syntax)
	data.Set("expiration", config.Expiration)

	return data
}

func defaultPoster() string {
	u, err := user.Current()
	if err != nil {
		return "anonymous"
	}

	return u.Username
}
