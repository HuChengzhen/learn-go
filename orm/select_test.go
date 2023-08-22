package orm

// func TestSlectorBuilder(t *testing.T) {
// 	db, err := Open()
// 	require.NoError(t, err)
// 	testCases := []struct {
// 		name      string
// 		builder   QueryBuilder
// 		wantQuery *Query
// 		wantErr   error
// 	}{
// 		{
// 			name:    "no from",
// 			builder: NewSelector[TestModel](db),
// 			wantQuery: &Query{
// 				SQL:  "SELECT * FROM `test_model`;",
// 				Args: nil,
// 			},
// 		},
// 		{
// 			name:    "from",
// 			builder: (NewSelector[TestModel](db)).From("test_model"),
// 			wantQuery: &Query{
// 				SQL:  "SELECT * FROM test_model;",
// 				Args: nil,
// 			},
// 		},

// 		{
// 			name:    "empty from",
// 			builder: (NewSelector[TestModel](db)).From(""),
// 			wantQuery: &Query{
// 				SQL:  "SELECT * FROM `test_model`;",
// 				Args: nil,
// 			},
// 		},
// 		{
// 			name:    "where",
// 			builder: (NewSelector[TestModel](db)).Where(C("Age").Eq(18)),
// 			wantQuery: &Query{
// 				SQL:  "SELECT * FROM `test_model`WHERE `Age`=?;",
// 				Args: []any{18},
// 			},
// 		},
// 		{
// 			name:    "or",
// 			builder: (NewSelector[TestModel](db)).Where(Not(C("Age").Eq(18))),
// 			wantQuery: &Query{
// 				SQL:  "SELECT * FROM `test_model`WHERE NOT `Age`=?;",
// 				Args: []any{18},
// 			},
// 		},

// 		{
// 			name:    "or",
// 			builder: (NewSelector[TestModel](db)).Where(Not(C("Age").Eq(18)).And(C("name").Eq("asdf"))),
// 			wantQuery: &Query{
// 				SQL:  "SELECT * FROM `test_model`WHERE NOT `Age`=?;",
// 				Args: []any{18},
// 			},
// 		},
// 	}

// 	for _, testCase := range testCases {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			build, err := testCase.builder.Build()

// 			assert.Equal(t, testCase.wantErr, err)
// 			if err != nil {
// 				return
// 			}
// 			assert.Equal(t, testCase.wantQuery, build)
// 		})
// 	}
// }

// type TestModel struct {
// 	Id        int64
// 	FirstName string
// 	Age       int8
// 	LastName  *sql.NullString
// }
