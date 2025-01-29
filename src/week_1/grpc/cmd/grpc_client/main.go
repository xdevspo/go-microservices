package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	desc "github.com/xdevspo/go-microservices/week_1/grpc/pkg/note_v1"
)

const (
	address = "localhost:50051"
	noteID  = 12
)

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("error closing connection: %v", err)
		}
	}(conn)

	c := desc.NewNoteV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &desc.GetRequest{Id: noteID})
	if err != nil {
		fmt.Printf("could not get note: %v", err)
	}

	log.Printf(color.RedString("Note info:\n"), color.GreenString("%+v", r.GetNote()))
}
