package api

import (
	"bank-service/db/sqlc"
	"bank-service/util"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(context.Background(), arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, ErrorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusCreated, rsp)
}

type loginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginUserResponse struct {
	User                  userResponse `json:"user"`
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiredAt  time.Time    `json:"access_token_expired_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiredAt time.Time    `json:"refresh_token_expired_at"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	user, err := server.store.GetUser(context.Background(), req.Username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	if err := util.CheckPassword(req.Password, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiredAt:    refreshTokenPayload.ExpiredAt,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	rsp := loginUserResponse{
		User:                  newUserResponse(user),
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiredAt:  accessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: refreshTokenPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}
