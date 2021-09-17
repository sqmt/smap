package smap

import (
    "reflect"
    "testing"
)

func getMapStrAny(data map[string]interface{}, safe ...bool) *MapStrAny {
    m := NewMapStrAny(safe...)
    for i, i2 := range data {
        m.Set(i, i2)
    }

    return m
}

func TestNewMapStrAny(t *testing.T) {
    type args struct {
        safe []bool
    }
    tests := []struct {
        name string
        args args
        want *MapStrAny
    }{
        {name: "unsafe", want: &MapStrAny{any: &MapAny{safe: false, data: map[interface{}]interface{}{}}}, args: args{}},
        {name: "safe", want: &MapStrAny{any: NewMapAny(true)}, args: args{[]bool{true}}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := NewMapStrAny(tt.args.safe...); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("NewMapStrAny() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestMapStrAny_All(t *testing.T) {
    data := map[string]interface{}{"1": 1, "2": 2, "3": "3"}

    tests := []struct {
        name string
        a    *MapStrAny
        want map[string]interface{}
    }{
        {name: "unsafe empty", a: getMapStrAny(nil), want: map[string]interface{}{}},
        {name: "unsafe not empty", a: getMapStrAny(data), want: data},
        {name: "safe empty", a: getMapStrAny(nil, true), want: map[string]interface{}{}},
        {name: "safe not empty", a: getMapStrAny(data, true), want: data},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.a.All(); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("All() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestMapStrAny_Get(t *testing.T) {
    type args struct {
        key string
    }
    tests := []struct {
        name    string
        a       *MapStrAny
        args    args
        wantVal interface{}
        wantOk  bool
        wantErr bool
    }{
        {name: "unsafe not found", a: getMapStrAny(map[string]interface{}{"a": 1}), args: args{"test"}, wantVal: nil, wantOk: false, wantErr: false},
        {name: "unsafe found", a: getMapStrAny(map[string]interface{}{"test": 1}), args: args{"test"}, wantVal: 1, wantOk: true, wantErr: false},
        {name: "safe not found", a: getMapStrAny(map[string]interface{}{"a": 1}, true), args: args{"test"}, wantVal: nil, wantOk: false, wantErr: false},
        {name: "safe found", a: getMapStrAny(map[string]interface{}{"test": 1}, true), args: args{"test"}, wantVal: 1, wantOk: true, wantErr: false},
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

func TestMapStrAny_Has(t *testing.T) {
    type args struct {
        key string
    }
    tests := []struct {
        name  string
        a     *MapStrAny
        args  args
        wantB bool
    }{
        {name: "unsafe not found", a: getMapStrAny(map[string]interface{}{"a": 1}), args: args{"test"}, wantB: false},
        {name: "unsafe found", a: getMapStrAny(map[string]interface{}{"test": 1}), args: args{"test"}, wantB: true},
        {name: "safe not found", a: getMapStrAny(map[string]interface{}{"a": 1}, true), args: args{"test"}, wantB: false},
        {name: "safe found", a: getMapStrAny(map[string]interface{}{"test": 1}, true), args: args{"test"}, wantB: true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if gotB := tt.a.Has(tt.args.key); gotB != tt.wantB {
                t.Errorf("Has() = %v, want %v", gotB, tt.wantB)
            }
        })
    }
}

func TestMapStrAny_Keys(t *testing.T) {
    data := map[string]interface{}{"1": 1, "2": 2, "3": "3"}
    tests := []struct {
        name string
        a    *MapStrAny
        want []string
    }{
        {name: "unsafe empty", a: getMapStrAny(nil), want: []string{}},
        {name: "unsafe not empty", a: getMapStrAny(data), want: []string{"1", "2", "3"}},
        {name: "safe empty", a: getMapStrAny(nil, true), want: []string{}},
        {name: "safe not empty", a: getMapStrAny(data, true), want: []string{"1", "2", "3"}},
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

func TestMapStrAny_Remove(t *testing.T) {
    type args struct {
        keys []string
    }
    tests := []struct {
        name     string
        a        *MapStrAny
        args     args
        wantErr  bool
        wantKeys []string
    }{
        {name: "unsafe one", a: getMapStrAny(map[string]interface{}{"a": 1}), args: args{[]string{"test"}}, wantErr: false, wantKeys: []string{"a"}},
        {name: "unsafe two", a: getMapStrAny(map[string]interface{}{"test": 1, "b": 1}), args: args{[]string{"test", "b"}}, wantErr: false, wantKeys: []string{}},
        {name: "safe one", a: getMapStrAny(map[string]interface{}{"a": 1}, true), args: args{[]string{"test"}}, wantErr: false, wantKeys: []string{"a"}},
        {name: "safe two", a: getMapStrAny(map[string]interface{}{"test": 1, "b": 1}, true), args: args{[]string{"test", "b"}}, wantErr: false, wantKeys: []string{}},
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

func TestMapStrAny_Set(t *testing.T) {
    type args struct {
        key   string
        value interface{}
    }
    tests := []struct {
        name    string
        a       *MapStrAny
        args    args
        wantErr bool
    }{
        {name: "unsafe one", a: getMapStrAny(map[string]interface{}{}), args: args{"test", 1}, wantErr: false},
        {name: "unsafe two", a: getMapStrAny(map[string]interface{}{}), args: args{"test", []string{"test"}}, wantErr: false},
        {name: "safe one", a: getMapStrAny(map[string]interface{}{}, true), args: args{"test", 1}, wantErr: false},
        {name: "safe two", a: getMapStrAny(map[string]interface{}{}, true), args: args{"test", []string{"test"}}, wantErr: false},
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

func TestMapStrAny_Size(t *testing.T) {
    tests := []struct {
        name string
        a    *MapStrAny
        want int
    }{
        {name: "unsafe one", a: getMapStrAny(map[string]interface{}{"a": 1}), want: 1},
        {name: "unsafe two", a: getMapStrAny(map[string]interface{}{"test": 1, "b": 1}), want: 2},
        {name: "unsafe three", a: getMapStrAny(map[string]interface{}{}), want: 0},
        {name: "safe one", a: getMapStrAny(map[string]interface{}{"a": 1}), want: 1},
        {name: "safe two", a: getMapStrAny(map[string]interface{}{"test": 1, "b": 1}), want: 2},
        {name: "safe three", a: getMapStrAny(map[string]interface{}{}), want: 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.a.Size(); got != tt.want {
                t.Errorf("Size() = %v, want %v", got, tt.want)
            }
        })
    }
}

func BenchmarkNewMapStrAny_unsafe(b *testing.B) {
    for i := 0; i < b.N; i++ {
        NewMapStrAny()
    }
}

func BenchmarkNewMapStrAny_safe(b *testing.B) {
    for i := 0; i < b.N; i++ {
        NewMapStrAny(true)
    }
}

func BenchmarkMapStrAny_Set_unsafe(b *testing.B) {
    m := NewMapStrAny()
    for i := 0; i < b.N; i++ {
        m.Set("test", "test")
    }
}

func BenchmarkMapStrAny_Set_safe(b *testing.B) {
    m1 := NewMapStrAny(true)
    for i := 0; i < b.N; i++ {
        m1.Set("test", "test")
    }
}

func BenchmarkMapStrAny_Get_unsafe(b *testing.B) {
    m := NewMapStrAny()
    m.Set("test", "test")
    for i := 0; i < b.N; i++ {
        m.Get("test")
    }
}

func BenchmarkMapStrAny_Get_safe(b *testing.B) {
    m1 := NewMapStrAny(true)
    m1.Set("test", "test")
    for i := 0; i < b.N; i++ {
        m1.Get("test")
    }
}
