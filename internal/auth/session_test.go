package auth_test

import (
	"context"
	"testing"

	"github.com/defryheryanto/piggy-bank-backend/internal/auth"
	"github.com/stretchr/testify/assert"
)

func TestAuthFromContext(t *testing.T) {
	ctx := context.TODO()

	t.Run("should get current session", func(t *testing.T) {
		currentSession := &auth.AuthSession{
			UserID:   1,
			Username: "defryheryanto",
		}
		ctx = context.WithValue(ctx, auth.SessionKey, currentSession)

		session := auth.FromContext(ctx)
		assert.NotNil(t, session)
		assert.Equal(t, currentSession.UserID, session.UserID)
		assert.Equal(t, currentSession.Username, session.Username)
	})

	t.Run("return nil if current session invalid", func(t *testing.T) {
		currentSession := &auth.User{
			UserID:   1,
			Username: "defryheryanto",
		}
		ctx = context.WithValue(ctx, auth.SessionKey, currentSession)

		session := auth.FromContext(ctx)
		assert.Nil(t, session)
	})
}
