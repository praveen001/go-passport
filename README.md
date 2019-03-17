## Example: local strategy

Each strategy has it's own options, but a common and required option is the `Verify` function. Which has to verify the incoming credentials, and return a response (including a JWT token).

```go
// Creating a local strategy
func localStrategy() *local.Strategy {
	opt := &local.StrategyOptions{
		UsernameField: "username",
		PasswordField: "password",
		Verify: func(username, password string) *passport.Result {
            // check if username/password is valid

			token, err := getJWT(username)
			if err != nil {
				return &passport.Result{
					Info: err.Error(),
				}
            }

			return &passport.Result{
				Ok: true,
				Info: AuthResponse{
					Token: token,
					Email: username,
				},
			}
		},
	}

	return local.New(opt)
}
```

Once you have created your strategy, you can attach them to your router. This example uses chi router.

```go
func main() {
    r := chi.NewRouter()

    // Create passport instance
    p := passport.New(&passport.Options{
		Deserializer: deserializer,
    })

    // Attach the strategy, via a name
    p.Use("local", localStrategy())

    // Assign a endpoint
    r.Route("/auth", func(r chi.Router) {
		r.Get("/login", p.Authenticate("local", nil))
    })

    // Authenticated routes
    r.Group(func(r chi.Router) {
		r.Use(p.AuthRequired)
		r.Get("/home", HomeHandler)
	})

	http.ListenAndServe(":5000", r)
}
```

## Creating a Deserializer

Deserializer function is required to create a instance of passport. Responsibility of deserializer function is to verify the JWT, and return some data about the user, which is then available to other http.Handler in context with key `passport.Ctxkey`

This example uses "github.com/dgrijalva/jwt-go" to parse the JWT

```go
func deserializer(tokstr string) (interface{}, error) {
	token, err := jwt.Parse(tokstr, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		return nil, err
	}

	m := token.Claims.(jwt.MapClaims)
	return AuthResponse{
		Email: m["email"].(string),
	}, nil
}
```
