package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/vutranhoang1411/SimpleBank/db/sqlc"
)
type createAccountReq struct{
	Owner string `json:"owner" form:"owner" binding:"required"`
	Currency string `json:"currency" form:"currency" binding:"required,currency"`
}
func (server *Server)createAccount(ctx *gin.Context){
	var req createAccountReq;
	if err:=ctx.ShouldBind(&req);err!=nil{
		ctx.JSON(http.StatusBadRequest,handleError(err));
		return
	}
	arg:=db.CreateAccountParams{
		Owner: req.Owner,
		Currency: req.Currency,
		Balance: 0,
	}

	result,err:=server.store.CreateAccount(ctx,arg)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,handleError(err))
		return
	}
	ctx.JSON(http.StatusOK,result)
}
type getAccountByReq struct{
	ID int64 `json:"id" form:"id" uri:"id" binding:"required,numeric,min=0"`
}
func (server *Server)getAccountByID(ctx *gin.Context){
	var arg getAccountByReq;
	if err:=ctx.ShouldBindUri(&arg);err!=nil{
		ctx.JSON(http.StatusBadRequest,handleError(err));
		return
	}
	result,err:=server.store.GetAccount(ctx,arg.ID);
	if err!=nil{
		if (err==sql.ErrNoRows){
			ctx.JSON(http.StatusNotFound,handleError(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError,handleError(err))
		return
	}
	ctx.JSON(http.StatusOK,result)
}

type listAccountsReq struct{
	Owner string `json:"owner" form:"owner" binding:"required"`
	Limit int32 `json:"limit" form:"limit" binding:"required,numeric,min=0"`
	Offset int32 `json:"offset" form:"offset" binding:"required,numeric,min=0"`
}
func (server *Server)listAccounts(ctx *gin.Context){
	var arg listAccountsReq
	if err:=ctx.ShouldBindQuery(&arg);err!=nil{
		ctx.JSON(http.StatusBadRequest,handleError(err));
		return
	}
	result,err:=server.store.ListAccounts(ctx,db.ListAccountsParams{
		Owner: arg.Owner,
		Limit: arg.Limit,
		Offset: arg.Offset,
	})
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,handleError(err))
		return
	}
	ctx.JSON(http.StatusOK,result)
}
