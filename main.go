// Файл main.go содержит код главной функции сервера
package main

import (
	"log"

	"github.com/akionka/aviasales/internal/store/mysqlstore"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := sqlx.Connect("mysql", "root:password@(localhost)/aviacompany?parseTime=true&time_zone=%27GMT%27")
	if err != nil {
		log.Fatal(err)
	}

	store := mysqlstore.New(db)
	server := newServer(store)
	log.Println("Starting server on port :8080...")
	log.Fatal(server.start())
}
