package pastebin

import (
	"fmt"
	"net/url"
)

type Config struct {
	APIDevKey          string `ini:"api_dev_key" comment:"API key"`
	APIPastePrivate    int    `ini:"api_paste_private"`
	APIPasteExpireDate string `ini:"api_paste_expire_date"`
}

func DefaultConfig() *Config {
	return &Config{APIPasteExpireDate: "N"}
}

func (config *Config) FormValues() url.Values {
	data := url.Values{}

	data.Set("api_dev_key", config.APIDevKey)
	data.Set("api_option", "paste")
	data.Set("api_paste_private", fmt.Sprintf("%d", config.APIPastePrivate))
	data.Set("api_paste_expire_date", config.APIPasteExpireDate)

	return data
}
