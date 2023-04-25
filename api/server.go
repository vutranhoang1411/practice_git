package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/vutranhoang1411/SimpleBank/db/sqlc"
	"github.com/vutranhoang1411/SimpleBank/token"
	"github.com/vutranhoang1411/SimpleBank/util"
)
type Server struct{
	config util.Config
	store db.Store
	router *gin.Engine
	maker token.Maker
}

func NewServer(config util.Config,store db.Store) (*Server,error){
	server:=&Server{store:store}
	//config
	server.config=config;

	//token maker
	var err error
	server.maker,err=token.NewPasetoMaker(config.KeyString)
	if err!=nil{
		return nil,err
	}

	//router
	router:=gin.Default();
	if v,ok:=binding.Validator.Engine().(*validator.Validate);ok{
		v.RegisterValidation("currency",validCurrency)
		v.RegisterValidation("email",validEmail)
	}

	router.POST("/user/",server.createUser)
	router.POST("/user/login/",server.userLogin)
	//register router
	authorize_route:=router.Group("/").Use(token_validator(server.maker));

	authorize_route.POST("/account/",server.createAccount)
	authorize_route.GET("/account/:id",server.getAccountByID)
	authorize_route.GET("/account/",server.listAccounts)
	authorize_route.GET("/token/new/",server.newAcessToken)
	server.router=router
	return server,nil
}
func (server *Server) Start (addr string)error{
	return server.router.Run(addr)
}
func handleError(err error)gin.H{
	return gin.H{"error":err.Error()}
}