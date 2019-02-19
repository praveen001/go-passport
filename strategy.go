package passport

import "net/http"

// Result ..
type Result struct {
	Ok   bool
	Info interface{}
}

// CallbackFunc ..
type CallbackFunc func(r *Result)

// Strategy ..
type Strategy interface {
	Authenticate(w http.ResponseWriter, r *http.Request, c CallbackFunc)
}
