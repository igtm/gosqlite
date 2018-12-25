package main

import (
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

func main() {
	// DB
	db, err := sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		panic(err)
	}

	// Table作成
	_, err = db.Exec("CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, name TEXT)", nil)
	if err != nil {
		panic(err)
	}

	// insert
	err = insert(db, "hello")
	if err != nil {
		panic(err)
	}
	err = insert(db, "goodbye")
	if err != nil {
		panic(err)
	}

	// Select
	rows, err := db.Query("SELECT name FROM foo")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		n := ""
		rows.Scan(&n)

		fmt.Printf("%v\n", n)
	}

}

type inserter interface {
	Exec(string, ...interface{}) (sql.Result, error)
}

func insert(i inserter, s string) error {
	for {
		_, err := i.Exec("INSERT INTO foo (name) VALUES (?)", s)
		if err == nil {
			return nil
		}
		if err == sqlite3.ErrLocked || err == sqlite3.ErrBusy || err.Error() == "database table is locked" {
			continue
		}
		return err
	}
}
