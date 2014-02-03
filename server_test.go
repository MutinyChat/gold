package gold

import (
	"github.com/drewolson/testflight"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

var (
	handler = Handler{}
)

func TestSkin(t *testing.T) {
	testflight.WithServer(handler, func(r *testflight.Requester) {
		request, _ := http.NewRequest("GET", "/", nil)
		request.Header.Add("Accept", "text/html")
		response := r.Do(request)
		assert.Contains(t, response.Body, "<html")
		assert.Equal(t, 200, response.StatusCode)
	})
}

func TestOPTIONS(t *testing.T) {
	testflight.WithServer(handler, func(r *testflight.Requester) {
		request, _ := http.NewRequest("OPTIONS", "/", nil)
		request.Header.Add("Accept", "text/turtle")
		request.Header.Add("Access-Control-Request-Headers", "Triples")
		request.Header.Add("Access-Control-Request-Method", "PATCH")
		response := r.Do(request)
		assert.Empty(t, response.Body)
		assert.Equal(t, 200, response.StatusCode)
	})
}

func TestOPTIONSOrigin(t *testing.T) {
	origin := "http://localhost:1234"
	testflight.WithServer(handler, func(r *testflight.Requester) {
		request, _ := http.NewRequest("OPTIONS", "/", nil)
		request.Header.Add("Accept", "text/turtle")
		request.Header.Add("Origin", origin)
		response := r.Do(request)
		assert.Empty(t, response.Body)
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, response.RawResponse.Header.Get("Access-Control-Allow-Origin"), origin)
	})
}

func TestPOSTSPARQL(t *testing.T) {
	testflight.WithServer(handler, func(r *testflight.Requester) {
		request, _ := http.NewRequest("POST", "/abc", strings.NewReader("INSERT DATA { <a> <b> <c>, <c0> . }"))
		request.Header.Add("Content-Type", "application/sparql-update")
		response := r.Do(request)
		assert.Empty(t, response.Body)
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, response.RawResponse.Header.Get("Triples"), "2")
		request, _ = http.NewRequest("POST", "/abc", strings.NewReader("DELETE DATA { <a> <b> <c> . }"))
		request.Header.Add("Content-Type", "application/sparql-update")
		response = r.Do(request)
		assert.Empty(t, response.Body)
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, response.RawResponse.Header.Get("Triples"), "1")
	})
}

func TestPOSTTurtle(t *testing.T) {
	testflight.WithServer(handler, func(r *testflight.Requester) {
		response := r.Post("/abc", "text/turtle", "<a> <b> <c1> .")
		assert.Empty(t, response.Body)
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, response.RawResponse.Header.Get("Triples"), "2")
	})
	testflight.WithServer(handler, func(r *testflight.Requester) {
		response := r.Post("/abc", "text/turtle", "<a> <b> <c2> .")
		assert.Empty(t, response.Body)
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, response.RawResponse.Header.Get("Triples"), "3")
	})
	testflight.WithServer(handler, func(r *testflight.Requester) {
		request, _ := http.NewRequest("GET", "/abc", nil)
		request.Header.Add("Accept", "text/turtle")
		response := r.Do(request)
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, response.Body, "@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .\n\n<a>\n    <b> <c0>, <c1>, <c2> .\n\n")
	})
}

func TestPATCHJson(t *testing.T) {
	testflight.WithServer(handler, func(r *testflight.Requester) {
		request, _ := http.NewRequest("PATCH", "/abc", strings.NewReader(`{"a":{"b":[{"type":"uri","value":"c"}]}}`))
		request.Header.Add("Content-Type", "application/json")
		response := r.Do(request)
		assert.Empty(t, response.Body)
		assert.Equal(t, 200, response.StatusCode)
	})
	testflight.WithServer(handler, func(r *testflight.Requester) {
		request, _ := http.NewRequest("GET", "/abc", nil)
		request.Header.Add("Accept", "text/turtle")
		response := r.Do(request)
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, response.Body, "@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .\n\n<a>\n    <b> <c> .\n\n")
		assert.Equal(t, response.RawResponse.Header.Get("Triples"), "1")
	})
}

