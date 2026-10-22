use axum::{
    routing::post,
    Json, Router,
    http::StatusCode,
};
use serde::{Deserialize, Serialize};
use std::net::SocketAddr;
use tracing_subscriber;

// Estructura del reporte militar (lo que envía Locust)
#[derive(Debug, Deserialize, Serialize)]
struct WarReport {
    country: String,
    warplanes_in_air: i32,
    warships_in_water: i32,
    timestamp: String,
}

// Respuesta que devuelve la API
#[derive(Serialize)]
struct ApiResponse {
    status: String,
}

// Handler para POST /report
async fn receive_report(Json(report): Json<WarReport>) -> (StatusCode, Json<ApiResponse>) {
    println!("Reporte recibido: {:?}", report);

    //inicializamos el cliente asincrono de reqwest
    let client = reqwest::Client::new();

    let go_url = std::env::var("GO_SERVICE_URL").unwrap_or_else(|_| "http://localhost:8081/war-report".to_string());

    //reenviamos el JSON a la API de destino
    let response = client
        .post(&go_url)
        .json(&report)
        .send()
        .await;

    match response {
        Ok(res) => {
            println!("Go respondio con estado: {}", res.status());
            (StatusCode::OK, Json(ApiResponse { status: "Reporte reenviado exitosamente".into() }))
        }
        Err(err) => {
            eprintln!("Error al reenviar el reporte: {}", err);
            (StatusCode::INTERNAL_SERVER_ERROR, Json(ApiResponse { status: "Error al reenviar el reporte".into() }))
        }
    }
}


#[tokio::main]
async fn main() {
    // Inicializar logs
    tracing_subscriber::fmt::init();

    // Definir rutas
    let app = Router::new()
        .route("/report", post(receive_report));

    // Escuchar en el puerto 8080
    let addr = SocketAddr::from(([0, 0, 0, 0], 8080));
    println!("Servidor escuchando en {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}