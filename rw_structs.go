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

// funcion para leer el mbr
func read_MBR(path string) (MBR, bool) {
	mbr := MBR{}
	disk, err := os.Open(path)
	if err != nil {
		fmt.Println("Error abriendo el archivo")
		disk.Close()
		return mbr, false
	}
	_, err1 := disk.Seek(0, io.SeekStart)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero")
		disk.Close()
		return mbr, false
	}

	err2 := binary.Read(disk, binary.LittleEndian, &mbr)
	if err2 != nil {
		fmt.Println("Error al leer el mbr")
		disk.Close()
		return mbr, false
	}
	disk.Close()
	return mbr, true
}

// funcion para escribir el mbr
func write_MBR(mbr MBR, path string) bool {
	disk, err := os.OpenFile(path, os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error abriendo el archivo")
		disk.Close()
		return false
	}
	_, err1 := disk.Seek(int64(0), io.SeekStart)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero")
		disk.Close()
		return false
	}
	// Se escribe el mbr en el archivo
	err = binary.Write(disk, binary.LittleEndian, mbr)
	if err != nil {
		fmt.Println("Error al escribir el mbr")
		disk.Close()
		return false
	}
	disk.Close()
	return true

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
		disk.Close()
		return ebr, false
	}
	//se posiciona el puntero en la posicion del disco
	_, err1 := disk.Seek(int64(d), io.SeekStart)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero ")
		disk.Close()
		return ebr, false
	}
	//se lee la particion extendida
	err2 := binary.Read(disk, binary.LittleEndian, &ebr)
	if err2 != nil {
		fmt.Println("Error leyendo la particion extendida ")
		disk.Close()
		return ebr, false
	}
	disk.Close()
	return ebr, true
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

// funcion para leer el superbloque de una particion
func read_sb(path string, position [10]byte) (SuperBloque, bool) {
	sb := SuperBloque{}
	disk, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		disk.Close()
		return sb, false
	}
	//asignamos a c el string con el numero
	c := string(position[:])
	//se elimina los caracteres nulos
	c = strings.TrimRight(c, "\x00")
	//se convierte a entero
	d, err := strconv.Atoi(c)
	if err != nil {
		fmt.Println("Error al convertir a entero ")
		return sb, false
	}
	//se posiciona el puntero en la posicion del disco
	_, err1 := disk.Seek(int64(d), io.SeekStart)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero ")
		return sb, false
	}
	//se lee la particion extendida
	err2 := binary.Read(disk, binary.LittleEndian, &sb)
	if err2 != nil {
		fmt.Println("Error leyendo el superbloque ")
		return sb, false
	}
	disk.Close()
	return sb, true
}

// funcion para escribir el superbloque
func write_sb(sb SuperBloque, path string, position int) bool {
	disk, err := os.OpenFile(path, os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error abriendo el archivo, sb")
		disk.Close()
		return false
	}
	_, err1 := disk.Seek(int64(position), io.SeekStart)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero. sb")
		disk.Close()
		return false
	}
	// Se escribe el mbr en el archivo
	err = binary.Write(disk, binary.LittleEndian, sb)
	if err != nil {
		fmt.Println("Error al escribir el superbloque, sb")
		disk.Close()
		return false
	}
	disk.Close()
	return true
}

// funcion para leer el inodo
func read_inodo(path string, position int) (Inodo, bool) {
	inodo := Inodo{}
	disk, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		disk.Close()
		return inodo, false
	}
	//se posiciona el puntero en la posicion
	_, err1 := disk.Seek(int64(position), io.SeekStart)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero, read inodo ")
		return inodo, false
	}
	//se lee la particion extendida
	err2 := binary.Read(disk, binary.LittleEndian, &inodo)
	if err2 != nil {
		fmt.Println("Error leyendo el inodo ")
		return inodo, false
	}
	disk.Close()
	return inodo, true
}

