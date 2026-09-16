package main

import (
	"fmt"
	"os/exec"
)

func cargar_modulo(rutaModulo string) error{
	cmd := exec.Command("sudo", "insmod", rutaModulo)
	salida, err :=  cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error al cargar el módulo: %v, salida: %s", err, string(salida))
	}
	return nil
}

func descargar_Modulo(nombreModulo string)error{
	cmd := exec.Command("sudo", "rmmod", nombreModulo)
	salida, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error al descargar el módulo: %v, salida: %s", err, string(salida))
	}
	return nil
}