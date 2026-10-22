package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Struct para recibir el JSON del gRPC Server
type WarReportJSON struct {
	Country         string `json:"country"`
	WarplanesInAir  int32  `json:"warplanes_in_air"`
	WarshipsInWater int32  `json:"warships_in_water"`
	Timestamp       string `json:"timestamp"`
}

var channel *amqp.Channel
var queueName = "war_reports"

func main() {
	// Leer la URL de RabbitMQ desde variable de entorno
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://admin:admin123@localhost:5672/" // Fallback local
	}

	// Conectar a RabbitMQ con reintentos
	var conn *amqp.Connection
	var err error
	for i := 1; i <= 10; i++ {
		conn, err = amqp.Dial(rabbitURL)
		if err == nil {
			break
		}
		fmt.Printf("Intento %d: error al conectar a RabbitMQ: %v\n", i, err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("No se pudo conectar a RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Crear canal
	channel, err = conn.Channel()
	if err != nil {
		log.Fatalf("Error al crear canal: %v", err)
	}
	defer channel.Close()

	// Declarar la cola
	_, err = channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Error al declarar la cola: %v", err)
	}

	fmt.Println("RabbitMQ Writer conectado y cola declarada")

	// Servidor HTTP para recibir del gRPC Server
	http.HandleFunc("/publish", publishHandler)
	fmt.Println("RabbitMQ Writer (HTTP) escuchando en :8083")
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatalf("Error al servir HTTP: %v", err)
	}
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar JSON
	var report WarReportJSON
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	// Serializar de nuevo a JSON para publicar
	body, err := json.Marshal(report)
	if err != nil {
		http.Error(w, "Error al serializar JSON", http.StatusInternalServerError)
		return
	}

	// Publicar en RabbitMQ
	err = channel.Publish(
		"",        // exchange (vacío = default)
		queueName, // routing key (nombre de la cola)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Error al publicar en RabbitMQ: %v", err)
		http.Error(w, "Error al publicar en RabbitMQ", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Publicado en RabbitMQ: %s\n", string(body))

	// Responder al gRPC Server
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"published"}`))
}