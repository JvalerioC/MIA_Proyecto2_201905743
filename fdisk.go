package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

// funcion para el comando mkdir
func fdisk(params []string) {
	var size int
	var path string
	var unit byte
	var fit byte
	var type_ byte
	var name string

	//se recuperan los parametros y se asigna como corresponde
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
			} else if value == "B" {
				unit = 'B'
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
		} else if param == ">type" {
			value := array[1]
			if value == "P" {
				type_ = 'P'
			} else if value == "E" {
				type_ = 'E'
			} else if value == "L" {
				type_ = 'L'
			} else {
				fmt.Println("Error, el valor ingresado para el parametro type no es valido")
				return
			}
		} else if param == ">name" {
			name = array[1]
		} else {
			fmt.Println("Error, el parametro ingresado no es valido")
			return
		}
	}
	//se validan los parametros obligatorios
	if size == 0 || path == "" || name == "" {
		fmt.Println("Error, parametro obligatorio vacio")
		return
	}
	//se asignan los valores por default a los parametros opcionales, si estos estan vacios
	if unit == 0 {
		unit = 'K'
	}
	if fit == 0 {
		fit = 'W'
	}
	if type_ == 0 {
		type_ = 'P'
	}
	//se verfica que el disco exista si no existe retorna
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("Error, el disco no existe")
		return
	}
	//recuperamos el MBR
	mbr, flag := read_MBR(path)
	if !flag {
		fmt.Println("Error, no se pudo recuperar el mbr")
		return
	}

	//se verifica que el nombre de la particion no este en uso, 4 particiones
	for i := 0; i < 4; i++ {
		if name == string(mbr.Mbr_Partition[i].Part_name[:len(name)]) && mbr.Mbr_Partition[i].Part_name[len(name)] == 0 {
			fmt.Println("Error, el nombre dela particion ya esta en uso")
			return
		}
	}
	//se verifica que el nombre de la particion en las particiones logicas
	for i := 0; i < 4; i++ {
		//verificamos que exista extendida
		if mbr.Mbr_Partition[i].Part_type[0] == 'E' {
			//recuperamos la particion extendida
			extended, flag := read_ebr(path, mbr.Mbr_Partition[i].Part_start)
			if !flag {
				return
			}
			//recorremos las particiones logicas
			for extended.Part_next[0] != 0 {
				//se valida el nombre
				if name == string(extended.Part_name[:len(name)]) && extended.Part_name[len(name)] == 0 {
					fmt.Println("Error, el nombre dela particion ya esta en uso")
					return
				}
				//se lee la siguiente particion logica
				extended, flag = read_ebr(path, extended.Part_next)
				if !flag {
					fmt.Println("Ocurrio un error, ebr")
					return
				}

			}
		}
	}
	//en este punto el nombre de la particion no se ha registrado
	//ahora se valida cuantas particiones hay
	countPrimary := 0
	countExtended := 0
	countPart := 0
	for i := 0; i < 4; i++ {
		if mbr.Mbr_Partition[i].Part_type[0] == 'P' {
			countPrimary++
		} else if mbr.Mbr_Partition[i].Part_type[0] == 'E' {
			countExtended++
		}
	}
	countPart = countPrimary + countExtended

	//se valida que no se haya llegado al limite de particiones y que la particion no sea logica, si es asi retorna
	if countPart == 4 && type_ != 'L' {
		fmt.Println("Error, ya hay 4 particiones en el disco")
		return
	}
	//ahora se ve que tipo de particion es
	if type_ == 'P' {
		createPartition(size, unit, type_, path, name, fit)
	} else if type_ == 'E' && countExtended == 0 {
		createPartition(size, unit, type_, path, name, fit)
	} else if type_ == 'E' && countExtended == 1 {
		fmt.Println("Error, ya hay una particion extendida")
		return
	} else if type_ == 'L' && countExtended == 0 {
		fmt.Println("Error, no hay particiones extendidas")
		return
	} else if type_ == 'L' && countExtended == 1 {
		createLogic(size, unit, type_, path, name, fit)
	} else {
		fmt.Println("Error, no se porque llego aqui")
		return
	}
}

