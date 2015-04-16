package tagger

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
)

type Server struct {
	URL    string
	counts map[string]int
}

func (s *Server) doit(w http.ResponseWriter, r *http.Request) {
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

	resp, err := http.Get(u)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Write the page header
	w.Write([]byte(fmt.Sprintf(`<html><head><link rel="stylesheet" type="text/css" href="/static/source.css" media="screen" /></head><body>`)))

	// Parse the source HTML, output the decorated
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
				s.counts[name] = s.counts[name] + 1
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
	w.Write([]byte(fmt.Sprintf("<script>window.counts=")))
	err = json.NewEncoder(w).Encode(s.counts)
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
	s.counts = make(map[string]int)

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
