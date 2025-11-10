package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jxgzzztang/simplebank/util"
	"net/http"
	"time"
)

type RenewAccessTokenParams struct {
	RefreshAccessToken string `json:"access_token"`
}

type RenewAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) RenewAccessToken(ctx *gin.Context) {
	var params RenewAccessTokenParams
	if err := ctx.ShouldBindBodyWithJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshTokenPayload, isValid := util.ParseToken(params.RefreshAccessToken)
	if !isValid {
		validRefreshTokenError := errors.New("invalid refresh token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(validRefreshTokenError))
		return
	}

	session, err := server.store.GetSessions(ctx, refreshTokenPayload.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("session is locked")))
		return
	}

	if refreshTokenPayload.Username != session.Username {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("session username is no match")))
		return
	}

	if session.RefreshToken != params.RefreshAccessToken {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("session refresh token is no match")))
		return
	}

	if time.Now().After(session.ExpiresAt.Time) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("session is expired")))
		return
	}

	accessToken, accessTokenPayload,  err := util.CreateToken(session.Username, util.Config.Jwt.ExpiresDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, RenewAccessTokenResponse{
		AccessToken: accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpiresAt.Time,
	})

}
