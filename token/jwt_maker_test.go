package token

import (
	"sampla_bank/util"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	role := util.DepositorRole
	duration := time.Minute

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, role, duration)

	require.NoError(t, err)
	require.Equal(t, payload.Username, username)

	parsedPayload, err := maker.VerifyToken(token)

	require.NoError(t, err)

	require.NotEmpty(t, token)

	require.NotEqual(t, token, "")

	require.Equal(t, username, parsedPayload.Username)

	require.Equal(t, role, parsedPayload.Role)

	require.WithinDuration(t, issuedAt, time.Unix(parsedPayload.IssuedAt, 0), time.Second)

	require.WithinDuration(t, expiresAt, time.Unix(parsedPayload.ExpiresAt, 0), time.Second)

	require.NotZero(t, parsedPayload.ID)

}

func TestJWTExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomOwner(), util.DepositorRole, -time.Minute)

	require.NoError(t, err)
	require.NotEmpty(t, payload.ID)

	payload, err = maker.VerifyToken(token)

	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestJWTInvalidToken(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), util.DepositorRole, time.Minute)

	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType) //
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32))

	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)

	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())

	require.Nil(t, payload)
}
