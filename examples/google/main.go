package main

import (
	"fmt"
	"net/http"

	passport "github.com/praveen001/go-passport"

	"github.com/go-chi/chi"
	"github.com/praveen001/go-passport/strategies/google"
)

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("login")
	w.Write([]byte("test"))
}

func success(w http.ResponseWriter, r *http.Request) {
	v := r.Context().Value(passport.CtxKey)
	w.Write([]byte(fmt.Sprintln(v)))
}

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("middleware")
		h.ServeHTTP(w, r)
	})
}

func googleStrategy() *google.Strategy {
	opt := &google.StrategyOptions{
		ClientID:     "813004830141-2leap7abkjf00bsofqcdn6a6po8l9348.apps.googleusercontent.com",
		ClientSecret: "A8ItEQFJpwu05wwsjoyytJSn",
		CallbackURL:  "http://localhost:5000/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Verify: func(accessToken, refreshToken string, profile *google.Profile) (bool, interface{}) {
			return true, profile
		},
	}

	return google.New(opt)
}

func main() {
	r := chi.NewRouter()

	p := passport.New(&passport.Options{})
	p.Use("google", googleStrategy())

	r.Group(func(r chi.Router) {
		r.Get("/auth/google", p.Authenticate("google", nil))
		r.Get("/auth/google/callback", p.Authenticate("google", nil))

		r.Group(func(r chi.Router) {
			r.Use(passport.AuthRequired)
			r.Get("/success", success)
		})
	})

	http.ListenAndServe(":5000", r)
}
