package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	refreshTokenPayload, err := server.tokenMaker.ValidateToken(req.RefreshToken)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshTokenPayload.ID)

	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			ctx.JSON(http.StatusNotFound, ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	if session.IsBlocked {
		err := errors.New("session is block")
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}

	if session.Username != refreshTokenPayload.Username {
		err := fmt.Errorf("mismatch username %s vs %s", session.Username, refreshTokenPayload.Username)
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatch refresh token %s vs %s", session.RefreshToken, req.RefreshToken)
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}

	if time.Now().After(session.ExpiredAt) {
		err := errors.New("expired session")
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(session.Username, server.config.AccessTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	})
}
