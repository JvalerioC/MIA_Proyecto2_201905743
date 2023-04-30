package main

import (
	"fmt"
	"os"
	"strings"
)

// funcion para el comando mkdir
func rmdisk(params []string) {
	var path string
	for i := 0; i < len(params); i++ {
		array := strings.Split(params[i], "=")
		param := strings.ToLower(array[0])
		if param == ">path" {
			path = array[1]
		} else {
			fmt.Println("Error, el parametro ingresado no es valido")
			cadRespuesta += "Error, el parametro ingresado no es valido\n"
			return
		}
	}
	//se verifica que el path no este vacio
	if path == "" {
		fmt.Println("Error, parametro obligatorio vacio")
		cadRespuesta += "Error, parametro obligatorio vacio\n"
		return
	}
	//se verfica que el disco exista si no existe retorna
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("Error, el disco no existe")
		cadRespuesta += "Error, el disco no existe\n"
		return
	}

	/* //se muestra un mensaje de confirmacion y se valida el caracter ingresado
	fmt.Println("Confirmar que desea eliminar el disco? (Y/N): ")
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "Y" {
		fmt.Println("No se eliminara el disco")
		return
	} */

	//se elimina el disco
	err = os.Remove(path)
	if err != nil {
		fmt.Println("Error, no se pudo eliminar el disco")
		cadRespuesta += "Error, no se pudo eliminar el disco\n"
		return
	}
	fmt.Println("Disco eliminado correctamente")
	cadRespuesta += "Disco eliminado correctamente\n"

}
