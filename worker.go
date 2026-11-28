package main

import (
	"database/sql"
	"log"
	"sync"
	"time"
)

func jalankanMonitor(db *sql.DB){
	log.Println("Worker mulai,interval 30 detik")
	cekSemua(db)

	ticker := time.NewTicker(30*time.Second)
	defer ticker.Stop()
	for range ticker.C{
		cekSemua(db)
	}
}

const maxConcurrent = 10

func cekSemua(db *sql.DB){
	monitors, err := getAllMonitors(db)
	if err!=nil{
		log.Println("worker: gagal ambil daftar monitor: ",err)
		return 
	}
	mulai:=time.Now()
	log.Printf("worker : mengecek %d website",len(monitors))

	var wg sync.WaitGroup
	sem :=make (chan struct{},maxConcurrent)


	for _,m:= range monitors{
		wg.Add(1)
		sem <- struct{}{}

		go func (m Monitor){
			defer wg.Done()
			defer func(){<-sem}()
		
		hasil:=cekWebsite(m.URL)
		if err:=saveLog(db,m.ID,hasil);err!=nil{
			log.Printf("worker: gagal simpan log %s: %v",m.Name,err)
			return
		}
		log.Printf("worker: %s -> %s (%d) %dms",m.Name,hasil.Status,hasil.StatusCode,hasil.LatencyMs)
	}(m)
}
wg.Wait()
log.Printf("worker : putaran seleasi dalam %dms", time.Since(mulai).Milliseconds())
}