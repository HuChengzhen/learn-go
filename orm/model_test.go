package orm

import (
	"github.com/stretchr/testify/assert"
	"learn_geektime_go/orm/internal/errs"
	"testing"
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
		},
		{
			name:      "error",
			entity:    TestModel{},
			wantModel: nil,
			wantErr:   errs.ErrPointerOnly,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := parseModel(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
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
