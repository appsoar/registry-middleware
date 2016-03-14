package session

import (
	"crypto/rand"
	"encoding/base64"
	//	"error"
	"fmt"
	"io"
	"net/http"
	"net/url"
	//_ "scheduler/session/provider"
	"sync"
	"time"
)

var (
	//globalSessions *session.Manager
	provides = make(map[string]Provider)
)

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{} //获取session存储的键值对.比如当前实现中"username"为sessionID绑定的用户名
	Delete(key interface{}) error
	SessionID() string
}

type Manager struct {
	cookieName  string     //private cookiename
	lock        sync.Mutex //protects session
	provider    Provider
	maxlifetime int64
}

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTIme int64)
}

//获取全局唯一的session id
func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

//检测是否已经有某个session和当前的来访用户发生了关联,如果没有则创建,
//有则返回.*通过session.value可以获取session相关的参数.比如说sessionid对应的用户名
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	//获取requset中session相关的cookie信息
	cookie, err := r.Cookie(manager.cookieName)
	//cookie值为空,没有关联session
	if err != nil || cookie.Value == "" {
		sid := manager.sessionId()
		session, _ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{
			Name:     manager.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(manager.maxlifetime),
		}
		//增加cookie到回应头
		http.SetCookie(w, &cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}
	return
}

func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	//获取session相关的cookie信息
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.provider.SessionDestroy(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{
			Name:     manager.cookieName,
			Path:     "/",
			HttpOnly: true,
			Expires:  expiration,
			MaxAge:   -1,
		}

		http.SetCookie(w, &cookie)
	}
}

func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxlifetime)
	time.AfterFunc(time.Duration(manager.maxlifetime), func() { manager.GC() })
}

//创建新的session管理器
func NewManager(providerName, cookieName string, maxlifetime int64) (*Manager, error) {
	provider, ok := provides[providerName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import ?)", providerName)
	}

	return &Manager{provider: provider, cookieName: cookieName, maxlifetime: maxlifetime}, nil
}

func Register(name string, provide Provider) {
	if provide == nil {
		panic("session: Register provide is nil")
	}

	if _, dup := provides[name]; dup {
		panic("session: Register called twice for provider" + name)
	}
	provides[name] = provide
}

/*
func init() {
	var err error
	globalSessions, err = NewManager("memory", "gosessionid", 3600)
	if err != nil {
		panic(err)
	}
	go globalSessions.GC()
}
*/
