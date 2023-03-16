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
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
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

func Hello(req *proto_api.HelloRequest) (*proto_api.HelloResponse, error) {
	message := "Hello, "
	if req.Title != "" {
		message += req.Title + " "
	}
	message += req.Name + "!"
	return &proto_api.HelloResponse{Message: message, Method: req.Method}, nil
}

func (s *greetingService) GetHello(ctx context.Context, req *proto_api.HelloRequest) (*proto_api.HelloResponse, error) {
	return Hello(req)
}
func (s *greetingService) PostHello(ctx context.Context, req *proto_api.HelloRequest) (*proto_api.HelloResponse, error) {
	return Hello(req)
}
func (s *greetingService) DeleteHello(ctx context.Context, req *proto_api.HelloRequest) (*proto_api.HelloResponse, error) {
	return Hello(req)
}
func (s *greetingService) PutHello(ctx context.Context, req *proto_api.HelloRequest) (*proto_api.HelloResponse, error) {
	return Hello(req)
}
func (s *greetingService) PatchHello(ctx context.Context, req *proto_api.HelloRequest) (*proto_api.HelloResponse, error) {
	return Hello(req)
}

func (s *greetingService) Error(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, fmt.Errorf("server error occur")
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

	s := grpc.NewServer(
		grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			requestMD, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				requestMD = metadata.MD{}
			}
			log.Printf("api: %#v\n", info.FullMethod)
			log.Printf("request metadata: %#v\n", requestMD)
			log.Printf("request body: %#v\n", req)
			res, err := handler(ctx, req)
			responseHeaderMD := metadata.New(map[string]string{"header-metadata": "header:" + info.FullMethod})
			if err := grpc.SetHeader(ctx, responseHeaderMD); err != nil {
				return nil, err
			}
			log.Printf("response header metadata: %#v\n", responseHeaderMD)
			responseTrailerMD := metadata.New(map[string]string{"trailer-metadata": "trailer:" + info.FullMethod})
			if err := grpc.SetTrailer(ctx, responseTrailerMD); err != nil {
				return nil, err
			}
			log.Printf("response trailer metadata: %#v\n", responseTrailerMD)
			log.Printf("response body: %#v\n", res)
			log.Printf("error: %#v\n", err)
			return res, err
		}),
	)
	proto_api.RegisterGreetingServiceServer(s, &greetingService{})
	reflection.Register(s)

	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
