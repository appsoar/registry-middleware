package handler

import (
	"encoding/json"
	//"scheduler/client"
	"fmt"
	"net/http"
	"scheduler/client/database"
	"scheduler/errjson"
	"scheduler/log"
)

func getAccounts(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
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

func addAccounts(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		return
	}

	log.Logger.Info(user + "add new user")
	decoder := json.NewDecoder(r.Body)
	var ui database.UserInfo
	err = decoder.Decode(&ui)
	if err != nil {
		panic(err)
	}
	//检测数据合法性?
	if len(ui.Id) == 0 || len(ui.Password) == 0 {
		err = errjson.NewNotValidEntityError("invalid account")
		return
	}
	//检测是否已经注册
	//待完成
	/*对ui的数据进行处理*/
	_, err = globalClient.AddUserAccount(ui)
	if err != nil {
		err = errjson.NewInternalServerError(err.Error())
		return
	}
	return
}
