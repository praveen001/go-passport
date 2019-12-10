package passport

import (
	"context"
	"net/http"
)

// AuthRequired ..
func (p *Passport) AuthRequired(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok := r.Header.Get("authorization")
		if tok == "" {
			tok = r.URL.Query().Get("token")
		}

		if tok == "" {
			w.WriteHeader(403)
			return
		}

		info, err := p.Options.Deserializer(tok)
		if err != nil {
			w.WriteHeader(403)
		}

		ctx := context.WithValue(r.Context(), AuthCtxKey, info)
		if tok != "" {
			h.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		w.WriteHeader(403)
	})
}
