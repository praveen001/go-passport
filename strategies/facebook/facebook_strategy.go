package facebook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2/facebook"

	"github.com/praveen001/go-passport"
	"golang.org/x/oauth2"
)

const (
	profileURL = "https://graph.facebook.com/me"
)

// Strategy for facebook
type Strategy struct {
	Options *StrategyOptions
}

// StrategyOptions .
type StrategyOptions struct {
	CallbackURL  string
	ClientID     string
	ClientSecret string
	Scopes       []string
	Fields       []string
	Verify       func(accessToken, refreshToken string, profile *Profile) (ok bool, info interface{})
}

// Profile ..
type Profile struct {
	ID        string  `json:"id,omitempty"`
	Email     string  `json:"email,omitempty"`
	Name      string  `json:"name,omitempty"`
	Picture   picture `json:"picture,omitempty"`
	FirstName string  `json:"first_name,omitempty"`
	LastName  string  `json:"last_name,omitempty"`

	Error graphAPIError `json:"error"`
}

type picture struct {
	Data pictureData `json:"data"`
}

type pictureData struct {
	Height       int    `json:"height"`
	IsSilhouette bool   `json:"is_silhouette"`
	URL          string `json:"url"`
	Width        int    `json:"width"`
}

type graphAPIError struct {
	Message   string `json:"message"`
	Type      string `json:"type"`
	Code      int    `json:"code"`
	FBTraceID string `json:"fbtrace_id"`
}

// New ..
func New(opt *StrategyOptions) *Strategy {
	return &Strategy{
		Options: opt,
	}
}

// Authenticate ..
func (f *Strategy) Authenticate(w http.ResponseWriter, r *http.Request) *passport.Result {
	config := oauth2.Config{
		ClientID:     f.Options.ClientID,
		ClientSecret: f.Options.ClientSecret,
		Endpoint:     facebook.Endpoint,
		RedirectURL:  f.Options.CallbackURL,
		Scopes:       f.Options.Scopes,
	}

	if err := r.FormValue("error"); err != "" {
		return &passport.Result{
			Info: err,
		}
	}

	if code := r.FormValue("code"); code != "" {
		token, _ := config.Exchange(r.Context(), r.FormValue("code"))

		url := fmt.Sprintf("%s?fields=%s&access_token=%s", profileURL, strings.Join(f.Options.Fields, ","), token.AccessToken)
		res, err := http.Get(url)
		if err != nil {
			return &passport.Result{
				Ok:   false,
				Info: err.Error(),
			}
		}

		profile := Profile{}
		json.NewDecoder(res.Body).Decode(&profile)

		if profile.Error.Type != "" {
			return &passport.Result{
				Ok:   false,
				Info: profile.Error,
			}
		}

		ok, info := f.Options.Verify(token.AccessToken, token.RefreshToken, &profile)
		return &passport.Result{
			Ok:   ok,
			Info: info,
		}
	}

	url := config.AuthCodeURL("State", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	return &passport.Result{}
}
