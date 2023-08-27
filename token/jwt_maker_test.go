package token

import (
	"testing"
	"time"

	"bank-service/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

func TestJwtMaker(t *testing.T) {
	maker, err := NewJwtMaker(util.GenerateRandomString(32))

	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, paylaod, err := maker.CreateToken(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, paylaod)
	require.Equal(t, paylaod.Username, username)

	payload, err := maker.ValidateToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
}

func TestExpiredJwtMaker(t *testing.T) {
	maker, err := NewJwtMaker(util.GenerateRandomString(32))

	require.NoError(t, err)
	require.NotEmpty(t, maker)

	token, payload, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.ValidateToken(token)

	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidTokenJWT(t *testing.T) {
	payload, err := NewPayload(util.GenerateRandomString(10), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)

	require.NoError(t, err)

	maker, err := NewJwtMaker(util.GenerateRandomString(32))
	require.NoError(t, err)

	payload, err = maker.ValidateToken(token)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidToken)
	require.Nil(t, payload)
}
