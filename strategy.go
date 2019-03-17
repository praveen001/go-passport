package passport

import (
	"net/http"
)

// Result all strategies should return this.
type Result struct {
	StrategyName string

	Ok    bool
	Error bool

	Info interface{}
}

// Strategy all strategies must implement this.
type Strategy interface {
	Authenticate(w http.ResponseWriter, r *http.Request, cb func(res *Result))
}
