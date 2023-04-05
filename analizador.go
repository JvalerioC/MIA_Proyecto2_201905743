package main

import (
	"fmt"
	"regexp"
	"strings"
)

func analizador(entrada string) {
	//expresion regular para separar por espacios sin incluir los espacios de que estan entre comillas
	regex := regexp.MustCompile(`\S+="[^"]*"|\S+`)
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
	//se pasa a minusculas la primera posicion de la entrada para saber el comando
	result[0] = strings.ToLower(result[0])

	//aqui podemos empezar a analizar el primer valor de la entrada
	if result[0] == "mkdisk" {
		mkdisk(result[1:])
	} else if result[0] == "rmdisk" {
		rmdisk(result[1:])
	} else if result[0] == "fdisk" {
		fdisk(result[1:])
	} else if result[0] == "mount" {
		mount(result[1:])
	} else if result[0] == "rep" {
		rep(result[1:])
	} else if result[0] == "mkfs" {
		mkfs(result[1:])
	} else if result[0] == "login" {
		login(result[1:])
	} else if result[0] == "logout" {
		logout()
	} else if result[0] == "mkgrp" {
		mkgrp(result[1:])
	} else if result[0] == "rmgrp" {
		rmgrp(result[1:])
	} else if result[0] == "mkusr" {
		mkusr(result[1:])
	} else if result[0] == "rmusr" {
		rmusr(result[1:])
	} else {
		fmt.Println("El comando ingresado no es valido")
		return
	}
}
