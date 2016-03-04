package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"registry/client"
	"registry/debug"
)

var (
	opts client.ClientOpts
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

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	/*将结构体转换成json*/
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}

}

func ListImages(w http.ResponseWriter, r *http.Request) {
	resp := RespError{
		Type:    "error",
		Status:  500,
		Code:    "500 Internal Server Error",
		Message: "server cannot get image lists",
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

func DeleteImageTag(w http.ResponseWriter, r *http.Request) {
	resp := RespError{
		Type:    "error",
		Status:  500,
		Code:    "500 Internal Server Error",
		Message: "server cannot get tag list",
	}

	vars := mux.Vars(r)
	image := vars["image"]
	tag := vars["tag"]
	debug.Print(image, tag)

	err := client.DeleteImage(opts, image, tag)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		resp.Message = err.Error()
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "")

}

func ListImageDigest(w http.ResponseWriter, r *http.Request) {
	resp := RespError{
		Type:    "error",
		Status:  500,
		Code:    "500 Internal Server Error",
		Message: "server cannot get tag list",
	}

	vars := mux.Vars(r)
	image := vars["image"]
	tag := vars["tag"]
	debug.Print(image, tag)

	digest, err := client.GetImageDigest(opts, image, tag)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		resp.Message = err.Error()
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}
	//	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, digest)

}

func init() {

	opts.Url = os.Getenv("REGISTRY_URL")
	if len(opts.Url) == 0 {
		panic("missing REGISTRY_URL")
	}
	/*
		opts.AccessKey = os.Getenv("REGISTRY_ACCESS_KEY")
		if len(opts.Url) == 0 {
			panic("missing AccessKey")
		}
		opts.AccessKey = os.Getenv("REGISTRY_SECRET_KEY")
		if len(opts.Url) == 0 {
			panic("missing SecretKey")
		}*/
}
