package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vutranhoang1411/SimpleBank/token"
)
const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)
func token_validator(maker token.Maker) gin.HandlerFunc{
	return func(ctx *gin.Context){
		authorization:=ctx.GetHeader(authorizationHeaderKey);
		if len(authorization)==0{
			ctx.JSON(http.StatusUnauthorized,handleError(errors.New("Authorization token is not provided")))
			ctx.Abort()
			return;
		}
		fields:=strings.Fields(authorization);
		if len(fields)<2{
			ctx.JSON(http.StatusUnauthorized,handleError(errors.New("Invalid access token")));
			ctx.Abort()
			return
		}

		if strings.ToLower(fields[0])!=authorizationTypeBearer{
			ctx.JSON(http.StatusUnauthorized,handleError(errors.New("Invalid access token")));
			ctx.Abort()
			return
		}

		payload,err:=maker.VerifyToken(fields[1]);
		if (err!=nil){
			ctx.JSON(http.StatusUnauthorized,handleError(errors.New("Invalid access token")));
			ctx.Abort()
			return
		}
		ctx.Set(authorizationPayloadKey,payload);
		ctx.Next();
	}
}