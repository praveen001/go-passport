package passport

import (
	"net/http"
)

// Result all strategies should return this.
type Result struct {
	StrategyName string
	Code         int
	Message      string
	Data         interface{}
}

// Strategy all strategies must implement this.
type Strategy interface {
	Authenticate(w http.ResponseWriter, r *http.Request, cb func(res *Result))
}
