package passport

import (
	"context"
	"net/http"
)

// Options ..
type Options struct {
	Session         bool
	SuccessRedirect string
	FailureRedirect string
}

// Result ..
type Result struct {
	Ok   bool
	Info interface{}
}

// CallbackFunc ..
type CallbackFunc func(r *Result)

// Authenticate ..
func Authenticate(s Strategy, opt *Options, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Authenticate(w, r, func(res *Result) {
			if h == nil {
				res.withDefaultHandler(w, r)
				return
			}

			res.withCustomHandler(w, r, h)
		})
	}
}

func (res *Result) withCustomHandler(w http.ResponseWriter, r *http.Request, h http.HandlerFunc) {
	ctx := context.WithValue(r.Context(), PassportCtxKey, res)

	h.ServeHTTP(w, r.WithContext(ctx))
}

func (res *Result) withDefaultHandler(w http.ResponseWriter, r *http.Request) {
	if res.Ok {
		res.success(w, r)
	} else {
		res.failure(w, r)
	}
}

func (res *Result) success(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from default handler"))
}

func (res *Result) failure(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Go to hell"))
}
