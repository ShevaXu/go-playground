package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/Shevaxu/playground/grpc/contacts"
)

const addr = ":15001"

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	cl := pb.NewContactsManagerClient(conn)
	ctx := context.Background()

	phones := append(make([]*pb.PhoneNumber, 0), &pb.PhoneNumber{
		Number: "123456789",
	})
	req := &pb.AddPersonRequest{
		Name:  "sheva",
		Email: "sheva@email.com",
		Phone: phones,
	}
	resp, err := cl.AddPerson(ctx, req)
	if err != nil {
		log.Fatalf("could not add: %v", err)
	}
	log.Println(resp)

	req2 := &pb.GetPersonRequest{
		Id: "1",
	}
	resp2, err := cl.GetPerson(ctx, req2)
	if err != nil {
		log.Fatalf("could not get: %v", err)
	}
	log.Println(resp2)
}
