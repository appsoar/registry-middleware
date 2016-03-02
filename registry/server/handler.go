package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RespError struct {
	Type    string `json:"type"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	resp := RespError{
		Type:    "error",
		Status:  404,
		Code:    "404 notfound",
		Message: "The specified page not found",
	}

	/*
		b, err := json.Marshal(resp)
		if err != nil {
			fmt.Println("json err:", err)
			return
		}
	*/

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	/*将结构体转换成json*/
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}

}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func ListImage(w http.ResponseWriter, r *http.Request) {

}
