package server

import (
	"context"
	"log"
	pb "project/proto/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	mongoClient *mongo.Client
	pb.UnimplementedUserServiceServer
}

func NewServer(mc *mongo.Client) pb.UserServiceServer {
	return &Server{
		mongoClient: mc,
	}
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Printf("Received User Creation Request %v \n", req.GetUser())
	InsertedID, err := s.CreateUserHelper(ctx, s.mongoClient, req)

	if err != nil {
		return nil, err
	}

	res := &pb.CreateUserResponse{
		User:       req.GetUser(),
		InsertedId: InsertedID.(primitive.ObjectID).Hex(),
	}
	return res, nil
}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("Received User Finding Request %s \n", req.GetName())
	resp, err := s.GetUserHelper(ctx, s.mongoClient, req)

	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{User: resp}, nil

}

func (s *Server) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
    log.Printf("Received Users Finding Request")
	protoUsers, err := s.GetUsersHelper(ctx, s.mongoClient, req)

	if err != nil {
		return nil, err
	}
    
	return &pb.GetUsersResponse{ Userlist: protoUsers}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
    log.Printf("Received User Geletion Request: %s \n", req.GetName())
	count, err := s.DeleteUserHelper(ctx, s.mongoClient, req)
	if err != nil {
		return nil, err
	}
    
	res:=&pb.DeleteUserResponse{
		Count: int32(count),
	}
	return res, nil
}
func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	log.Printf("Received User Updation Request: %s \n", req.GetName())
    count, err:= s.UpdateUserHelper(ctx,s.mongoClient, req)
	if err != nil {
		return nil, err
	}

	res:= &pb.UpdateUserResponse{
	  Count: int32(count),
	}
	return res, nil
}

// GetUser fetches User details from mongodb
