package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	carnet = "202110378"
	vm     = "VM1"
	port   = ":5001"
)

// IPs de las otras VMs
const (
	vm2IP = "192.168.122.101"
	vm3IP = "192.168.122.245"
)

var httpClient = &http.Client{Timeout: 3 * time.Second}

type HealthResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	VM        string `json:"VM"`
	Carnet    string `json:"carnet"`
}

type CallResponse struct {
	Apiname    string `json:"apiname"`
	Message    string `json:"message"`
	Connection bool   `json:"connection"`
	Carnet     string `json:"carnet"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{
		Status:    "UP",
		Message:   "API1 is Ready",
		Timestamp: time.Now().Format(time.RFC3339),
				  VM:        vm,
				  Carnet:    carnet,
	})
}

func callAPI2Handler(w http.ResponseWriter, r *http.Request) {
	resp, err := httpClient.Get(fmt.Sprintf("http://%s:5002/health", vm2IP))
	if err == nil && resp.StatusCode == 200 {
		var health HealthResponse
		json.NewDecoder(resp.Body).Decode(&health)
		resp.Body.Close()
		if health.Status == "UP" {
			json.NewEncoder(w).Encode(CallResponse{
				Apiname:    "API2",
				Message:    "The API2 located on the VM2 is working",
				Connection: true,
				Carnet:     carnet,
			})
			return
		}
	}
	json.NewEncoder(w).Encode(CallResponse{
		Apiname:    "API2",
		Message:    "ERROR: The API2 located on the VM2 is not working",
		Connection: false,
		Carnet:     carnet,
	})
}

func callAPI3Handler(w http.ResponseWriter, r *http.Request) {
	resp, err := httpClient.Get(fmt.Sprintf("http://%s:5003/health", vm3IP))
	if err == nil && resp.StatusCode == 200 {
		var health HealthResponse
		json.NewDecoder(resp.Body).Decode(&health)
		resp.Body.Close()
		if health.Status == "UP" {
			json.NewEncoder(w).Encode(CallResponse{
				Apiname:    "API3",
				Message:    "The API3 located on the VM3 is working",
				Connection: true,
				Carnet:     carnet,
			})
			return
		}
	}
	json.NewEncoder(w).Encode(CallResponse{
		Apiname:    "API3",
		Message:    "ERROR: The API3 located on the VM3 is not working",
		Connection: false,
		Carnet:     carnet,
	})
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api1/202110378/call-api2", callAPI2Handler)
	http.HandleFunc("/api1/202110378/call-api3", callAPI3Handler)

	fmt.Println("API1 corriendo en puerto 5001")
	http.ListenAndServe(port, nil)
}
