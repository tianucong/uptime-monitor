package main

import "time"

type Monitor struct{
	ID int `json:"id"`
	Name string `json:"name"`
	URL string `json:"url"`
	CreatedAt time.Time `json:"created_at"`

	LastStatus *string `json:"last_status"`
	LastStatusCode *string `json:"last_status_code"`
	LastLatencyMs *string `json:"last_latency_ms"`
	LastCheckedAt *time.Time `json:"last_checked_at"`
}

type CreateMonitorRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}