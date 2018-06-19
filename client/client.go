package main

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "MicroService/microServices/pb"

)

func main() {
	conn,err := grpc.Dial(fmt.Sprintf("localhost:%v",9093), grpc.WithInsecure())
	if err != nil{
		log.Fatal(err)
	}

	defer conn.Close()

	// Create a gRPC server client.
	client := pb.NewDemoServiceClient(conn)
	// Call “SayHello” method and wait for response from gRPC Server.
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "Test"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)

}