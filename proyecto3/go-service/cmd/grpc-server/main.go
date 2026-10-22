package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "go-service/proto"
)

type server struct {
	pb.UnimplementedWarReportServiceServer
}

// Struct para enviar al RabbitMQ Writer
type WarReportJSON struct {
	Country         string `json:"country"`
	WarplanesInAir  int32  `json:"warplanes_in_air"`
	WarshipsInWater int32  `json:"warships_in_water"`
	Timestamp       string `json:"timestamp"`
}

// Mapeo de enum a string
var countryToString = map[pb.Countries]string{
	pb.Countries_usa: "USA",
	pb.Countries_rus: "RUS",
	pb.Countries_chn: "CHN",
	pb.Countries_esp: "ESP",
	pb.Countries_gtm: "GTM",
}

func (s *server) SendReport(ctx context.Context, req *pb.WarReportRequest) (*pb.WarReportResponse, error) {
	fmt.Printf("=== gRPC Server: reporte recibido ===\n")
	fmt.Printf("País: %v\n", req.Country)
	fmt.Printf("Aviones en aire: %d\n", req.WarplanesInAir)
	fmt.Printf("Barcos en mar: %d\n", req.WarshipsInWater)
	fmt.Printf("Timestamp: %s\n", req.Timestamp)
	fmt.Printf("====================================\n")

	// Leer la URL del RabbitMQ Writer
	writerURL := os.Getenv("RABBITMQ_WRITER_URL")
	if writerURL == "" {
		writerURL = "http://localhost:8083/publish" // Fallback local
	}

	// Convertir el enum a string
	countryStr := countryToString[req.Country]
	if countryStr == "" {
		countryStr = "UNKNOWN"
	}

	// Crear el JSON para el writer
	reportJSON := WarReportJSON{
		Country:         countryStr,
		WarplanesInAir:  req.WarplanesInAir,
		WarshipsInWater: req.WarshipsInWater,
		Timestamp:       req.Timestamp,
	}

	// Serializar y enviar
	body, err := json.Marshal(reportJSON)
	if err != nil {
		log.Printf("Error al serializar: %v", err)
		return &pb.WarReportResponse{Status: "error"}, nil
	}

	resp, err := http.Post(writerURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error al llamar al RabbitMQ Writer: %v", err)
		return &pb.WarReportResponse{Status: "error"}, nil
	}
	defer resp.Body.Close()

	fmt.Printf("RabbitMQ Writer respondió con estado: %d\n", resp.StatusCode)

	return &pb.WarReportResponse{Status: "ok"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterWarReportServiceServer(s, &server{})
	reflection.Register(s)

	fmt.Println("gRPC Server escuchando en :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}