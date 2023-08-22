package unsafe

import (
	"errors"
	"reflect"
	"unsafe"
)

type UnsafeAccessor struct {
	fields  map[string]FieldMeta
	address unsafe.Pointer
}

func NewUnsafeAccessor(entity any) *UnsafeAccessor {
	typ := reflect.TypeOf(entity)

	numField := typ.NumField()

	fields := make(map[string]FieldMeta, numField)

	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		fields[fd.Name] = FieldMeta{
			offset: fd.Offset,
			typ:    fd.Type,
		}
	}
	val := reflect.ValueOf(entity)
	return &UnsafeAccessor{
		fields:  fields,
		address: val.UnsafePointer(),
	}
}

func (u *UnsafeAccessor) Field(field string) (any, error) {
	fd, ok := u.fields[field]
	if !ok {
		return nil, errors.New("非法字段")
	}

	fdAddress := unsafe.Pointer(uintptr(u.address) + fd.offset)

	//
	//return *(*int)(fdAddress), nil
	return reflect.NewAt(fd.typ, fdAddress).Elem().Interface(), nil
}

func (u *UnsafeAccessor) SetField(field string, val any) error {
	fd, ok := u.fields[field]
	if !ok {
		return errors.New("非法字段")
	}

	fdAddress := unsafe.Pointer(uintptr(u.address) + fd.offset)
	//*(*int)(fdAddress) = val.(int)
	reflect.NewAt(fd.typ, fdAddress).Elem().Set(reflect.ValueOf(val))
	return nil
}

type FieldMeta struct {
	offset uintptr
	typ    reflect.Type
}
