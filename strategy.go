package passport

import "net/http"

// Strategy ..
type Strategy interface {
	Authenticate(w http.ResponseWriter, r *http.Request, c CallbackFunc)
}
