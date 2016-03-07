package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
	"scheduler/errors"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	resp := NotFoundError("The specified page not found")

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	/*将结构体转换成json*/
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}

}

func SetImagesProperty(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		//这里应该返回一个500,主机出错
		panic(err)
	}
	//返回422无法处理的对象
	//还要检测镜像是否存在
	//需要加锁
	if len(r.Form["name"][0]) == 0 {
		resp := NoValidEntityError("name cannot be empty")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		//返回状态码422,未在net/http中实现,使用自定义的422
		w.WriteHeader(ErrorNotValidEntity)
		/*将结构体转换成json*/
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

}

func Login(w http.ResponseWriter, r *http.Request) {
	err := r.PostForm()

	if err != nil {
		panic(err)
	}

	if len(r.Form["username"][0]) == 0 || len(r.Form["password"]) == 0 {
		resp := NoValidEntityError("name cannot be empty")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		//返回状态码422,未在net/http中实现,使用自定义的422
		w.WriteHeader(ErrorNotValidEntity)
		/*将结构体转换成json*/
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}
	//连接数据库进行验证
	//验证通过
	//先不用cookie
	//	cookie := http.Cookie{Name:""}

}
