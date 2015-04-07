package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type NavTimingReport struct {
	Details   NavTimingDetails `json:"nav-timing" statName:"navTiming"`
	Page      string           `json:"page-uri" statName:"pageUri"`
	Referer   string           `statName:"referer"`
	UserAgent string           `statName:"userAgent`
}

type NavTimingDetails struct {
	DNS      int64 `json:"dns" statName:"dns"`
	Connect  int64 `json:"connect" statName:"connect"`
	TTFB     int64 `json:"ttfb" statName:"ttfb"`
	BasePage int64 `json:"basePage" statName:"basePage"`
	FrontEnd int64 `json:"frontEnd" statName:"frontEnd"`
}

// for Navigation Timing API
func NavTimingHandler(recorders []Recorder) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var timing NavTimingReport

		if req.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			w.Header().Set("Allow", "POST")
			return
		}

		if err := json.NewDecoder(req.Body).Decode(&timing); err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		// You could consider this a flaw, but we don't send the stat anywhere
		// if it can't go to one of the recorders.
		for _, recorder := range recorders {
			if !recorder.validStat(timing.Page) {
				http.Error(w, "Invalid page-uri passed", http.StatusNotAcceptable)
				return
			}
		}

		// for each recorder we're sending all the NavTimingDetails stats
		t := reflect.TypeOf(timing.Details)
		v := reflect.ValueOf(timing.Details)
		for i := 0; i < v.NumField(); i++ {
			for _, recorder := range recorders {
				stat := recorder.cleanURI(timing.Page) + t.Field(i).Tag.Get("statName")
				val := v.Field(i).Int()
				recorder.pushStat(stat, val)
			}
		}
	}
}

type JsErrorReport struct {
	PageURI     string  `json:"page-uri"`
	QueryString string  `json:"query-string"`
	Details     JsError `json:"js-error"`
	ReportTime  time.Time
}

type JsError struct {
	UserAgent   string `json:"user-agent"`
	ErrorType   string `json:"error-type"`
	Description string `json:"description"`
}

func JsErrorReportHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var jsError JsErrorReport

		if req.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			w.Header().Set("Allow", "POST")
			return
		}

		if err := json.NewDecoder(req.Body).Decode(&jsError); err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		jsError.ReportTime = time.Now().UTC()

		// do something smart with the error
		dets, _ := json.Marshal(jsError)
		log.Println(strings.Join(req.Header["X-Real-Ip"], ""), "encountered a javascript error:", string(dets))
	}
}

type CSPReport struct {
	Details    CSPDetails `json:"csp-report" statName:"cspReport"`
	ReportTime time.Time  `statName:"dateTime"`
}

type CSPDetails struct {
	DocumentUri       string `json:"document-uri" statName:"documentUri" validate:"min=1,max=200"`
	Referrer          string `json:"referrer" statName:"referrer" validate:"max=200"`
	BlockedUri        string `json:"blocked-uri" statName:"blockedUri" validate:"max=200"`
	ViolatedDirective string `json:"violated-directive" statName:"violatedDirective" validate:"min=1,max=200,regexp=^[a-z0-9 '/\\*\\.:;-]+$"`
}

// for Content Security Policy
func CSPReportHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var report CSPReport

		if req.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			w.Header().Set("Allow", "POST")
			return
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(body, &report); err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		report.ReportTime = time.Now().UTC()

		// if validationError := validator.Validate(report); validationError != nil {
		// 	log.Println("Request failed validation:", validationError)
		// 	log.Println("Failed with report:", report)
		// 	http.Error(w, "Unable to validate JSON", http.StatusBadRequest)
		// 	return
		// }

		// do something smart with the report
		log.Println("policy violation from", strings.Join(req.Header["X-Real-Ip"], ""), "was:", report.Details.ViolatedDirective)
	}
}
