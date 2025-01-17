package tagger

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Server struct {
	URL string
}

func (s *Server) doit(w http.ResponseWriter, r *http.Request) {
	counts := make(map[string]int)

	err := error(nil)

	defer func() {
		if err != nil && err != io.EOF {
			_, _ = w.Write([]byte(err.Error()))
			log.Printf("Deferred error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
	}()

	err = r.ParseForm()
	if err != nil {
		return
	}

	u := r.Form.Get("url")
	if u == "" {
		err = errors.New("Must provide an url")
		return
	}

	url, err := url.Parse(u)
	if err != nil {
		return
	}

	if url.Scheme == "" {
		url.Scheme = "http"
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Write the page header
	// This HTML cannot go into a template since the data is never loaded into memory
	w.Write([]byte(fmt.Sprintf(`
<html>
	<head>
		<link rel="stylesheet" type="text/css" href="/static/source.css" media="screen" />
		<link rel="stylesheet" type="text/css" href="/static/normalize.css" media="screen" />
	</head>
	<body><pre>`)))

	// Parse the source HTML, output the wrapped HTML
	z := html.NewTokenizer(resp.Body)

SCAN:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			err = z.Err()
			if err == io.EOF {
				break SCAN
			}
			return

		case html.StartTagToken, html.EndTagToken:
			nbytes, _ := z.TagName()
			name := string(nbytes)

			_, err = w.Write([]byte(fmt.Sprintf("<span class=%q>", name)))
			if err != nil {
				return
			}

			raw := string(z.Raw())
			_, err = w.Write([]byte(html.EscapeString(raw)))
			if err != nil {
				return
			}

			_, err = w.Write([]byte(fmt.Sprintf("</span>")))
			if err != nil {
				return
			}

			if tt == html.StartTagToken {
				counts[name] = counts[name] + 1
			}

		default:
			raw := string(z.Raw())
			_, err = w.Write([]byte(html.EscapeString(raw)))
			if err != nil {
				return
			}
		}
	}

	// Write the page footer
	w.Write([]byte(fmt.Sprintf("</pre><script>window.counts=")))
	err = json.NewEncoder(w).Encode(counts)
	if err != nil {
		return
	}
	w.Write([]byte(fmt.Sprintf(";</script>")))
	w.Write([]byte(fmt.Sprintf("</body></html>")))
}

func (s *Server) healthcheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Healthcheck: server healthy")
	w.Write([]byte(fmt.Sprintf("Alive and healthy")))
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../static/tagger.html")
}

func (s *Server) Run() {
	http.HandleFunc("/doit", s.doit)
	http.HandleFunc("/healthcheck", s.healthcheck)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	//	http.Handle("/static/", http.FileServer(http.Dir("../static")))
	http.HandleFunc("/", s.handleRoot)

	log.Printf("Starting server on %q", s.URL)
	err := http.ListenAndServe(s.URL, nil)
	if err != nil {
		log.Fatal(err)
	}
}
