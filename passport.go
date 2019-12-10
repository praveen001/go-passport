package passport

import (
	"context"
	"encoding/json"
	"net/http"
)

// Passport .
type Passport struct {
	Options *Options
}

// Options for passport
type Options struct {
	strategies   map[string]Strategy
	Deserializer func(s string) (interface{}, error)
}

// New creates a new passport instance
func New(opt *Options) *Passport {
	opt.strategies = make(map[string]Strategy)
	return &Passport{opt}
}

// Use registers a strategy by name
//
// Later, the registered strategy can be used by calling Authenticate() method
func (p *Passport) Use(name string, s Strategy) {
	p.Options.strategies[name] = s
}

// Authenticate calls `Strategy.Authenticate` method of registered strategies, and checks the `passport.Result` returned by it.
//
// The result is stored in the request context with `passport.CtxKey` as key.
//
// If a handler is provided, it is invoked, otherwise a `DefaultHandler` will be used
func (p *Passport) Authenticate(name string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		s, ok := p.Options.strategies[name]
		if !ok {
			w.WriteHeader(404)
			return
		}

		s.Authenticate(w, r, func(res *Result) {
			res.StrategyName = name

			if h == nil {
				if err := json.NewEncoder(w).Encode(res.Data); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(res.Code)
				return
			}

			ctx := context.WithValue(r.Context(), ResultCtxKey, res)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
