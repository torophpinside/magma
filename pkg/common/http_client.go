package common

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

var httpClientTor *http.Client

func NewHttpClient() *http.Client {
	if httpClientTor != nil {
		return httpClientTor
	}

	up, err := url.Parse("socks5://127.0.0.1:9050")
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to parse proxy URL: %v\n", err))
	}

	tr := &http.Transport{Proxy: http.ProxyURL(up)}
	c := &http.Client{Transport: tr}

	httpClientTor = c

	return httpClientTor
}
