package tcpio

import (
	"errors"
	"fmt"
	"reflect"
)

type caller struct {
	Func       reflect.Value
	Args       []reflect.Type
	NeedSocket bool
}

func newCaller(f interface{}) (*caller, error) {
	// 获取传入接口的类型
	fv := reflect.ValueOf(f)
	// 如果传入的不是方式异常
	if fv.Kind() != reflect.Func {
		return nil, fmt.Errorf("f is not func")
	}
	ft := fv.Type()
	// 如果方法体不需要传入参数则直接返回
	if ft.NumIn() == 0 {
		return &caller{
			Func: fv,
		}, nil
	}
	// 取出需要传递的参数
	args := make([]reflect.Type, ft.NumIn())
	for i, n := 0, ft.NumIn(); i < n; i++ {
		args[i] = ft.In(i)
	}
	// 判断第一个传入的参数名称是否是socket
	needSocket := false
	if args[0].Name() == "Socket" {
		args = args[1:]
		needSocket = true
	}
	// 返回对象
	return &caller{
		Func:       fv,
		Args:       args,
		NeedSocket: needSocket,
	}, nil
}

func (c *caller) GetArgs() []interface{} {
	ret := make([]interface{}, len(c.Args))
	for i, argT := range c.Args {
		if argT.Kind() == reflect.Ptr {
			argT = argT.Elem()
		}
		v := reflect.New(argT)
		ret[i] = v.Interface()
	}
	return ret
}

func (c *caller) Call(so Socket, args []interface{}) []reflect.Value {
	var a []reflect.Value
	diff := 0
	if c.NeedSocket {
		diff = 1
		a = make([]reflect.Value, len(args)+1)
		a[0] = reflect.ValueOf(so)
	} else {
		a = make([]reflect.Value, len(args))
	}

	if len(args) != len(c.Args) {
		return []reflect.Value{reflect.ValueOf([]interface{}{}), reflect.ValueOf(errors.New("Arguments do not match"))}
	}

	for i, arg := range args {
		v := reflect.ValueOf(arg)
		if c.Args[i].Kind() != reflect.Ptr {
			if v.IsValid() {
				v = v.Elem()
			} else {
				v = reflect.Zero(c.Args[i])
			}
		}
		a[i+diff] = v
	}

	return c.Func.Call(a)
}