// funcion para leer el bloque de carpeta
func read_b_carpeta(path string, position int) (B_Carpeta, bool) {
	bloque := B_Carpeta{}
	disk, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		disk.Close()
		return bloque, false
	}
	//se posiciona el puntero en la posicion
	_, err1 := disk.Seek(int64(position), io.SeekStart)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero, read bloque carpeta ")
		disk.Close()
		return bloque, false
	}
	//se lee la particion extendida
	err2 := binary.Read(disk, binary.LittleEndian, &bloque)
	if err2 != nil {
		fmt.Println("Error leyendo el bloque carpeta ")
		disk.Close()
		return bloque, false
	}
	disk.Close()
	return bloque, true
}

// funcion para leer el bloque de archivo
func read_b_archivo(path string, position int) (B_Archivo, bool) {
	bloque := B_Archivo{}
	disk, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		disk.Close()
		return bloque, false
	}
	//se posiciona el puntero en la posicion
	_, err1 := disk.Seek(int64(position), io.SeekStart)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero, read bloque archivo ")
		disk.Close()
		return bloque, false
	}
	//se lee la particion extendida
	err2 := binary.Read(disk, binary.LittleEndian, &bloque)
	if err2 != nil {
		fmt.Println("Error leyendo el bloque archivo ")
		disk.Close()
		return bloque, false
	}
	disk.Close()
	return bloque, true
}

// funcion para escribir un bloque archivo
func write_b_archivo(barchivo B_Archivo, position int) bool {

	disk, err := os.OpenFile(ItemLogin.LoginItem.Path, os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error abriendo el archivo, ba")
		disk.Close()
		return false
	}
	//se posiciona el puntero en la posicion del primer bloque libre

	_, err1 := disk.Seek(int64(position), io.SeekStart)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero, ba")
		disk.Close()
		return false
	}
	// Se escribe el mbr en el archivo
	err = binary.Write(disk, binary.LittleEndian, barchivo)
	if err != nil {
		fmt.Println("Error al escribir el bloque archivo, ba")
		disk.Close()
		return false
	}
	disk.Close()
	return true
}

// funcion para escribir un bloque carpeta
func write_b_carpeta(bcarpeta B_Carpeta, position int) bool {

	disk, err := os.OpenFile(ItemLogin.LoginItem.Path, os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error abriendo el archivo, ba")
		disk.Close()
		return false
	}
	//se posiciona el puntero en la posicion del primer bloque libre

	_, err1 := disk.Seek(int64(position), io.SeekStart)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero, ba")
		disk.Close()
		return false
	}
	// Se escribe el mbr en el archivo
	err = binary.Write(disk, binary.LittleEndian, bcarpeta)
	if err != nil {
		fmt.Println("Error al escribir el bloque archivo, ba")
		disk.Close()
		return false
	}
	disk.Close()
	return true
}

// funcion para escribir un inodo
func write_inodo(inodo Inodo, position int) bool {
	disk, err := os.OpenFile(ItemLogin.LoginItem.Path, os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error abriendo el archivo, ba")
		disk.Close()
		return false
	}
	//se posiciona el puntero en la posicion del primer bloque libre

	_, err1 := disk.Seek(int64(position), io.SeekStart)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero, ba")
		disk.Close()
		return false
	}
	// Se escribe el mbr en el archivo
	err = binary.Write(disk, binary.LittleEndian, inodo)
	if err != nil {
		fmt.Println("Error al escribir el bloque archivo, ba")
		disk.Close()
		return false
	}
	disk.Close()
	return true
}

