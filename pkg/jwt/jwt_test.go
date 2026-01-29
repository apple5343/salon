package jwt

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func TestJwt(t *testing.T) {
	t.Parallel()
	const ttl = time.Hour
	const secret = "secret"
	userID := gofakeit.UUID()
	token, err := GenerateToken(UserInfo{ID: userID}, []byte(secret), ttl)
	require.NoError(t, err)

	claims, err := VerifyToken(token, []byte(secret))
	require.NoError(t, err)

	require.Equal(t, claims.ID, userID)

	_, err = VerifyToken(token, []byte("wrong"))
	require.Error(t, err)

	_, err = VerifyToken("wrong", []byte(secret))
	require.Error(t, err)
}
