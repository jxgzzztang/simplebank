package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/jxgzzztang/simplebank/db/sqlc"
	"github.com/jxgzzztang/simplebank/util"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	SessionID             pgtype.UUID      `json:"session_id"`
	UserInfo              UserInfoResponse `json:"user"`
	AccessToken           string           `json:"access_token"`
	AccessTokenExpiredAt  time.Time        `json:"access_token_expired_at"`
	RefreshToken          string           `json:"refresh_token"`
	RefreshTokenExpiredAt time.Time        `json:"refresh_token_expired_at"`
}

func (server *Server) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if err := util.CheckPassword(req.Password, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	token, accessTokenPayload, err := util.CreateToken(req.Username, util.Config.Jwt.ExpiresDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := util.CreateToken(user.Username, util.Config.Jwt.RefreshDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session := db.CreateSessionsParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    pgtype.Timestamptz{Time: refreshTokenPayload.ExpiresAt.Time, Valid: true},
	}

	_, err = server.store.CreateSessions(ctx, session)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := UserResponse{
		SessionID:             refreshTokenPayload.ID,
		UserInfo:              CreateUserInfoResponse(user),
		AccessToken:           token,
		AccessTokenExpiredAt:  accessTokenPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: refreshTokenPayload.ExpiresAt.Time,
	}

	ctx.JSON(http.StatusOK, resp)

}
