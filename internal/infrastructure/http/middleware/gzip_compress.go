package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func GZIPCompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isRequestCompressed(r) {
			decompressed, err := decompressRequest(r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			r.Body = decompressed
		}

		if !isAcceptsCompression(r) {
			next.ServeHTTP(w, r)
			return
		}

		/* if !canCompressContent(r) {
			next.ServeHTTP(w, r)
			return
		} */

		gz, _ := gzip.NewWriterLevel(w, gzip.BestSpeed)
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzResponseWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

type gzResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func isRequestCompressed(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
}

func decompressRequest(r *http.Request) (io.ReadCloser, error) {
	gzr, err := gzip.NewReader(r.Body)
	if err != nil {
		return nil, err
	}

	return gzr, nil
}

func isAcceptsCompression(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

func canCompressContent(r *http.Request) bool {
	contentType := r.Header.Get("Accept")

	return strings.Contains(contentType, "json") ||
		strings.Contains(contentType, "html") ||
		strings.Contains(contentType, "xml") ||
		strings.Contains(contentType, "text")
}
