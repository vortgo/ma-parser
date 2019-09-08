package tor

import (
	"crypto/tls"
	"github.com/sycamoreone/orc/control"
	"github.com/vortgo/ma-parser/logger"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type Client struct {
	http.Client
}

var lastUpdTimestamp int64
var mutex = &sync.Mutex{}

const userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36"

func init() {
	lastUpdTimestamp = time.Now().Unix() - 1000
}

func NewClient() *Client {
	return configureClient()
}

func (client *Client) MakeGetRequest(requestUrl string) *http.Response {
	log := logger.New()
	retry := func(requestUrl string) *http.Response {
		client.RenewIP()
		return client.MakeGetRequest(requestUrl)
	}

	req, _ := http.NewRequest("GET", requestUrl, nil)

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil || resp == nil {
		log.SetData(logger.Data{
			"request_url": requestUrl,
			"culprit":     "Requester",
		}).Warningf("Failed GET request - %s", err)

		return retry(requestUrl)
	}

	if resp.StatusCode != http.StatusOK {

		log.SetData(logger.Data{
			"culprit":     "Requester",
			"request_url": requestUrl,
			"status_code": resp.StatusCode,
		}).Warning("Failed GET request - status <> 200")

		return retry(requestUrl)
	}

	return resp
}

func (client *Client) RenewIP() {
	time.Sleep(1 * time.Second)
	var log = logger.New()
	currentTimestamp := time.Now().Unix()
	mutex.Lock()
	if currentTimestamp > lastUpdTimestamp+10 {
		c, err := control.Dial(os.Getenv("TOR_CONTROL_URL"))
		if err != nil {
			log.SetData(logger.Data{
				"culprit": "Tor",
				"command": "set connection",
			}).Warning(err)

			client.RenewIP()
			return
		}

		err = c.Auth("secret-password-tor")
		if err != nil {
			log.SetData(logger.Data{
				"culprit": "Tor",
				"command": "auth",
			}).Warning(err)
			return
		}

		err = c.Signal(control.SignalNewNym)
		if err != nil {
			log.SetData(logger.Data{
				"culprit": "Tor",
				"command": "signal",
			}).Warning(err)
		}

		time.Sleep(1 * time.Second)
		lastUpdTimestamp = time.Now().Unix()
	}

	*client = *NewClient()
	mutex.Unlock()
}

func configureClient() *Client {
	var log = logger.New()
	tbProxyURL, err := url.Parse(os.Getenv("TOR_PROXY_URL"))
	if err != nil {
		log.Warningf("Failed to parse proxy URL: %v/n", err)
	}

	tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
	if err != nil {
		log.Warningf("Failed to obtain proxy dialer: %v\n", err)
	}

	tbTransport := &http.Transport{Dial: tbDialer.Dial}
	tbTransport.MaxIdleConns = 100
	tbTransport.MaxIdleConnsPerHost = 100
	tbTransport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	tbTransport.TLSHandshakeTimeout = 15 * time.Second
	timeout := time.Duration(30 * time.Second)
	return &Client{http.Client{Transport: tbTransport, Timeout: timeout}}
}
