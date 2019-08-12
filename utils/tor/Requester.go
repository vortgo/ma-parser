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

func init() {
	lastUpdTimestamp = time.Now().Unix() - 1000
}

func NewClient() *Client {
	return configureClient()
}

func (client *Client) MakeGetRequest(requestUrl string) *http.Response {
	var log = logger.New()
	retry := func(requestUrl string) *http.Response {
		client.RenewIP()
		return client.MakeGetRequest(requestUrl)
	}

	resp, err := client.Get(requestUrl)
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
	log.Info("RenewIP")
	currentTimestamp := time.Now().Unix()
	mutex.Lock()
	log.Info("RenewIP lock")
	if currentTimestamp > lastUpdTimestamp+10 {
		log.Info("RenewIP currentTimestamp > lastUpdTimestamp+10")
		c, err := control.Dial(os.Getenv("TOR_CONTROL_URL"))
		if err != nil {
			log.SetData(logger.Data{
				"culprit": "Tor",
				"command": "set connection",
			}).Warning(err)

			client.RenewIP()
			return
		}

		log.Info("RenewIP controll success")

		err = c.Auth("secret-password-tor")
		if err != nil {
			log.SetData(logger.Data{
				"culprit": "Tor",
				"command": "auth",
			}).Warning(err)
			return
		}

		log.Info("RenewIP auth success")

		err = c.Signal(control.SignalNewNym)
		if err != nil {
			log.SetData(logger.Data{
				"culprit": "Tor",
				"command": "signal",
			}).Warning(err)
		}

		log.Info("RenewIP signal success")
		time.Sleep(1 * time.Second)
		lastUpdTimestamp = time.Now().Unix()

		log.Info("RenewIP ud timestamp")
	}

	*client = *NewClient()
	log.Info("RenewIP set new client")
	mutex.Unlock()

	log.Info("RenewIP - SUCCESS")
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
	tbTransport.TLSHandshakeTimeout = 10 * time.Second
	timeout := time.Duration(30 * time.Second)
	return &Client{http.Client{Transport: tbTransport, Timeout: timeout}}
}
