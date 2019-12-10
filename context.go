package passport

type contextKey struct {
	name string
}

// ResultCtxKey is the context key where passport sets the authentication result.
var (
	ResultCtxKey = &contextKey{"passport_result"}
	AuthCtxKey   = &contextKey{"passport_auth"}
)
