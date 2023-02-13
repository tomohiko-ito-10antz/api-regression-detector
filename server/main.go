package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	proto_api "github.com/Jumpaku/api-regression-detector/server/gen/proto/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	grpcPort := 50051
	grpcHostPort := fmt.Sprintf(":%d", grpcPort)
	go runGRPCGateway(grpcHostPort)
	runGRPCServer(grpcPort)
}

type greetingService struct {
	proto_api.UnimplementedGreetingServiceServer
}

func (s *greetingService) SayHello(ctx context.Context, in *proto_api.HelloRequest) (*proto_api.HelloResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &proto_api.HelloResponse{Message: "Hello " + in.Name}, nil
}

func runGRPCGateway(grpcHostPort string) {
	gwPort := 80
	mux := runtime.NewServeMux()
	err := proto_api.RegisterGreetingServiceHandlerFromEndpoint(context.Background(), mux,
		grpcHostPort,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		})
	if err != nil {
		log.Fatalf("failed to register handler: %v", err)
	}

	log.Printf("grpc gateway listening at %v", gwPort)
	err = http.ListenAndServe(fmt.Sprintf(":%d", gwPort), mux)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}

func runGRPCServer(port int) {
	grpcHostPort := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", grpcHostPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto_api.RegisterGreetingServiceServer(s, &greetingService{})
	reflection.Register(s)

	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
