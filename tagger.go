package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	// "time"
	// "net/http/httputil"
)

var counts = make(map[string]int)

func doit(w http.ResponseWriter, r *http.Request) {
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
	w.Write([]byte(fmt.Sprintf("<html><body>")))

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

			_, err = w.Write(z.Raw())
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
			_, err = w.Write(z.Raw())
			if err != nil {
				return
			}
		}
	}

	// Write the page footer
	w.Write([]byte(fmt.Sprintf("<script>window.counts =")))
	err = json.NewEncoder(w).Encode(counts)
	if err != nil {
		return
	}
	w.Write([]byte(fmt.Sprintf("; </script>")))
	w.Write([]byte(fmt.Sprintf("</body></html>")))
}

func test(w http.ResponseWriter, r *http.Request) {
	s := `
<html>
	<head>
		<title>Test page</title>
	</head>
	<body>
		<article id="content">
			<section class="section-1">
				<section class="nested-section">
				<section>
			<section>
			<hr />
			<div class="div-1">
				<div class="nested-div">
				<div>
			<div>
		</article>
	</body>
</html>
`
	w.Write([]byte(s))
}

func main() {
	http.HandleFunc("/doit", doit)
	http.HandleFunc("/test", test)
	//TODO a flag for port
	// go func() {
	err := http.ListenAndServe(":8580", nil)
	if err != nil {
		log.Fatal(err)
	}
	// }()

	/*
		// time.Sleep(1 * time.Second)
		cl := http.Client{}
		res, err := cl.Get("http://localhost:8580/test")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		bb, err := httputil.DumpResponse(res, true)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s", bb)
	*/
}