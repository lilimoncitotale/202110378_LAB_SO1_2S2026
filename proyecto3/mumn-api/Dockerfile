# para ahorrar espacio se utiliza el multiStage build, primero se construye la imagen con rust y luego se copia el binario a una imagen más ligera

# Use the official Rust image from Docker Hub
FROM rust:latest AS builder

# Set the working directory inside the container
WORKDIR /app
# Copy the project files into the container
COPY Cargo.toml Cargo.lock ./
COPY src ./src

# Build the project
RUN cargo build --release

# Etapa 2: Ejecutar
FROM debian:bookworm-slim

WORKDIR /app

# Instalar librerías necesarias
RUN apt-get update && apt-get install -y \
    libssl3 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copiar el binario
COPY --from=builder /app/target/release/mumn-api ./mumn-api

EXPOSE 8080
CMD ["./mumn-api"]