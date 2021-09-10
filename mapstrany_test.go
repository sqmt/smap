package smap_test

import (
    "reflect"
    "testing"

    "github.com/sqmt/blu-core/library/smap"
)

type mock struct {
    key   string
    value interface{}
}

func TestMapStrAny_All(t *testing.T) {
    tests := []struct {
        name string
        want map[string]interface{}
        mock []mock
    }{
        {name: "empty", want: map[string]interface{}{}},
        {name: "all", want: map[string]interface{}{"test1": 1, "test2": 2}, mock: []mock{{"test1", 1}, {"test2", 2}}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := smap.NewMapStrAny()
            for _, m := range tt.mock {
                c.Set(m.key, m.value)
            }
            if got := c.All(); !reflect.DeepEqual(got, tt.want) {
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
        name string
        args args
        want interface{}
        mock []mock
    }{
        {name: "empty", args: args{"test"}, want: nil},
        {name: "not found", args: args{"test"}, want: nil, mock: []mock{{"test1", 1}, {"test2", 2}}},
        {name: "found", args: args{"test1"}, want: 1, mock: []mock{{"test1", 1}, {"test2", 2}}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := smap.NewMapStrAny()
            for _, m := range tt.mock {
                c.Set(m.key, m.value)
            }
            if got := c.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Get() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestMapStrAny_Has(t *testing.T) {
    type args struct {
        key string
    }
    tests := []struct {
        name string
        args args
        want bool
        mock []mock
    }{
        {name: "not found", args: args{"test"}, want: false},
        {name: "found", args: args{"test"}, want: true, mock: []mock{{"test", 1}, {"test2", 2}}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := smap.NewMapStrAny()
            for _, m := range tt.mock {
                c.Set(m.key, m.value)
            }
            if got := c.Has(tt.args.key); got != tt.want {
                t.Errorf("Has() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestMapStrAny_Keys(t *testing.T) {
    tests := []struct {
        name string
        want []string
        mock []mock
    }{
        {name: "empty", want: []string{}},
        {name: "all", want: []string{"test", "test2"}, mock: []mock{{"test", 1}, {"test2", 2}}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := smap.NewMapStrAny()
            for _, m := range tt.mock {
                c.Set(m.key, m.value)
            }
            if got := c.Keys(); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Keys() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestMapStrAny_Load(t *testing.T) {
    type args struct {
        key string
    }
    tests := []struct {
        name   string
        args   args
        wantI  interface{}
        wantOk bool
        mock   []mock
    }{
        {name: "empty", args: args{"test"}, wantI: nil, wantOk: false},
        {name: "found", args: args{"test"}, wantI: 1, wantOk: true, mock: []mock{{"test", 1}, {"test2", 2}}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := smap.NewMapStrAny()
            for _, m := range tt.mock {
                c.Set(m.key, m.value)
            }
            gotI, gotOk := c.Load(tt.args.key)
            if !reflect.DeepEqual(gotI, tt.wantI) {
                t.Errorf("Load() gotI = %v, want %v", gotI, tt.wantI)
            }
            if gotOk != tt.wantOk {
                t.Errorf("Load() gotOk = %v, want %v", gotOk, tt.wantOk)
            }
        })
    }
}

func TestMapStrAny_Remove(t *testing.T) {
    type args struct {
        key []string
    }
    tests := []struct {
        name string
        args args
        mock []mock
        want int
    }{
        {name: "remove empty", args: args{}, mock: []mock{}},
        {name: "remove empty1", args: args{key: []string{"test"}}, mock: []mock{}},
        {name: "remove one", args: args{key: []string{"test"}}, mock: []mock{{"test", 1}, {"test2", 2}}, want: 1},
        {name: "remove more", args: args{key: []string{"test", "test2"}}, mock: []mock{{"test", 1}, {"test2", 2}}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := smap.NewMapStrAny()
            for _, m := range tt.mock {
                c.Set(m.key, m.value)
            }
            c.Remove(tt.args.key...)
            if got := c.Size(); got != tt.want {
                t.Errorf("Remove() failed, got %v, want %v", got, tt.want)
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
        name string
        args args
        want interface{}
    }{
        {name: "set", args: args{"test", 1}, want: 1},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := smap.NewMapStrAny()
            c.Set(tt.args.key, tt.args.value)
            if got := c.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Set() faild, got %v, want %v", got, tt.want)
            }
        })
    }
}

func TestMapStrAny_Size(t *testing.T) {
    tests := []struct {
        name string
        want int
        mock []mock
    }{
        {name: "size empty", want: 0},
        {name: "size 1", want: 1, mock: []mock{{"test", 1}}},
        {name: "size 2", want: 2, mock: []mock{{"test", 1}, {"test2", 2}}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := smap.NewMapStrAny()
            for _, m := range tt.mock {
                c.Set(m.key, m.value)
            }
            if got := c.Size(); got != tt.want {
                t.Errorf("Size() = %v, want %v", got, tt.want)
            }
        })
    }
}

func BenchmarkNewMapStrAny(b *testing.B) {
    for i := 0; i < b.N; i++ {
        smap.NewMapStrAny()
    }
}

func BenchmarkMapStrAny_Set(b *testing.B) {
    s := smap.NewMapStrAny()
    for i := 0; i < b.N; i++ {
        s.Set("test", i)
    }
}
func BenchmarkMapStrAny_Get(b *testing.B) {
    s := smap.NewMapStrAny()
    s.Set("test", 1)
    for i := 0; i < b.N; i++ {
        s.Get("test")
    }
}