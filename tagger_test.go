package tagger

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

type testData struct {
	t    *testing.T
	data []struct {
		in, out string
	}
}

func (td testData) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		td.t.Fatal(err)
	}

	i, err := strconv.Atoi(r.Form.Get("i"))
	if err != nil {
		td.t.Fatal(err)
	}

	if i >= len(td.data) {
		td.t.Fatalf("i %v out of range %v", i, len(td.data))
	}

	w.Write([]byte(td.data[i].in))
}

func TestTaggerTest(t *testing.T) {
	td := testData{
		t: t,
		data: []struct{ in, out string }{
			{
				in:  `<html>Hello</html>`,
				out: `<html><head><linkrel="stylesheet"type="text/css"href="/static/source.css"media="screen"/><linkrel="stylesheet"type="text/css"href="/static/normalize.css"media="screen"/></head><body><pre><spanclass="html">&lt;html&gt;</span>Hello<spanclass="html">&lt;/html&gt;</span></pre><script>window.counts={"html":1};</script></body></html>`,
			},
			{
				in: ` <html> <head> <title>Test page</title> </head> <body> <article id="content">
					<section class="section-1"> <section class="nested-section"> <section> <section>
					<hr /> <div class="div-1"> <div class="nested-div"> <div> <div> </article>
					</body> </html> `,
				out: `<html><head><linkrel="stylesheet"type="text/css"href="/static/source.css"media="screen"/><linkrel="stylesheet"type="text/css"href="/static/normalize.css"media="screen"/></head><body><pre><spanclass="html">&lt;html&gt;</span><spanclass="head">&lt;head&gt;</span><spanclass="title">&lt;title&gt;</span>Testpage<spanclass="title">&lt;/title&gt;</span><spanclass="head">&lt;/head&gt;</span><spanclass="body">&lt;body&gt;</span><spanclass="article">&lt;articleid=&#34;content&#34;&gt;</span><spanclass="section">&lt;sectionclass=&#34;section-1&#34;&gt;</span><spanclass="section">&lt;sectionclass=&#34;nested-section&#34;&gt;</span><spanclass="section">&lt;section&gt;</span><spanclass="section">&lt;section&gt;</span>&lt;hr/&gt;<spanclass="div">&lt;divclass=&#34;div-1&#34;&gt;</span><spanclass="div">&lt;divclass=&#34;nested-div&#34;&gt;</span><spanclass="div">&lt;div&gt;</span><spanclass="div">&lt;div&gt;</span><spanclass="article">&lt;/article&gt;</span><spanclass="body">&lt;/body&gt;</span><spanclass="html">&lt;/html&gt;</span></pre><script>window.counts={"article":1,"body":1,"div":4,"head":1,"html":1,"section":4,"title":1};</script></body></html>`,
			},
		},
	}

	s := Server{
		URL: ":8586",
	}
	go s.Run()

	ts := httptest.NewServer(td)

	for i, tt := range td.data {
		url := fmt.Sprintf("http://%s/doit?url=%s?i=%d", s.URL, ts.URL, i)
		resp, err := http.Get(url)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		replace := func(s string) string {
			s = strings.Replace(s, " ", "", -1)
			s = strings.Replace(s, "\n", "", -1)
			s = strings.Replace(s, "\t", "", -1)
			return s
		}
		got := replace(string(bytes))
		expected := replace(tt.out)
		if got != expected {
			t.Fatalf("%d:\n\tgot:\t%s\nexpected:\t%s", i, got, expected)
		}
	}
}
