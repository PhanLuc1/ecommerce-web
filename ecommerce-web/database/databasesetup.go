package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func DBSet() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/ecommerce?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Susscessfully connected to MYSQL")
	return db
}

var Client *sql.DB = DBSet()

func UserData(DB *sql.DB, collectionName string) *sql.Rows {
	rows, err := DB.Query("SELECT * FROM " + collectionName)
	if err != nil {
		return nil
	}
	return rows
}
func ProductData(DB *sql.DB, collectionName string) *sql.Rows {
	rows, err := DB.Query("SELECT * FROM " + collectionName)
	if err != nil {
		return nil
	}
	return rows
}
func CartItemData(DB *sql.DB, collectionName string) *sql.Rows {
	rows, err := DB.Query("SELECT * FROM " + collectionName)
	if err != nil {
		return nil
	}
	return rows
}
