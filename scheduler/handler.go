package scheduler

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"net/http"
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

type schedulerHandler func(http.ResponseWriter, *http.Request) error

func (fn schedulerHandler) ServeHttp(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		errJsonReturn(w, r, err)
		//转换成respError
		//类型断言处理
	}
	//否则返回正确信息
}

func getRequestUser(w http.ResponseWriter, r *http.Request) (string, error) {
	sess := globalSessions.SessionStart(w, r)
	strI := sess.Get("username")
	if strI == nil {
		log.Logger.Warn("invalid logout request")
		globalSessions.SessionDestroy(w, r)
		e := NewUnauthorizedError("user doesn't login")
		errJsonReturn(w, r, e)
	}

	username, ok := strI.(string)
	if !ok {
		errStr := "session username key/value pair are not string"
		log.Logger.Error(errStr)
		panic(errStr)
	}

	return username, nil
}

//错误状态JSON返回
func errJsonReturn(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	switch e := err.(type) {
	//为什么以下类型断言错误类型不能放在同一个case中,会报e.resp错误
	//类型断言时如果多个类型放在switch case,go语言不知道使用哪个
	//因此go会使用原来的类型(这里是error)
	case UnauthorizedError:

		w.WriteHeader(e.resp.Status)
		if err := json.NewEncoder(w).Encode(e.resp); err != nil {
			panic(err)
		}
	case NotFoundError:
		w.WriteHeader(e.resp.Status)
		if err := json.NewEncoder(w).Encode(e.resp); err != nil {
			panic(err)
		}
	case NotValidEntityError:
		w.WriteHeader(e.resp.Status)
		if err := json.NewEncoder(w).Encode(e.resp); err != nil {
			panic(err)
		}

	case InternalServerError:
		w.WriteHeader(e.resp.Status)
		if err := json.NewEncoder(w).Encode(e.resp); err != nil {
			panic(err)
		}
	default:
		panic("not json return error")
	}
}

//无效url请求
func NotFound(w http.ResponseWriter, r *http.Request) {
	e := NewNotFoundError("specified page not found")
	errJsonReturn(w, r, e)
}

func SetImagesProperty(w http.ResponseWriter, r *http.Request) {

	user, err := getRequestUser(w, r)
	if err != nil {
		e := NewUnauthorizedError("user doesn't login")
		errJsonReturn(w, r, e)
	}
	log.Debug(user + "set image")

	err = r.ParseForm()
	if err != nil {
		//这里应该返回一个500,主机出错
		e := NewInternalServerError("request form parse fail:" + err.Error())
		errJsonReturn(w, r, e)
	}
	//还要检测镜像是否存在
	//需要加锁
	if len(r.Form["name"][0]) == 0 {
		e := NewNotValidEntityError("name cannot be empty")
		errJsonReturn(w, r, e)
	}

}

func GetImageProperty(w http.ResponseWriter, r *http.Request) {

	user, err := getRequestUser(w, r)
	if err != nil {
		e := NewUnauthorizedError("user doesn't login")
		errJsonReturn(w, r, e)
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
		e := NewUnauthorizedError("invalid username or password ")
		errJsonReturn(w, r, e)
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
		e := NewUnauthorizedError("invalid username or password ")
		errJsonReturn(w, r, e)
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

//返回系统信息
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

func GetUserStatus(ws *websocket.Conn) {
	for {
		/*使用channel来控制更新频率*/
		/*新增或者移除用户时,发送channel,触发websocket写入*/
		/*新增或移除镜像时,发送channel,触发websocket写入*/
		/*新增或者移除namespace时,发送channel,触发websocket写入*/
		/*只有以上动作触发时才会触发更新前端数据*/
		/*
			go func() {
				select {
				case <-UserChannel: //获取image数据,更新
				case <-ImageChannel: client.GetUserStatus()
				case <-NamespaceChannel: client.GetUserStaus
				}
			}()
			userStatus := <-UserStatusChannel
			b, err := json.Marshal(userStatus)
			if err != nil {
				panic(err)
			}

			if err = websocket.Message.Send(ws, string(b)); err != nil {
				panic(err)
			}*/
		//暂时不用channel
		//使用轮询,10s一次查询
	}

}

type LogStruct struct {
	Log string `json:"log"`
}

//需要修改
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

func DeleteImage(w http.ResponseWriter, r *http.Request) {
	user, err := getRequestUser(w, r)
	if err != nil {
		e := NewUnauthorizedError("user doesn't login")
		errJsonReturn(w, r, e)
	}
	log.Debug(user + "delete images")

	/*提取镜像名,Tag*/
	/*获取锁,删除*/
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
