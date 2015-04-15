package main

import (
	"flag"
	"github.com/vintik/tagger"
)

var s = tagger.Server{}

func init() {
	flag.StringVar(&s.URL, "url", ":8080", "url to listen on")
}

func main() {
	flag.Parse()
	s.Run()
}
