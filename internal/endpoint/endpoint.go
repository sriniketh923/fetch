package endpoint

import (
	"fetch-interview/internal/instrumentation"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jhunt/go-log"
	"github.com/prometheus/client_golang/prometheus"
)

type Endpoint struct {
	name    string            `yaml:"name"`
	url     string            `yaml:"url"`
	method  string            `yaml:"string,omitempty"`
	headers map[string]string `yaml:"headers,omitempty"`
	body    string            `yaml:"body,omitempty"`
}

// Initialize the endpoint struct
func InitializeEndpoint(name, url, method, body string, headers map[string]string) Endpoint {
	return Endpoint{
		name:    name,
		url:     url,
		method:  method,
		headers: headers,
		body:    body,
	}
}

func (ep Endpoint) PerformHealthCheck(metrics instrumentation.Instrumentation) bool {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	req, err := http.NewRequest(ep.GetMethod(), ep.GetURL(), strings.NewReader(ep.GetBody()))
	if err != nil {
		log.Errorf("Error creating request for %s and url %s: %s", ep.GetName(), ep.GetURL(), err)
		return false
	}

	for key, value := range ep.GetHeaders() {
		req.Header.Add(key, value)
	}

	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start)

	if err != nil {
		log.Errorf("Error performing request '%s': %s", ep.GetName(), err)
		return false
	}
	defer resp.Body.Close()
	metrics.ApiLatency.With(prometheus.Labels{"name": ep.GetName(), "url": ep.GetURL()}).Observe(float64(latency.Milliseconds()))
	metrics.ApiCalls.With(prometheus.Labels{"name": ep.GetName(), "url": ep.GetURL(), "status_code": strconv.Itoa(resp.StatusCode)}).Inc()
	log.Debugf("Name: %s  StatusCode: %d  Latency (ms): %d", ep.GetName(), resp.StatusCode, latency.Milliseconds())
	return resp.StatusCode >= 200 && resp.StatusCode < 300 && latency < 500*time.Millisecond
}

func (ep Endpoint) GetName() string {
	return ep.name
}

func (ep Endpoint) GetURL() string {
	return ep.url
}

func (ep Endpoint) GetMethod() string {
	return ep.method
}

func (ep Endpoint) GetHeaders() map[string]string {
	return ep.headers
}

func (ep Endpoint) GetBody() string {
	return ep.body
}
