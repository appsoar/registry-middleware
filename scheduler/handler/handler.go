package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"scheduler/client"
	"scheduler/errjson"
	"scheduler/log"
	"scheduler/session"
	_ "scheduler/session/provider"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

var (
	globalSessions   *session.Manager
	globalClient     *client.Client
	globalLoginedMap *LoginedUser
)

//这里并没有处理,用户来自CLI的情况。
//CLI没有cookie,用户应当是通过Access/Secret的方式进行的
//目前想到的方法是getRequestUser失败后,
//从request请求中获取授权信息(即access/secret)
//进行权限判断之类的e.
func getRequestUser(w http.ResponseWriter, r *http.Request) (string, error) {
	sess := globalSessions.SessionStart(w, r)
	strI := sess.Get("username")
	if strI == nil {
		//处理从CLI发送的请求
		log.Logger.Warn(r.RemoteAddr + " :invalid request")
		globalSessions.SessionDestroy(w, r)
		e := errjson.NewUnauthorizedError(r.RemoteAddr + " :user doesn't login")
		return "", e
	}

	username, ok := strI.(string)
	if !ok {
		errStr := "session username key/value pair are not string"
		log.Logger.Error(errStr)
		panic(errStr)
	}

	return username, nil
}

/*请求错误退出时,采用Json格式进行返回错误信息*/
func errJsonReturn(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	switch e := err.(type) {
	//为什么以下类型断言错误类型不能放在同一个case中,会报e.Resp错误
	//类型断言时如果多个类型放在switch case,go语言不知道使用哪个
	//因此go会使用原来的类型(这里是error)
	case errjson.UnauthorizedError:
		w.WriteHeader(e.Status)
		if err := json.NewEncoder(w).Encode(e); err != nil {
			panic(err)
		}
	case errjson.NotFoundError:
		w.WriteHeader(e.Status)
		if err := json.NewEncoder(w).Encode(e); err != nil {
			panic(err)
		}
	case errjson.NotValidEntityError:
		w.WriteHeader(e.Status)
		if err := json.NewEncoder(w).Encode(e); err != nil {
			panic(err)
		}

	case errjson.InternalServerError:
		w.WriteHeader(e.Status)
		if err := json.NewEncoder(w).Encode(e); err != nil {
			panic(err)
		}

	case errjson.ErrForbidden:
		w.WriteHeader(e.Status)
		if err := json.NewEncoder(w).Encode(e); err != nil {
			panic(err)
		}
	default:
		panic("not json return error")
	}
}

func withJsonReturn(fn func(http.ResponseWriter, *http.Request) error, w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		errJsonReturn(w, r, err)
		return
	}
	jsonReturn(w, r)
	return
}

//和Fprintf(w,"")冲突,会触发multiple header write错误
func jsonReturn(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	//	w.WriteHeader(http.StatusOK)
	return
}

//返回系统信息
func GetSysInfo(ws *websocket.Conn) {
	defer func() {
		ws.Close()
	}()
	for {
		sysinfo, err := globalClient.GetSysInfo()
		if err != nil {
			panic(err)
		}

		b, err := json.Marshal(sysinfo)
		if err != nil {
			//panic(err)
			log.Logger.Error(err.Error())
			time.Sleep(2 * time.Second)
			continue
		}

		if err = websocket.Message.Send(ws, string(b)); err != nil {
			//panic(err)
			//			log.Logger.Error(err.Error())
		}
		time.Sleep(2 * time.Second)

	}
}

func GetUserStats(ws *websocket.Conn) {
	defer func() {
		ws.Close()
	}()
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
		us, err := globalClient.GetUserStats()
		if err != nil {
			//panic(err)
			//	log.Logger.Error(err.Error())
			continue
		}
		//		b, err := json.Marshal(us)
		if err := websocket.Message.Send(ws, string(us)); err != nil {
			//	panic(err)
			//			log.Logger.Error(err.Error())
		}
		time.Sleep(10 * time.Second)
	}

}

