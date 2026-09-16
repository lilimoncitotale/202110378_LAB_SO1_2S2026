package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
	"fmt"
)

func main() {
	rutaKO := "/home/limon/proyecto2/kernel_module/continfo.ko"
	nombreModulo := "continfo"

	//cargamos el modulo
	fmt.Println("Cargando módulo de kernel...")
	if err := cargar_modulo(rutaKO); err != nil {
		log.Fatalf("Error al cargar el módulo: %v", err)
	}
	fmt.Println("Módulo de kernel cargado exitosamente.")

	//garantizar la descarga al salir
	defer func(){
		fmt.Println("Descargando módulo de kernel...")
		if err := descargar_Modulo(nombreModulo); err != nil {
			log.Printf("Error al descargar el módulo: %v", err)
		}else{
			fmt.Println("Módulo de kernel descargado exitosamente.")
		}
	}()
    log.Println("=== Daemon de telemetría iniciado ===")
    
    // Aquí irá la lógica del daemon
	resultado, err := leer_archivo("/proc/continfo_pr2_so1_202110378")
	if err != nil {
		log.Fatalf("Error al leer el archivo: %v", err)
	}

	fmt.Println("datos del kernel obtenidos")
	fmt.Print(resultado)
    // 1. Inicializar servicios (Grafana, Valkey)
    // 2. Cargar módulo de kernel
    // 3. Configurar cronjob
    // 4. Cargar programa eBPF
    // 5. Loop principal
    
    // Canal para señales del sistema
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    // Esperar señal para terminar
    <-sigChan
    log.Println("=== Daemon detenido ===")
}