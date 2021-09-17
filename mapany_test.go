package smap

import (
    "reflect"
    "sync"
    "testing"
)

func getMapAny(data map[interface{}]interface{}, safe ...bool) *MapAny {
    m := NewMapAny(safe...)
    for i, i2 := range data {
        m.Set(i, i2)
    }
    return m
}

func TestNewMapAny(t *testing.T) {
    type args struct {
        safe []bool
    }
    tests := []struct {
        name string
        args args
        want *MapAny
    }{
        {name: "unsafe", want: &MapAny{safe: false, data: map[interface{}]interface{}{}}},
        {name: "safe", args: args{safe: []bool{true}}, want: &MapAny{safe: true, mutex: new(sync.RWMutex), data: map[interface{}]interface{}{}}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := NewMapAny(tt.args.safe...); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("NewMapAny() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestMapAny_All(t *testing.T) {
    data := map[interface{}]interface{}{1: 1, 2: 2, "3": "3"}
    tests := []struct {
        name string
        a    *MapAny
        want map[interface{}]interface{}
    }{
        {name: "unsafe empty", a: getMapAny(nil), want: map[interface{}]interface{}{}},
        {name: "unsafe not empty", a: getMapAny(data), want: data},
        {name: "safe empty", a: getMapAny(nil, true), want: map[interface{}]interface{}{}},
        {name: "safe not empty", a: getMapAny(data, true), want: data},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.a.All(); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("All() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestMapAny_Get(t *testing.T) {
    data := map[interface{}]interface{}{"test": 1, 1: 3}
    type args struct {
        key interface{}
    }
    tests := []struct {
        name    string
        a       *MapAny
        args    args
        wantVal interface{}
        wantOk  bool
        wantErr bool
    }{
        {name: "unsafe not found", a: getMapAny(data), args: args{"a"}, wantVal: nil, wantOk: false, wantErr: false},
        {name: "unsafe found", a: getMapAny(data), args: args{"test"}, wantVal: 1, wantOk: true, wantErr: false},
        {name: "unsafe error", a: getMapAny(data), args: args{[]string{"test"}}, wantVal: nil, wantOk: false, wantErr: true},
        {name: "safe not found", a: getMapAny(data, true), args: args{"a"}, wantVal: nil, wantOk: false, wantErr: false},
        {name: "safe found", a: getMapAny(data, true), args: args{"test"}, wantVal: 1, wantOk: true, wantErr: false},
        {name: "safe error", a: getMapAny(data, true), args: args{[]string{"test"}}, wantVal: nil, wantOk: false, wantErr: true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            gotVal, gotOk, err := tt.a.Get(tt.args.key)
            if (err != nil) != tt.wantErr {
                t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(gotVal, tt.wantVal) {
                t.Errorf("Get() gotVal = %v, want %v", gotVal, tt.wantVal)
            }
            if gotOk != tt.wantOk {
                t.Errorf("Get() gotOk = %v, want %v", gotOk, tt.wantOk)
            }
        })
    }
}

func TestMapAny_Has(t *testing.T) {
    data := map[interface{}]interface{}{"test": 1, 1: 3}
    type args struct {
        key interface{}
    }
    tests := []struct {
        name  string
        a     *MapAny
        args  args
        wantB bool
    }{
        {name: "unsafe not exits", a: getMapAny(nil), args: args{"a"}, wantB: false},
        {name: "unsafe not exits", a: getMapAny(data), args: args{"a"}, wantB: false},
        {name: "unsafe exits", a: getMapAny(data), args: args{"test"}, wantB: true},
        {name: "unsafe error", a: getMapAny(data), args: args{[]string{"test"}}, wantB: false},
        {name: "safe not exits", a: getMapAny(data, true), args: args{"a"}, wantB: false},
        {name: "safe not exits", a: getMapAny(nil, true), args: args{"a"}, wantB: false},
        {name: "safe exits", a: getMapAny(data, true), args: args{"test"}, wantB: true},
        {name: "safe error", a: getMapAny(data, true), args: args{[]string{"test"}}, wantB: false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if gotB := tt.a.Has(tt.args.key); gotB != tt.wantB {
                t.Errorf("Has() = %v, want %v", gotB, tt.wantB)
            }
        })
    }
}

func TestMapAny_Keys(t *testing.T) {
    data := map[interface{}]interface{}{1: 1, 2: 2, "3": "3"}
    tests := []struct {
        name string
        a    *MapAny
        want []interface{}
    }{
        {name: "unsafe empty", a: getMapAny(nil), want: []interface{}{}},
        {name: "unsafe not empty", a: getMapAny(data), want: []interface{}{1, 2, "3"}},
        {name: "safe empty", a: getMapAny(nil, true), want: []interface{}{}},
        {name: "safe not empty", a: getMapAny(data, true), want: []interface{}{1, 2, "3"}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := tt.a.Keys()
            if len(got) != len(tt.want) {
                t.Errorf("Keys() = %v, want %v", got, tt.want)
            }
            if err := tt.a.Remove(tt.want...); err != nil || len(tt.a.Keys()) != 0 {
                t.Errorf("Keys() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestMapAny_Remove(t *testing.T) {
    type args struct {
        keys []interface{}
    }
    tests := []struct {
        name     string
        a        *MapAny
        args     args
        wantErr  bool
        wantKeys []interface{}
    }{
        {name: "unsafe one", a: getMapAny(map[interface{}]interface{}{"a": 1}), args: args{[]interface{}{"test"}}, wantErr: false, wantKeys: []interface{}{"a"}},
        {name: "unsafe two", a: getMapAny(map[interface{}]interface{}{"test": 1, "b": 1}), args: args{[]interface{}{"test", "b"}}, wantErr: false, wantKeys: []interface{}{}},
        {name: "unsafe three", a: getMapAny(map[interface{}]interface{}{"test": 1}), args: args{[]interface{}{[]string{"test"}}}, wantErr: true, wantKeys: []interface{}{"test"}},
        {name: "safe one", a: getMapAny(map[interface{}]interface{}{"a": 1}, true), args: args{[]interface{}{"test"}}, wantErr: false, wantKeys: []interface{}{"a"}},
        {name: "safe two", a: getMapAny(map[interface{}]interface{}{"test": 1, "b": 1}, true), args: args{[]interface{}{"test", "b"}}, wantErr: false, wantKeys: []interface{}{}},
        {name: "safe three", a: getMapAny(map[interface{}]interface{}{"test": 1}, true), args: args{[]interface{}{[]string{"test"}}}, wantErr: true, wantKeys: []interface{}{"test"}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if err := tt.a.Remove(tt.args.keys...); (err != nil) != tt.wantErr {
                t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
            }
            if keys := tt.a.Keys(); !reflect.DeepEqual(keys, tt.wantKeys) {
                t.Errorf("Remove() error = key not removed, wantKeys %v, got %v", tt.wantKeys, keys)
            }
        })
    }
}

func TestMapAny_Set(t *testing.T) {
    type args struct {
        key   interface{}
        value interface{}
    }
    tests := []struct {
        name    string
        a       *MapAny
        args    args
        wantErr bool
    }{
        {name: "unsafe one", a: getMapAny(map[interface{}]interface{}{}), args: args{"test", 1}, wantErr: false},
        {name: "unsafe two", a: getMapAny(map[interface{}]interface{}{}), args: args{"test", []string{"test"}}, wantErr: false},
        {name: "unsafe three", a: getMapAny(map[interface{}]interface{}{}), args: args{[]string{"test"}, []string{"test"}}, wantErr: true},
        {name: "safe one", a: getMapAny(map[interface{}]interface{}{}, true), args: args{"test", 1}, wantErr: false},
        {name: "safe two", a: getMapAny(map[interface{}]interface{}{}, true), args: args{"test", []string{"test"}}, wantErr: false},
        {name: "safe three", a: getMapAny(map[interface{}]interface{}{}, true), args: args{[]string{"test"}, []string{"test"}}, wantErr: true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if err := tt.a.Set(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
                t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
            }
            if v, ok, err := tt.a.Get(tt.args.key); !tt.wantErr && (!ok || err != nil || !reflect.DeepEqual(v, tt.args.value)) {
                t.Errorf("Set failed, Get() error = %v, ok %v wantVal %v got %v", err, ok, tt.args.value, v)
            }
        })
    }
}

func TestMapAny_Size(t *testing.T) {
    tests := []struct {
        name string
        a    *MapAny
        want int
    }{
        {name: "unsafe one", a: getMapAny(map[interface{}]interface{}{"a": 1}), want: 1},
        {name: "unsafe two", a: getMapAny(map[interface{}]interface{}{"test": 1, "b": 1}), want: 2},
        {name: "unsafe three", a: getMapAny(map[interface{}]interface{}{}), want: 0},
        {name: "safe one", a: getMapAny(map[interface{}]interface{}{"a": 1}), want: 1},
        {name: "safe two", a: getMapAny(map[interface{}]interface{}{"test": 1, "b": 1}), want: 2},
        {name: "safe three", a: getMapAny(map[interface{}]interface{}{}), want: 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.a.Size(); got != tt.want {
                t.Errorf("Size() = %v, want %v", got, tt.want)
            }
        })
    }
}

func BenchmarkNewMapAny_unsafe(b *testing.B) {
    for i := 0; i < b.N; i++ {
        NewMapAny()
    }
}

func BenchmarkNewMapAny_safe(b *testing.B) {
    for i := 0; i < b.N; i++ {
        NewMapAny(true)
    }
}

func BenchmarkMapAny_Set_unsafe(b *testing.B) {
    m := NewMapAny()
    for i := 0; i < b.N; i++ {
        m.Set("test", "test")
    }
}

func BenchmarkMapAny_Set_safe(b *testing.B) {
    m1 := NewMapAny(true)
    for i := 0; i < b.N; i++ {
        m1.Set("test", "test")
    }
}

func BenchmarkMapAny_Get_unsafe(b *testing.B) {
    m := NewMapAny()
    m.Set("test", "test")
    for i := 0; i < b.N; i++ {
        m.Get("test")
    }
}

func BenchmarkMapAny_Get_safe(b *testing.B) {
    m1 := NewMapAny(true)
    m1.Set("test", "test")
    for i := 0; i < b.N; i++ {
        m1.Get("test")
    }
}
