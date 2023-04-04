package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
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
