package passport

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

// Passport .
type Passport struct {
	Options *Options
}

// OutputFormat ..
type OutputFormat int

// JSON .
const (
	JSON OutputFormat = iota
	XML
	String
)

// Options for passport
type Options struct {
	strategies   map[string]Strategy
	Serializer   func(info interface{}) string
	Deserializer func(s string) (info interface{})
	OutputFormat OutputFormat
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
			h = p.DefaultHandler
		}

		h.ServeHTTP(w, r.WithContext(ctx))

	}
}

// DefaultHandler returns the info object returned by the strategy after authentication.
//
// If authentication has failed, it returns a 403 status response.
func (p *Passport) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value(CtxKey).(*Result)

	if res.Ok {
		w.WriteHeader(200)
		p.output(w, res.Info)
		return
	}

	w.WriteHeader(403)
}

func (p *Passport) output(w http.ResponseWriter, o interface{}) {
	switch p.Options.OutputFormat {
	case JSON:
		if err := json.NewEncoder(w).Encode(o); err != nil {
			log.Println("Unable to format output to JSON")
		}

	case XML:
		if err := xml.NewEncoder(w).Encode(o); err != nil {
			log.Println("Unable to format output to XML")
		}

	case String:
		fmt.Fprintf(w, "%v", o)
	}
}
