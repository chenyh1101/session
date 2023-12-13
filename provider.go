package session

//session 管理接口 管理的是一个session对象

type IProvider interface {
	Init(sid string) (ISession, error)
	Read(sid string) (ISession, error)
	Destroy(sid string) error
	GC(maxAge int64)
}

var providers = make(map[string]IProvider)

// Register 根据provider管理器名称 获取管理器名字（不同存储方式和管理方式会有不同的管理器）

func Register(name string, provider IProvider) {
	if provider == nil {
		panic("provider register:provider is null")
	}
	if _, ok := providers[name]; ok {
		panic("provider register:provider already exists")
	}
	providers[name] = provider
}
