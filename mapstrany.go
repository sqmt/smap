package smap

import "sync"

type MapStrAny struct {
    mutex *sync.RWMutex
    data  map[string]interface{}
}

func NewMapStrAny() *MapStrAny {
    return &MapStrAny{
        mutex: new(sync.RWMutex),
        data:  make(map[string]interface{}),
    }
}

func (c *MapStrAny) Keys() []string {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    keys := make([]string, 0)
    for s := range c.data {
        keys = append(keys, s)
    }
    return keys
}

func (c *MapStrAny) Size() int {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    return len(c.data)
}

func (c *MapStrAny) All() map[string]interface{} {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    return c.data
}

func (c *MapStrAny) Has(key string) bool {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    _, ok := c.data[key]

    return ok
}

func (c *MapStrAny) Get(key string) interface{} {
    if c.Has(key) {
        return c.data[key]
    }

    return nil
}

func (c *MapStrAny) Set(key string, value interface{}) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.data[key] = value
}

func (c *MapStrAny) Remove(key ...string) {
    if len(key) == 0 {
        return
    }
    c.mutex.Lock()
    defer c.mutex.Unlock()
    for _, s := range key {
        delete(c.data, s)
    }
}

func (c *MapStrAny) Load(key string) (i interface{}, ok bool) {
    i = c.Get(key)

    return i, i != nil
}