func createLogic(size int, unit byte, type_ byte, path string, name string, fit byte) {
	//recuperamos el MBR
	mbr, flag := read_MBR(path)
	if !flag {
		fmt.Println("Error, no se pudo recuperar el mbr")
		return
	}
	flag_e := false
	for i := 0; i < 4; i++ {
		//verificamos que exista extendida, para trabajar en ella
		if mbr.Mbr_Partition[i].Part_type[0] == 'E' {
			//se guarda el tamano de la particion en una variable
			c1 := string(mbr.Mbr_Partition[i].Part_size[:])
			c1 = strings.TrimRight(c1, "\x00")
			part_size, err := strconv.Atoi(c1)
			if err != nil {
				fmt.Println("Error al convertir a entero, tamaño del particion")
				return
			}
			//se calcula el tamano de la particion
			finalsize := 0
			if unit == 'K' {
				finalsize = size * 1024
			} else if unit == 'M' {
				finalsize = size * 1024 * 1024
			} else if unit == 'B' {
				finalsize = size
			}

			//recuperamos la particion extendida
			extended, flag := read_ebr(path, mbr.Mbr_Partition[i].Part_start)
			if !flag {
				fmt.Println("Ocurrio un error con el ebr")
				return
			}
			//recorremos las particiones logicas
			for extended.Part_next[0] != 0 {
				//se lee la siguiente particion logica
				extended, flag = read_ebr(path, extended.Part_next)
				if !flag {
					fmt.Println("Ocurrio un error con el ebr")
					return
				}

			}
			position := 0
			position_last := 0
			if extended.Part_size[0] == 0 {
				//se calcula la posicion de la particion
				c := string(mbr.Mbr_Partition[i].Part_start[:])
				c = strings.TrimRight(c, "\x00")
				d, err := strconv.Atoi(c)
				if err != nil {
					fmt.Println("Error al convertir a entero, posicion 1")
					return
				}
				position = d
				if (position + finalsize) > part_size {
					fmt.Println("Error, no es posible ingresar particion, no hay espacio disponible")
					return
				}
			} else {
				//se calcula la posicion de la particion
				c := string(extended.Part_start[:])
				c2 := string(extended.Part_size[:])
				c = strings.TrimRight(c, "\x00")
				c2 = strings.TrimRight(c2, "\x00")
				d, err := strconv.Atoi(c)
				if err != nil {
					fmt.Println("Error al convertir a entero, posicion 2")
					return
				}
				d2, err := strconv.Atoi(c2)
				if err != nil {
					fmt.Println("Error al convertir a entero, tamaño ")
					return
				}
				position = d + d2
				position_last = d
				if (position + finalsize) > part_size {
					fmt.Println("Error, no es posible ingresar particion, no hay espacio disponible")
					return
				}
			}
			//se crea el ebr
			ebr := EBR{}
			//se asignan los valores
			copy(ebr.Part_name[:], name)
			ebr.Part_fit[0] = fit
			ebr.Part_status[0] = 0
			num_to_string := strconv.Itoa(finalsize)
			copy(ebr.Part_size[:], num_to_string)
			num_to_string = strconv.Itoa(position)
			copy(ebr.Part_start[:], num_to_string)
			flag1 := write_ebr(ebr, path, position)
			if !flag1 {
				fmt.Println("Error, no se pudo escribir el ebr")
				return
			}
			//se actualiza el ebr anterior
			if position_last != 0 {
				num_to_string = strconv.Itoa(position)
				copy(extended.Part_next[:], num_to_string)
				flag1 = write_ebr(extended, path, position_last)
				if !flag1 {
					fmt.Println("Error, no se pudo escribir el ebr")
					return
				}
			}
			// si todo termino bien
			flag_e = true

		}
	}
	if flag_e {
		fmt.Println("Particion logica creada exitosamente")
	}
}

// funcion para escribir un ebr en una particion
func write_ebr(ebr EBR, path string, position int) bool {
	disk, err := os.OpenFile(path, os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error abriendo el archivo")
		disk.Close()
		return false
	}
	_, err1 := disk.Seek(int64(position), io.SeekStart)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero")
		disk.Close()
		return false
	}
	// Se escribe el mbr en el archivo
	err = binary.Write(disk, binary.LittleEndian, ebr)
	if err != nil {
		fmt.Println("Error al escribir el ebr")
		disk.Close()
		return false
	}
	disk.Close()
	return true
}

