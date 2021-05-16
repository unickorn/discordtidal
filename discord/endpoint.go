package discord

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// EndpointInterface is an interface for all interactions that have a method and a url.
// The JSON payload is the object "inheriting"(?) this interface, if passed to Do with hasBody=true.
type EndpointInterface interface {
	url() string
	method() string
}

// Do sends an EndpointInterface to Discord with the appropriate headers and returns the response.
func Do(e EndpointInterface, hasBody bool) (*http.Response, error) {
	var req *http.Request
	var err error

	if hasBody {
		j, err := json.Marshal(e)
		if err != nil {
			panic(err)
		}
		req, err = http.NewRequest(e.method(), e.url(), bytes.NewReader(j))
	} else {
		req, err = http.NewRequest(e.method(), e.url(), nil)
	}
	if err != nil {
		panic(err)
	}

	req.Header.Set("Host", "discord.com")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Authorization", config.Token)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", config.UserAgent)

	r, err := http.DefaultClient.Do(req)
	return r, err
}
