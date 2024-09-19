package main

import (
	"fetch-interview/internal/config"
	"fetch-interview/internal/endpoint"
	"fetch-interview/internal/instrumentation"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/jhunt/go-log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

type DomainChecks struct {
	TotalChecks int
	UpChecks    int

	mu sync.Mutex
}

func main() {
	// Read from config file
	var (
		configFile = kingpin.Flag(
			"config",
			"Configuration file path.",
		).Default("config.yaml").String()

		logLevel = kingpin.Flag(
			"log-level",
			"Set log level to debug, warn, info or error.",
		).Default("info").String()

		metricsPort = kingpin.Flag(
			"metrics",
			"Port to expose prometheus metrics.",
		).Default("8080").Int()
	)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	// Set up log level
	log.SetupLogging(log.LogConfig{
		Type:  "console",
		Level: *logLevel,
	})

	// Initialize config
	config, err := config.NewConfig(*configFile)
	if err != nil {
		log.Errorf("Error reading configuration file: %s", err)
		os.Exit(1)
	}

	metrics := instrumentation.NewInstrumentation()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(fmt.Sprintf(":%d", *metricsPort), nil)
		log.Debugf("Prometheus server listening on port %d", *metricsPort)
	}()

	var endpoints []endpoint.Endpoint
	for _, conf := range config {
		endpoints = append(endpoints, endpoint.InitializeEndpoint(conf.Name, conf.URL, conf.Method, conf.Body, conf.Headers))
	}

	// Run checks every 15 seconds
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	checks := make(map[string]*DomainChecks)

	var wg sync.WaitGroup
	for {
		for _, ep := range endpoints {
			url, _ := url.Parse(ep.GetURL())
			host := url.Host
			if _, found := checks[host]; !found {
				checks[host] = &DomainChecks{}
			}

			// Do checks in parallel. Use locks for synchronizing read/write to the same dictionary key.
			wg.Add(1)
			go func() {
				defer wg.Done()
				isUp := ep.PerformHealthCheck(metrics)
				checks[host].mu.Lock()
				checks[host].TotalChecks += 1
				if isUp {
					checks[host].UpChecks += 1
				}
				checks[host].mu.Unlock()
			}()
		}

		wg.Wait()
		logAvailability(checks)

		select {
		case <-ticker.C:
			continue
		}
	}
}

func logAvailability(checks map[string]*DomainChecks) {
	for domain, check := range checks {
		check.mu.Lock()
		log.Debugf("%s:  up=%d total=%d", domain, check.UpChecks, check.TotalChecks)
		result := float64(check.UpChecks) / float64(check.TotalChecks)
		avaibality := 100 * result
		log.Infof("%s has %d%% availability", domain, int(math.Round(avaibality)))
		check.mu.Unlock()
	}
}
