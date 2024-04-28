package api

import (
	"fmt"
	"net/http"

	db "github.com/rakeshdr543/go-bank-app/db/sqlc"
	"github.com/rakeshdr543/go-bank-app/token"

	"github.com/gin-gonic/gin"
)

type createMoneyTransferRequest struct {
	From_Account int64  `json:"from_account" binding:"required"`
	To_Account   int64  `json:"to_account" binding:"required"`
	Amount       int64  `json:"amount" binding:"required,gt=1"`
	Currency     string `json:"currency" binding:"required,currency"`
}

func (server *Server) createMoneyTransfer(ctx *gin.Context) {
	var req createMoneyTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := server.validateAccount(ctx, req.From_Account, req.Currency)

	if !valid {
		return
	}

	owner := ctx.MustGet(authorizationPayloadKey).(*token.Payload).Username

	if fromAccount.Owner != owner {
		err := fmt.Errorf("account %d does not belong to the authenticated user", req.From_Account)
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validateAccount(ctx, req.To_Account, req.Currency)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.From_Account,
		ToAccountID:   req.To_Account,
		Amount:        req.Amount,
	}

	account, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) validateAccount(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		if err.Error() == "no rows in result set" {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}
	if account.Currency != currency {
		err := fmt.Errorf(`Currency mismatch for account id ` + string(rune(accountId)) + ` and currency ` + currency + ` is not allowed`)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return account, false
	}

	return account, true

}
