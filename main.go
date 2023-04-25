package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
	"github.com/vutranhoang1411/SimpleBank/api"
	db "github.com/vutranhoang1411/SimpleBank/db/sqlc"
	"github.com/vutranhoang1411/SimpleBank/gapi"
	pb "github.com/vutranhoang1411/SimpleBank/pb/proto"
	"github.com/vutranhoang1411/SimpleBank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)
const(

)
func main(){
	config,err:=util.LoadConfig("./");
	if (err!=nil){
		log.Fatal(err);
	}
	conn,err:=sql.Open(config.DBDriver,config.DBSource);
	if err!=nil{
		log.Fatal(err);
	}
	store:=db.NewStore(conn)
	startGRPCServer(config,store)

}
func startGRPCServer(config util.Config,store db.Store){
	//new server
	server,err:=gapi.NewServer(config,store)
	if err!=nil{
		log.Fatal(err)
	}
	//register grpc server
	grpcServer:=grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer,server)
	reflection.Register(grpcServer)

	listener,err:=net.Listen("tcp",config.GrpcServerAddress)
	log.Print("Server running on port: ",config.GrpcServerAddress);
	err=grpcServer.Serve(listener)
	if err!=nil{
		log.Fatal(err)
	}
	
}
func startHTTPServer(config util.Config,store db.Store){
	server,err:=api.NewServer(config,store)
	if err!=nil{
		log.Fatal(err);
	}
	err=server.Start(config.HttpServerAddress);
	if err!=nil{
		log.Fatal(err)
	}
	log.Print("Server running on port: ",config.HttpServerAddress);
}