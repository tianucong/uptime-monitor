package main

import (
	"net/http"
	"time"
)

type CheckResult struct{
	Status string
	StatusCode int 
	LatencyMs int64
}

var httpClient=&http.Client{
	Timeout:5*time.Second,
}

func cekWebsite(url string)CheckResult{
	start:=time.Now()

	resp,err:=httpClient.Get(url)
	latency:=time.Since(start).Milliseconds()

	if err!= nil{
		return CheckResult{
			Status : "DOWN",
			StatusCode : 0,
			LatencyMs: latency,

		}
	}
	defer resp.Body.Close()

	status:="UP"
	if resp.StatusCode>=400{
		status="DOWN"
	}

	return CheckResult{
		Status:status,
		StatusCode:resp.StatusCode,
		LatencyMs:latency,

	}
}