package main

import (
	"database/sql"
	"fmt"

	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
)

const (
	//insertItemQuery = "insert into items ('name', 'price', 'description') values (?,?,?)"
	insertItemQuery  = "insert into items (name, price, description) values (?,?,?)"
	insertItemQuery2 = "insert into items (name, price, description) values ('brownie',240,'sizzling')"
	selectItemQuery  = "select * from items where id=?"
	selectItemQuery2 = "select * from items where id=5"
)

var db sql.DB

/*
func init() {
}
*/
func useQuery(db *sql.DB, query string, params ...interface{}) error {
	if len(params) > 0 {
		rs, err := db.Query(query, params...)
		if err != nil {
			return err
		}
		rs.Close()
		return nil
	}
	rs, err := db.Query(query)
	if err != nil {
		return err
	}
	for rs.Next() {
		fmt.Println("just iterating")
	}

	defer rs.Close()

	return nil
}

func useExec(db *sql.DB, query string, params ...interface{}) error {
	if len(params) > 0 {
		_, err := db.Exec(query, params)
		if err != nil {
			return err
		}
	}
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func usePrepared(db *sql.DB, ch chan error) error {
	stmt, err := db.Prepare(insertItemQuery)
	defer func() {
		stmt.Close()
	}()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	_, err = stmt.Exec("brownie", 240, "sizzling")
	if err != nil {
		ch <- err
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func main() {
	//db, err := sql.Open("sqlite3", "file:/media/sf_alok/swiggy.db?cache=shared")
	db, err := sql.Open("mysql", "vm:vm@tcp(192.168.56.101:3306)/items")
	ch := make(chan error)
	fmt.Println(db.Stats())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for i := 0; i < 300; i++ {
		/*
			err := useQuery(db, selectItemQuery2)
			if err != nil {
				fmt.Println(err.Error())
			}
		*/

		/*
			err := useQuery(db, selectItemQuery, 5)
			if err != nil {
				fmt.Println(err.Error())
			}
		*/
		/*
			err = useExec(db, insertItemQuery2)
			if err != nil {
				fmt.Println(err.Error())
			}
		*/

		/*
			err = useExec(db, insertItemQuery, "brownie", 240, "sizzling")
			if err != nil {
				fmt.Println(err.Error())
			}
		*/

		go usePrepared(db, ch)
		/*
			if err != nil {
				fmt.Println(err.Error())
			}
		*/
	}
	fmt.Println(db.Stats())
	fmt.Println("Test End")
	fmt.Println((<-ch).Error())
}
