package passport

import (
	"context"
	"fmt"
	"net/http"
)

// Passport .
type Passport struct {
	Options *Options
}

// Options for passport
type Options struct {
	strategies   map[string]Strategy
	Serializer   func(info interface{}) string
	Deserializer func(s string) (info interface{})
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

		res := s.Authenticate(w, r)
		res.StrategyName = name
		ctx := context.WithValue(r.Context(), CtxKey, res)

		if h == nil {
			h = DefaultHandler
		}

		h.ServeHTTP(w, r.WithContext(ctx))

	}
}

// DefaultHandler returns the info object returned by the strategy after authentication.
//
// If authentication has failed, it returns a 403 status response.
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value(CtxKey).(*Result)

	if res.Ok {
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintln(res)))
		return
	}

	w.WriteHeader(403)
}
