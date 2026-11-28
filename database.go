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

func saveLog(db *sql.DB,monitorID int, hasil CheckResult) error{
	_,err:=db.Exec(
		"INSERT INTO monitor_logs (monitor_id,status,status_code,latency_ms) VALUES (?,?,?,?)",
		monitorID,hasil.Status,hasil.StatusCode,hasil.LatencyMs,
	)
	return err
}

func getAllMonitors (db *sql.DB)([]Monitor,error){
	rows,err := db.Query("SELECT id,name,url FROM monitors")
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	monitors := []Monitor{}
	for rows.Next(){
		var m Monitor 
		if err := rows.Scan(&m.ID,&m.Name,&m.URL);err != nil{
			return nil,err
		}
		monitors=append(monitors,m)
	}
	if err := rows.Err(); err != nil{
		return nil,err
	}
	return monitors,nil
}