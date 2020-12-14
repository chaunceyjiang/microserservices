package handlers

import (
	"encoding/json"
	"io"
)

// swagger:response Resp
// Resp 返回统一封装
type Resp struct {
	// Required: true
	// Code 状态码
	Code int `json:"code"`
	// Required: true
	// Message 状态码
	Message string `json:"message"`
	// Required: true
	// Data 状态码
	Data interface{} `json:"data,omitempty"`
}

func (r *Resp) ToJSON(w io.Writer) error {
	enc:=json.NewEncoder(w)
	return enc.Encode(r)
}
