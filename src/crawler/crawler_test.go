package crawler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			// Home page with links to other pages
			fmt.Fprint(w, "<html><body>")
			fmt.Fprint(w, `<a href="/1">Link 1</a>`)
			fmt.Fprint(w, `<a href="/2">Link 2</a>`)
			fmt.Fprint(w, `<a href="/3">Link 3</a>`)
			fmt.Fprint(w, "</body></html>")
		case "/1":
			// Page 1 with a link to the home page
			fmt.Fprint(w, "<html><body>")
			fmt.Fprint(w, `<a href="/">Home</a>`)
			fmt.Fprint(w, "</body></html>")
		case "/2":
			// Page 2 with a link to an external domain
			fmt.Fprint(w, "<html><body>")
			fmt.Fprint(w, `<a href="http://www.external.com">External</a>`)
			fmt.Fprint(w, "</body></html>")
		case "/3":
			// Page 3 with a link to itself
			fmt.Fprint(w, "<html><body>")
			fmt.Fprint(w, `<a href="/3">Self</a>`)
			fmt.Fprint(w, "</body></html>")
		}
	}))
}

// idea stolen from https://stackoverflow.com/a/10476304/4916416
func captureOutput(fun func()) string {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fun()

	outChannel := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outChannel <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outChannel

	return string(out)
}

func TestCrawler(t *testing.T) {
	// Arrange
	server := createTestServer()
	defer server.Close()

	// Act
	capturedOutput := captureOutput(func() {
		c := NewCrawler()
		serverUrlParsed, _ := url.Parse(server.URL)
		c.Crawl(server.URL, serverUrlParsed.Host)
	})

	// Assert
	// Visited: http://127.0.0.1:49368
	// Links found:
	// - http://127.0.0.1:49368/1
	// - http://127.0.0.1:49368/2
	// - http://127.0.0.1:49368/3
	expectedOutput := fmt.Sprintf("Visited: %s\n", server.URL)
	expectedOutput += "Links found:\n"
	expectedOutput += fmt.Sprintf("- %s/1\n", server.URL)
	expectedOutput += fmt.Sprintf("- %s/2\n", server.URL)
	expectedOutput += fmt.Sprintf("- %s/3\n", server.URL)
	// Visited: http://127.0.0.1:49368/1
	// Links found:
	// - http://127.0.0.1:49368/
	expectedOutput += fmt.Sprintf("Visited: %s/1\n", server.URL)
	expectedOutput += "Links found:\n"
	expectedOutput += fmt.Sprintf("- %s/\n", server.URL)
	// Visited: http://127.0.0.1:49368/
// Links found:
// - http://127.0.0.1:49368/1
// - http://127.0.0.1:49368/2
// - http://127.0.0.1:49368/3
	expectedOutput += fmt.Sprintf("Visited: %s/\n", server.URL)
	expectedOutput += "Links found:\n"
	expectedOutput += fmt.Sprintf("- %s/1\n", server.URL)
	expectedOutput += fmt.Sprintf("- %s/2\n", server.URL)
	expectedOutput += fmt.Sprintf("- %s/3\n", server.URL)
	// Visited: http://127.0.0.1:49368/2
	// No links found!
	expectedOutput += fmt.Sprintf("Visited: %s/2\n", server.URL)
	expectedOutput += "No links found!\n"
	// Visited: http://127.0.0.1:49368/3
	// Links found:
	// - http://127.0.0.1:49368/3
	expectedOutput += fmt.Sprintf("Visited: %s/3\n", server.URL)
	expectedOutput += "Links found:\n"
	expectedOutput += fmt.Sprintf("- %s/3\n", server.URL)
	assert.Equal(t, expectedOutput, capturedOutput, "Unexpected output!")
}
