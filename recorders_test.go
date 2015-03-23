package main

import (
	"testing"
)

var (
	recorder = StatsDRecorder{}
)

var validStatTests = []struct {
	Input    string
	Expected bool
}{
	{"/", true},
	{"foo/bar", true},
	{"/foo/bar", true},
	{"/foo/bar/", true},
	{"/foo/bar///", false},
	{"/about.html", true},
	{"/company/about.png.php", true},
}

func TestStatsdValidStat(t *testing.T) {
	for _, tt := range validStatTests {
		if ret := recorder.validStat(tt.Input); ret != tt.Expected {
			t.Errorf("input was %v and expected %v, but got %v", tt.Input, tt.Expected, ret)
		}
	}
}

var uriTests = []struct {
	Input    string
	Expected string
}{
	{"/", "index."},
	{"foo/bar", "foo.bar."},
	{"/foo/bar", "foo.bar."},
	{"/foo/bar/", "foo.bar.index."},
	{"/about.html", "about."},
	{"/company/about.png.php", "company.about."},
}

func TestCleanURI(t *testing.T) {
	for _, tt := range uriTests {
		if ret := recorder.cleanURI(tt.Input); ret != tt.Expected {
			t.Errorf("input was %v and expected %v, but got %v", tt.Input, tt.Expected, ret)
		}
	}
}
