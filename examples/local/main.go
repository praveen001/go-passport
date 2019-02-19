package main

import (
	"fmt"
	"net/http"

	passport "github.com/praveen001/go-passport"

	"github.com/go-chi/chi"
	"github.com/praveen001/go-passport/strategies/local"
)

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("login")
	w.Write([]byte("test"))
}

func success(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome"))
}

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("middleware")
		h.ServeHTTP(w, r)
	})
}

func localStrategy() *local.Strategy {
	opt := &local.StrategyOptions{
		UsernameField: "username",
		PasswordField: "password",
		Verify: func(username, password string) (bool, interface{}) {
			return true, nil
		},
	}

	return local.New(opt)
}

func main() {
	r := chi.NewRouter()

	opt := &passport.Options{
		Session:         false,
		SuccessRedirect: "/success",
		FailureRedirect: "/failure",
	}

	r.Group(func(r chi.Router) {
		r.Post("/login", passport.Authenticate(localStrategy(), opt, nil))

		r.Group(func(r chi.Router) {
			r.Use(middleware)
			r.Get("/success", success)
		})
	})

	http.ListenAndServe(":5000", r)
}
