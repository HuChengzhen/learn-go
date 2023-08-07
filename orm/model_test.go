package orm

import (
	"learn_geektime_go/orm/internal/errs"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Register(t *testing.T) {
	testCases := []struct {
		name      string
		entity    any
		wantModel *Model
		opts      []ModelOption
		wantErr   error
	}{
		{
			name:   "test model",
			entity: &TestModel{},
			wantModel: &Model{
				tableName: "test_model",
				fields: map[string]*Field{
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
			name:   "test model opts",
			entity: &TestModel{},
			wantModel: &Model{
				tableName: "test_model_t",
				fields: map[string]*Field{
					"Id": {
						colName: "id_t",
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
			opts: []ModelOption{ModelWithTableName("test_model_t"), ModelWithColumnName("Id", "id_t")},
		},
		{
			name:    "test model opts err",
			entity:  &TestModel{},
			wantErr: errs.NewErrUnknownField("Id_df"),
			opts:    []ModelOption{ModelWithTableName("test_model_t"), ModelWithColumnName("Id_df", "id_t")},
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
			m, err := r.Register(tc.entity, tc.opts...)
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
		wantModel *Model
		wantErr   error
		cacheSize int
	}{
		{
			name:   "testModel",
			entity: &TestModel{},
			wantModel: &Model{
				tableName: "test_model",
				fields: map[string]*Field{
					"Id": &Field{
						colName: "id",
					},
					"FirstName": &Field{
						colName: "first_name",
					},
					"LastName": &Field{
						colName: "last_name",
					},
					"Age": &Field{
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
			wantModel: &Model{
				tableName: "tag_table",
				fields: map[string]*Field{
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
			wantModel: &Model{
				tableName: "tag_table",
				fields: map[string]*Field{
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
			m, err := r.Get(tc.entity)
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

type CustomTableName struct {
}

func (c CustomTableName) TableName() string {
	return "custom_table_name_t"
}

type CustomTableNamePtr struct {
}

func (c *CustomTableNamePtr) TableName() string {
	return "custom_table_name_t"
}

func TestCustomTableName(t *testing.T) {
	testCases := []struct {
		name      string
		entity    any
		wantModel *Model
		wantErr   error
		cacheSize int
	}{
		{
			name:   "CustomTableName",
			entity: &CustomTableName{},
			wantModel: &Model{
				tableName: "custom_table_name_t",
				fields:    map[string]*Field{},
			},
		},
		{
			name:   "CustomTableNamePtr",
			entity: &CustomTableNamePtr{},
			wantModel: &Model{
				tableName: "custom_table_name_t",
				fields:    map[string]*Field{},
			},
		},
	}

	r := newRegistry()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Get(tc.entity)
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
