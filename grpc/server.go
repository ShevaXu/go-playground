package main

import (
	"log"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "github.com/Shevaxu/playground/grpc/contacts"
)

const addr = ":15001"

// server is used to implement contacts.ContaceManagerServer.
type server struct{}

// AddPerson implements contacts.ContaceManagerServer
func (s *server) AddPerson(ctx context.Context, req *pb.AddPersonRequest) (*pb.AddPersonResponse, error) {
	log.Println(req)
	return &pb.AddPersonResponse{
		Id: "1",
	}, nil
}

// GetPerson implements contacts.ContaceManagerServer
func (s *server) GetPerson(ctx context.Context, req *pb.GetPersonRequest) (*pb.GetPersonResponse, error) {
	log.Println(req)
	return &pb.GetPersonResponse{
		Name: "sheva",
		Id:   "1",
	}, nil
}

func authorize(ctx context.Context) error {
	// code from the authorize() function:
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "retrieving metadata failed")
	}

	// so far, just check its existence
	_, ok = md["authorization"]
	if !ok {
		return status.Errorf(codes.InvalidArgument, "no auth details supplied")
	}

	// passed
	return nil
}

// unaryInterceptor handles per RPC call
// with logging and authorization
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, h grpc.UnaryHandler) (interface{}, error) {
	// timing
	t := time.Now()
	defer log.Printf("Handle request %s in %s\n", info.FullMethod, time.Since(t))

	// auth first
	if err := authorize(ctx); err != nil {
		return nil, err
	}

	// let the original handler handle it
	return h(ctx, req)
}

func main() {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	cred, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("faild to create credentials")
	}

	// functional options
	s := grpc.NewServer(grpc.Creds(cred),
		grpc.UnaryInterceptor(unaryInterceptor))

	pb.RegisterContactsManagerServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("Serving Contacts Manager @", addr)
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
