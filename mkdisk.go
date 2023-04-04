package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// funcion para el comando mkdir
func mkdisk(params []string) {
	var size int
	var path string
	var unit byte
	var fit byte

	for i := 0; i < len(params); i++ {
		array := strings.Split(params[i], "=")
		param := strings.ToLower(array[0])
		if param == ">size" {
			value, err := strconv.Atoi(array[1])
			if err != nil {
				fmt.Println("Error, el valor ingresado para el parametro size no es valido")
				return
			}
			if value <= 0 {
				fmt.Println("Error, el valor ingresado para el parametro size debe ser mayor a 0")
				return
			}
			size = value
		} else if param == ">unit" {
			value := array[1]
			if value == "M" {
				unit = 'M'
			} else if value == "K" {
				unit = 'K'
			} else {
				fmt.Println("Error, el valor ingresado para el parametro unit no es valido")
				return
			}
		} else if param == ">path" {
			path = array[1]
		} else if param == ">fit" {
			value := array[1]
			if value == "BF" {
				fit = 'B'
			} else if value == "FF" {
				fit = 'F'
			} else if value == "WF" {
				fit = 'W'
			} else {
				fmt.Println("Error, el valor ingresado para el parametro fit no es valido")
				return
			}
		} else {
			fmt.Println("El parametro ingresado no es valido")
			return
		}
	}
	//ya obtenidos los parametros se hacen validaciones
	if size == 0 || path == "" {
		fmt.Println("Error, parametro obligatorio vacio")
		return
	}
	if unit == 0 {
		unit = 'M'
	}
	if fit == 0 {
		fit = 'F'
	}
	//creamos las carpetas y subcarpetas si estas no existen
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)

	if err != nil {
		fmt.Println("Error creando el directorio")
		return
	}
	//verificamos si el archivo existe o no
	_, err1 := os.Stat(path)
	if os.IsNotExist(err1) {
		//fmt.Println()
	} else {
		fmt.Println("Error, el disco ya existe")
		return
	}
	//creamos el archivo
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creando el archivo")
		return
	}
	//se calcula el tamaño del archivo
	finalSize := size * 1024
	if unit == 'M' {
		finalSize = finalSize * 1024
	}
	data := make([]byte, finalSize)

	//se escribe en el archivo el arreglo de bytes
	_, err2 := file.Write(data)
	if err2 != nil {
		fmt.Println("Error al escribir en el archivo")
		return
	}
	file.Close()
	//creado el archivo crearemos el mbr y se escribira en el archivo

	//creamos el mbr
	mbr := MBR{}
	//se copia el tamaño al atributo del mbr
	num_to_string := strconv.Itoa(finalSize)
	copy(mbr.Mbr_tamano[:], num_to_string)
	//se copia el fit al atributo del mbr
	mbr.Dsk_fit[0] = fit
	//se crea la fecha de creacion y se asigna
	date := time.Now()
	formatted := date.Format("02/01/2006 15:04:05")
	copy(mbr.Mbr_fecha_creacion[:], formatted)
	//se crea un numero aleatorio entre 0 y 9999
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source).Intn(9999)
	//se copia el numero aleatorio al atributo del mbr
	num_to_string = strconv.Itoa(random)
	copy(mbr.Mbr_dsk_signature[:], num_to_string)

	//se llama la funcion encargada de escribir el mbr
	flag := write_MBR(mbr, path)
	if flag {
		fmt.Println("Disco creado exitosamente")
		//intentaremos recuperar el mbr para saber si se esta escribiendo bien
		/* mbr2 := read_MBR(path)
		fmt.Println("Tamaño del disco: ", string(mbr2.Mbr_tamano[:]))
		fmt.Println("Fit del disco: ", string(mbr2.Dsk_fit[:]))
		fmt.Println("Fecha de creacion: ", string(mbr2.Mbr_fecha_creacion[:]))
		fmt.Println("Disk Signature: ", string(mbr2.Mbr_dsk_signature[:])) */

	} else {
		fmt.Println("Error al crear el disco")
	}

}