// funcion para crear particiones primarias y extendidas
func createPartition(size int, unit byte, type_ byte, path string, name string, fit byte) {
	//se crea la particion
	partition := Partition{}
	//se asignan los valores a la particion
	copy(partition.Part_name[:], name)
	partition.Part_type[0] = type_
	partition.Part_fit[0] = fit
	partition.Part_status[0] = '0'
	finalsize := 0
	if unit == 'K' {
		finalsize = size * 1024
	} else if unit == 'M' {
		finalsize = size * 1024 * 1024
	} else if unit == 'B' {
		finalsize = size
	}
	num_to_string := strconv.Itoa(finalsize)
	copy(partition.Part_size[:], num_to_string)

	//se recupera el MBR
	mbr, flag := read_MBR(path)
	if !flag {
		fmt.Println("Error, no se pudo recuperar el mbr")
		return
	}
	//se guarda el tamano del disco en una variable
	c1 := string(mbr.Mbr_tamano[:])
	c1 = strings.TrimRight(c1, "\x00")
	disk_size, err := strconv.Atoi(c1)
	if err != nil {
		fmt.Println("Error al convertir a entero, tamaño del disco")
		return
	}

	//se busca la posicion donde se pueda ingresar
	position := int(unsafe.Sizeof(mbr))
	for i := 0; i < 4; i++ {
		//se verifica que la particion este libre
		if mbr.Mbr_Partition[i].Part_size[0] == 0 {
			//se escribe un ebr en la primera posicion de una particion extendida
			if type_ == 'E' {
				ebr := EBR{}
				file, err := os.OpenFile(path, os.O_RDWR, 0664)
				if err != nil {
					fmt.Println("Error al abrir el archivo, ebr")
					file.Close()
					return
				}
				_, err1 := file.Seek(int64(position), io.SeekStart)
				if err1 != nil {
					fmt.Println("Error posicionando el puntero")
					return
				}
				// Se escribe el mbr en el archivo
				err = binary.Write(file, binary.LittleEndian, ebr)
				if err != nil {
					fmt.Println("Error al escribir el ebr")
					return
				}
				file.Close()

			}
			num_to_string := strconv.Itoa(position)
			copy(partition.Part_start[:], num_to_string)
			//se asigna la particion
			mbr.Mbr_Partition[i] = partition
			//se escribe el MBR
			flag1 := write_MBR(mbr, path)
			if !flag1 {
				fmt.Println("Error, no se pudo recuperar el mbr")
				return
			}

			break

		} else {
			c := string(mbr.Mbr_Partition[i].Part_size[:])
			c = strings.TrimRight(c, "\x00")
			d, err := strconv.Atoi(c)
			if err != nil {
				fmt.Println("Error al convertir a entero, posicion 3")
				return
			}
			position += d
			if (position + finalsize) > disk_size {
				fmt.Println("Error, no es posible ingresar particion, no hay espacio disponible")
				return
			}

		}
	}

	fmt.Println("Particion creada exitosamente")
}

// funcion para leer particiones extendidas
func read_ebr(path string, position [10]byte) (EBR, bool) {
	ebr := EBR{}
	disk, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		disk.Close()
		return ebr, false
	}
	//asignamos a c el string con el numero
	c := string(position[:])
	//se elimina los caracteres nulos
	c = strings.TrimRight(c, "\x00")
	//se convierte a entero
	d, err := strconv.Atoi(c)
	if err != nil {
		fmt.Println("Error al convertir a entero ")
		return ebr, false
	}
	//se posiciona el puntero en la posicion del disco
	_, err1 := disk.Seek(int64(d), io.SeekStart)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero ")
		return ebr, false
	}
	//se lee la particion extendida
	err2 := binary.Read(disk, binary.LittleEndian, &ebr)
	if err2 != nil {
		fmt.Println("Error leyendo la particion extendida ")
		return ebr, false
	}
	disk.Close()
	return ebr, true
}
