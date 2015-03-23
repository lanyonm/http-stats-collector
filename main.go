// Application to collect HTTP stats and send them to various storage.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/cactus/go-statsd-client/statsd"
)

// run the server
func main() {
	var (
		port          = flag.Int("port", 8080, "Server listen port.")
		statsHostPort = flag.String("statsHostPort", "127.0.0.1:8125", "host:port of statsd server")
		statsPrefix   = flag.String("statsPrefix", "http-stats-collector", "the prefix used when stats are sent to statsd")
	)
	flag.Parse()

	// Create recorders and pass those to handlers.
	// Not sure if this is the best way to continue, but it's a point along
	// the process.
	var client *statsd.Client
	client, err := statsd.New(*statsHostPort, *statsPrefix)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	recorders := []Recorder{StatsDRecorder{client}}
	http.HandleFunc("/nav-timing", NavTimingHandler(recorders))
	http.HandleFunc("/csp-report", CSPReportHandler())

	log.Println("http-stats-collector: listening on port", *port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
