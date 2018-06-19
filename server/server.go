package main

import (
	"net"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/grpc"

	pb "MicroService/microServices/pb"
	"MicroService/microServices"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)


var (
	// Create a metrics registry.
	//创建一个指标注册表。
	reg = prometheus.NewRegistry()

	// Create some standard server metrics.
	//创建一些标准服务器指标。
	grpcMetrics = grpc_prometheus.NewServerMetrics()

	// Create a customized counter metric.
	//创建一个自定义计数器指标
	//记录say_hello的请求数
	customizedCounterMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
		Name: "demo_server_say_hello_method_handle_count",
		Help: "Total number of RPCs handled on the server.",
	}, []string{"name"})

)

func init() {
	// Register standard server metrics and customized metrics to registry.
	//向注册表注册标准服务器指标和自定义指标。
	reg.MustRegister(grpcMetrics, customizedCounterMetric)
	customizedCounterMetric.WithLabelValues("Test")
}


func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9093))
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
	}

	defer lis.Close()

	// Create a HTTP server for prometheus.
	//为prometheus创建一个HTTP服务器。
	httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
	Addr: fmt.Sprintf("0.0.0.0:%d", 9092)}

	// Create a gRPC Server with gRPC interceptor.
	//使用gRPC拦截器创建gRPC服务器。
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
	)

	// Create a new api server.
	demoServer := microServices.NewDemoServer()

	// Register your service.
	//创建一个新的api服务
	pb.RegisterDemoServiceServer(grpcServer, demoServer)

	// Initialize all metrics.
	//初始化所有度量标准。
	grpcMetrics.InitializeMetrics(grpcServer)

	// Start your http server for prometheus.
	//为prometheus启动你的http服务器
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	// Start your gRPC server.
	//启动你的gRPC服务器
	log.Fatal(grpcServer.Serve(lis))
}