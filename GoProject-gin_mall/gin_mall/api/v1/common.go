package v1

import (
	"example.com/unicorn-acc/serializer"
	"fmt"
	"github.com/goccy/go-json"
)

// ErrorResponse 因为error不是json格式，在contorller层发送错误返回错误信息时需要返回json格式的err
func ErrorResponse(err error) serializer.Response {
	//if ve, ok := err.(validator.ValidationErrors); ok {
	//	for _, e := range ve {
	//		field := conf.T(fmt.Sprintf("Field.%s", e.Field))
	//		tag := conf.T(fmt.Sprintf("Tag.Valid.%s", e.Tag))
	//		return serializer.Response{
	//			Status: 400,
	//			Msg:    fmt.Sprintf("%s%s", field, tag),
	//			Error:  fmt.Sprint(err),
	//		}
	//	}
	//}
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return serializer.Response{
			Status: 400,
			Msg:    "JSON类型不匹配",
			Error:  fmt.Sprint(err),
		}
	}

	return serializer.Response{
		Status: 400,
		Msg:    "参数错误",
		Error:  fmt.Sprint(err),
	}
}