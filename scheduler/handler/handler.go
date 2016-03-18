package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"net/http"
	"scheduler/client"
	"scheduler/errjson"
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
		e := errjson.NewUnauthorizedError("user doesn't login")
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
		w.WriteHeader(e.Resp.Status)
		if err := json.NewEncoder(w).Encode(e.Resp); err != nil {
			panic(err)
		}
	case errjson.NotFoundError:
		w.WriteHeader(e.Resp.Status)
		if err := json.NewEncoder(w).Encode(e.Resp); err != nil {
			panic(err)
		}
	case errjson.NotValidEntityError:
		w.WriteHeader(e.Resp.Status)
		if err := json.NewEncoder(w).Encode(e.Resp); err != nil {
			panic(err)
		}

	case errjson.InternalServerError:
		w.WriteHeader(e.Resp.Status)
		if err := json.NewEncoder(w).Encode(e.Resp); err != nil {
			panic(err)
		}
	default:
		panic("not json return error")
	}
}

func jsonReturn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return
}

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

/*UI传递过来的密码是先经过Md5加密过的,服务端再进行bcrypt加密.*/
/*这里存在一个问题:当不通过UI,而通过客户端连接,用户输入的是不经过、
  mk5加密过的明文密码.会无法通过密码验证.
*/

func login(w http.ResponseWriter, r *http.Request) (err error) {

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
		err = errjson.NewUnauthorizedError("invalid username or password ")
		//	errJsonReturn(w, r, e)
		return
	}
	/*
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(info.password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}*/

	ui, err := globalClient.GetUserAccount(info.Username)
	if err != nil {
		log.Logger.Error(err.Error())
		err = errjson.NewInternalServerError("server can't get userinfo")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(ui.Password), []byte(info.Password))
	if err != nil {
		err = errjson.NewUnauthorizedError("incorrect password")
		return
	}

	sess := globalSessions.SessionStart(w, r)
	sess.Set("username", info.Username)
	log.Logger.Info("%s Login", info.Username)
	return
}

