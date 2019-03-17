package google

import (
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/praveen001/go-passport"
)

const (
	profileURL = "https://www.googleapis.com/oauth2/v2/userinfo"
)

// Strategy ..
type Strategy struct {
	Options *StrategyOptions
}

// StrategyOptions ..
type StrategyOptions struct {
	CallbackURL  string
	ClientID     string
	ClientSecret string
	Scopes       []string
	Verify       func(accessToken, refreshToken string, profile *Profile) *passport.Result
}

// Profile ..
type Profile struct {
	Email         string `json:"email"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

// New ..
func New(opt *StrategyOptions) *Strategy {
	return &Strategy{
		Options: opt,
	}
}

// Authenticate ..
func (g *Strategy) Authenticate(w http.ResponseWriter, r *http.Request, cb func(*passport.Result)) {
	config := oauth2.Config{
		ClientID:     g.Options.ClientID,
		ClientSecret: g.Options.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  g.Options.CallbackURL,
		Scopes:       g.Options.Scopes,
	}

	if err := r.FormValue("error"); err != "" {
		cb(&passport.Result{
			Info: err,
		})
		return
	}

	if code := r.FormValue("code"); code != "" {
		token, err := config.Exchange(r.Context(), r.FormValue("code"))
		if err != nil {
			cb(&passport.Result{
				Info: err.Error(),
			})
			return
		}

		res, err := http.Get(profileURL + "?access_token=" + token.AccessToken)
		if err != nil {
			cb(&passport.Result{
				Info: err.Error(),
			})
			return
		}

		profile := Profile{}
		json.NewDecoder(res.Body).Decode(&profile)

		cb(g.Options.Verify(token.AccessToken, token.RefreshToken, &profile))
		return
	}

	url := config.AuthCodeURL("State", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
