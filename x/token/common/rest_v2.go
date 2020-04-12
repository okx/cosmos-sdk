package common

import (
	"encoding/json"
	"net/http"
	"strconv"

	//jsoniter "github.com/json-iterator/go"
)

const (
	// common error
	ErrorMissingRequiredParam ErrorCodeV2 = 60001
	ErrorInvalidParam         ErrorCodeV2 = 60002
	ErrorServerException      ErrorCodeV2 = 60003
	ErrorDataNotExist         ErrorCodeV2 = 60004

	// account error
	ErrorInvalidAddress ErrorCodeV2 = 61001

	// order error
	ErrorOrderNotExist        ErrorCodeV2 = 62001
	ErrorInvalidCurrency      ErrorCodeV2 = 62002
	ErrorEmptyInstrumentId    ErrorCodeV2 = 62003
	ErrorInstrumentIdNotExist ErrorCodeV2 = 62004
)

func DefaultErrorMessageV2(code ErrorCodeV2) (message string) {
	switch code {
	case ErrorMissingRequiredParam:
		message = "missing required param"
	case ErrorInvalidParam:
		message = "invalid request param"
	case ErrorServerException:
		message = "internal server error"
	case ErrorDataNotExist:
		message = "data not exists"

	case ErrorInvalidAddress:
		message = "invalid address"

	case ErrorOrderNotExist:
		message = "order not exists"
	case ErrorInvalidCurrency:
		message = "invalid currency"
	case ErrorEmptyInstrumentId:
		message = "instrument_id is empty"
	case ErrorInstrumentIdNotExist:
		message = "instrument_id not exists"
	default:
		message = "unknown error"
	}
	return
}

type ErrorCodeV2 int

func (code ErrorCodeV2) Code() string {
	return strconv.Itoa(int(code))
}

func (code ErrorCodeV2) Message() string {
	return DefaultErrorMessageV2(code)
}

type ResponseErrorV2 struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func HandleErrorResponseV2(w http.ResponseWriter, statusCode int, errCode ErrorCodeV2) {
	response, _ := json.Marshal(ResponseErrorV2{
		Code:    errCode.Code(),
		Message: errCode.Message(),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(response)
}

func HandleSuccessResponseV2(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}
//
//func HandleResponseV2(w http.ResponseWriter, data []byte, err error) {
//	if err != nil {
//		HandleErrorResponseV2(w, http.StatusInternalServerError, ErrorServerException)
//		return
//	}
//	if len(data) == 0 {
//		HandleErrorResponseV2(w, http.StatusBadRequest, ErrorDataNotExist)
//	}
//
//	HandleSuccessResponseV2(w, data)
//}
//
//func JsonMarshalV2(v interface{}) ([]byte, error) {
//	var jsonV2 = jsoniter.Config{
//		EscapeHTML:             true,
//		SortMapKeys:            true,
//		ValidateJsonRawMessage: true,
//		TagKey:                 "v2",
//	}.Froze()
//
//	return jsonV2.MarshalIndent(v, "", "  ")
//}
//
//func JsonUnMarshalV2(data []byte, v interface{}) error {
//	var jsonV2 = jsoniter.Config{
//		EscapeHTML:             true,
//		SortMapKeys:            true,
//		ValidateJsonRawMessage: true,
//		TagKey:                 "v2",
//	}.Froze()
//
//	return jsonV2.Unmarshal(data, v)
//}
