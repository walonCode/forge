package middleware

import "net/http"

// maxRequestBodyBytes caps the size of a request body the API will read (1 MB).
const maxRequestBodyBytes = 1 << 20

// MaxBodyBytes limits every request body to maxRequestBodyBytes. Reads past the
// limit fail, so handlers decoding oversized bodies return an error instead of
// buffering unbounded input.
func MaxBodyBytes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
		}
		next.ServeHTTP(w, r)
	})
}
