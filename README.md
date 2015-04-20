# tagger

A web application that allows for the fetching of another page's source code. Enter a valid url and get the source, a list of tags on the page, and the counts of each tag. Click on a count to highlight the tags in the source.

# Running the server

Haven't tested this on my local machine (all development done on the ec2 box).
- Install go: https://golang.org/doc/install
- `cd server`
- `go build server.go`
  - if build fails due to missing packages `go get PACKAGE_NAME` whatever packages need to be installed and rebuild
- `sudo ./server -p ":8080"`

## Code design

The solution that I chose to go with, was to have a go server that would have two endpoints:

1) GET / - load up the main page and DOM of the application

2) POST /doit?url={{url}} - a request that passes the url in the query string and fetches the DOM server side. The DOM gets parsed in a stream, and passed off to the tokenizer (golang.org/x/net/html package) which wraps each tag with a span. Each span has a className of the wrapped tagName. 

The frontend will POST the search requests to an iframe and load the response in there. This is done due to the /doit response being a stream, and letting the browser do the heavy work of parsing sizeable pages and displaying them. This keeps the memory overhead on the frontend to a minimum. 


## TODOs

- Found out how to stream data into a go template. Would be a much cleaner implementation to use a go template than to post to an iframe. This will resolve the bookmarking/routing of the frontend
- Rename /doit endpoint to something more descriptive
- Decide whether frontend should move to jquery to simplify the code (quite a lot of document.createElement)
- Clean up some styling and alignments
- Fix relative pathing issue to serve static assets. Currently the path is relative to server.go, and the app won't work if started in a different directory of the project.
