package pastebin

type Config struct {
	APIDevKey       string `ini:"api_dev_key" comment:"API key"`
	APIPastePrivate int    `ini:"api_paste_private"`
}

func DefaultConfig() *Config {
	return &Config{}
}
