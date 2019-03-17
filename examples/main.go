package main

import (
	"encoding/json"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/praveen001/go-passport"
	"github.com/praveen001/go-passport/strategies/facebook"
	"github.com/praveen001/go-passport/strategies/google"
	"github.com/praveen001/go-passport/strategies/local"
)

//
const (
	Secret = "secret"
)

func getJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})

	return token.SignedString([]byte(Secret))
}

// HomeHandler .
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	info := r.Context().Value(passport.CtxKey)
	auth := info.(AuthResponse)
	json.NewEncoder(w).Encode(auth)
}

// AuthResponse .
type AuthResponse struct {
	Email string `json:"email"`
	Token string `json:"token,omitempty"`
}

func deserializer(tokstr string) (interface{}, error) {
	token, err := jwt.Parse(tokstr, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		return nil, err
	}

	m := token.Claims.(jwt.MapClaims)
	return AuthResponse{
		Email: m["email"].(string),
	}, nil
}

func localStrategy() *local.Strategy {
	opt := &local.StrategyOptions{
		UsernameField: "username",
		PasswordField: "password",
		Verify: func(username, password string) *passport.Result {
			token, err := getJWT(username)
			if err != nil {
				return &passport.Result{
					Info: err.Error(),
				}
			}

			return &passport.Result{
				Ok: true,
				Info: AuthResponse{
					Token: token,
					Email: username,
				},
			}
		},
	}

	return local.New(opt)
}

func googleStrategy() *google.Strategy {
	opt := &google.StrategyOptions{
		ClientID:     "813004830141-2leap7abkjf00bsofqcdn6a6po8l9348.apps.googleusercontent.com",
		ClientSecret: "A8ItEQFJpwu05wwsjoyytJSn",
		CallbackURL:  "http://localhost:5000/auth/google",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Verify: func(accessToken, refreshToken string, profile *google.Profile) *passport.Result {
			token, err := getJWT(profile.Email)
			if err != nil {
				return &passport.Result{
					Info: err.Error(),
				}
			}

			return &passport.Result{
				Ok: true,
				Info: AuthResponse{
					Token: token,
					Email: profile.Email,
				},
			}
		},
	}

	return google.New(opt)
}

func facebookStrategy() *facebook.Strategy {
	opt := &facebook.StrategyOptions{
		ClientID:     "2377736872283552",
		ClientSecret: "dfdcbc53fc36d8735a28bca4d04b4d76",
		CallbackURL:  "http://localhost:5000/auth/facebook",
		Scopes:       []string{"public_profile", "email"},
		Fields:       []string{"id", "name", "email", "picture.type(large)"},
		Verify: func(accessToken, refreshToken string, profile *facebook.Profile) *passport.Result {
			token, err := getJWT(profile.Email)
			if err != nil {
				return &passport.Result{
					Info: err.Error(),
				}
			}

			return &passport.Result{
				Ok: true,
				Info: AuthResponse{
					Token: token,
					Email: profile.Email,
				},
			}
		},
	}

	return facebook.New(opt)
}

func main() {
	r := chi.NewRouter()
	p := passport.New(&passport.Options{
		Deserializer: deserializer,
	})

	p.Use("local", localStrategy())
	p.Use("google", googleStrategy())
	p.Use("facebook", facebookStrategy())

	r.Route("/auth", func(r chi.Router) {
		r.Get("/login", p.Authenticate("local", nil))
		r.Get("/google", p.Authenticate("google", nil))
		r.Get("/facebook", p.Authenticate("facebook", nil))
	})

	r.Group(func(r chi.Router) {
		r.Use(p.AuthRequired)
		r.Get("/home", HomeHandler)
	})

	http.ListenAndServe(":5000", r)
}
