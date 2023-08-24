package rpc

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

func InitClientProxy(service Service) error {

	return setFunField(service, nil)
}

func setFunField(service Service, p Proxy) error {
	if service == nil {
		return errors.New("rpc: not support nil")
	}
	val := reflect.ValueOf(service)

	typ := reflect.TypeOf(service)

	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return errors.New("rpc: only support one level pointer")
	}

	val = val.Elem()
	typ = typ.Elem()

	numField := typ.NumField()
	for i := 0; i < numField; i++ {
		fieldTyp := typ.Field(i)
		fieldVal := val.Field(i)

		if fieldVal.CanSet() {
			fn := func(args []reflect.Value) (results []reflect.Value) {

				ctx := args[0].Interface().(context.Context)

				req := &Request{
					ServiceName: service.Name(),
					MethodName:  fieldTyp.Name,
					Args:        args[1].Interface(),
				}

				var p Proxy

				resp, err := p.Invoke(ctx, req)
				if err != nil {
					return []reflect.Value{reflect.Zero(fieldTyp.Type.Out(0)), reflect.ValueOf(err)}
				}

				respVal := reflect.ValueOf(resp)

				fmt.Println(respVal)
				return []reflect.Value{reflect.Zero(fieldTyp.Type.Out(0)), reflect.ValueOf(((*error)(nil)))}
			}

			fnVal := reflect.MakeFunc(fieldTyp.Type, fn)

			fieldVal.Set(fnVal)
		}
	}
}
