package local

import (
	"encoding/json"
	"log"
	"net/http"

	passport "github.com/praveen001/go-passport"
)

// Strategy ..
type Strategy struct {
	Options *StrategyOptions
}

// StrategyOptions ..
type StrategyOptions struct {
	UsernameField string
	PasswordField string
	Verify        func(username, password string) (ok bool, info interface{})
}

// New ..
func New(opt *StrategyOptions) *Strategy {
	return &Strategy{
		Options: opt,
	}
}

// Authenticate ..
func (l *Strategy) Authenticate(w http.ResponseWriter, r *http.Request) *passport.Result {
	body := make(map[string]string)

	// Read username, password
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("Unable to decode request body", err.Error())
		return &passport.Result{Ok: false}
	}

	// If username/password is not present, return 400
	username, hasUsername := body[l.Options.UsernameField]
	password, hasPassword := body[l.Options.PasswordField]
	if !hasUsername || !hasPassword {
		log.Println("Missing credentials")
		return &passport.Result{Ok: false}
	}

	// Call verify
	ok, info := l.Options.Verify(username, password)

	return &passport.Result{
		Ok:   ok,
		Info: info,
	}

}
