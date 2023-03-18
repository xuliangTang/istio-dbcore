package main

import (
	"context"
	"dbcore/pbfiles"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
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

	structPb, err := structpb.NewStruct(map[string]interface{}{
		"name": "txl",
		"age":  18,
	})

	params := &pbfiles.SimpleParam{
		Params: structPb,
	}

	req := &pbfiles.ExecRequest{Name: "adduser", Params: params}
	rsp := &pbfiles.ExecResponse{}
	err = client.Invoke(context.Background(),
		"/DBService/Exec", req, rsp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rsp.RowsAffected, rsp.Select.AsMap())
	//for _, item := range rsp.Result {
	//	fmt.Println(item.AsMap())
	//}
}
