package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/vutranhoang1411/SimpleBank/db/sqlc"
	"github.com/vutranhoang1411/SimpleBank/token"
	"github.com/vutranhoang1411/SimpleBank/util"
)
var(
	WrongLoginInfo=errors.New("Wrong email or password");
)
type createUserRequest struct{
	Name     string `json:"name" form:"name" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=8"`
}
func (server *Server) createUser(ctx *gin.Context){
	var reqBody createUserRequest;
	if err:=ctx.ShouldBind(&reqBody);err!=nil{
		ctx.JSON(http.StatusBadRequest,handleError(err))
		return;
	}
	password,err:=util.HashPassword(reqBody.Password);
	if (err!=nil){
		ctx.JSON(http.StatusInternalServerError,handleError(err));
		return
	}
	
	user,err:=server.store.CreateUser(ctx,db.CreateUserParams{
		ID:uuid.NewString()[2:18],
		Name:reqBody.Name,
		Email:reqBody.Email,
		Password: password,
	})
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,handleError(err));
		return
	}
	ctx.JSON(http.StatusOK,user)
}
type loginRequest struct{
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=8"`

}
type loginResponse struct{
	AccessToken string `json:"access_token"`
	AccessTokenExpiredAt time.Time `json:"access_token_expired_at"`
	RefreshToken string `json:"refresh_token"`
	RefreshTokenExpiredAt time.Time `json:"refresh_token_expired_at"`
}
func (server *Server)userLogin(ctx *gin.Context){
	//get request body
	var reqBody loginRequest;
	if err:=ctx.ShouldBind(&reqBody);err!=nil{
		ctx.JSON(http.StatusBadRequest,handleError(err));
		return;
	}

	user,err:=server.store.GetUserByEmail(ctx,reqBody.Email);
	if (err!=nil){
		if err==sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound,handleError(WrongLoginInfo));
			return
		}
		ctx.JSON(http.StatusInternalServerError,handleError(err))
		return
	}
	if !util.CheckPasswordHash(user.Password,reqBody.Password){
		ctx.JSON(http.StatusBadRequest,handleError(WrongLoginInfo));
		return;
	}
	accessToken,payload,err:=server.maker.CreateToken(user.Email,server.config.AccessTokenDuration);
	if (err!=nil){
		ctx.JSON(http.StatusInternalServerError,handleError(err))
		return
	}
	refreshToken,rfPayload,err:=server.maker.CreateToken(user.Email,server.config.RefreshTokenDuration);
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,handleError(err))
		return
	}
	server.store.CreateSession(ctx,db.CreateSessionParams{
		ID:rfPayload.ID,  //use token's payload id
		UserEmail: user.Email,
		RefreshToken: refreshToken,
		IsBlocked: false,
		ExpiresAt: rfPayload.ExpiredAt,
	})

	rsp:=loginResponse{
		AccessToken: accessToken,
		AccessTokenExpiredAt: payload.ExpiredAt,
		RefreshToken: refreshToken,
		RefreshTokenExpiredAt: rfPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK,rsp);
}
type newAccessTokenRes struct{
	AccessToken string `json:"access_token"`
	AccessTokenExpiredAt time.Time `json:"access_token_expired_at"`
}
func (server *Server)newAcessToken(ctx *gin.Context){
	//get payload data
	payloadValue,exist:=ctx.Get(authorizationPayloadKey)
	
	if !exist{
		ctx.JSON(http.StatusInternalServerError,handleError(errors.New("Something wrong with the server, please try again later")))
		return
	}

	rfPayload,ok:=payloadValue.(*token.Payload)
	if !ok{
		ctx.JSON(http.StatusInternalServerError,handleError(errors.New("Something wrong with the server, please try again later")))
		return
	}

	//get sesion on db
	session,err:=server.store.GetSession(ctx,rfPayload.ID)
	if err!=nil{
		ctx.JSON(http.StatusUnauthorized,handleError(errors.New("Can't find the session info, unauthorized!")))
		return
	}
	if session.IsBlocked{
		ctx.JSON(http.StatusUnauthorized,handleError(errors.New("The session has been blocked, please sign in again")))
		return
	}

	//return new access token
	accessToken,payload,err:=server.maker.CreateToken(session.UserEmail,server.config.AccessTokenDuration)

	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,handleError(errors.New("Something wrong with the server, please try again later")))
		return
	}
	ctx.JSON(http.StatusOK,newAccessTokenRes{
		AccessToken: accessToken,
		AccessTokenExpiredAt: payload.ExpiredAt,
	})

}