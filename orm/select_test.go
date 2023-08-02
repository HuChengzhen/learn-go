package orm

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlectorBuilder(t *testing.T) {
	testCases := []struct {
		name      string
		builder   QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name:    "no from",
			builder: &Selector[TestModel]{},
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel`;",
				Args: nil,
			},
		},
		{
			name:    "no from",
			builder: (&Selector[TestModel]{}).From("test_model"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
		},

		{
			name:    "empty from",
			builder: (&Selector[TestModel]{}).From(""),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel`;",
				Args: nil,
			},
		},
		{
			name:    "where",
			builder: (&Selector[TestModel]{}).Where(C("Age").Eq(18)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel`WHERE `Age`=?;",
				Args: []any{18},
			},
		},
		{
			name:    "or",
			builder: (&Selector[TestModel]{}).Where(Not(C("Age").Eq(18))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel`WHERE NOT `Age`=?;",
				Args: []any{18},
			},
		},

		{
			name:    "or",
			builder: (&Selector[TestModel]{}).Where(Not(C("Age").Eq(18)).And(C("name").Eq("asdf"))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel`WHERE NOT `Age`=?;",
				Args: []any{18},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			build, err := testCase.builder.Build()

			assert.Equal(t, testCase.wantErr, err)
			if err != nil {

			}
			assert.Equal(t, testCase.wantQuery, build)
		})
	}
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
