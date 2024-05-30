package server

import (
	"context"
	"encoding/json"
	"log"
	app "project/config"
	"project/mongodb"
	pb "project/proto/user"
	model "project/schema"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) CreateUserHelper(ctx context.Context, mg *mongo.Client, req *pb.CreateUserRequest) (interface{}, error) {
	userInfo := &model.User{}
	userInfo.ConvertToSchema(req.GetUser())
	res, err := mongodb.InsertOne(ctx, s.mongoClient, app.DbName, app.CollectionName, userInfo)
	if err != nil {
		log.Printf("Error in inserting the document")
		return "", err
	}
	return res.InsertedID, err
}

func (s *Server) GetUserHelper(ctx context.Context, mg *mongo.Client, req *pb.GetUserRequest) (*pb.User, error) {
	username := req.GetName()
	filter := bson.M{"name": username}
	res := mongodb.FindOne(ctx, s.mongoClient, app.DbName, app.CollectionName, filter)

	user := &model.User{}
	err := res.Decode(user)
	if err != nil {
		log.Printf("Not Found")
		return nil, err
	}
	return user.ConvertToProto(), nil
}

func (s *Server) GetUsersHelper(ctx context.Context, mg *mongo.Client, req *pb.GetUsersRequest) ([]*pb.User, error) {
	filter := bson.D{{}}
	cur, err := mongodb.FindMany(ctx, s.mongoClient, app.DbName, app.CollectionName, filter)
	if err != nil {
		log.Printf("Not Found any")
		return nil, err
	}

	users := &[]model.User{}
	if err = cur.All(ctx, users); err != nil {
		log.Printf("Failed to decode users: %v", err)
		return nil, err
	}

	protoUsers := make([]*pb.User, len(*users))
	for i, user := range *users {
		protoUsers[i] = user.ConvertToProto()
	}
	return protoUsers, nil
}

func (s *Server) DeleteUserHelper(ctx context.Context, mg *mongo.Client, req *pb.DeleteUserRequest) (int, error) {
	username := req.GetName()
	filter := bson.M{"name": username}

	res, err := mongodb.DeleteOne(ctx, s.mongoClient, app.DbName, app.CollectionName, filter)
	if err != nil {
		log.Printf("Error somewhere")
		return 0, err
	}

	return int(res.DeletedCount), nil
}

func (s *Server) UpdateUserHelper(ctx context.Context, mg *mongo.Client, req *pb.UpdateUserRequest) (int, error) {
	username := req.GetName()
	filter := bson.M{"name": username}

	var updateParams map[string]interface{}
	userBytes, _ := json.Marshal(req.GetUser())
	json.Unmarshal(userBytes, &updateParams)

	fields:= bson.M{"$set": updateParams}

	res, err:= mongodb.UpdateOne(ctx, s.mongoClient, app.DbName, app.CollectionName, filter,fields)
	
	if err != nil {
		log.Printf("Failed to update")
		return 0, err
	}
    return int(res.ModifiedCount), nil
}
