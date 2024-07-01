package main

import (
	"context"
	"log"
	"time"

	pb "example.com/TrainReservatioSystem/proto"

	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewTrainServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Purchase a ticket
	purchaseReq := &pb.PurchaseTicketRequest{
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
	}
	purchaseRes, err := c.PurchaseTicket(ctx, purchaseReq)
	if err != nil {
		log.Fatalf("could not purchase ticket: %v", err)
	}
	log.Printf("PurchaseTicket: %v", purchaseRes)

	// Get the receipt
	receiptReq := &pb.GetReceiptRequest{
		ReceiptId: purchaseRes.ReceiptId,
	}
	receiptRes, err := c.GetReceipt(ctx, receiptReq)
	if err != nil {
		log.Fatalf("could not get receipt: %v", err)
	}
	log.Printf("GetReceipt: %v", receiptRes)

	// View users by section
	viewReq := &pb.ViewUsersBySectionRequest{
		Section: "A",
	}
	viewRes, err := c.ViewUsersBySection(ctx, viewReq)
	if err != nil {
		log.Fatalf("could not view users: %v", err)
	}
	log.Printf("ViewUsersBySection: %v", viewRes)

	// Remove a user
	removeReq := &pb.RemoveUserRequest{
		Email: "john.doe@example.com",
	}
	removeRes, err := c.RemoveUser(ctx, removeReq)
	if err != nil {
		log.Fatalf("could not remove user: %v", err)
	}
	log.Printf("RemoveUser: %v", removeRes)

	// Modify a seat
	modifyReq := &pb.ModifySeatRequest{
		Email: "john.doe@example.com",
		NewSeat: &pb.Seat{
			Section: "B",
			Number:  2,
		},
	}
	modifyRes, err := c.ModifySeat(ctx, modifyReq)
	if err != nil {
		log.Fatalf("could not modify seat: %v", err)
	}
	log.Printf("ModifySeat: %v", modifyRes)
}
