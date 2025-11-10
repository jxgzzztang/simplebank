package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/jxgzzztang/simplebank/db/sqlc"
	"github.com/jxgzzztang/simplebank/util"
)

type TransferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required"`
	ToAccountID   int64 `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}


// Transfer godoc
// @Summary      Transfer
// @Description  transfer
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param transferData body TransferRequest true "参数"
// @Success      200  {object} 	db.TransferTxResult
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /transfer [post]
func (server *Server) Transfer(ctx *gin.Context)  {
	var req TransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fromAccount, isValid := server.validateCurrency(ctx, req.FromAccountID, req.Currency)
	if !isValid {
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*util.TokenPayload)

	if payload.Username != fromAccount.Owner {
		err := errors.New("from account does not have permission to transfer")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	_, isValid = server.validateCurrency(ctx, req.ToAccountID, req.Currency)

	if !isValid {
		return
	}

	createTransfer := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID: req.ToAccountID,
		Amount: req.Amount,
	}

	err, result := server.store.TransferTx(ctx, createTransfer)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validateCurrency(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {

	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}
	if account.Currency != currency {
		customErr := fmt.Errorf("valid currency is [%s], transfer [%s]", currency, account.Currency)

		ctx.JSON(http.StatusBadRequest, errorResponse(customErr))
		return account, false
	}
	return account, true

}