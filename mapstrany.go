package smap

type MapStrAny struct {
    any *MapAny
}

// NewMapStrAny 创建一个MapStrAny对象
// 默认是非并发安全的，如果传入safe参数为true则开启并发安全锁
func NewMapStrAny(safe ...bool) *MapStrAny {
    return &MapStrAny{
        any: NewMapAny(safe...),
    }
}

// Get 获取value
func (a *MapStrAny) Get(key string) (val interface{}, ok bool, err error) {
    return a.any.Get(key)
}

// Set 设置k/v
func (a *MapStrAny) Set(key string, value interface{}) (err error) {
    return a.any.Set(key, value)
}

// Has 检查key是否存在
func (a *MapStrAny) Has(key string) (b bool) {
    return a.any.Has(key)
}

// Remove 移除单个或多个key
func (a *MapStrAny) Remove(keys ...string) (err error) {
    ks := make([]interface{}, 0)
    for _, key := range keys {
        ks = append(ks, key)
    }
    return a.any.Remove(ks...)
}

// Keys 获取所有key
func (a *MapStrAny) Keys() []string {
    ks := make([]string, 0)
    for _, k := range a.any.Keys() {
        ks = append(ks, k.(string))
    }
    return ks
}

// Size 获取数据长度
func (a *MapStrAny) Size() int {
    return a.any.Size()
}

// All 获取所有数据
func (a *MapStrAny) All() map[string]interface{} {
    kvs := make(map[string]interface{})
    for i, v := range a.any.All() {
        kvs[i.(string)] = v
    }

    return kvs
}
