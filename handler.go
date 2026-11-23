package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func writeJSON(w http.ResponseWriter,status int, data any){
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func createMonitorHandler(db *sql.DB)http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request){

		var req CreateMonitorRequest
		if err:=json.NewDecoder(r.Body).Decode(&req);err!=nil{
			writeJSON(w, http.StatusBadRequest,map[string]string{
				"error":"format JSON tidak valid",
			})
			return
		}

		req.Name=strings.TrimSpace(req.Name)
		req.URL=strings.TrimSpace(req.URL)

		if req.Name==""||req.URL==""{
			writeJSON(w, http.StatusBadRequest,map[string]string{
				"error":"nama dan url wajib diisi",
			})
			return
		}

		parsed,err:=url.ParseRequestURI(req.URL)
		if err !=nil||(parsed.Scheme!="http"&&parsed.Scheme!="https")||parsed.Host==""{
			writeJSON(w, http.StatusBadRequest,map[string]string{
				"error": "url harus diawali http//: atau https://",
			})
			return
		}

		result, err:=db.Exec(
			"INSERT INTO monitors (name,url) VALUES(?,?)",
			req.Name,req.URL,
		)

		if err!=nil{
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error":"gagal menyimpan ke database",
			})
			return
		}
		id,_:=result.LastInsertId()

		writeJSON(w,http.StatusCreated,Monitor{
			ID :int (id),
			Name:req.Name,
			URL:req.URL,
			CreatedAt: time.Now(),
		})
	}


}

func listMonitorHandler(db *sql.DB)http.HandlerFunc{
	return func (w http.ResponseWriter,r *http.Request){
		rows, err:=db.Query(`
		SELECT  m.id, m.name,m.url,m.created_at,
		l.status,l.status_code,l.latency_ms,l.checked_at
		FROM monitors ,
		LEFT JOIN monitor_logs l
		On l.id=(SELECT MAX (id) FROM monitor_logs WHERE monitor_id=m.id)
		ORDER BY M.id DESC

		`)
		if err!=nil{
			log.Println("gagal query monitors:",err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "gagal membaca data",
			})
			return
		}
		defer rows.Close()
		monitors :=[]Monitor{}

		for rows.Next(){
			var m Monitor 

			if err:=rows.Scan(&m.ID, &m.Name, &m.URL, &m.CreatedAt,&m.LastStatus,&m.LastStatusCode,&m.LastLatencyMs,&m.LastCheckedAt,); err!=nil{
				log.Println("gagal scan baris:",err)
				writeJSON(w, http.StatusInternalServerError, map[string]string{
					"error": "gagal membaca data",
				})
				return
			}
			monitors=append(monitors,m)
		}
		if err:=rows.Err();err!=nil{
			log.Println("error saat iterasi baris:",err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{
					"error": "gagal membaca data",
			})
			return
		}
		writeJSON(w,http.StatusOK,monitors)
	}
}

func checkMonitorHandler(db *sql.DB)http.HandlerFunc{
	return func (w http.ResponseWriter,r *http.Request){
		id,err:=strconv.Atoi(r.PathValue("id"))
		if err != nil{
			writeJSON(w, http.StatusBadRequest,map[string]string{
				"error":"id harus berupa angka",
			})
			return
		}
		var target string
		err=db.QueryRow("SELECT url FROM monitors WHERE id=?",id).Scan(&target)
		if errors.Is(err,sql.ErrNoRows){
			writeJSON(w,http.StatusNotFound,map[string]string{
				"error":"monitor tidak ditemukan",
			})
			return
		}
		if err != nil{
			log.Println("gagal ambil monitor:",err)
			writeJSON(w, http.StatusInternalServerError,map[string]string{
				"error":"gagal membaca data",
			})
			return 
		}
		hasil:=cekWebsite(target)

		if err:=saveLog(db,id,hasil);err!=nil{
			log.Println("gagal simpan log:",err)
			writeJSON(w, http.StatusInternalServerError,map[string]string{
				"error":"gagal menyimpan hasil",
			})
			return
		}

		writeJSON(w,http.StatusOK,map[string]any{
			"url":target,
			"status":hasil.Status,
			"status_code": hasil.StatusCode,
			"latency_ms":hasil.LatencyMs,
		})
	}
}