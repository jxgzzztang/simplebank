package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jxgzzztang/simplebank/util"
)

const (
	authorizationHeader = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "payloadKey"
)

func authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader(authorizationHeader)

		if len(token) == 0 {
			err := errors.New("no authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(token)

		if len(fields) < 2 {
			err := errors.New("invalid authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		fmt.Println(authorizationType, authorizationTypeBearer)
		if authorizationType != authorizationTypeBearer {
			err := errors.New("invalid authorization header type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]

		payload, ok := util.ParseToken(accessToken);
		fmt.Println(payload)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("token is invalid")))
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}