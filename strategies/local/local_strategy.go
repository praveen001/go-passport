package local

import (
	"encoding/json"
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
	Verify        func(username, password string) *passport.Result
}

// New ..
func New(opt *StrategyOptions) *Strategy {
	return &Strategy{
		Options: opt,
	}
}

// Authenticate ..
func (l *Strategy) Authenticate(w http.ResponseWriter, r *http.Request, cb func(*passport.Result)) {
	body := make(map[string]string)

	// Read username, password
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		cb(&passport.Result{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// If username/password is not present, return 400
	username, hasUsername := body[l.Options.UsernameField]
	password, hasPassword := body[l.Options.PasswordField]
	if !hasUsername || !hasPassword {
		cb(&passport.Result{
			Code:    http.StatusBadGateway,
			Message: "Missing credentials",
		})
		return
	}

	// Call verify
	cb(l.Options.Verify(username, password))
}
