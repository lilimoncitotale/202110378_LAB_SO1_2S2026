package main

import (
	
	"os"
	"bufio"
)

func leer_archivo(ruta string)(string, error){
	archivo, err := os.Open(ruta)
	if err != nil {
		return "", err
	}
	defer archivo.Close()

	var contenido string
	scanner := bufio.NewScanner(archivo)

	//leer linea por línea
	for scanner.Scan(){
		contenido += scanner.Text() + "\n"
	}

	//validar si hubo error al leer el archivo
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return contenido, nil
}