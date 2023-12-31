package api

import (
	"bank-service/token"
	"context"
	"database/sql"
	"errors"
	"net/http"

	"bank-service/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,oneof=USD EUR KZT RUB"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	authPayload, ok := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if !ok {
		err := errors.New("mismatch token payload")
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}

	args := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(context.Background(), args)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, ErrorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	account, err := server.store.GetAccount(context.Background(), req.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	authPayload, ok := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if !ok {
		err := errors.New("mismatch token payload")
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}

	if authPayload.Username != account.Owner {
		err := errors.New("accounts doesn't belong to the user")
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type ListAccountRequest struct {
	Page int32 `form:"page" binding:"required,min=1"`
	Size int32 `form:"size" binding:"required,min=1"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	authPayload, ok := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if !ok {
		err := errors.New("mismatch token payload")
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}

	args := db.GetListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.Size,
		Offset: (req.Page - 1) * req.Size,
	}
	accounts, err := server.store.GetListAccounts(context.Background(), args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	if len(accounts) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"accounts": accounts})
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

type DeleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req DeleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	err := server.store.DeleteAccount(context.Background(), req.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": req.ID})
}

type UpdateAccountRequest struct {
	ID      int64 `json:"id"`
	Balance int64 `json:"balance" binding:"required"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var req UpdateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	args := db.UpdateAccountParams{
		ID:      req.ID,
		Balance: req.Balance,
	}
	account, err := server.store.UpdateAccount(context.Background(), args)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
