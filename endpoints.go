package microServices

import (
	"golang.org/x/net/context"

	 pb "MicroService/microServices/pb"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	SayHelloEndpoint  endpoint.Endpoint
}


func (e Endpoints)SayHello(ctx context.Context,i interface{})(interface{},error){
	return e.SayHelloEndpoint(ctx,i)
}

func MakeEndpoints(s Service) Endpoints{
	return Endpoints{
		SayHelloEndpoint: MakeSayHelloEndpoints(s),
	}
}

func MakeSayHelloEndpoints(s Service) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.HelloRequest)
		return s.SayHello(ctx,req)
	}
}


