package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	app "project/config"
	user_pb "project/proto/user"
	server "project/server"

	"project/mongodb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {

	//Mongo Connection
	mongoClient, err, ctx, cancelFunc := mongodb.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database %v \n", err)
	}
	defer mongodb.Close(mongoClient, ctx, cancelFunc)
    
	//Create a listener
	_, err = net.Listen("tcp", "localhost"+app.Serverport)
	if err != nil {
		log.Fatalf("Error in listening %v \n", err)
	}
	fmt.Printf("Server listening on localhost: %s \n", app.Serverport)

	//Register grpc server
	grpcServer := grpc.NewServer()
	userServer := server.NewServer(mongoClient)
	user_pb.RegisterUserServiceServer(grpcServer, userServer)
	fmt.Println("grpc Server Registered Successfully")

	//Register grpc gateway mux
	gwMux := runtime.NewServeMux()
	user_pb.RegisterUserServiceHandlerServer(context.Background(), gwMux, userServer)
 
	err = http.ListenAndServe("localhost"+app.Serverportmux, gwMux)
	if err != nil {
		log.Fatalf("Error in listening mux %v \n", err)
	}

}
