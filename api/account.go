package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/jxgzzztang/simplebank/db/sqlc"
	"github.com/jxgzzztang/simplebank/util"
	"net/http"
)

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,oneof=CNY USD"`
}

// CreateAccount godoc
// @Summary      CreateAccount
// @Description  create a account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Success      200  {object} 	db.Account
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /createAccount [post]
func (server *Server) CreateAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	parsePayload := ctx.MustGet(authorizationPayloadKey).(*util.TokenPayload)

	accountArg := db.CreateAccountParams{
		Owner: parsePayload.Username,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(ctx, accountArg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "owner_currency_key", "accounts_owner_fkey":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}
	ctx.JSON(http.StatusOK, account)
}


type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// GetAccount godoc
// @Summary      GetAccount
// @Description  get a account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param 		id  path  int true "Account ID"
// @Success      200  {object} 	db.Account
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /createAccount/{id} [get]
func (server *Server) GetAccount(ctx *gin.Context) {
	var req GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*util.TokenPayload)

	if payload.Username != account.Owner {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid owner")))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type ListAccountsRequest struct {
	PageSize int32 `form:"pageSize" binding:"required,min=1,max=10"`
	PageNumber int32 `form:"pageNumber" binding:"required,min=1"`
}

// ListAccount godoc
// @Summary      ListAccounts
// @Description  get a accounts list
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param 		pageSize  query  int true "page size"
// @Param 		pageNumber query int true "page number"
// @Success      200  {array} 	db.Account
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /listAccounts [get]
func (server *Server) ListAccounts(ctx *gin.Context) {
	var req ListAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*util.TokenPayload)

	listArg := db.ListAccountParams{
		Owner: payload.Username,
		Limit: req.PageSize,
		Offset: (req.PageNumber - 1) * req.PageSize,
	}
	accounts, err := server.store.ListAccount(ctx, listArg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}