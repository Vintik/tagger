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
				in: `<html>Hello</html>`,
				out: `<html><body><span class="html"><html></span>Hello<span class="html">
					</html></span><script>window.counts={"html":1};</script></body></html>`,
			},
			{
				in: ` <html> <head> <title>Test page</title> </head> <body> <article id="content">
					<section class="section-1"> <section class="nested-section"> <section> <section>
					<hr /> <div class="div-1"> <div class="nested-div"> <div> <div> </article>
					</body> </html> `,
				out: `<html><body><spanclass="html"><html></span><spanclass="head"><head></span><spanclass="title"><title></span>Testpage<spanclass="title"></title></span><spanclass="head"></head></span><spanclass="body"><body></span><spanclass="article"><articleid="content"></span><spanclass="section"><sectionclass="section-1"></span><spanclass="section"><sectionclass="nested-section"></span><spanclass="section"><section></span><spanclass="section"><section></span><hr/><spanclass="div"><divclass="div-1"></span><spanclass="div"><divclass="nested-div"></span><spanclass="div"><div></span><spanclass="div"><div></span><spanclass="article"></article></span><spanclass="body"></body></span><spanclass="html"></html></span><script>window.counts={"article":1,"body":1,"div":4,"head":1,"html":2,"section":4,"title":1};</script></body></html>`,
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
