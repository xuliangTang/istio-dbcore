package core

import (
	"context"
	"dbcore/helpers"
	"dbcore/pbfiles"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
)

type DBService struct {
	pbfiles.UnimplementedDBServiceServer
}

func (this *DBService) Query(ctx context.Context, in *pbfiles.QueryRequest) (*pbfiles.QueryResponse, error) {
	api := SysConfig.FindAPI(in.Name)
	if api == nil {
		return nil, status.Error(codes.Unavailable, "error api name")
	}

	ret, err := api.Query(in.Params) // 返回值是一个map[string]interface{}
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	// 把 map 转化为 StructList
	structList, err := helpers.MapListToStructList(ret)
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &pbfiles.QueryResponse{
		Message: "success",
		Result:  structList,
	}, nil
}

func (this *DBService) Exec(ctx context.Context, in *pbfiles.ExecRequest) (*pbfiles.ExecResponse, error) {
	api := SysConfig.FindAPI(in.Name)
	if api == nil {
		return nil, status.Error(codes.Unavailable, "error api name")
	}

	// 判断客户端是否超时或主动取消
	if code, yes := helpers.ContextIsBroken(ctx); yes {
		return nil, status.Error(code, ctx.Err().Error())
	}

	rowsAffected, selectKey, err := api.ExecBySql(in.Params)
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	// 把 map[string]interface{} 转化为 protobuf struct
	selectKeyValue, _ := helpers.MapToStruct(selectKey)

	return &pbfiles.ExecResponse{
		Message:      "success",
		RowsAffected: rowsAffected,
		Select:       selectKeyValue,
	}, nil
}

func (this *DBService) Tx(c pbfiles.DBService_TxServer) error {
	tx := GormDB.Begin()

	for {
		txRequest, err := c.Recv()
		if err != nil {
			if err == io.EOF {
				log.Println("EOF")
				tx.Commit()
				return nil
			}

			log.Println("错误0:", err)
			tx.Rollback()
			return err
		}

		// 坑：client的CloseSend() 没有实现EOF标志，要么在服务端检查是否关闭连接，要么在客户端使用cancelContext
		/*if err = c.Context().Err(); err != nil {
			log.Printf("Client closed stream: %v", err)
			return nil
		}*/

		api := SysConfig.FindAPI(txRequest.Name)
		if api == nil {
			tx.Rollback()
			return status.Error(codes.Unavailable, "error api name")
		}

		api.SetDb(tx) // 设置db连接对象
		ret := make(map[string]interface{})
		if txRequest.Type == "query" {
			queryResult, err := api.Query(txRequest.Params)
			if err != nil {
				log.Println("错误1:", err)
				tx.Rollback()
				return err
			}

			// ret["query"] = queryResult
			// 取出来的内容是 []map[string]interface 需要转换为[]interface{}，否则 MapToStruct 会出错
			ret["query"] = helpers.MapListToInterfaceList(queryResult)
		} else {
			rowsAffected, selectKey, err := api.ExecBySql(txRequest.Params)
			if err != nil {
				log.Println("错误2:", err)
				tx.Rollback()
				return err
			}
			ret["exec"] = []interface{}{rowsAffected, selectKey}
		}

		m, _ := helpers.MapToStruct(ret)
		err = c.Send(&pbfiles.TxResponse{
			Result: m,
		})
		if err != nil {
			log.Println("错误3:", err)
			tx.Rollback()
			return err
		}
	}
}
