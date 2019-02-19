package passport

type contextKey struct {
	name string
}

// PassportCtxKey ..
var (
	PassportCtxKey = &contextKey{"passport"}
)