type LogLine struct {
	Lines LogStruct `json:"lines"`
}

type LogStruct struct {
	Time   string `json:"time"`
	User   string `json:"user"`
	Level  string `json:"level"`
	Detail string `json:"detail"`
}

func getLog(w http.ResponseWriter, r *http.Request) (err error) {
	/*
		content, err := ioutil.ReadFile("/home/kiongf/registry-middleware/src/scheduler/handler/logs.json")
		if err != nil {
			err = errjson.NewInternalServerError("read logs.json fail:" + err.Error())
			return
		}
	*/

	vars := mux.Vars(r)
	line_offset, ok := vars["line_offset"]
	if !ok {
		panic("invalid")
	}
	content, err := globalClient.GetLog(line_offset)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, string(content))
	return

}

func GetLog(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(getLog, w, r)
}

func getIfs(w http.ResponseWriter, r *http.Request) (err error) {
	ifs, err := globalClient.GetNetIfs()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, string(ifs))
	return
}

func GetIfs(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(getIfs, w, r)
}

func getIfStat(w http.ResponseWriter, r *http.Request) (err error) {
	vars := mux.Vars(r)
	netif, ok := vars["netif"]
	if !ok {
		err := errjson.NewErrForbidden("netif must set")
		return err
	}
	ifs, err := globalClient.GetNetIfStat(netif)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, string(ifs))
	return
}

func GetIfStat(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(getIfStat, w, r)
}

/*为了解决请求返回信息冗余的问题,合并http请求控制器最后路径为errJsonReturn或jsonReturn
参考go web编程错误处理章节.由于集成了gorilla/mux,限制控制器类型为http.HandlerFunc
因此使用此方法.下一步,尝试更改gorilla/mux代码?*/
/*登录控制器*/
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(login, w, r)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(logout, w, r)
}

//无效url请求
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	e := errjson.NewNotFoundError("specified page not found")
	errJsonReturn(w, r, e)
}

func GetAllNsHandler(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(GetAllNs, w, r)
}

func GetSpecNsHandler(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(getSpecNs, w, r)
}

func UpdateNs(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(updateNs, w, r)
}

func DeleteNs(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(deleteNs, w, r)
}

func GetReposHandler(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(getRepos, w, r)
}
func GetNsReposHandler(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(getNsRepos, w, r)
}

func GetUserReposHandler(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(getUserRepos, w, r)
}

func ListRepoTagsHandler(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(listRepoTags, w, r)
}

func GetTagImageHandler(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(getTagImage, w, r)
}

func GetAccounts(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(getAccounts, w, r)
}

func GetUserAccount(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(getUserAccount, w, r)
}

func AddAccount(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(addAccount, w, r)
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(updateAccount, w, r)
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(deleteAccount, w, r)
}

func GetNsUgroup(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(getNsUgroup, w, r)
}

func AddUgroup(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(addUgroup, w, r)
}

func DeleteUgroup(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(deleteUgroup, w, r)
}

func UpdateUgroup(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(updateUgroup, w, r)
}

func GetUgroup(w http.ResponseWriter, r *http.Request) {
	withJsonReturn(updateUgroup, w, r)
}

/*记录已登录用户*/
type LoginedUser struct {
	user map[string]int
	m    *sync.RWMutex
}

func init() {
	var err error
	globalLoginedMap = &LoginedUser{
		user: make(map[string]int),
		m:    new(sync.RWMutex),
	}

	//所有的客户端请求通过globalClient完成
	globalClient, err = client.NewClient()
	if err != nil {
		log.Logger.Debug("fail to create scheduler client:%s", err.Error())
		panic(err)
	}

	//创建一个全局的session管理器,session存储方式为内存,cookie名为gosessionid
	globalSessions, err = session.NewManager("memory", "gosessionid", 3600000)
	if err != nil {
		log.Logger.Debug("fail to create session manager:%s", err.Error())
		panic(err)
	}
	go globalSessions.GC()
}
