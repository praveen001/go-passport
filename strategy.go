package passport

import (
	"net/http"
)

// Result all strategies should return this.
type Result struct {
	Ok           bool
	StrategyName string
	Info         interface{}
}

// Strategy all strategies must implement this.
type Strategy interface {
	Authenticate(w http.ResponseWriter, r *http.Request) *Result
}
