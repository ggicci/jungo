package handlers

import (
	"compress/flate"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func foo(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("foo"))
}

func TestCompresser(t *testing.T) {
	rd := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:8090/foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept-Encoding", "gzip;q=0.3,deflate;q=0.9")
	Compresser(http.HandlerFunc(foo)).ServeHTTP(rd, req)

	t.Logf("status code: %d", rd.Code)
	if rd.Code != 200 {
		t.Fatalf("status code should be 200, but got %d", rd.Code)
	}

	t.Logf("headers: %v", rd.HeaderMap)
	if ce := rd.HeaderMap["Content-Encoding"]; len(ce) == 0 || ce[0] != "deflate" {
		t.Fatalf("Content-Encoding in header should be deflate, but got %s", ce)
	}

	content, err := ioutil.ReadAll(flate.NewReader(rd.Body))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("body: %s", content)
	if sContent := string(content); sContent != "foo" {
		t.Fatalf("body content after decompression should be \"foo\", but got %q", sContent)
	}
}
