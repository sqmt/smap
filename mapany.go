package smap

import (
    "fmt"
    "sync"
)

type MapAny struct {
    safe  bool
    mutex *sync.RWMutex
    data  map[interface{}]interface{}
}

// NewMapAny 创建一个MapAny对象
// 默认是非并发安全的，如果传入safe参数为true则开启并发安全锁
func NewMapAny(safe ...bool) *MapAny {
    m := &MapAny{}
    if len(safe) > 0 {
        m.safe = safe[0]
    }
    if m.safe {
        m.mutex = new(sync.RWMutex)
    }
    m.data = make(map[interface{}]interface{})
    return m
}

// errorRecover 捕获panic异常信息
func errorRecover(err *error) {
    if r := recover(); r != nil {
        *err = fmt.Errorf("%v", r)
    }
}

// lock 加锁
func (a *MapAny) lock() {
    if a.mutex != nil {
        a.mutex.Lock()
    }
}

// unlock 解锁
func (a *MapAny) unlock() {
    if a.mutex != nil {
        a.mutex.Unlock()
    }
}

// Get 获取value
func (a *MapAny) Get(key interface{}) (val interface{}, ok bool, err error) {
    a.lock()
    defer a.unlock()
    defer errorRecover(&err)
    if v, ok := a.data[key]; ok {
        return v, ok, nil
    }

    return nil, false, err
}

// Set 设置k/v
func (a *MapAny) Set(key interface{}, value interface{}) (err error) {
    a.lock()
    defer a.unlock()
    defer errorRecover(&err)
    a.data[key] = value

    return err
}

// Has 检查key是否存在
func (a *MapAny) Has(key interface{}) (b bool) {
    a.lock()
    defer a.unlock()
    defer func(b *bool) {
        if r := recover(); r != nil {
            *b = false
        }
    }(&b)
    if _, ok := a.data[key]; ok {
        return ok
    }

    return b
}

// Remove 移除单个或多个key
func (a *MapAny) Remove(keys ...interface{}) (err error) {
    if len(keys) == 0 {
        return
    }
    a.lock()
    defer a.unlock()
    defer errorRecover(&err)
    for _, key := range keys {
        delete(a.data, key)
    }

    return err
}

// Keys 获取所有key
func (a *MapAny) Keys() []interface{} {
    a.lock()
    defer a.unlock()
    keys := make([]interface{}, 0)
    for s := range a.data {
        keys = append(keys, s)
    }
    return keys
}

// Size 获取数据长度
func (a *MapAny) Size() int {
    a.lock()
    defer a.unlock()

    return len(a.data)
}

// All 获取所有数据
func (a *MapAny) All() map[interface{}]interface{} {
    a.lock()
    defer a.unlock()
    return a.data
}
