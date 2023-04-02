package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// funcion para el comando mkdir
func mount(params []string) {
	var path string
	var name string
	for i := 0; i < len(params); i++ {
		array := strings.Split(params[i], "=")
		param := strings.ToLower(array[0])
		if param == ">path" {
			path = array[1]
		} else if param == ">name" {
			name = array[1]
		} else {
			fmt.Println("Error, el parametro ingresado no es valido")
			return
		}
	}
	if name == "" || path == "" {
		fmt.Println("Error, parametro obligatorio vacio")
		return
	}
	//verificamos si el disco existe o no
	_, err1 := os.Stat(path)
	if os.IsNotExist(err1) {
		fmt.Println("Error, el disco no existe")
		return
	}
	mbr, flag := read_MBR(path)

	if !flag {
		fmt.Println("Error, no se pudo leer el MBR")
		return
	}
	//verificamos si el disco ya esta montado
	count_status_disk := 0
	i_mount := itemMount{}
	for i := 0; i < 4; i++ {
		name_ := string(mbr.Mbr_Partition[i].Part_name[:])
		if name_ == name {
			if mbr.Mbr_Partition[i].Part_status[0] == '1' {
				fmt.Println("Error, la particion ya esta montada")
				return
			}
			i_mount.Part = mbr.Mbr_Partition[i]
			mbr.Mbr_Partition[i].Part_status[0] = '1'
		}
		if mbr.Mbr_Partition[i].Part_status[0] == '1' {
			count_status_disk++
		}
	}
	// se valida que si se recupero bien la particion
	if i_mount.Part.Part_size[0] == '0' {
		fmt.Println("no existe una particion  con el nombre ingresado")
		return
	}
	letter := "A"
	if count_status_disk == 0 {
		letter = "A"
	} else if count_status_disk == 1 {
		letter = "B"
	} else if count_status_disk == 2 {
		letter = "C"
	} else if count_status_disk == 3 {
		letter = "D"
	}

	count_disk := 0
	new_flag := false
	for _, item := range PartMount {
		if item.Path == path {
			count_disk = item.Number
			new_flag = true
		}
	}
	if !new_flag {
		disc_counter++
		count_disk = disc_counter
	}
	//reescribimos el mbr
	flagg := write_MBR(mbr, path)
	if !flagg {
		fmt.Println("Error, no se pudo escribir el MBR, mount")
		return
	}
	//se asignan los valores
	i_mount.Id = "43" + strconv.Itoa(count_disk) + letter
	i_mount.Path = path
	i_mount.Number = count_disk
	PartMount = append(PartMount, i_mount)
	fmt.Println("Particion montada exitosamente")
	fmt.Println()

	//se imprimen las particiones montadas
	fmt.Println(" Particiones Montadas ")
	fmt.Println("----------------------")
	for _, item := range PartMount {
		fmt.Println(item.Id)
	}
}
