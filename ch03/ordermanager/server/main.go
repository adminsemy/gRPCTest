package main

import (
	"context"
	"fmt"
	"io"
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
	port           = ":50051"
	orderBatchSize = 3
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

func (s *server) UpdateOrders(stream orderManager.OrderManager_UpdateOrdersServer) error {
	ordersStr := "Update Order is: "
	for {
		order, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&wrappers.StringValue{Value: "Orders processed " + ordersStr})
		}
		s.orders[order.Id] = order
		log.Print("Order ID ", order.Id, ": updated!")
		ordersStr += order.Id + ", "
		time.Sleep(time.Second)
	}
}

func (s *server) ProcessOrders(stream orderManager.OrderManager_ProcessOrdersServer) error {
	var combinedShipmentMap = make(map[string]*orderManager.CombinedShipment)
	batchMarker := 1
	for {
		orderId, err := stream.Recv()
		if err == io.EOF {
			for _, comb := range combinedShipmentMap {
				stream.Send(comb)
			}
			return nil
		}
		if err != nil {
			return err
		}
		destination := s.orders[orderId.GetValue()].Destination
		shipment, ok := combinedShipmentMap[destination]
		if ok {
			ord := s.orders[orderId.GetValue()]
			shipment.OrdersList = append(shipment.OrdersList, ord)
			combinedShipmentMap[destination] = shipment
		} else {
			comShip := orderManager.CombinedShipment{
				Id:     "cmb - " + (s.orders[orderId.GetValue()].Destination),
				Status: "Processed",
			}
			ord := s.orders[orderId.GetValue()]
			comShip.OrdersList = append(comShip.OrdersList, ord)
			combinedShipmentMap[destination] = &comShip
			log.Print(len(comShip.OrdersList), comShip.GetId())
		}
		if batchMarker == orderBatchSize {
			for _, comb := range combinedShipmentMap {
				log.Printf("Shipping: %v -> %v", comb.Id, len(comb.OrdersList))
				if err := stream.Send(comb); err != nil {
					return err
				}
			}
			batchMarker = 0
			combinedShipmentMap = make(map[string]*orderManager.CombinedShipment)
		} else {
			batchMarker++
		}
	}
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
