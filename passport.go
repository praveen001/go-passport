package passport

import (
	"context"
	"net/http"
)

// Options for passport
type Options struct {
	Session         bool
	SuccessRedirect string
	FailureRedirect string
}

// Authenticate ..
func Authenticate(s Strategy, opt *Options, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Authenticate(w, r, func(res *Result) {
			if h == nil {
				res.withDefaultHandler(w, r, opt)
				return
			}

			res.withCustomHandler(w, r, h)
		})
	}
}

func (res *Result) withCustomHandler(w http.ResponseWriter, r *http.Request, h http.HandlerFunc) {
	ctx := context.WithValue(r.Context(), CtxKey, res)

	h.ServeHTTP(w, r.WithContext(ctx))
}

func (res *Result) withDefaultHandler(w http.ResponseWriter, r *http.Request, opt *Options) {
	if res.Ok {
		res.success(w, r, opt)
	} else {
		res.failure(w, r, opt)
	}
}

func (res *Result) success(w http.ResponseWriter, r *http.Request, opt *Options) {
	http.Redirect(w, r, opt.SuccessRedirect, http.StatusTemporaryRedirect)
}

func (res *Result) failure(w http.ResponseWriter, r *http.Request, opt *Options) {
	http.Redirect(w, r, opt.FailureRedirect, http.StatusTemporaryRedirect)
}
