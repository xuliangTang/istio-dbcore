package main

import (
	"context"
	"dbcore/pbfiles"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func main() {
	client, err := grpc.DialContext(context.Background(),
		"localhost:8080",
		grpc.WithInsecure(),
	)

	if err != nil {
		log.Fatal(err)
	}
	req := &pbfiles.QueryRequest{Name: "userlist"}
	rsp := &pbfiles.QueryResponse{}
	err = client.Invoke(context.Background(),
		"/DBService/Query", req, rsp)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range rsp.Result {
		fmt.Println(item.AsMap())
	}
}
