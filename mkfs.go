package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func mkfs(params []string) {
	var id string
	var type_ string

	// se recorren los parametros
	for i := 0; i < len(params); i++ {
		array := strings.Split(params[i], "=")
		param := strings.ToLower(array[0])
		if param == ">id" {
			id = array[1]
		} else if param == ">type" {
			value := array[1]
			if value == "Full" {
				type_ = "Full"
			} else {
				fmt.Println("Error, el valor ingresado para el parametro type no es valido")
				return
			}
		} else {
			fmt.Println("Error, el parametro ingresado no es valido")
			return
		}
	}
	//se validan los parametros
	if id == "" {
		fmt.Println("Error, parametro obligatorio vacio")
		return
	}
	if type_ == "" {
		type_ = "Full"
	}
	//se busca la particion montada
	item := itemMount{}
	for i := 0; i < len(PartMount); i++ {
		if PartMount[i].Id == id {
			item = PartMount[i]
			break
		}
	}
	//se valida que se haya encontrado la particion montada
	if item.Id == "" {
		fmt.Println("Error, no existe una particion montada con el id ingresado")
		return
	}
	//se comienza con el formateo de la particion
	sb := SuperBloque{}
	//se asignan los valores al superbloque
	sb.S_filesystem_type[0] = '2'
	//se calcula el numero de estructuras
	c1 := string(item.Part.Part_size[:])
	c1 = strings.TrimRight(c1, "\x00")
	size_partition, err1 := strconv.Atoi(c1)
	if err1 != nil {
		fmt.Println("Error, no se pudo obtener el tamaño de la particion")
		return
	}
	//este resukltado ya esta
	n := (size_partition - int(unsafe.Sizeof(SuperBloque{}))) / (4 + int(unsafe.Sizeof(Inodo{})) + (3 * int(unsafe.Sizeof(B_Carpeta{}))))
	copy(sb.S_inodes_count[:], strconv.Itoa(n))
	copy(sb.S_blocks_count[:], strconv.Itoa(n*3))
	copy(sb.S_free_blocks_count[:], strconv.Itoa((n*3)-2))
	copy(sb.S_free_inodes_count[:], strconv.Itoa(n-2))
	date := time.Now()
	formatted := date.Format("02/01/2006 15:04:05")
	copy(sb.S_mtime[:], formatted)
	sb.S_mnt_count[0] = '1'
	sb.S_magic[0] = '6'
	sb.S_magic[1] = '2'
	sb.S_magic[2] = '1'
	sb.S_magic[3] = '6'
	sb.S_magic[4] = '7'
	//vamo a probar esta nueva forma de copiar
	copy(sb.S_inode_size[:], strconv.Itoa(int(unsafe.Sizeof(Inodo{}))))
	copy(sb.S_block_size[:], strconv.Itoa(int(unsafe.Sizeof(B_Carpeta{}))))
	sb.S_first_ino[0] = '0'
	sb.S_first_blo[0] = '0'
	c1 = string(item.Part.Part_start[:])
	c1 = strings.TrimRight(c1, "\x00")
	start_partition, err1 := strconv.Atoi(c1)
	if err1 != nil {
		fmt.Println("Error al convertir el inicio de la particion")
		return
	}
	copy(sb.S_bm_inode_start[:], strconv.Itoa(start_partition+int(unsafe.Sizeof(SuperBloque{}))))
	copy(sb.S_bm_block_start[:], strconv.Itoa(start_partition+int(unsafe.Sizeof(SuperBloque{}))+n))
	copy(sb.S_inode_start[:], strconv.Itoa(start_partition+int(unsafe.Sizeof(SuperBloque{}))+(4*n)))
	copy(sb.S_block_start[:], strconv.Itoa(start_partition+int(unsafe.Sizeof(SuperBloque{}))+(4*n)+(n*int(unsafe.Sizeof(Inodo{})))))
	/* final_position := start_partition + int(unsafe.Sizeof(SuperBloque{})) + (4 * n) + (n * int(unsafe.Sizeof(Inodo{}))) + ((3 * n) * int(unsafe.Sizeof(B_Carpeta{})))
	//para asegurarnos
	fmt.Print("Tamano de la particion: ")
	fmt.Println(size_partition)
	fmt.Print("inicio de la particion: ")
	fmt.Println(start_partition)
	fmt.Print("posicion final: ")
	fmt.Println(final_position)
	fmt.Print("Estructuras calculadas: ")
	fmt.Println(n)
	fmt.Print("tamano final: ")
	fmt.Println(int(unsafe.Sizeof(SuperBloque{})) + (4 * n) + (n * int(unsafe.Sizeof(Inodo{}))) + (3 * n * int(unsafe.Sizeof(B_Carpeta{})))) */

	//ahora se escribe el inodo raiz
	inodo := Inodo{}
	inodo.I_uid[0] = '1'
	inodo.I_gid[0] = '1'
	copy(inodo.I_atime[:], formatted)
	copy(inodo.I_ctime[:], formatted)
	copy(inodo.I_mtime[:], formatted)
	inodo.I_type[0] = '0'
	inodo.I_perm[0] = '6'
	inodo.I_perm[1] = '6'
	inodo.I_perm[2] = '4'

	sb.S_first_ino[0] = '1'

	//se crea el bloque de carpeta
	bc := B_Carpeta{}
	bc.B_content[0].B_inodo[0] = '0'
	bc.B_content[0].B_name[0] = '.'
	bc.B_content[1].B_inodo[0] = '0'
	bc.B_content[1].B_name[0] = '.'
	bc.B_content[1].B_name[1] = '.'
	copy(bc.B_content[2].B_name[:], "users.txt")
	bc.B_content[2].B_inodo[0] = '1'

	sb.S_first_blo[0] = '1'
	inodo.I_block[0] = [10]byte{'0', 0, 0, 0, 0, 0, 0, 0, 0, 0}

	//se crea el inodo users.txt
	inodo2 := Inodo{}
	inodo2.I_uid[0] = '1'
	inodo2.I_gid[0] = '1'

	copy(inodo2.I_atime[:], formatted)
	copy(inodo2.I_ctime[:], formatted)
	copy(inodo2.I_mtime[:], formatted)
	inodo2.I_type[0] = '1'
	inodo2.I_perm[0] = '6'
	inodo2.I_perm[1] = '6'
	inodo2.I_perm[2] = '4'
	inodo2.I_block[0] = [10]byte{'1', 0, 0, 0, 0, 0, 0, 0, 0, 0}

	sb.S_first_ino[0] = '2'

	//se crea el bloque de archivo
	ba := B_Archivo{}
	users := "1,G,root\n1,U,root,root,123\n"
	copy(ba.B_content[:], users)
	c1 = strconv.Itoa(len(users))
	copy(inodo2.I_size[:], c1)
	copy(inodo.I_size[:], c1)

	sb.S_first_blo[0] = '2'

	//----se escribe el superbloque, inodos y bloques
	disk_partition, err := os.OpenFile(item.Path, os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error al abrir el disco, mkfs")
		return
	}
	//se escribe el superbloque
	_, err = disk_partition.Seek(int64(start_partition), io.SeekStart)
	if err != nil {
		fmt.Println("Error posicionando el puntero, mkfs")
		return
	}

	err = binary.Write(disk_partition, binary.LittleEndian, sb)
	if err != nil {
		fmt.Println("Error al escribir el superbloque, mkfs")
		return
	}
	//se escriben los inodos
	_, err = disk_partition.Seek(int64(start_partition+int(unsafe.Sizeof(SuperBloque{}))+(4*n)), io.SeekStart)
	if err != nil {
		fmt.Println("Error posicionando el puntero inodo, mkfs")
		return
	}
	err = binary.Write(disk_partition, binary.LittleEndian, inodo)
	if err != nil {
		fmt.Println("Error al escribir el inodo, mkfs")
		return
	}
	err = binary.Write(disk_partition, binary.LittleEndian, inodo2)
	if err != nil {
		fmt.Println("Error al escribir el inodo2, mkfs")
		return
	}
	//se escriben los bloques
	_, err = disk_partition.Seek(int64(start_partition+int(unsafe.Sizeof(SuperBloque{}))+(4*n)+(n*int(unsafe.Sizeof(Inodo{})))), io.SeekStart)
	if err != nil {
		fmt.Println("Error posicionando el puntero bloque, mkfs")
		return
	}
	err = binary.Write(disk_partition, binary.LittleEndian, bc)
	if err != nil {
		fmt.Println("Error al escribir el bloque carpeta, mkfs")
		return
	}
	err = binary.Write(disk_partition, binary.LittleEndian, ba)
	if err != nil {
		fmt.Println("Error al escribir el bloque archivo, mkfs")
		return
	}
	//se escriben los datos al bitmap de inodos y bloques
	_, err = disk_partition.Seek(int64(start_partition+int(unsafe.Sizeof(SuperBloque{}))), io.SeekStart)
	if err != nil {
		fmt.Println("Error posicionando el puntero bitmap inodos, mkfs")
		return
	}
	//se escribe el bitmap de inodos
	err = binary.Write(disk_partition, binary.LittleEndian, '1')
	if err != nil {
		fmt.Println("Error al escribir el bitmap de inodos, mkfs")
		return
	}
	err = binary.Write(disk_partition, binary.LittleEndian, '1')
	if err != nil {
		fmt.Println("Error al escribir el bitmap de inodos, mkfs")
		return
	}
	//se posiciona en el bitmap de bloques
	_, err = disk_partition.Seek(int64(start_partition+int(unsafe.Sizeof(SuperBloque{}))+(4*n)+(n*int(unsafe.Sizeof(Inodo{})))+(3*n*int(unsafe.Sizeof(B_Carpeta{})))), io.SeekStart)
	if err != nil {
		fmt.Println("Error posicionando el puntero bitmap bloques, mkfs")
		return
	}
	//se escribe el bitmap de bloques
	err = binary.Write(disk_partition, binary.LittleEndian, '1')
	if err != nil {
		fmt.Println("Error al escribir el bitmap de bloques, mkfs")
		return
	}
	err = binary.Write(disk_partition, binary.LittleEndian, '1')
	if err != nil {
		fmt.Println("Error al escribir el bitmap de bloques, mkfs")
		return
	}
	disk_partition.Close()
	fmt.Println("Particion formateada correctamente")

}
