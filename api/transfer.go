package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/vutranhoang1411/SimpleBank/db/sqlc"
)	

type transferRequest struct{
	FromAccountID int64 `form:"from_account_id" binding:"required,min=1"`
	ToAccountID int64	`form:"to_account_id" binding:"required,min=1"`
	Amount int64	`form:"amount" binding:"required,gt=0"`
	Currency string `form:"currency" binding:"required,currency"`
}
func (server *Server) createTransfer(ctx *gin.Context){
	//get request body
	var reqBody transferRequest;
	if err:=ctx.ShouldBind(&reqBody);err!=nil{
		ctx.JSON(http.StatusBadRequest,handleError(err));
		return
	}
	if !server.validateAccount(ctx,reqBody.FromAccountID,reqBody.Currency){
		return
	}
	if !server.validateAccount(ctx,reqBody.ToAccountID,reqBody.Currency){
		return
	}
	//
	result,err:=server.store.TransferTx(ctx,db.TransferTxParams{
		FromAccountID: reqBody.FromAccountID,
		ToAccountID: reqBody.ToAccountID,
		Amount: reqBody.Amount,
	})
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,handleError(err))
		return
	}
	ctx.JSON(http.StatusOK,result);
}
func (server *Server)validateAccount(ctx *gin.Context,accountID int64, currency string)bool{
	account,err:=server.store.GetAccount(ctx,accountID);
	if (err!=nil){
		if err==sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound,handleError(err))
			return false;
		}
		ctx.JSON(http.StatusInternalServerError,handleError(err))
		return false;
	}
	if account.Currency!=currency{
		err:=fmt.Errorf("Currency type mismatch: %v %v",account.Currency,currency);
		ctx.JSON(http.StatusBadRequest,handleError(err))
		return false;
	}
	return true;
}