package main

import (
	"context"
	"io"
	"log"
	"time"

	orderManager "ordermagager/client/ecommerce/proto"

	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connct server %v: %v", address, err)
	}
	defer conn.Close()
	c := orderManager.NewOrderManagerClient(conn)
	id := &wrappers.StringValue{Value: "3"}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	order, err := c.GetOrder(ctx, id)
	if err != nil {
		log.Fatalf("Could not Order ID %v: %v", id.Value, err)
	}
	log.Println("Order - ", order)
	searchStream, _ := c.SearchOrders(ctx, id)

	for {
		searchOrder, err := searchStream.Recv()
		if searchOrder == nil {
			break
		}
		log.Println("Search result: ", searchOrder)
		if err == io.EOF {
			break
		}
	}
	log.Println("UPDATE ORDERS 1,2,3")
	updateStream, err := c.UpdateOrders(ctx)
	if err != nil {
		log.Fatalf("%v.UpdateOrders(_) = _, %v", c, err)
	}

	updOrder1 := orderManager.Order{
		Id:          "1",
		Items:       []string{"1", "2", "3"},
		Description: "First Order",
		Price:       8.0,
		Destination: "What is this",
	}
	updOrder2 := orderManager.Order{
		Id:          "2",
		Items:       []string{"1", "2", "3"},
		Description: "Second Order",
		Price:       16.0,
		Destination: "What is this",
	}
	updOrder3 := orderManager.Order{
		Id:          "3",
		Items:       []string{"1", "2", "3"},
		Description: "Third Order",
		Price:       32.0,
		Destination: "What is this",
	}
	if err := updateStream.Send(&updOrder1); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateStream, &updOrder1, err)
	}
	if err := updateStream.Send(&updOrder2); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateStream, &updOrder2, err)
	}
	if err := updateStream.Send(&updOrder3); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateStream, &updOrder3, err)
	}

	updateRes, err := updateStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndResv() got error %v, want %v", updateStream, err, nil)
	}
	log.Printf("Updates Orders Res: %s", updateRes)

}
