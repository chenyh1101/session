package memory

import (
	"container/list"
	session "http_session"
	"sync"
	"time"
)

type Memory struct {
	mutex sync.Mutex               //互斥锁
	list  *list.List               //用于GC
	data  map[string]*list.Element //用于存储在内存
}

func (m *Memory) Init(sid string) (session.ISession, error) {
	//加锁
	m.mutex.Lock()
	defer m.mutex.Unlock()
	//创建session
	store := &Store{sid: sid, data: make(map[interface{}]interface{}, 0), time: time.Now()}
	elem := m.list.PushBack(store)
	m.data[sid] = elem
	return store, nil
}
func (m *Memory) Read(sid string) (session.ISession, error) {
	ele, ok := m.data[sid]
	if ok {
		return ele.Value.(*Store), nil
	}
	return m.Init(sid)
}
func (m *Memory) Destroy(sid string) error {
	ele, ok := m.data[sid]
	if ok {
		delete(m.data, sid)
		m.list.Remove(ele)
	}
	return nil
}
func (m *Memory) GC(maxAge int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for {
		ele := m.list.Back()
		if ele == nil {
			break
		}
		Session := ele.Value.(*Store)
		if Session.time.Unix()+maxAge >= time.Now().Unix() {
			break
		}
		m.list.Remove(ele)
		delete(m.data, Session.sid)
	}
}
func (m *Memory) Update(sid string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	ele, ok := m.data[sid]
	if ok {
		ele.Value.(*Store).time = time.Now()
		m.list.MoveToFront(ele)
	}
	return nil
}

var memory = &Memory{list: list.New()}

func init() {
	memory.data = make(map[string]*list.Element, 0)
	session.Register("memory", memory)
}
