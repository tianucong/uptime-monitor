package main

import "time"

type Monitor struct{
	ID int `json:"id"`
	Name string `json:"name"`
	URL string `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateMonitorRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}