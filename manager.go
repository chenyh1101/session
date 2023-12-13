package session

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Manager struct {
	mutex      sync.Mutex
	provider   IProvider //session 存储
	cookieName string
	maxAge     int64
}

func NewManager(providerName, cookieName string, maxAge int64) *Manager {
	provider, ok := providers[providerName]
	if !ok {
		return nil
	}
	return &Manager{cookieName: cookieName, provider: provider, maxAge: maxAge}
}

// SessionID 生成全局唯一的sessionID 用于识别每一个用户
func (m *Manager) SessionID() string {
	buf := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(buf)
}

// Start 根据当前请求中的COOKIE判断是否存在有效的Session，不存在则创建
func (m *Manager) Start(w http.ResponseWriter, r *http.Request) ISession {
	//添加互斥锁
	m.mutex.Lock()
	defer m.mutex.Unlock()
	//获取cookie
	cookie, err := r.Cookie(m.cookieName)
	log.Printf("%v", cookie)
	if err != nil || cookie.Value == "" {
		//创建sessionID
		sid := m.SessionID()
		//Session初始化
		session, _ := m.provider.Init(sid)
		//设置Cookie到Response
		http.SetCookie(w, &http.Cookie{
			Name:     m.cookieName,
			Value:    url.QueryEscape(cookie.Value),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(m.maxAge),
		})
		return session
	} else {
		//从Cookie 获取Session
		sid, _ := url.QueryUnescape(cookie.Value)
		//获取Session
		Session, _ := m.provider.Read(sid)
		return Session
	}

}

// Destroy 注销Session
func (m *Manager) Destroy(w http.ResponseWriter, r *http.Request) {
	//从请求中获取Cookie值
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}
	//添加互斥锁
	m.mutex.Lock()
	defer m.mutex.Unlock()
	//销毁Session内存
	_ = m.provider.Destroy(cookie.Value)
	//设置客户端Cookie立即过期
	http.SetCookie(w, &http.Cookie{
		Name:     m.cookieName,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Now(),
	})
}

// GC 销毁Session
func (m *Manager) GC() {
	//添加互斥锁
	m.mutex.Lock()
	defer m.mutex.Unlock()
	//设置过期时间销毁Session
	m.provider.GC(m.maxAge)
	//添加计时器当Session超时自动销毁
	time.AfterFunc(time.Duration(m.maxAge), func() {
		m.GC()
	})
}
