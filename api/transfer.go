package api

import (
	"fmt"
	"net/http"
	db "sampla_bank/db/sqlc"

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

	if !server.validateCurrencyAccount(ctx, req.From_Account, req.Currency) {
		return
	}

	if !server.validateCurrencyAccount(ctx, req.To_Account, req.Currency) {
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

func (server *Server) validateCurrencyAccount(ctx *gin.Context, accountId int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		if err.Error() == "no rows in result set" {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	if account.Currency != currency {
		err := fmt.Errorf(`Currency mismatch for account id ` + string(rune(accountId)) + ` and currency ` + currency + ` is not allowed`)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return false
	}

	return true

}
