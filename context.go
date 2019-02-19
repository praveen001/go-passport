package passport

type contextKey struct {
	name string
}

// CtxKey is the context key where passport sets the authentication result.
var (
	CtxKey = &contextKey{"passport"}
)
