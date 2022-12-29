package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/podhmo/nullable"
)

type Foo struct {
	ID   int    `db:"id" json:"id"` // pk
	Name string `db:"name" json:"name"`

	NickName nullable.Type[string] `db:"nickname" json:"nickname"`
	// NickName sql.NullString `db:"nickname" json:"nickname"`
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
func run() error {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := `
	create table foo (id integer not null primary key, name text not null, nickname text);
	`
	if _, err := db.Exec(sqlStmt); err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into foo(id, name, nickname) values(?, ?, ?)")
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()
	for i := 0; i < 5; i++ {
		if i%2 == 0 {
			_, err = stmt.Exec(i, fmt.Sprintf("こんにちは世界%03d", i), nil)
		} else {
			_, err = stmt.Exec(i, fmt.Sprintf("こんにちは世界%03d", i), "ko")
		}
		if err != nil {
			return fmt.Errorf("insert: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	rows, err := db.Query("select id, name, nickname from foo")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var ob Foo
		err = rows.Scan(&ob.ID, &ob.Name, &ob.NickName)
		if err != nil {
			return fmt.Errorf("scan: %w", err)
		}
		fmt.Printf("%#+v\n", ob)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}
