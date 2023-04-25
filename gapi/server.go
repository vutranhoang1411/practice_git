package gapi

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
	db "github.com/vutranhoang1411/SimpleBank/db/sqlc"
	pb "github.com/vutranhoang1411/SimpleBank/pb/proto"
	"github.com/vutranhoang1411/SimpleBank/token"
	"github.com/vutranhoang1411/SimpleBank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)
type Server struct{
	pb.UnimplementedSimpleBankServer
	config util.Config
	store db.Store
	maker token.Maker
}

func NewServer(config util.Config,store db.Store) (pb.SimpleBankServer,error){
	server:=&Server{store:store}
	//config
	server.config=config;

	//token maker
	var err error
	server.maker,err=token.NewPasetoMaker(config.KeyString)
	if err!=nil{
		return nil,err
	}
	return server,err
}
func (server Server) CreateUser(ctx context.Context,rq *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	password,err:=util.HashPassword(rq.GetPassword());
	if (err!=nil){
		return nil,status.Errorf(codes.Internal,"Unable to hash the password")
	}

	user,err:=server.store.CreateUser(ctx,db.CreateUserParams{
		ID:uuid.NewString()[2:18],
		Name:rq.GetUsername(),
		Email:rq.GetEmail(),
		Password: password,
	})
	if err!=nil{
		if pqError,ok:=err.(*pq.Error);ok{
			switch pqError.Code.Name(){
			case "unique_violation":
				return nil,status.Errorf(codes.AlreadyExists,"User already exist")
			}
		}
		return nil,status.Errorf(codes.Internal,"Can't create new user")
	}

	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:user.ID,
			Username: user.Name,
			Email:user.Email,
			Password: user.Password,
			CreatedAt: &timestamppb.Timestamp{Seconds: int64(user.CreatedAt.Second())},
		},
	},err
	
	// return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}