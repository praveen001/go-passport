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
	AuthorizationURL string
	CallbackURL      string
	ClientID         string
	ClientSecret     string
	TokenURL         string
	Scopes           []string
	Verify           func(accessToken, refreshToken string, profile *Profile) (ok bool, info interface{})
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
func (g *Strategy) Authenticate(w http.ResponseWriter, r *http.Request) *passport.Result {
	config := oauth2.Config{
		ClientID:     g.Options.ClientID,
		ClientSecret: g.Options.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  g.Options.CallbackURL,
		Scopes:       g.Options.Scopes,
	}

	q := r.URL.Query()

	if err := q.Get("error"); err != "" {
		return &passport.Result{
			Info: err,
		}
	}

	if code := q.Get("code"); code != "" {
		token, _ := config.Exchange(r.Context(), q.Get("code"))

		res, _ := http.Get(profileURL + "?access_token=" + token.AccessToken)
		profile := Profile{}
		json.NewDecoder(res.Body).Decode(&profile)

		ok, info := g.Options.Verify(token.AccessToken, token.RefreshToken, &profile)
		return &passport.Result{
			Ok:   ok,
			Info: info,
		}
	}

	url := config.AuthCodeURL("State", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	return nil
}
