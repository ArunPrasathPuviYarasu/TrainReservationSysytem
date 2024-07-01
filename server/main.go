package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	pb "example.com/TrainReservatioSystem/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTrainServiceServer
	receipts    map[string]*pb.PurchaseTicketResponse
	seats       map[string][]*pb.PurchaseTicketResponse
	seatCounter int32
	mu          sync.Mutex
}

func newServer() *server {
	return &server{
		receipts:    make(map[string]*pb.PurchaseTicketResponse),
		seats:       map[string][]*pb.PurchaseTicketResponse{"A": {}, "B": {}},
		seatCounter: 1,
	}
}

func (s *server) PurchaseTicket(ctx context.Context, req *pb.PurchaseTicketRequest) (*pb.PurchaseTicketResponse, error) {
	if req.User.FirstName == "" || req.User.LastName == "" || req.User.Email == "" {
		return nil, grpc.Errorf(grpc.Code(errors.New("invalid argument")), "User information is incomplete")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	receiptID := fmt.Sprintf("receipt_%d", len(s.receipts)+1)
	section := "A"
	if len(s.seats["A"]) > len(s.seats["B"]) {
		section = "B"
	}
	seat := &pb.Seat{Section: section, Number: s.seatCounter}
	s.seatCounter++

	receipt := &pb.PurchaseTicketResponse{
		ReceiptId: receiptID,
		From:      "London",
		To:        "France",
		User:      req.User,
		PricePaid: 20.0,
		Seat:      seat,
	}

	s.receipts[receiptID] = receipt
	s.seats[section] = append(s.seats[section], receipt)

	return receipt, nil
}

func (s *server) GetReceipt(ctx context.Context, req *pb.GetReceiptRequest) (*pb.GetReceiptResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	receipt, ok := s.receipts[req.ReceiptId]
	if !ok {
		return nil, grpc.Errorf(grpc.Code(errors.New("not found")), "Receipt not found")
	}

	return &pb.GetReceiptResponse{
		From:      receipt.From,
		To:        receipt.To,
		User:      receipt.User,
		PricePaid: receipt.PricePaid,
		Seat:      receipt.Seat,
	}, nil
}

func (s *server) ViewUsersBySection(ctx context.Context, req *pb.ViewUsersBySectionRequest) (*pb.ViewUsersBySectionResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users, ok := s.seats[req.Section]
	if !ok {
		return nil, grpc.Errorf(grpc.Code(errors.New("not found")), "Section not found")
	}

	userSeatAllocations := make([]*pb.UserSeatAllocation, 0)
	for _, receipt := range users {
		userSeatAllocations = append(userSeatAllocations, &pb.UserSeatAllocation{
			User: receipt.User,
			Seat: receipt.Seat,
		})
	}

	return &pb.ViewUsersBySectionResponse{UserSeatAllocations: userSeatAllocations}, nil
}

func (s *server) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	if req.Email == "" {
		return &pb.RemoveUserResponse{Success: false, Message: "Email is required"}, grpc.Errorf(grpc.Code(errors.New("invalid argument")), "Email is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for section, receipts := range s.seats {
		for i, receipt := range receipts {
			if receipt.User.Email == req.Email {
				s.seats[section] = append(receipts[:i], receipts[i+1:]...)
				delete(s.receipts, receipt.ReceiptId)
				return &pb.RemoveUserResponse{Success: true, Message: "User removed successfully"}, nil
			}
		}
	}

	return &pb.RemoveUserResponse{Success: false, Message: "User not found"}, grpc.Errorf(grpc.Code(errors.New("not found")), "User not found")
}

func (s *server) ModifySeat(ctx context.Context, req *pb.ModifySeatRequest) (*pb.ModifySeatResponse, error) {
	if req.Email == "" || req.NewSeat.Section == "" || req.NewSeat.Number <= 0 {
		return &pb.ModifySeatResponse{Success: false, Message: "Invalid email or seat information"}, grpc.Errorf(grpc.Code(errors.New("invalid argument")), "Invalid email or seat information")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for section, receipts := range s.seats {
		for i, receipt := range receipts {
			if receipt.User.Email == req.Email {
				s.seats[section] = append(receipts[:i], receipts[i+1:]...)
				receipt.Seat.Section = req.NewSeat.Section
				receipt.Seat.Number = req.NewSeat.Number
				s.seats[req.NewSeat.Section] = append(s.seats[req.NewSeat.Section], receipt)
				return &pb.ModifySeatResponse{Success: true, Message: "Seat modified successfully"}, nil
			}
		}
	}

	return &pb.ModifySeatResponse{Success: false, Message: "User not found"}, grpc.Errorf(grpc.Code(errors.New("not found")), "User not found")
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTrainServiceServer(s, newServer())

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
