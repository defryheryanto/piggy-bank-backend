package auth

import "context"

type key string

const SessionKey = key("sessionKey")

func FromContext(ctx context.Context) *AuthSession {
	session, ok := ctx.Value(SessionKey).(*AuthSession)
	if !ok {
		return nil
	}

	return session
}
