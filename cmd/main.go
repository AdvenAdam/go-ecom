package main

import (
	"fmt"
	"log"

	"github.com/AdvenAdam/go-ecom/cmd/api"
	"github.com/AdvenAdam/go-ecom/config"
	"github.com/AdvenAdam/go-ecom/db"
	"github.com/go-sql-driver/mysql"
)

func main() {
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Net:                  "tcp",
		Addr:                 config.Envs.DBHost,
		DBName:               config.Envs.DBName,
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}
	port := fmt.Sprintf(":%s", config.Envs.Port)
	server := api.NewAPIServer(port, db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
