package main

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func connectDB(dsn string)(*sql.DB,error){
	db,err:=sql.Open("mysql",dsn)
	if err != nil{
		return nil,err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5*time.Minute)

	if err:=db.Ping();err!=nil{
		return nil,err
	}
	return db, nil
}