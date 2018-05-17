package microServices

import (
	"fmt"

	"golang.org/x/net/context"
	pb "MicroService/microServices/pb"
)

type demoServiceServer struct {}


func NewDemoServer() *demoServiceServer{
	return &demoServiceServer{}
}


func (s *demoServiceServer) SayHello(ctx context.Context,request *pb.HelloRequest)(*pb.HelloResponse,error){
	return &pb.HelloResponse{Message:fmt.Sprintf("Hello %s",request.Name)},nil
}

