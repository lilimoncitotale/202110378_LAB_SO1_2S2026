package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"context"
)

type WarReport struct {
	Country         string `json:"country"`
	WarplanesInAir  int    `json:"warplanes_in_air"`
	WarshipsInWater int    `json:"warships_in_water"`
	Timestamp       string `json:"timestamp"`
}

var ctx = context.Background()

func main() {
	// 1. Conectar a RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://admin:admin123@localhost:5672/"
	}

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

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error al abrir canal: %v", err)
	}
	defer ch.Close()

	// Declarar la cola (por si no existe)
	q, err := ch.QueueDeclare(
		"war_reports",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Error al declarar cola: %v", err)
	}

	// 2. Conectar a Valkey (Redis)
	valkeyAddr := os.Getenv("VALKEY_ADDR")
	if valkeyAddr == "" {
		valkeyAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: valkeyAddr,
	})
	// Probar conexión
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("No se pudo conectar a Valkey: %v", err)
	}
	fmt.Println("Conectado a RabbitMQ y Valkey")

	// 3. Consumir mensajes
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Error al consumir: %v", err)
	}

	fmt.Println("Esperando mensajes...")
	for msg := range msgs {
		var report WarReport
		if err := json.Unmarshal(msg.Body, &report); err != nil {
			log.Printf("Error al parsear JSON: %v", err)
			continue
		}
		fmt.Printf("Reporte recibido: %+v\n", report)
		processReport(rdb, report)
	}
}

func processReport(rdb *redis.Client, report WarReport) {
	// 1. Actualizar contadores globales (max, min, total)
	updateGlobalStats(rdb, "max_warplanes", int64(report.WarplanesInAir), true)
	updateGlobalStats(rdb, "min_warplanes", int64(report.WarplanesInAir), false)
	updateGlobalStats(rdb, "max_warships", int64(report.WarshipsInWater), true)
	updateGlobalStats(rdb, "min_warships", int64(report.WarshipsInWater), false)
	rdb.Incr(ctx, "global:total_reports")

	// 2. Actualizar contadores por país
	countryKey := "country:" + report.Country
	rdb.IncrBy(ctx, countryKey+":total_warplanes", int64(report.WarplanesInAir))
	rdb.IncrBy(ctx, countryKey+":total_warships", int64(report.WarshipsInWater))
	rdb.Incr(ctx, countryKey+":report_count")
	rdb.Set(ctx, countryKey+":last_warplanes", report.WarplanesInAir, 0)
	rdb.Set(ctx, countryKey+":last_warships", report.WarshipsInWater, 0)
	rdb.Set(ctx, countryKey+":last_timestamp", report.Timestamp, 0)

	// 3. Frecuencia de valores para la moda
	rdb.HIncrBy(ctx, "frequency:warplanes", strconv.Itoa(report.WarplanesInAir), 1)
	rdb.HIncrBy(ctx, "frequency:warships", strconv.Itoa(report.WarshipsInWater), 1)
}

func updateGlobalStats(rdb *redis.Client, key string, value int64, isMax bool) {
	// Obtener el valor actual
	currentStr, err := rdb.Get(ctx, "global:"+key).Result()
	if err == redis.Nil {
		// No existe, crear
		rdb.Set(ctx, "global:"+key, value, 0)
		return
	} else if err != nil {
		log.Printf("Error al obtener %s: %v", key, err)
		return
	}

	current, _ := strconv.ParseInt(currentStr, 10, 64)
	if isMax && value > current {
		rdb.Set(ctx, "global:"+key, value, 0)
	} else if !isMax && value < current {
		rdb.Set(ctx, "global:"+key, value, 0)
	}
}