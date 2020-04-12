package common

import "encoding/json"

type BaseResponse struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	DetailMsg string      `json:"detail_msg"`
	Data      interface{} `json:"data"`
}

func GetErrorResponse(code int, msg, detailMsg string) *BaseResponse {
	return &BaseResponse{
		Code:      code,
		DetailMsg: detailMsg,
		Msg:       msg,
		Data:      nil,
	}
}

func GetErrorResponseJson(code int, msg, detailMsg string) []byte {
	res, _ := json.Marshal(BaseResponse{
		Code:      code,
		DetailMsg: detailMsg,
		Msg:       msg,
		Data:      nil,
	})
	return res
}

func GetBaseResponse(data interface{}) *BaseResponse {
	return &BaseResponse{
		Code:      0,
		Msg:       "",
		DetailMsg: "",
		Data:      data,
	}
}

type ParamPage struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Total   int `json:"total"`
}

type ListDataRes struct {
	Data      interface{} `json:"data"`
	ParamPage ParamPage   `json:"param_page"`
}

type ListResponse struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	DetailMsg string      `json:"detail_msg"`
	Data      ListDataRes `json:"data"`
}

func GetListResponse(total, page, perPage int, data interface{}) *ListResponse {
	return &ListResponse{
		Code:      0,
		Msg:       "",
		DetailMsg: "",
		Data: ListDataRes{
			Data:      data,
			ParamPage: ParamPage{page, perPage, total},
		},
	}
}

func GetEmptyListResponse(total, page, perPage int) *ListResponse {
	return &ListResponse{
		Code:      0,
		Msg:       "",
		DetailMsg: "",
		Data: ListDataRes{
			Data:      []string{},
			ParamPage: ParamPage{page, perPage, total},
		},
	}
}
