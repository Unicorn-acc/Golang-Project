package serializer

import "example.com/unicorn-acc/pkg/e"

// 通用返回类(基础序列化器)
type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Error  string      `json:"error"`
}

// DataList 带有总数的Data结构
type DataList struct {
	Item  interface{} `json:"item"`
	Total uint        `json:"total"`
}

// TokenData 带有token的Data结构
type TokenData struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

func BuildListResponse(items interface{}, total uint) Response {
	return Response{
		Status: 200,
		Data: DataList{
			Item:  items,
			Total: total,
		},
		Msg: "ok",
	}
}

// 每次验证err != nil 进行返回时都需要手动写这部分代码，直接抽取
func Result(code int) Response {
	return Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
