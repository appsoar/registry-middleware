package scheduler

import (
	"encoding/json"
	"fmt"
	//	"github.com/gorilla/mux"
	//	"html/template"
	"net/http"
	//	"regexp"
	//	"scheduler/errors"
	//"scheduler/session/session"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"scheduler/client"
	"scheduler/log"
	"scheduler/session"
	_ "scheduler/session/provider"
	"time"
)

var (
	globalSessions *session.Manager
	globalClient   *client.Client
)

func getRequestUser(w http.ResponseWriter, r *http.Request) (string, error) {
	sess := globalSessions.SessionStart(w, r)
	strI := sess.Get("username")
	if strI == nil {
		log.Logger.Warn("invalid logout request")
		globalSessions.SessionDestroy(w, r)
		return "", &UnloginUserError{}
	}

	username, ok := strI.(string)
	if !ok {
		errStr := "session username key/value pair are not string"
		log.Logger.Error(errStr)
		panic(errStr)
	}

	return username, nil
}

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
	//
	user, err := getRequestUser(w, r)
	if err != nil {
		if _, ok := err.(UnloginUserError); ok {
			//返回错误,error json
			return
		}
	}

	r.ParseForm()
	if len(r.Form["image"]) == 0 {
		errStr := "image name is empty"
		log.Logger.Error(errStr)
		return
	}
	/*从数据库获取数据*/
	log.Logger.Info("%s get image[%s] ", user, r.Form["image"])
	/**/

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
	log.Logger.Debug(string(body))

	//转换成byte
	err = json.Unmarshal(body, &info)
	if err != nil {
		panic(err)
	}
	if len(info.Username) == 0 || len(info.Password) == 0 {
		log.Logger.Info("invalid username or password")
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
	log.Logger.Info("%s Login", info.Username)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	//*如果是一个未登录用户调用了该方法
	//SessionManager会根据session中有无设置username来确定
	//是以登录用户，还是未登录用户
	sess := globalSessions.SessionStart(w, r)
	strI := sess.Get("username")
	if strI == nil {
		log.Logger.Warn("invalid logout request")
		globalSessions.SessionDestroy(w, r)
		return
	}

	username, ok := strI.(string)
	if !ok {
		errStr := "session username key/value pair are not string"
		log.Logger.Error(errStr)
		panic(errStr)
	}
	log.Logger.Debug("%s logout", username) //打印当前登录用户的用户名

	globalSessions.SessionDestroy(w, r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "")
}

//router.Handler(websocket.Handler(GetSystemInfo)
//这里返回有所有的系统相关数据？
func GetSysInfo(ws *websocket.Conn) {
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

type LogStruct struct {
	Log string `json:"log"`
}

func GetLog(ws *websocket.Conn) {
	for {
		infos := <-log.LogChannel
		logStruct := LogStruct{Log: infos}
		log.Logger.Info(infos)

		b, err := json.Marshal(logStruct)
		if err != nil {
			panic(err)
		}
		if err := websocket.Message.Send(ws, string(b)); err != nil {
			panic(err)
		}

	}
}

func Test(w http.ResponseWriter, r *http.Request) {
	/*	images, err := globalClient.ListImages()

		if err != nil {
			panic(err)
		}*/
	//cookie, _ := r.Cookie("gosessionid")
	sess := globalSessions.SessionStart(w, r)

	log.Logger.Debug(sess.Get("username").(string)) //打印当前登录用户的用户名

}

/*
func ShowLogin(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	t, _ := template.ParseFiles(loginPage)
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, sess.Get("username"))

}*/

func init() {
	var err error
	//所有的客户端请求通过globalClient完成
	globalClient, err = client.NewClient()
	if err != nil {
		log.Logger.Debug("fail to create scheduler client:%s", err.Error())
		panic(err)
	}
	//创建一个全局的session管理器,session存储方式为内存,cookie名为gosessionid
	globalSessions, err = session.NewManager("memory", "gosessionid", 3600)
	if err != nil {
		log.Logger.Debug("fail to create session manager:%s", err.Error())
		panic(err)
	}
	go globalSessions.GC()
}