func TestPUTTurtle(t *testing.T) {
	testflight.WithServer(handler, func(r *testflight.Requester) {
		response := r.Put("/abc", "text/turtle", "<d> <e> <f> .")
		assert.Empty(t, response.Body)
		assert.Equal(t, 201, response.StatusCode)
	})
	testflight.WithServer(handler, func(r *testflight.Requester) {
		request, _ := http.NewRequest("GET", "/abc", nil)
		request.Header.Add("Accept", "text/turtle")
		response := r.Do(request)
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, response.Body, "@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .\n\n<d>\n    <e> <f> .\n\n")
	})
}

func TestMKCOL(t *testing.T) {
	testflight.WithServer(handler, func(r *testflight.Requester) {
		request, _ := http.NewRequest("MKCOL", "/abc", nil)
		response := r.Do(request)
		assert.Equal(t, 500, response.StatusCode)
	})
	testflight.WithServer(handler, func(r *testflight.Requester) {
		request, _ := http.NewRequest("MKCOL", "/_folder", nil)
		response := r.Do(request)
		assert.Empty(t, response.Body)
		assert.Equal(t, 201, response.StatusCode)
	})
	testflight.WithServer(handler, func(r *testflight.Requester) {
		response := r.Post("/_folder", "text/turtle", "<a> <b> <c>.")
		assert.Equal(t, 500, response.StatusCode)
	})
	testflight.WithServer(handler, func(r *testflight.Requester) {
		request, _ := http.NewRequest("GET", "/_folder", nil)
		response := r.Do(request)
		assert.Equal(t, 501, response.StatusCode)
	})
}

func TestHEAD(t *testing.T) {
	testflight.WithServer(handler, func(r *testflight.Requester) {
		request, _ := http.NewRequest("HEAD", "/abc", nil)
		response := r.Do(request)
		assert.Empty(t, response.Body)
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, response.RawResponse.Header.Get("Triples"), "1")
	})
}

func TestStreaming(t *testing.T) {
	*stream = true
	testflight.WithServer(handler, func(r *testflight.Requester) {
		response := r.Get("/abc")
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, response.Body, "@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .\n\n<d>\n    <e> <f> .\n\n")
	})
	*stream = false
}

func TestDELETE(t *testing.T) {
	testflight.WithServer(handler, func(r *testflight.Requester) {
		response := r.Delete("/abc", "", "")
		assert.Empty(t, response.Body)
		assert.Equal(t, 200, response.StatusCode)
	})
	testflight.WithServer(handler, func(r *testflight.Requester) {
		response := r.Get("/abc")
		assert.Equal(t, 404, response.StatusCode)
	})
}

func TestDELETEFolder(t *testing.T) {
	testflight.WithServer(handler, func(r *testflight.Requester) {
		response := r.Delete("/_folder", "", "")
		assert.Empty(t, response.Body)
		assert.Equal(t, 200, response.StatusCode)
	})
	testflight.WithServer(handler, func(r *testflight.Requester) {
		response := r.Get("/_folder")
		assert.Equal(t, 404, response.StatusCode)
	})
	testflight.WithServer(handler, func(r *testflight.Requester) {
		response := r.Delete("/_folder", "", "")
		assert.Empty(t, response.Body)
		assert.Equal(t, 404, response.StatusCode)
	})
	testflight.WithServer(handler, func(r *testflight.Requester) {
		response := r.Delete("/", "", "")
		assert.Equal(t, 500, response.StatusCode)
	})
}

func TestInvalidMethod(t *testing.T) {
	testflight.WithServer(handler, func(r *testflight.Requester) {
		request, _ := http.NewRequest("TEST", "/test", nil)
		response := r.Do(request)
		assert.Equal(t, 405, response.StatusCode)
	})
}

func TestInvalidAccept(t *testing.T) {
	testflight.WithServer(handler, func(r *testflight.Requester) {
		request, _ := http.NewRequest("GET", "/test", nil)
		request.Header.Add("Accept", "text/csv")
		response := r.Do(request)
		assert.Equal(t, 406, response.StatusCode)
	})
}

func TestInvalidContent(t *testing.T) {
	testflight.WithServer(handler, func(r *testflight.Requester) {
		response := r.Post("/test", "text/csv", "a\tb\tc\n")
		assert.Equal(t, 415, response.StatusCode)
	})
}
