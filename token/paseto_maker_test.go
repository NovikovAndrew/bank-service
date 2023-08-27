package token

import (
	"testing"
	"time"

	"bank-service/util"

	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.GenerateRandomString(32))

	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, payload, err := maker.CreateToken(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.ValidateToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
}

func TestExpiredPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.GenerateRandomString(32))

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

func TestInvalidTokenPaseto(t *testing.T) {
	key := []byte(util.GenerateRandomString(32))
	payload, err := NewPayload(util.GenerateRandomString(10), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	token, err := paseto.NewV2().Encrypt(key, payload, nil)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	maker, err := NewPasetoMaker(util.GenerateRandomString(32))
	require.NoError(t, err)

	payload, err = maker.ValidateToken(token)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidToken)
	require.Nil(t, payload)
}
