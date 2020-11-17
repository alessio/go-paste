// Package pastebin wraps the basic functions of the Pastebin API and exposes a
// Go API.
package pastebin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	pasteErrors "github.com/bearbin/go-paste/internal/errors"
)

const baseURL = "https://pastebin.com/"

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
	data := p.config.FormValues()
	data.Set("api_paste_code", text)
	data.Set("api_paste_name", title)

	resp, err := http.PostForm(baseURL+"api/api_post.php", data)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%w: %s", pasteErrors.ErrPutFailed, string(respBody))
	}

	return p.StripURL(string(respBody)), nil
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
		return "", fmt.Errorf("%w: %s", pasteErrors.ErrGetFailed, string(respBody))
	}

	return string(respBody), nil
}

// StripURL returns the paste ID from a pastebin URL.
func (p *Pastebin) StripURL(url string) string {
	return strings.ReplaceAll(url, baseURL, "")
}

// WrapID returns the pastebin URL from a paste ID.
func (p *Pastebin) WrapID(id string) string {
	return "http://pastebin.com/" + id
}
