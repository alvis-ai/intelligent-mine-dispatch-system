package response

import "net/http"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

func Success(data interface{}) Response {
	return Response{Code: 0, Message: "success", Data: data}
}

func Fail(code int, msg string) Response {
	return Response{Code: code, Message: msg}
}

func Page(list interface{}, total int64, page, pageSize int) Response {
	return Success(PageResult{List: list, Total: total, Page: page, PageSize: pageSize})
}

var MsgMap = map[int]string{
	http.StatusOK:                  "success",
	http.StatusBadRequest:          "invalid request",
	http.StatusUnauthorized:        "unauthorized",
	http.StatusForbidden:           "forbidden",
	http.StatusNotFound:            "not found",
	http.StatusInternalServerError: "internal error",
}
