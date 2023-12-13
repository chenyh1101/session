package session

//session内部操作
//设置key-value 获取key对应的value 获取一个sessionID

type ISession interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}
