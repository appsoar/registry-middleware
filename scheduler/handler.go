package scheduler

import (
	"encoding/json"
	"fmt"
	//	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	//	"regexp"
	//	"scheduler/errors"
	//"scheduler/session/session"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"scheduler/client"
	"scheduler/session"
	_ "scheduler/session/provider"
	"time"
)

const (
	loginPage = "login.gtpl"
)

var (
	globalSessions *session.Manager
	globalClient   *client.Client
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
		resp := NotValidEntityError("name cannot be empty")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		//返回状态码422,未在net/http中实现,使用自定义的422
		w.WriteHeader(ErrorNotValidEntity)
		/*将结构体转换成json*/
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

}

func GetImageProperty(w http.ResponseWriter, r *http.Request) {
	//	fmt.Printf(w, "")
}

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {

	var info LoginInfo
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))

	//转换成byte
	err = json.Unmarshal(body, &info)
	if err != nil {
		panic(err)
	}
	if len(info.Username) == 0 || len(info.Password) == 0 {
		resp := NotValidEntityError("invalid username or password ")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		//返回状态码422,未在net/http中实现,使用自定义的422
		w.WriteHeader(ErrorNotValidEntity)
		/*将结构体转换成json*/
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}
	/*未实现,密码验证*/
	sess := globalSessions.SessionStart(w, r)
	sess.Set("username", info.Username)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "")
}

//router.Handler(websocket.Handler(GetSystemInfo)
//这里返回有所有的系统相关数据？
func GetSysInfo(ws *websocket.Conn) {
	//	var err error
	//cpuchan := make(chan SystemInfo)
	//memchan := make(chan .....)
	//diskchan := make(chan....)
	for {
		sysinfo, err := globalClient.GetSysInfo()
		if err != nil {
			panic(err)
		}

		b, err := json.Marshal(sysinfo)
		if err != nil {
			panic(err)
		}

		if err = websocket.Message.Send(ws, string(b)); err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)

	}

}

func ShowLogin(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	t, _ := template.ParseFiles(loginPage)
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, sess.Get("username"))

}

func init() {
	var err error
	//所有的客户端请求通过globalClient完成
	globalClient, err = client.NewClient()
	if err != nil {
		panic(err)
	}
	//创建一个全局的session管理器,session存储方式为内存,cookie名为gosessionid
	globalSessions, err = session.NewManager("memory", "gosessionid", 3600)
	if err != nil {
		panic(err)
	}
	go globalSessions.GC()
}
