package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func initDB() {
	var err error
	dsn := "root:02020202@tcp(192.168.178.1:3306)/StudentService?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal("数据库连接失败", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("数据库连接测试失败", err)
	}

	log.Println("数据库连接成功", err)
}
