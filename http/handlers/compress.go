package handlers

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"

	httput "github.com/ggicci/jungo/http"
)

type compressResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (crw *compressResponseWriter) Write(bs []byte) (int, error) {
	if h := crw.Header(); h.Get("Content-Type") == "" {
		h.Set("Content-Type", http.DetectContentType(bs))
	}
	return crw.Writer.Write(bs)
}

func Compresser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Decide encoding.
		rr := httput.NewRequest(r)
		var enc string
		for _, enc = range rr.AcceptEncodings() {
			if enc == "*" {
				enc = "identity"
			}
			// Now, only support gzip and deflate compression methods.
			if enc == "identity" || enc == "gzip" || enc == "deflate" {
				break
			}
		}

		// As per rfc-2616 section 14.3. If the server can't encode the content
		// according to the accept encodings list, it should send 406 error.
		if enc == "" {
			http.Error(rw, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
			return
		}

		rw.Header().Set("Content-Encoding", enc)

		var output io.Writer = rw
		// Switch to compress writer.
		switch enc {
		case "gzip":
			output, _ = gzip.NewWriterLevel(rw, gzip.BestSpeed)
		case "deflate":
			output, _ = flate.NewWriter(rw, flate.BestSpeed)
		}

		// To close the writer.
		defer func(toclose interface{}) {
			if closer, ok := output.(io.Closer); ok {
				closer.Close()
			}
		}(output)

		crw := &compressResponseWriter{
			Writer:         output,
			ResponseWriter: rw,
		}
		next.ServeHTTP(crw, r)
	})
}
