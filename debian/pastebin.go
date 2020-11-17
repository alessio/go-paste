package debian

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const baseURL = "https://paste.debian.net/"
const Name = "paste.debian.net"

var (
	// ErrPutFailed is returned when a paste could not be uploaded to pastebin.
	ErrPutFailed = errors.New("pastebin put failed")
	// ErrGetFailed is returned when a paste could not be fetched from pastebin.
	ErrGetFailed = errors.New("pastebin get failed")
)

// Pastebin represents an instance of the pastebin service.
type Pastebin struct {
	config *Config
}

func New(cfg *Config) *Pastebin {
	if cfg != nil {
		return &Pastebin{config: cfg}
	}

	return &Pastebin{config: DefaultConfig()}
}

// Put uploads text to Pastebin with optional title returning the ID or an error.
func (p *Pastebin) Put(text, title string) (id string, err error) {
	data := url.Values{}
	// Required values.
	data.Set("poster", p.config.Poster)
	data.Set("lang", "-1")
	data.Set("expire", "-1")
	data.Set("code", text)

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // or maybe the error from the request
		},
	}

	resp, err := client.PostForm(baseURL, data)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 302 {
		return resp.Header.Get("Location"), nil
	}

	return "", ErrPutFailed
}

// Get returns the text inside the paste identified by ID.
func (p *Pastebin) Get(id string) (text string, err error) {
	resp, err := http.DefaultClient.Get(baseURL + "raw.php?i=" + id)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%w: %s", ErrGetFailed, string(respBody))
	}

	return string(respBody), nil
}

// StripURL returns the paste ID from a pastebin URL.
func (p *Pastebin) StripURL(url string) string {
	return strings.ReplaceAll(url, baseURL, "")
}

// WrapID returns the pastebin URL from a paste ID.
func (p *Pastebin) WrapID(id string) string {
	return baseURL + id
}
