package passport

import (
	"encoding/json"
	"log"
	"net/http"
)

// LocalStrategy ..
type LocalStrategy struct {
	Options *LocalStrategyOptions
}

// LocalStrategyOptions ..
type LocalStrategyOptions struct {
	UsernameField string
	PasswordField string
	Verify        func(username, password string) (ok bool, info interface{})
}

// NewLocalStrategy ..
func NewLocalStrategy(opt *LocalStrategyOptions) *LocalStrategy {
	return &LocalStrategy{
		Options: opt,
	}
}

// Authenticate ..
func (l *LocalStrategy) Authenticate(w http.ResponseWriter, r *http.Request, next CallbackFunc) {
	body := make(map[string]string)

	// Read username, password
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("Unable to decode request body", err.Error())
		return
	}

	// If username/password is not present, return 400
	username, hasUsername := body[l.Options.UsernameField]
	password, hasPassword := body[l.Options.PasswordField]
	if !hasUsername || !hasPassword {
		log.Println("Missing credentials")
		return
	}

	// Call verify
	ok, info := l.Options.Verify(username, password)

	res := &Result{
		Ok:   ok,
		Info: info,
	}
	next(res)
}
