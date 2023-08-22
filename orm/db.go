package orm

import "database/sql"

type DBOption func(*DB)

type DB struct {
	r  *registry
	db *sql.DB
}

func Open(driver string, dataSourceName string, opts ...DBOption) (*DB, error) {
	res, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}
	db := &DB{
		r:  newRegistry(),
		db: res,
	}
	for _, opt := range opts {
		opt(db)
	}
	return db, nil
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		r:  &registry{},
		db: db,
	}

	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

// MustNewDB 创建一个 DB，如果失败则会 panic
// 我个人不太喜欢这种
func MustOpenDB(driver string, dataSourceName string, opts ...DBOption) *DB {
	db, err := Open(driver, dataSourceName, opts...)
	if err != nil {
		panic(err)
	}
	return db
}
