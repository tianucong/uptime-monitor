package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func main(){
	dsn:="root:@tcp(127.0.0.1:3306)/uptime_monitor?parseTime=true&charset=utf8mb4"

	db,err:=connectDB(dsn)
	if err!=nil{
		log.Fatal("Gagal menghubungkan ke database:",err)
	}

	defer db.Close()
	log.Println("Berhasil connnect ke database")

	go jalankanMonitor(db)

	mux:=http.NewServeMux()
	mux.HandleFunc("/health",healthHandler(db))
	mux.HandleFunc("GET /api/monitors",listMonitorHandler(db))
	mux.HandleFunc("POST /api/monitors",createMonitorHandler(db))
	mux.HandleFunc("POST /api/monitors/{id}/check",checkMonitorHandler(db))

	log.Println("Server jalan di http://localhost:8080")
	err=http.ListenAndServe(":8080",mux)
	if err!=nil{
		log.Fatal(err)
	}
}

func healthHandler(db *sql.DB)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		statusDB:="ok"
		if err:=db.Ping();err!=nil{
			statusDB="error"
		}
	
		w.Header().Set("Content-Type","application/json")
		json.NewEncoder(w).Encode(map[string]string{
		"status":"ok",
		"database":statusDB,
	})
}
}