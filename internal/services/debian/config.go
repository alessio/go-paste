package debian

import (
	"net/url"
	"os/user"
)

type Config struct {
	Poster   string `ini:"poster"`
	Language string `ini:"language"`
	Expire   string `ini:"expire"`
	Private  string `ini:"private"`
	Wrap     string `ini:"wrap"`
}

func DefaultConfig() *Config {
	return &Config{Poster: defaultPoster(), Language: "1", Expire: "-1"}
}

func (config *Config) FormValues() url.Values {
	data := url.Values{}

	data.Set("poster", config.Poster)
	data.Set("lang", config.Language)
	data.Set("expire", config.Expire)
	data.Set("private", config.Private)
	data.Set("wrap", config.Wrap)

	return data
}

func defaultPoster() string {
	u, err := user.Current()
	if err != nil {
		return "anonymous"
	}

	return u.Username
}
