package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/jxgzzztang/simplebank/db/sqlc"
	"github.com/jxgzzztang/simplebank/util"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserInfoResponse struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	CreateAt time.Time `json:"create_at"`
	PasswordChangedAt time.Time `json:"password_change_at"`
}

func CreateUserInfoResponse(user db.User) UserInfoResponse {
	return UserInfoResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
		CreateAt: user.CreatedAt.Time,
		PasswordChangedAt: user.PasswordChangedAt.Time,
	}
}

// CreateUser godoc
// @Summary      CreateUser
// @Description  create a user
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Success      200  {object} 	db.User
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /createUser [post]
func (server *Server) CreateUser(ctx *gin.Context)  {
	var reqParams CreateUserRequest

	if err := ctx.ShouldBindJSON(&reqParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashPassword, err := util.HashPassword(reqParams.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userArg := db.CreateUserParams{
		Username: reqParams.Username,
		HashedPassword: hashPassword,
		FullName: reqParams.FullName,
		Email: reqParams.Email,
	}

	user, err := server.store.CreateUser(ctx, userArg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := CreateUserInfoResponse(user)
	ctx.JSON(http.StatusOK, resp)
}