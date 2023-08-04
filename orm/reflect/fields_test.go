package reflect

import (
	"errors"
	"learn_geektime_go/orm/reflect/types"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterateFields(t *testing.T) {
	type User struct {
		Name string
		age  int
	}

	testCases := []struct {
		name    string
		entity  any
		wantErr error
		wantRes map[string]any
	}{
		{
			name: "user",
			entity: User{
				Name: "Tom",
				age:  10,
			},

			wantErr: nil,
			wantRes: map[string]any{
				"Name": "Tom",
				"age":  0,
			},
		},
		{
			name: "pointer",
			entity: &User{
				Name: "Tom",
				age:  10,
			},

			wantErr: nil,
			wantRes: map[string]any{
				"Name": "Tom",
				"age":  0,
			},
		},

		{
			name:    "basic",
			entity:  1,
			wantErr: errors.New("不支持的类型"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fields, err := IterateFields(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, fields)
		})
	}
}

func TestSetField(t *testing.T) {
	type User struct {
		Name string
	}
	testCases := []struct {
		name       string
		entity     any
		field      string
		newValue   any
		wantEntity any
		wantErr    error
	}{
		{
			name: "struct",
			entity: &User{
				Name: "Tom",
			},
			field:    "Name",
			newValue: "Jerry",
			wantEntity: &User{
				Name: "Jerry",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SetField(tc.entity, tc.field, tc.newValue)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			assert.Equal(t, tc.wantEntity, tc.entity)
		})
	}
}

func TestSetBasic(t *testing.T) {
	var i = 0
	p := &i

	reflect.ValueOf(p).Elem().Set(reflect.ValueOf(2))

	assert.Equal(t, 2, i)
}

func TestIterateFunc(t *testing.T) {
	testCases := []struct {
		name    string
		entity  any
		wantRes map[string]FuncInfo
		wantErr error
	}{
		{
			name:   "struct",
			entity: types.NewUser("Tom", 18),
			wantRes: map[string]FuncInfo{
				"GetAge": FuncInfo{
					Name:        "GetAge",
					InputTypes:  []reflect.Type{reflect.TypeOf(types.User{})},
					OutputTypes: []reflect.Type{reflect.TypeOf(0)},
					Result:      []any{18},
				},
				//"ChangeName": FuncInfo{
				//	Name:        "ChangeName",
				//	InputTypes:  []reflect.Type{reflect.TypeOf("")},
				//	OutputTypes: []reflect.Type{},
				//	Result:      []any{},
				//},
			},
		},
		{
			name:   "pointer",
			entity: types.NewUserPointer("Tom", 18),
			wantRes: map[string]FuncInfo{
				"ChangeName": FuncInfo{
					Name:        "ChangeName",
					InputTypes:  []reflect.Type{reflect.TypeOf(&types.User{}), reflect.TypeOf("")},
					OutputTypes: []reflect.Type{},
					Result:      []any{},
				},
				"GetAge": {
					Name: "GetAge",
					// 下标 0 的指向接收器
					InputTypes:  []reflect.Type{reflect.TypeOf(&types.User{})},
					OutputTypes: []reflect.Type{reflect.TypeOf(0)},
					Result:      []any{18},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := IterateFunc(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
