package main

import (
	"context"
	"dbcore/pbfiles"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func makeParam(m map[string]interface{}) *pbfiles.SimpleParam {
	paramStruct, err := structpb.NewStruct(m)
	checkError(err)

	return &pbfiles.SimpleParam{
		Params: paramStruct,
	}
}

func main() {
	client, err := grpc.DialContext(context.Background(),
		"localhost:8080",
		grpc.WithInsecure(),
	)
	checkError(err)

	c := pbfiles.NewDBServiceClient(client)

	txClient, err := c.Tx(context.Background()) // 执行事务
	checkError(err)

	addUserParam := makeParam(map[string]interface{}{
		"name": "lain",
		"age":  17,
	})
	err = txClient.Send(&pbfiles.TxRequest{Name: "adduser", Params: addUserParam, Type: "exec"})
	checkError(err)

	addUseRsp, err := txClient.Recv()
	checkError(err)
	ret := addUseRsp.Result.AsMap()
	userId := ret["exec"].([]interface{})[1].(map[string]interface{})["user_id"]
	fmt.Println(ret)

	addScoreParam := makeParam(map[string]interface{}{
		"userId": userId,
		"score":  10,
	})
	err = txClient.Send(&pbfiles.TxRequest{Name: "adduserscore", Params: addScoreParam, Type: "exec"})
	checkError(err)

	addScoreRsp, err := txClient.Recv()
	checkError(err)
	fmt.Println(addScoreRsp.Result.AsMap())

	err = txClient.CloseSend()
	checkError(err)
	log.Fatal("结束")

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