// funcion para escribir en el bitmap de inodos
func write_bitmap_inodos() bool {
	disk, err := os.OpenFile(ItemLogin.LoginItem.Path, os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error abriendo el archivo, ba")
		disk.Close()
		return false
	}
	//se posiciona el puntero en la posicion del primer bloque libre
	sb, flag := read_sb(ItemLogin.LoginItem.Path, ItemLogin.LoginItem.Part.Part_start)
	if !flag {
		fmt.Println("Error leyendo el superbloque")
		disk.Close()
		return false
	}
	//se obtiene el primer inodo libre y la posicion del bitmap de inodos
	c1 := string(sb.S_first_ino[:])
	c1 = strings.TrimRight(c1, "\x00")
	posicion_free_ino, err1 := strconv.Atoi(c1)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero, ba")
		disk.Close()
		return false
	}
	c1 = string(sb.S_bm_inode_start[:])
	c1 = strings.TrimRight(c1, "\x00")
	posicion_bitmap_inodos, err1 := strconv.Atoi(c1)

	if err1 != nil {
		fmt.Println("Error posicionando el puntero, ba")
		disk.Close()
		return false
	}
	_, err = disk.Seek(int64(posicion_bitmap_inodos+(int(unsafe.Sizeof(B_Archivo{}))*posicion_free_ino)), io.SeekStart)
	if err != nil {
		fmt.Println("Error posicionando el puntero bitmap inodos, mkfs")
		return false
	}

	err = binary.Write(disk, binary.LittleEndian, '1')
	if err != nil {
		fmt.Println("Error al escribir el bitmap de bloques, mkfs")
		return false
	}
	//se actualiza el superbloque
	c1 = string(sb.S_free_inodes_count[:])
	c1 = strings.TrimRight(c1, "\x00")
	free_inodes_count, err1 := strconv.Atoi(c1)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero, ba")
		return false
	}
	free_inodes_count -= 1
	posicion_free_ino += 1
	copy(sb.S_inode_start[:], []byte(strconv.Itoa(posicion_free_ino)))
	copy(sb.S_free_inodes_count[:], []byte(strconv.Itoa(free_inodes_count)))
	c1 = string(ItemLogin.LoginItem.Part.Part_start[:])
	c1 = strings.TrimRight(c1, "\x00")
	posicion_superbloque, err1 := strconv.Atoi(c1)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero, ba")
		return false
	}
	flag = write_sb(sb, ItemLogin.LoginItem.Path, posicion_superbloque)
	if !flag {
		fmt.Println("Error escribiendo el superbloque")
		return false
	}
	disk.Close()
	return true
}

// funcion para escribir en el bitmap de bloques
func write_bitmap_bloques() bool {
	disk, err := os.OpenFile(ItemLogin.LoginItem.Path, os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error abriendo el archivo, ba")
		disk.Close()
		return false
	}
	//se posiciona el puntero en la posicion del primer bloque libre
	sb, flag := read_sb(ItemLogin.LoginItem.Path, ItemLogin.LoginItem.Part.Part_start)
	if !flag {
		fmt.Println("Error leyendo el superbloque")
		disk.Close()
		return false
	}
	//se obtiene el primer bloque libre y la posicion del bitmap de bloques
	c1 := string(sb.S_first_blo[:])
	c1 = strings.TrimRight(c1, "\x00")
	posicion_free_block, err1 := strconv.Atoi(c1)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero, ba")
		disk.Close()
		return false
	}
	c1 = string(sb.S_bm_block_start[:])
	c1 = strings.TrimRight(c1, "\x00")
	posicion_bitmap_bloques, err1 := strconv.Atoi(c1)

	if err1 != nil {
		fmt.Println("Error posicionando el puntero, ba")
		disk.Close()
		return false
	}
	_, err = disk.Seek(int64(posicion_bitmap_bloques+posicion_free_block), io.SeekStart)
	if err != nil {
		fmt.Println("Error posicionando el puntero bitmap bloques, mkfs")
		return false
	}

	err = binary.Write(disk, binary.LittleEndian, '1')
	if err != nil {
		fmt.Println("Error al escribir el bitmap de bloques, mkfs")
		return false
	}
	//se actualiza el superbloque
	c1 = string(sb.S_free_blocks_count[:])
	c1 = strings.TrimRight(c1, "\x00")
	free_blocks_count, err1 := strconv.Atoi(c1)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero, ba")
		return false
	}
	free_blocks_count -= 1
	posicion_free_block += 1
	copy(sb.S_first_blo[:], []byte(strconv.Itoa(posicion_free_block)))
	copy(sb.S_free_blocks_count[:], []byte(strconv.Itoa(free_blocks_count)))
	c1 = string(ItemLogin.LoginItem.Part.Part_start[:])
	c1 = strings.TrimRight(c1, "\x00")
	posicion_superbloque, err1 := strconv.Atoi(c1)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero, ba")
		return false
	}
	flag = write_sb(sb, ItemLogin.LoginItem.Path, posicion_superbloque)
	if !flag {
		fmt.Println("Error escribiendo el superbloque")
		return false
	}

	disk.Close()
	return true
}
