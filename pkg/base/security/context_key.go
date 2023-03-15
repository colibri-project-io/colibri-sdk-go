package security

type contextKey string

func (c contextKey) String() string {
	return "context key " + string(c)
}

const (
	contextKeyAuthenticationContext = contextKey("authentication-context-key")
)
