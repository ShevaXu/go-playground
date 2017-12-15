package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

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
	s := grpc.NewServer(grpc.Creds(cred))

	pb.RegisterContactsManagerServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("Serving Contacts Manager @", addr)
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
