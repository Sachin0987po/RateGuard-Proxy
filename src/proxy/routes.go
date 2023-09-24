package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/proxy-server-rateLimiter/config"
	"github.com/proxy-server-rateLimiter/ratelimiter"
)

var customTransport = http.DefaultTransport
var targetIP = "127.0.0.1"
var targetPort = "8081"

func getTargetUrl(r *http.Request) *url.URL {
	return &url.URL{
		Scheme:   "http",
		Host:     fmt.Sprintf("%s:%s", targetIP, targetPort),
		Path:     r.URL.Path,
		RawQuery: r.URL.RawQuery,
	}
}

func isApiUnderRateLimit(endpoint config.Endpoint, key string) bool {
	key = key + "-" + fmt.Sprintf("%d", endpoint.Id)
	return ratelimiter.RateLimiterHandler(key, endpoint)
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {

	ep := config.Endpoint{}
	if !config.GetEndpointDetail(r.URL.Path, &ep) {
		http.Error(w, "Endpoint not found", http.StatusBadRequest)
		return
	}

	key := r.Header.Get("api-key")
	if key == "" {
		http.Error(w, "API key is missing", http.StatusUnauthorized)
		return
	}

	if !isApiUnderRateLimit(ep, key) {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	targetURL := getTargetUrl(r)
	proxyReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	// Copy the headers from the original request to the proxy request
	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	// Send the proxy request using the custom transport
	resp, err := customTransport.RoundTrip(proxyReq)
	if err != nil {
		http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the headers from the proxy response to the original response
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Set the status code of the original response to the status code of the proxy response
	w.WriteHeader(resp.StatusCode)

	// Copy the body of the proxy response to the original response
	io.Copy(w, resp.Body)
}
