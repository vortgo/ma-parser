package utils

import (
	"net/http"
)

var CustomHttpClient *http.Client

func init() {
	// Customize the Transport to have larger connection pool
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer := defaultRoundTripper.(*http.Transport)
	defaultTransport := *defaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
	defaultTransport.MaxIdleConns = 100
	defaultTransport.MaxIdleConnsPerHost = 100

	CustomHttpClient = &http.Client{Transport: &defaultTransport}
}
