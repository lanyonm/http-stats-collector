package main

import (
	"log"
	"strings"

	"github.com/cactus/go-statsd-client/statsd"
)

// This interface is used to represent the functionality of clients/objects
// that send information to be persisted or further processed.
type Recorder interface {
	pushStat(stat string, value int64) bool
	cleanURI(input string) string
	validStat(stat string) bool
}

type StatsDRecorder struct {
	Client *statsd.Client
}

// Push stats to StatsD.
// This assumes that the data being written is always timing data and we are
// always collecting all the samples.
func (statsd StatsDRecorder) pushStat(stat string, value int64) bool {
	err := statsd.Client.Timing(stat, value, 1.0)
	if err != nil {
		log.Fatal("there was an error sending the statsd timing", err)
		return false
	}
	return true
}

// The valid page-uri checker for StatsD.  We don't want to accept anything
// that the storage would have trouble handing.
func (statsd StatsDRecorder) validStat(stat string) bool {
	return !strings.ContainsAny(stat, "&#") && strings.Index(stat, "//") == -1
}

// Clean up the page-uri.
// This will strip the leading and trailing /'s and replace the rest with '.'
func (statsd StatsDRecorder) cleanURI(input string) string {
	//replace rightmost / with "index"
	ret := input
	if last := len(ret) - 1; last >= 0 && ret[last] == '/' {
		ret = ret + "index"
	}
	ret = strings.TrimLeft(ret, "/")
	ret = strings.Split(ret, ".")[0]
	ret = strings.Replace(ret, "/", ".", -1)
	ret = strings.ToLower(ret)
	if len(ret) > 0 {
		ret = ret + "."
	}
	return ret
}
