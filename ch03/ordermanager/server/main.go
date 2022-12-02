package main

import (
	"context"
	"fmt"
	"log"
	"net"
	orderManager "ordermagager/server/ecommerce/proto"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	port = ":50051"
)

type server struct {
	orders map[string]*orderManager.Order
	orderManager.UnimplementedOrderManagerServer
}

func (s *server) GetOrder(ctx context.Context, id *wrappers.StringValue) (*orderManager.Order, error) {
	order, ok := s.orders[id.Value]
	if ok {
		return order, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "Order %v not found", id.Value)
}

func (s *server) SearchOrders(searchQuery *wrappers.StringValue,
	stream orderManager.OrderManager_SearchOrdersServer) error {
	for key, order := range s.orders {
		log.Print(key, order)
		for _, itemStr := range order.Items {
			log.Print(itemStr)
			if strings.Contains(itemStr, searchQuery.Value) {
				err := stream.Send(order)
				if err != nil {
					return fmt.Errorf("error send message to stream: %v", err)
				}
				log.Println("Matching Order Found: ", key)
			}
		}
		time.Sleep(time.Second)
	}
	return nil
}

func main() {
	serv := &server{}
	serv.orders = map[string]*orderManager.Order{
		"1": {Id: "1", Items: []string{"1", "2", "3"}, Description: "First", Price: 4.0, Destination: "..."},
		"2": {Id: "2", Items: []string{"2", "3", "1"}, Description: "Twice", Price: 12.0, Destination: "..."},
		"3": {Id: "3", Items: []string{"3", "1", "2"}, Description: "Third", Price: 15.0, Destination: "..."},
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen port %v - %v", port, err)
	}
	s := grpc.NewServer()
	orderManager.RegisterOrderManagerServer(s, serv)
	log.Println("Start listener server on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
