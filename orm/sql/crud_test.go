package sql

import (
	"context"
	"database/sql"
	"learn_geektime_go/orm"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestDB(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	defer db.Close()
	db.Ping()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	_, err = db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS test_model(
		id INTEGER PRIMARY KEY,
		first_name TEXT NOT NULL,
		age INTEGER,
		last_name TEXT NOT NULL
	)
	`)

	require.NoError(t, err)

	result, err := db.ExecContext(ctx, `INSERT INTO test_model VALUES(?, ?, ?, ?)`, 1, "Tom", 18, "Jerry")
	require.NoError(t, err)

	i, err := result.RowsAffected()
	require.NoError(t, err)
	t.Logf("rows affected %d", i)
	i2, err := result.LastInsertId()
	require.NoError(t, err)
	t.Logf("last insertId %d", i2)

	row := db.QueryRowContext(ctx, `SELECT id, first_name, age, last_name FROM test_model WHERE id = ?`, 1)
	tm := orm.TestModel{}
	err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	require.NoError(t, err)
	t.Logf("tm %v", tm)

	rows, err := db.QueryContext(ctx, `SELECT id, first_name, age, last_name FROM test_model WHERE id = ?`, 1)
	require.NoError(t, err)
	for rows.Next() {
		row.Scan()
	}
	cancel()
}
