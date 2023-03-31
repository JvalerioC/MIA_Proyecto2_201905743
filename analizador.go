package main

import (
	"fmt"
	"regexp"
	"strings"
)

func analizador(entrada string) {
	//expresion regular para separar por espacios sin incluir los espacios de que estan entre comillas
	regex := regexp.MustCompile(`(\S+"[^"]+"|\S+)`)
	//se buscan todas las coincidencias de la entrada
	matches := regex.FindAllStringSubmatch(entrada, -1)

	// Creamos un slice para almacenar los elementos sin comillas dobles
	result := make([]string, len(matches))
	//sustituimos cada valor del arreglo las commilas y se asignan al otro
	for i, match := range matches {
		element := strings.ReplaceAll(match[0], `"`, "")
		//fmt.Println(element)
		result[i] = element
	}
	result[0] = strings.ToLower(result[0])

	//aqui podemos empezar a analizar el primer valor de la entrada
	if result[0] == "mkdisk" {
		mkdisk(result[1:])
	} else if result[0] == "rmdisk" {
		rmdisk(result)
	} else if result[0] == "fdisk" {
		fdisk(result)
	} else if result[0] == "mount" {
		mount(result)
	} else {
		fmt.Println("El comando ingresado no es valido")
		return
	}
}
