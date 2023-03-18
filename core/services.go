package core

import (
	"context"
	"dbcore/helpers"
	"dbcore/pbfiles"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	ret, err := api.Exec(in.Params)
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &pbfiles.ExecResponse{
		Message:      "success",
		RowsAffected: ret,
	}, nil
}
