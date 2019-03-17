package main

import (
	"fmt"
	"net/http"

	passport "github.com/praveen001/go-passport"

	"github.com/go-chi/chi"
	"github.com/praveen001/go-passport/strategies/facebook"
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

func facebookStrategy() *facebook.Strategy {
	opt := &facebook.StrategyOptions{
		ClientID:     "2377736872283552",
		ClientSecret: "dfdcbc53fc36d8735a28bca4d04b4d76",
		CallbackURL:  "http://localhost:5000/auth/facebook",
		Scopes:       []string{"public_profile", "email"},
		Fields:       []string{"id", "name", "email", "picture.type(large)"},
		Verify: func(accessToken, refreshToken string, profile *facebook.Profile) (bool, interface{}) {
			return true, profile
		},
	}

	return facebook.New(opt)
}

func main() {
	r := chi.NewRouter()

	p := passport.New(&passport.Options{})
	p.Use("facebook", facebookStrategy())

	r.Group(func(r chi.Router) {
		r.Get("/auth/facebook", p.Authenticate("facebook", nil))

		r.Group(func(r chi.Router) {
			r.Use(passport.AuthRequired)
			r.Get("/success", success)
		})
	})

	http.ListenAndServe(":5000", r)
}
