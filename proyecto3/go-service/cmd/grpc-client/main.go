package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "go-service/proto"

)

//struct para re4cibir el json del api rest

type WarReportJSON struct {
	Country         string `json:"country"`
	WarplanesInAir  int32  `json:"warplanes_in_air"`
	WarshipsInWater int32  `json:"warships_in_water"`
	Timestamp       string `json:"timestamp"`
}

//mapeo de string a enum del proto

var countryMap = map[string]pb.Countries{
	"USA":    pb.Countries_usa,
	"RUS":    pb.Countries_rus,
	"CHN":    pb.Countries_chn,
	"ESP":    pb.Countries_esp,
	"GTM":    pb.Countries_gtm,
}

//cliente gRPC global
var grpcClient pb.WarReportServiceClient

func main(){
	//leer la direccion del gRPC server desde la variable de entorno
	grpcServerAddr := os.Getenv("GRPC_SERVER_ADDR")

	if grpcServerAddr == "" {
		grpcServerAddr = "localhost:50051" // valor por defecto
	}

	//conectar al gRPC server

	conn, err := grpc.Dial(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar al gRPC server: %v", err)
	}
	defer conn.Close()

	grpcClient = pb.NewWarReportServiceClient(conn)

	fmt.Println("gRPC Client conectado al servidor en", grpcServerAddr)

	//servidor http para recibir del api rest

	http.HandleFunc("/send-report", sendReportHandler)
	fmt.Println("HTTP Server escuchando en :8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("Error al iniciar el servidor HTTP: %v", err)
	}

}

func sendReportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}


	//decodificar el json del body
	var report WarReportJSON
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		http.Error(w, "Error al decodificar el JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	//validar pais
	countryEnum, ok := countryMap[report.Country]
	if !ok {
		http.Error(w, "País no válido: "+report.Country, http.StatusBadRequest)
		return
	}

	//llamar al gRPC server con el reporte
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := grpcClient.SendReport(ctx, &pb.WarReportRequest{
		Country:         countryEnum,
		WarplanesInAir:  report.WarplanesInAir,
		WarshipsInWater: report.WarshipsInWater,
		Timestamp:       report.Timestamp,
	})

	if err != nil {
		http.Error(w, "Error al enviar el reporte al gRPC server: "+err.Error(), http.StatusInternalServerError)
		return
	}


	//responder al api rest
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": resp.Status,
	})
}