func logout(w http.ResponseWriter, r *http.Request) (err error) {
	//*如果是一个未登录用户调用了该方法
	//SessionManager会根据session中有无设置username来确定
	//是以登录用户，还是未登录用户
	sess := globalSessions.SessionStart(w, r)
	strI := sess.Get("username")
	if strI == nil {
		log.Logger.Warn("invalid logout request")
		globalSessions.SessionDestroy(w, r)
		err = errjson.NewUnauthorizedError("invalid username or password ")
		//errJsonReturn(w, r, e)
		return
	}
	fmt.Println(strI)

	username, ok := strI.(string)
	if !ok {
		errStr := "session username key/value pair are not string"
		log.Logger.Error(errStr)
		panic(errStr)
	}
	log.Logger.Debug("%s logout", username) //打印当前登录用户的用户名

	globalSessions.SessionDestroy(w, r)
	return
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

func GetUserStats(ws *websocket.Conn) {
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
			panic(err)
		}
		b, err := json.Marshal(us)
		if err := websocket.Message.Send(ws, string(b)); err != nil {
			panic(err)
		}
		time.Sleep(10 * time.Second)
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
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}
	log.Debug(user + "delete images")

	/*提取镜像名,Tag*/
	/*获取锁,删除*/
}
func namespacesGet(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	log.Logger.Info(user + "get namespace info")
	nsJson, err := globalClient.GetNamespaces()
	if err != nil {
		err = errjson.NewInternalServerError("can't get ns info")
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func namespaceGetSpecific(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	vars := mux.Vars(r)
	ns := vars["namespace"]
	if len(ns) == 0 {
		err = errjson.NewNotValidEntityError("invalid namespace")
		return
	}

	log.Logger.Info(user + " get " + ns + " namespace info")

	nsJson, err := globalClient.GetSpecificNamespace(ns)
	_, err = globalClient.GetSpecificNamespace(ns)
	if err != nil {
		err = errjson.NewInternalServerError("can't get ns info")
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func getRepos(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	log.Logger.Info(user + " get repositories")

	nsJson, err := globalClient.GetRepositories()
	if err != nil {
		err = errjson.NewInternalServerError("can't get ns info")
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func getNsRepos(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	log.Logger.Info(user + " get repositories")

	vars := mux.Vars(r)
	ns := vars["namespace"]

	if len(ns) == 0 {
		err = errjson.NewNotValidEntityError("invalid namespace")
		return
	}

	nsJson, err := globalClient.GetNsRepos(ns)
	if err != nil {
		err = errjson.NewInternalServerError("can't get ns info")
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func getUserRepos(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	log.Logger.Info(user + " get repositories")

	vars := mux.Vars(r)
	ns := vars["user"]

	if len(ns) == 0 {
		err = errjson.NewNotValidEntityError("invalid namespace")
		return
	}

	nsJson, err := globalClient.GetUserRepos(ns)
	if err != nil {
		err = errjson.NewInternalServerError("can't get ns info")
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func listRepoTags(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	log.Logger.Info(user + " get repositories")

	vars := mux.Vars(r)
	name := vars["usernameOrNamespace"]
	repoName := vars["repoName"]

	if len(repoName) == 0 {
		err = errjson.NewNotValidEntityError("invalid namespace")
		return
	}

	nsJson, err := globalClient.ListRepoTags(name, repoName)
	if err != nil {
		err = errjson.NewInternalServerError("can't get repo info")
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func getTagImage(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	log.Logger.Info(user + " get repositories")

	vars := mux.Vars(r)
	name := vars["usernameOrNamespace"]
	repoName := vars["repoName"]
	tagName := vars["tagName"]

	if len(repoName) == 0 || len(tagName) == 0 {
		err = errjson.NewNotValidEntityError("invalid namespace")
		return
	}

	nsJson, err := globalClient.GetTagImage(name, repoName, tagName)
	if err != nil {
		err = errjson.NewInternalServerError("can't get repo info")
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func test(w http.ResponseWriter, r *http.Request) error {
	ui, err := globalClient.GetUserAccount("admin")
	if err != nil {
		log.Logger.Error(err.Error())
		return err
	}

	fmt.Println(ui)
	return nil

}

/*为了解决请求返回信息冗余的问题,合并http请求控制器最后路径为errJsonReturn或jsonReturn
参考go web编程错误处理章节.由于集成了gorilla/mux,限制控制器类型为http.HandlerFunc
因此使用此方法.下一步,尝试更改gorilla/mux代码?*/
/*登录控制器*/
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := login(w, r); err != nil {
		errJsonReturn(w, r, err)
		return
	}
	jsonReturn(w, r)
	return
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if err := logout(w, r); err != nil {
		errJsonReturn(w, r, err)
		return
	}
	jsonReturn(w, r)
	return
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	if err := test(w, r); err != nil {
		errJsonReturn(w, r, err)
		return
	}
	jsonReturn(w, r)
	return
}

//无效url请求
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	e := errjson.NewNotFoundError("specified page not found")
	errJsonReturn(w, r, e)
}

func NamespacesGetHandler(w http.ResponseWriter, r *http.Request) {
	if err := namespacesGet(w, r); err != nil {
		errJsonReturn(w, r, err)
		return
	}
	jsonReturn(w, r)
	return
}

func NamespaceSpecificGetHandler(w http.ResponseWriter, r *http.Request) {
	if err := namespaceGetSpecific(w, r); err != nil {
		errJsonReturn(w, r, err)
		return
	}
	jsonReturn(w, r)
	return
}

func GetReposHandler(w http.ResponseWriter, r *http.Request) {
	if err := getRepos(w, r); err != nil {
		errJsonReturn(w, r, err)
		return
	}
	jsonReturn(w, r)
	return
}
func GetNsReposHandler(w http.ResponseWriter, r *http.Request) {
	if err := getNsRepos(w, r); err != nil {
		errJsonReturn(w, r, err)
		return
	}
	jsonReturn(w, r)
	return
}

func GetUserReposHandler(w http.ResponseWriter, r *http.Request) {
	if err := getUserRepos(w, r); err != nil {
		errJsonReturn(w, r, err)
		return
	}
	jsonReturn(w, r)
	return
}

func ListRepoTagsHandler(w http.ResponseWriter, r *http.Request) {
	if err := listRepoTags(w, r); err != nil {
		errJsonReturn(w, r, err)
		return
	}
	jsonReturn(w, r)
	return
}

func GetTagImageHandler(w http.ResponseWriter, r *http.Request) {
	if err := getTagImage(w, r); err != nil {
		errJsonReturn(w, r, err)
		return
	}
	jsonReturn(w, r)
	return
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
