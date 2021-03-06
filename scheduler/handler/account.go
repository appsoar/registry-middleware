package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"scheduler/client/database"
	"scheduler/errjson"
	"scheduler/log"

	"scheduler/Godeps/_workspace/src/github.com/gorilla/mux"
	"scheduler/Godeps/_workspace/src/golang.org/x/crypto/bcrypt"
)

func checkDbErr(err1 error) (err error) {
	/*
		if e, ok := err1.(database.EDatabase); ok {
			switch e.Code {
			case database.EPermission,
				database.ENoRecord,
				database.EMissingId,
				database.EInvalidFilter,
				database.EIncompleteUserInfo,
				database.EUserExists,
				database.EIncompleteGroupInfo,
				database.EGroupExists,
				database.EInvalidNsInfo,
				database.ENsExists:
				err = errjson.NewErrForbidden(e.Msg)
			case database.EDbException,
				database.ENotInterface:
				err = errjson.NewInternalServerError(e.Msg)
			default:
				err = errjson.NewInternalServerError(err.Error())
			}
		}
	*/
	err = errjson.NewInternalServerError("Database Error:" + err1.Error())
	return

}

func getUserAccount(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError(r.RemoteAddr + " :user doesn't login")
		return
	}

	vars := mux.Vars(r)
	account := vars["account"]
	if len(account) == 0 {
		err = errjson.NewErrForbidden("invalid account argument")
		return
	}

	log.Logger.Info(r.RemoteAddr + " : " + user + " get user account info")
	nsJson, err := globalClient.GetUserAccount(account)
	if err != nil {
		err = errjson.NewInternalServerError(" can't get account")
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func getAccounts(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		return
	}

	log.Logger.Info(r.RemoteAddr + " : " + user + " get accounts")
	nsJson, err := globalClient.GetAccounts()
	if err != nil {
		err = checkDbErr(err)
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func updateAccount(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		return
	}

	log.Logger.Info(r.RemoteAddr + " : " + user + "update account")
	nsJson, err := globalClient.UpdateUserAccount()
	if err != nil {
		err = checkDbErr(err)
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func deleteAccount(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		return
	}

	log.Logger.Info(r.RemoteAddr + " : " + user + "delete account")

	vars := mux.Vars(r)
	uid, ok := vars["user_id"]
	if !ok {
		panic("user_id missing")
	}

	nsJson, err := globalClient.DeleteUserAccount(uid)
	if err != nil {
		err = checkDbErr(err)
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func addAccount(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		return
	}
	//检查权限

	log.Logger.Info(r.RemoteAddr + " : " + user + " add new user")
	decoder := json.NewDecoder(r.Body)
	var ui database.UserInfo
	err = decoder.Decode(&ui)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v:%v\n", "new acccount", ui)
	//检测数据合法性?
	if len(ui.Id) == 0 || len(ui.Password) == 0 {
		log.Logger.Error(r.RemoteAddr + ":" + "new account's Id or Password are empty")
		err = errjson.NewErrForbidden("invalid account")
		return
	}

	/*bcrypt加密*/
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(ui.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	ui.Password = string(hashedPassword)
	log.Logger.Debug("encrypted password:" + ui.Password)

	//检测是否已经注册
	//待完成
	/*对ui的数据进行处理*/
	_, err = globalClient.AddUserAccount(ui)
	if err != nil {
		err = checkDbErr(err)
		return
	}
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
/*如何解决验证用户是否已经登录? 创建一个内存对象链表?*/
/*无法通过session!*/
func login(w http.ResponseWriter, r *http.Request) (err error) {

	var info LoginInfo
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}
	log.Logger.Debug(r.RemoteAddr + " : " + string(body))

	//转换成byte
	err = json.Unmarshal(body, &info)
	if err != nil {
		panic(err)
	}
	if len(info.Username) == 0 || len(info.Password) == 0 {
		log.Logger.Info(r.RemoteAddr + ":invalid username or password")
		err = errjson.NewUnauthorizedError("invalid username or password ")
		return
	}
	/*
		globalLoginedMap.m.RLock()
		fmt.Println(globalLoginedMap.user)
		if _, ok := globalLoginedMap.user[info.Username]; ok {
			globalLoginedMap.m.RUnlock()
			log.Logger.Debug(info.Username + " have loggined")
			err = errjson.NewErrForbidden("User have loggin")

			return
		}
		globalLoginedMap.m.RUnlock()

	*/
	ui, err := globalClient.GetUserAccountDecoded(info.Username)
	if err != nil {
		log.Logger.Error(err.Error())
		err = checkDbErr(err)
		return
	}

	//比较加密后的数据
	err = bcrypt.CompareHashAndPassword([]byte(ui.Password), []byte(info.Password))
	if err != nil {
		err = errjson.NewUnauthorizedError("incorrect password")
		return
	}
	/*
		globalLoginedMap.m.Lock()
		globalLoginedMap.user[info.Username] = 1
		defer globalLoginedMap.m.Unlock()
	*/

	sess := globalSessions.SessionStart(w, r)
	sess.Set("username", info.Username)
	log.Logger.Info(r.RemoteAddr+" : "+"%s Login", info.Username)
	return
}

func logout(w http.ResponseWriter, r *http.Request) (err error) {
	/*检测是否已登录用户发送的请求*/
	sess := globalSessions.SessionStart(w, r)
	strI := sess.Get("username")
	if strI == nil {
		log.Logger.Warn(r.RemoteAddr + " :invalid logout request")
		globalSessions.SessionDestroy(w, r)
		err = errjson.NewUnauthorizedError("user don't login")
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
	log.Logger.Debug(r.RemoteAddr+":%s logout", username) //打印当前登录用户的用户名
	globalLoginedMap.m.Lock()
	delete(globalLoginedMap.user, username)
	defer globalLoginedMap.m.Unlock()

	globalSessions.SessionDestroy(w, r)
	return
}
