package main


import(
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	"os"
)

//definimos la estructura 

type WarReportRequest struct {

	Country string `json:"country"`
	WarplanesInAir int `json:"warplanes_in_air"`
	WarshipsInWater int `json:"warships_in_water"`
	Timestamp string `json:"timestamp"`
}

//funcion para manejar la solicitud POST
func warReportHandler(w http.ResponseWriter, r *http.Request){
	//validamos que sea un metodo post
	fmt.Println("Solicitud POST recibida en /war-report")
	if r.Method != http.MethodPost{
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	//decodificamos el cuerpo de la solicitud en la estructura WarReportRequest
	var report WarReportRequest

	//decodificamos el cuerpo de la solicitud en la estructura WarReportRequest
	err := json.NewDecoder(r.Body).Decode(&report)
	if err != nil{
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}


	//procesamos los datos recibidos (en este caso, simplemente los imprimimos en la consola)
	fmt.Printf("Reporte de guerra recibido: %+v\n", report)

	//leer la url del gRPC client desde la variable de entorno
	grpcClientAddr := os.Getenv("GRPC_CLIENT_URL")
	if grpcClientAddr == "" {
		grpcClientAddr = "http://localhost:8082/send-report" // valor por defecto
	}

	//reenviar el reporte a la url del gRPC client
	body, err := json.Marshal(report)
	if err != nil {
		http.Error(w, "Error al codificar el reporte", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(grpcClientAddr, "application/json", bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Error al enviar el reporte al gRPC client", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	//leer la respuesta del gRPC client
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Respuesta del gRPC client: %s\n", string(respBody))

	//devolver la respuesta a Rust
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))



}

func main(){
	http.HandleFunc("/war-report", warReportHandler)
	fmt.Println("Servidor escuchando en http://localhost:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil{
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
	}
}