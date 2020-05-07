package utils

import (
	"crypto/tls"
	"net/http"
	"time"
)

type Client struct {
	http.Client
}

func NewClient() *Client {
	return configureClient()
}

func (client *Client) MakeGetRequest(requestUrl string) *http.Response {
	retry := func(requestUrl string) *http.Response {
		return client.MakeGetRequest(requestUrl)
	}

	req, _ := http.NewRequest("GET", requestUrl, nil)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; rv:60.0) Gecko/20100101 Firefox/60.0")
	req.Header.Set("accept", "*/*")
	req.Header.Del("Accept-Encoding")

	resp, err := client.Do(req)
	if err != nil || resp == nil || resp.StatusCode == http.StatusForbidden {
		time.Sleep(time.Duration(10) * time.Second)
		return retry(requestUrl)
	}

	time.Sleep(time.Duration(2) * time.Second)
	return resp
}

func configureClient() *Client {
	tbTransport := &http.Transport{DisableCompression: true}
	tbTransport.MaxIdleConns = 100
	tbTransport.MaxIdleConnsPerHost = 100
	tbTransport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	tbTransport.TLSHandshakeTimeout = 15 * time.Second
	timeout := time.Duration(30 * time.Second)
	return &Client{http.Client{Transport: tbTransport, Timeout: timeout}}
}
