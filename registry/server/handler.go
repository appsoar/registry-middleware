package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"registry/client"
)

var (
	opts = client.ClientOpts{
		Url: "http://192.168.2.110:5000",
	}
)

type RespError struct {
	Type    string `json:"type"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	resp := RespError{
		Type:    "error",
		Status:  404,
		Code:    "404 not found",
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
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprintf(w, id)
}

func ListImage(w http.ResponseWriter, r *http.Request) {
	resp := RespError{
		Type:    "error",
		Status:  500,
		Code:    "500 Internal Server Error",
		Message: "server cannot get image lists",
	}
	opts := client.ClientOpts{
		Url: "http://192.168.2.110:5000",
	}

	content, err := client.ListRepositoriesPagination(opts, 0)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(content))
}

func ListImageTags(w http.ResponseWriter, r *http.Request) {
	resp := RespError{
		Type:    "error",
		Status:  500,
		Code:    "500 Internal Server Error",
		Message: "server cannot get tag list",
	}
	vars := mux.Vars(r)
	image := vars["image"]
	fmt.Println(image)
	content, err := client.ListImageTags(opts, image)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(content))
}

func init() {
}
