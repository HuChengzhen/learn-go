package orm

import (
	"learn_geektime_go/orm/internal/errs"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseModel(t *testing.T) {
	testCases := []struct {
		name      string
		entity    any
		wantModel *model
		wantErr   error
	}{
		{
			name:   "test model",
			entity: &TestModel{},
			wantModel: &model{
				tableName: "test_model",
				fields: map[string]*field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
					"Age": {
						colName: "age",
					},
				},
			},
		},
		{
			name:      "error",
			entity:    TestModel{},
			wantModel: nil,
			wantErr:   errs.ErrPointerOnly,
		},
	}

	r := &registry{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.parseModel(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
		})
	}
}

func TestRegistry_get(t *testing.T) {
	testCases := []struct {
		name      string
		entity    any
		wantModel *model
		wantErr   error
		cacheSize int
	}{
		{
			name:   "testModel",
			entity: &TestModel{},
			wantModel: &model{
				tableName: "test_model",
				fields: map[string]*field{
					"Id": &field{
						colName: "id",
					},
					"FirstName": &field{
						colName: "first_name",
					},
					"LastName": &field{
						colName: "last_name",
					},
					"Age": &field{
						colName: "age",
					},
				},
			},
			wantErr:   nil,
			cacheSize: 1,
		},
		{
			name: "tag",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column=first_name_t"`
				}

				return new(TagTable)
			}(),
			wantModel: &model{
				tableName: "tag_table",
				fields: map[string]*field{
					"FirstName": {
						colName: "first_name_t",
					},
				},
			},
		},
		{
			name: "empty tag",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column="`
				}

				return new(TagTable)
			}(),
			wantModel: &model{
				tableName: "tag_table",
				fields: map[string]*field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
		},
		{
			name: "no tag",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column"`
				}

				return new(TagTable)
			}(),
			wantErr: errs.NewErrInvalidTagContent("column"),
		},
	}

	r := newRegistry()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.get(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
			typ := reflect.TypeOf(tc.entity)
			m2, ok := r.models.Load(typ)
			assert.True(t, ok)
			assert.Equal(t, tc.wantModel, m2)
		})
	}
}

func TestCamelToSnake(t *testing.T) {
	type args struct {
		camel string
	}
	tests := []struct {
		name      string
		args      args
		wantSnake string
	}{
		{
			name: "ID",
			args: args{
				camel: "ID",
			},
			wantSnake: "id",
		},

		{
			name: "Id",
			args: args{
				camel: "Id",
			},
			wantSnake: "id",
		},
		{
			name: "AbCd",
			args: args{
				camel: "AbCd",
			},
			wantSnake: "ab_cd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantSnake, CamelToSnake(tt.args.camel), "CamelToSnake(%v)", tt.args.camel)
		})
	}
}
