package helpers

import (
	"google.golang.org/protobuf/types/known/structpb"
)

func MapListToInterfaceList(m []map[string]interface{}) []interface{} {
	ret := make([]interface{}, len(m))
	for i, item := range m {
		ret[i] = item
	}
	return ret
}

func MapToStruct(m map[string]interface{}) (*structpb.Struct, error) {
	if m == nil {
		return nil, nil
	}

	s, err := structpb.NewStruct(m)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func MapListToStructList(m []map[string]interface{}) ([]*structpb.Struct, error) {
	ret := make([]*structpb.Struct, len(m))
	iList := MapListToInterfaceList(m)
	vList, err := structpb.NewList(iList) // 把 map list 转为 ListValue
	if err != nil {
		return nil, err
	}

	for i, item := range vList.GetValues() {
		ret[i] = item.GetStructValue() // 把 value 转为 pb struct
	}

	return ret, nil
}
