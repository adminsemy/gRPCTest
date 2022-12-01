package main

import (
	"context"
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
	id := &wrappers.StringValue{Value: "1"}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	order, err := c.GetOrder(ctx, id)
	if err != nil {
		log.Fatalf("Could not Order ID %v", id.Value)
	}
	log.Println("Order - ", order)
}
