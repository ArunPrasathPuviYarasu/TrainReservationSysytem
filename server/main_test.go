package main

import (
	"context"
	"testing"

	pb "example.com/TrainReservatioSystem/proto"
)

func TestPurchaseTicket(t *testing.T) {
	s := newServer()

	req := &pb.PurchaseTicketRequest{
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
	}

	res, err := s.PurchaseTicket(context.Background(), req)
	if err != nil {
		t.Fatalf("PurchaseTicket failed: %v", err)
	}

	if res.User.FirstName != "John" || res.User.LastName != "Doe" || res.User.Email != "john.doe@example.com" {
		t.Errorf("PurchaseTicket returned incorrect user data")
	}

	if res.PricePaid != 20.0 {
		t.Errorf("PurchaseTicket returned incorrect price")
	}
}

func TestGetReceipt(t *testing.T) {
	s := newServer()

	req := &pb.PurchaseTicketRequest{
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
	}

	res, err := s.PurchaseTicket(context.Background(), req)
	if err != nil {
		t.Fatalf("PurchaseTicket failed: %v", err)
	}

	receiptReq := &pb.GetReceiptRequest{
		ReceiptId: res.ReceiptId,
	}

	receiptRes, err := s.GetReceipt(context.Background(), receiptReq)
	if err != nil {
		t.Fatalf("GetReceipt failed: %v", err)
	}

	if receiptRes.User.FirstName != "John" || receiptRes.User.LastName != "Doe" || receiptRes.User.Email != "john.doe@example.com" {
		t.Errorf("GetReceipt returned incorrect user data")
	}

	if receiptRes.PricePaid != 20.0 {
		t.Errorf("GetReceipt returned incorrect price")
	}
}

func TestViewUsersBySection(t *testing.T) {
	s := newServer()

	req := &pb.PurchaseTicketRequest{
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
	}

	_, err := s.PurchaseTicket(context.Background(), req)
	if err != nil {
		t.Fatalf("PurchaseTicket failed: %v", err)
	}

	sectionReq := &pb.ViewUsersBySectionRequest{
		Section: "A",
	}

	sectionRes, err := s.ViewUsersBySection(context.Background(), sectionReq)
	if err != nil {
		t.Fatalf("ViewUsersBySection failed: %v", err)
	}

	if len(sectionRes.UserSeatAllocations) == 0 {
		t.Errorf("ViewUsersBySection returned no users")
	}
}

func TestRemoveUser(t *testing.T) {
	s := newServer()

	req := &pb.PurchaseTicketRequest{
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
	}

	_, err := s.PurchaseTicket(context.Background(), req)
	if err != nil {
		t.Fatalf("PurchaseTicket failed: %v", err)
	}

	removeReq := &pb.RemoveUserRequest{
		Email: "john.doe@example.com",
	}

	removeRes, err := s.RemoveUser(context.Background(), removeReq)
	if err != nil {
		t.Fatalf("RemoveUser failed: %v", err)
	}

	if !removeRes.Success {
		t.Errorf("RemoveUser failed: %v", removeRes.Message)
	}
}

func TestModifySeat(t *testing.T) {
	s := newServer()

	req := &pb.PurchaseTicketRequest{
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
	}

	_, err := s.PurchaseTicket(context.Background(), req)
	if err != nil {
		t.Fatalf("PurchaseTicket failed: %v", err)
	}

	modifyReq := &pb.ModifySeatRequest{
		Email: "john.doe@example.com",
		NewSeat: &pb.Seat{
			Section: "B",
			Number:  2,
		},
	}

	modifyRes, err := s.ModifySeat(context.Background(), modifyReq)
	if err != nil {
		t.Fatalf("ModifySeat failed: %v", err)
	}

	if !modifyRes.Success {
		t.Errorf("ModifySeat failed: %v", modifyRes.Message)
	}
}
