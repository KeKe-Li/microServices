package microServices

import (
	"fmt"

	"golang.org/x/net/context"
	pb "MicroService/microServices/pb"
)

type Service interface {
	SayHello(ctx context.Context,request interface{})(interface{},error)
}

type DemoServiceServer struct {}

func NewDemoServer() *DemoServiceServer{
	return &DemoServiceServer{}
}


func (s *DemoServiceServer) SayHello(ctx context.Context,request *pb.HelloRequest)(*pb.HelloResponse,error){
	return &pb.HelloResponse{Message:fmt.Sprintf("Hello %s",request.Name)},nil
}
