package passport

import "net/http"

// AuthRequired ..
func AuthRequired(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok := r.Header.Get("authorization")
		if tok != "" {
			h.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(403)
	})
}
