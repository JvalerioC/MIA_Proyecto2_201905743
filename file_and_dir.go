package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func mkfile(params []string) {
	if len(PartMount) == 0 {
		fmt.Println("Error, No hay particiones montadas")
		cadRespuesta += "Error, No hay particiones montadas\n"
		return
	}
	if !ItemLogin.Iniciado {
		fmt.Println("Error, No hay una sesion iniciada")
		cadRespuesta += "Error, No hay una sesion iniciada\n"
		return
	}
	//se obtienen los parametros
	var path string
	var size int
	var content string
	var r bool

	for i := 0; i < len(params); i++ {
		array := strings.Split(params[i], "=")
		param := strings.ToLower(array[0])
		if param == ">path" {
			path = array[1]
		} else if param == ">size" {
			size, _ = strconv.Atoi(array[1])
		} else if param == ">cont" {
			content = array[1]
		} else if param == ">r" {
			if len(array) != 1 {
				fmt.Println("Error, parametro no valido para >r")
				cadRespuesta += "Error, parametro no valido para >r\n"
				return
			}
			r = true
		} else {
			fmt.Println("Error, el parametro ingresado no es valido")
			cadRespuesta += "Error, el parametro ingresado no es valido\n"
			return
		}
	}

	//se validan los parametros obligatorios
	if path == "" {
		fmt.Println("Error, parametro obligatorio vacio")
		cadRespuesta += "Error, parametro obligatorio vacio\n"
		return
	}

	//se separa la ruta en carpetas y se obtienen las variables
	array_path := strings.Split(path, "/")
	file := array_path[len(array_path)-1]
	array_path = array_path[1 : len(array_path)-1]

	//se verifica el tipo de archivo a crear
	if content != "" {
		createFileW(array_path, file, r, content)
		return
	}
	if size != 0 {
		if size < 0 {
			fmt.Println("Error, el tamaño debe ser mayor a cero")
			cadRespuesta += "Error, el tamaño debe ser mayor a cero\n"
			return
		}
		createFileS(array_path, file, r, size)
		return
	}
	createFileB(array_path, file, r)

}

func mkdir(params []string) {
	if len(PartMount) == 0 {
		fmt.Println("Error, No hay particiones montadas")
		cadRespuesta += "Error, No hay particiones montadas\n"
		return
	}
	if !ItemLogin.Iniciado {
		fmt.Println("Error, No hay una sesion iniciada")
		cadRespuesta += "Error, No hay una sesion iniciada\n"
		return
	}
	//se obtienen los parametros
	var path string
	var r bool

	for i := 0; i < len(params); i++ {
		array := strings.Split(params[i], "=")
		param := strings.ToLower(array[0])
		if param == ">path" {
			path = array[1]
		} else if param == ">r" {
			if len(array) != 1 {
				fmt.Println("Error, parametro no valido para >r")
				cadRespuesta += "Error, parametro no valido para >r\n"
				return
			}
			r = true
		} else {
			fmt.Println("Error, el parametro ingresado no es valido")
			cadRespuesta += "Error, el parametro ingresado no es valido\n"
			return
		}
	}

	//se validan los parametros obligatorios
	if path == "" {
		fmt.Println("Error, parametro obligatorio vacio")
		cadRespuesta += "Error, parametro obligatorio vacio\n"
		return
	}

	//se separa la ruta en carpetas
	array_path := strings.Split(path, "/")
	array_path = array_path[1:]

	createDir(array_path, r)
}

// funcion para crear carpetas
func createDir(array_path []string, r bool) {
	//se recupera el superbloque
	sb, flag := read_sb(ItemLogin.LoginItem.Path, ItemLogin.LoginItem.Part.Part_start)
	if !flag {
		fmt.Println("Error, no se pudo recuperar el superbloque, createdir")
		return
	}
	//se recupera el nodo raiz
	c := string(sb.S_inode_start[:])
	c = strings.TrimRight(c, "\x00")
	pos_inode_start, err := strconv.Atoi(c) //inicio de inodos
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, mkuser")
		return
	}
	//se lee el inodo raiz
	inodo, flag1 := read_inodo(ItemLogin.LoginItem.Path, pos_inode_start)
	if !flag1 {
		fmt.Println("Error, no se pudo recuperar el inodo, createdir")
		return
	}
	//se recorre el array para encontrar la carpeta o crearlas
	temp_inodo := inodo
	pos_inodo := pos_inode_start
	for i := 0; i < len(array_path); i++ {
		inodo = temp_inodo
		encounter := encounterInode(inodo, sb, array_path[i])
		if encounter == 0 {
			if r {

				//se crea el inodo padre
				pos_inodo = createInodeCarpeta(inodo, &sb, pos_inodo, array_path[i])
				if pos_inodo == 0 {
					//fmt.Println("Error, r, no se creo el inodo, mkdir 1")
					return
				}
				temp_inodo, flag1 = read_inodo(ItemLogin.LoginItem.Path, pos_inode_start+(int(unsafe.Sizeof(Inodo{}))*pos_inodo))
				if !flag1 {
					fmt.Println("Error, no se pudo recuperar el inodo, mkdir 1")
					return
				}

			} else if !r && len(array_path) == i+1 {

				//se crea la carpeta
				pos_inodo = createInodeCarpeta(inodo, &sb, pos_inodo, array_path[i])
				fmt.Println("entro al r negado ", pos_inodo)
				if pos_inodo == 0 {
					//fmt.Println("Error, r, no se creo el inodo, mkdir 2")
					return
				}
				temp_inodo, flag1 = read_inodo(ItemLogin.LoginItem.Path, pos_inode_start+(int(unsafe.Sizeof(Inodo{}))*pos_inodo))
				if !flag1 {
					fmt.Println("Error, no se pudo recuperar el inodo, mkdir 2")
					return
				}

			} else {
				fmt.Println("Error, la carpeta no existe, esta no se puede crear")
				return
			}
		} else {
			if len(array_path) == i+1 {
				fmt.Println("Error, la carpeta ya existe")
				cadRespuesta += "Error, la carpeta ya existe \n"
				return
			}
			//se recupera el inodo
			temp_inodo, flag1 = read_inodo(ItemLogin.LoginItem.Path, pos_inode_start+(int(unsafe.Sizeof(Inodo{}))*encounter))
			if !flag1 {
				fmt.Println("Error, no se pudo recuperar el inodo, mkdir 3")
				return
			}
			pos_inodo = encounter
		}
	}
	fmt.Println("Carpeta creada exitosamente")
	cadRespuesta += "Carpeta creada exitosamente\n"
}

// funcion para encontrar un inodo, retorna el numero de inodo que coincide
func encounterInode(inodo Inodo, sb SuperBloque, name string) int {

	//se recupera el inicio de los bloques
	c := string(sb.S_block_start[:])
	c = strings.TrimRight(c, "\x00")
	pos_block_start, err := strconv.Atoi(c) //inicio de bloques
	if err != nil {
		fmt.Println("Error, no se pudo ller el inicio de bloques, encounterinode")
		return 0
	}
	encounter := false
	pos_inodo := 0
	if inodo.I_type[0] != '1' {
		for i := 0; i < 16; i++ {
			if inodo.I_block[i][0] == 0 {
				break
			}
			c := string(inodo.I_block[i][:])
			c = strings.TrimRight(c, "\x00")
			pos_block, err := strconv.Atoi(c) //pos block
			if err != nil {
				fmt.Println("Error, no se pudo obtener pos bloque, encounterInode")
				break
			}
			bc, flag := read_b_carpeta(ItemLogin.LoginItem.Path, pos_block_start+(int(unsafe.Sizeof(B_Carpeta{}))*pos_block))
			if !flag {
				fmt.Println("Error, no se pudo leer el bloque, encounterInode")
				break
			}
			for j := 0; j < 4; j++ {
				tname := string(bc.B_content[j].B_name[:])
				tname = strings.TrimRight(tname, "\x00")
				if tname == name {
					encounter = true
					c = string(bc.B_content[j].B_inodo[:])
					c = strings.TrimRight(c, "\x00")
					pos_inodo, err = strconv.Atoi(c) //pos inodo
					if err != nil {
						fmt.Println("Error, no se pudo obtener el no inodo, encounterInode")
						break
					}
					break
				}
			}
		}
	}
	if encounter {
		return pos_inodo
	} else {
		return 0
	}
}

// funcion para crear un inodo carpeta
func createInodeCarpeta(inode Inodo, sb *SuperBloque, pos_inodo int, name string) int {
	//se validan los permisos
	gid := string(inode.I_gid[:])
	gid = strings.TrimRight(gid, "\x00") //id de grupo

	uid := string(inode.I_uid[:])
	uid = strings.TrimRight(uid, "\x00") //id de usuario
	if ItemLogin.Grupo_id != gid && ItemLogin.User_id != uid {
		fmt.Println("Error, no se tienen los permisos para crear la carpeta")
		cadRespuesta += "Error, no se tienen los permisos para crear la carpeta\n"
		return 0
	}

	//se recupera el primer inodo libre
	c := string((*sb).S_first_ino[:])
	c = strings.TrimRight(c, "\x00")
	pos_first_inodo, err := strconv.Atoi(c) //posicion del primer inodo libre
	if err != nil {
		fmt.Println("Error, no se pudo obtener posicion de inodo, cic")
		return 0
	}
	//se recupera el primer bloque libre
	c = string((*sb).S_first_blo[:])
	c = strings.TrimRight(c, "\x00")
	pos_first_block, err := strconv.Atoi(c) //posicion del primer bloque libre
	if err != nil {
		fmt.Println("Error, no se pudo obtener posicion de bloque, cic")
		return 0
	}
	//se recupera el inicio de inodos
	c = string((*sb).S_inode_start[:])
	c = strings.TrimRight(c, "\x00")
	pos_inode_start, err := strconv.Atoi(c) //posicion del inicio de inodos
	if err != nil {
		fmt.Println("Error, no se pudo obtener inicio de inodo, cic")
		return 0
	}
	//se recupera el inicio de bloques
	c = string((*sb).S_block_start[:])
	c = strings.TrimRight(c, "\x00")
	pos_block_start, err := strconv.Atoi(c) //posicion del inicio de bloques
	if err != nil {
		fmt.Println("Error, no se pudo obtener inicio de bloque, cic")
		return 0
	}

	//asignamos al inodo anterior el bloque y al bloque el nuevo inodo
	for i := 0; i < 16; i++ {
		if inode.I_block[i][0] == 0 {
			bc := B_Carpeta{}
			//se asigna el bloque al inodo
			c = strconv.Itoa(pos_first_block)
			copy(inode.I_block[i][:], c)
			//se asigna el inodo al bloque
			c = strconv.Itoa(pos_first_inodo)
			copy(bc.B_content[0].B_inodo[:], c)
			//se asigna el nombre al bloque
			copy(bc.B_content[0].B_name[:], name)
			//se escribe el bloque
			flag := write_b_carpeta(bc, pos_block_start+(int(unsafe.Sizeof(B_Carpeta{}))*pos_first_block))
			if !flag {
				fmt.Println("Error, no se pudo escribir el bloque, cic")
				return 0
			}
			//se escribe el inodo
			flag = write_inodo(inode, pos_inode_start+(int(unsafe.Sizeof(Inodo{}))*pos_inodo))
			if !flag {
				fmt.Println("Error, no se pudo escribir el inodo, cic")
				return 0
			}
			//se escribe el bloque en el bitmap
			flag = write_bitmap_bloques()
			if !flag {
				fmt.Println("Error, no se pudo escribir el bitmap de bloques, cic")
				return 0
			}
			//se actualizan los valores del superbloque
			copy((*sb).S_first_blo[:], strconv.Itoa(pos_first_block+1))
			pos_first_block++
			c = string((*sb).S_free_blocks_count[:])
			c = strings.TrimRight(c, "\x00")
			free_blocks, err := strconv.Atoi(c) //posicion del inicio de bloques
			if err != nil {
				fmt.Println("Error, no se pudo obtener free blocks, cic")
				return 0
			}
			copy((*sb).S_free_blocks_count[:], strconv.Itoa(free_blocks-1))
			break
		} else {
			c := string(inode.I_block[i][:])
			c = strings.TrimRight(c, "\x00")
			pos_block, err := strconv.Atoi(c) //pos block
			if err != nil {
				fmt.Println("Error, no se pudo obtener pos bloque, cic")
				return 0
			}
			bc, flag := read_b_carpeta(ItemLogin.LoginItem.Path, pos_block_start+(int(unsafe.Sizeof(B_Carpeta{}))*pos_block))
			if !flag {
				fmt.Println("Error, no se pudo leer el bloque, cic")
				return 0
			}
			flag_e := false
			for j := 0; j < 4; j++ {
				if bc.B_content[j].B_name[0] == 0 {
					//se asigna el inodo al bloque
					c = strconv.Itoa(pos_first_inodo)
					copy(bc.B_content[j].B_inodo[:], c)
					//se asigna el nombre al bloque
					copy(bc.B_content[j].B_name[:], name)
					//se escribe el bloque
					flag := write_b_carpeta(bc, pos_block_start+(int(unsafe.Sizeof(B_Carpeta{}))*pos_block))
					if !flag {
						fmt.Println("Error, no se pudo escribir el bloque, cic")
						return 0
					}
					flag_e = true
					break
				}
			}
			if flag_e {
				break
			}
		}
	}

	//se crea el inodo
	date := time.Now()
	formatted := date.Format("02/01/2006 15:04:05")
	bi := Inodo{}
	bi.I_uid[0] = ItemLogin.User_id[0]
	bi.I_gid[0] = ItemLogin.Grupo_id[0]
	copy(bi.I_atime[:], formatted)
	copy(bi.I_ctime[:], formatted)
	copy(bi.I_mtime[:], formatted)
	bi.I_type[0] = '0'
	bi.I_perm[0] = '6'
	bi.I_perm[1] = '6'
	bi.I_perm[2] = '4'

	//se crea el bloque con los padres
	bc := B_Carpeta{}
	bc.B_content[0].B_name[0] = '.'
	copy(bc.B_content[0].B_inodo[:], "0")
	bc.B_content[1].B_name[0] = '.'
	bc.B_content[1].B_name[1] = '.'
	copy(bc.B_content[1].B_inodo[:], "0")

	c = string((*sb).S_free_blocks_count[:])
	c = strings.TrimRight(c, "\x00")
	free_blocks, err := strconv.Atoi(c)
	if err != nil {
		fmt.Println("Error, no se pudo obtener free blocks, cic")
		return 0
	}
	copy((*sb).S_free_blocks_count[:], strconv.Itoa(free_blocks-1))
	c = string((*sb).S_free_inodes_count[:])
	c = strings.TrimRight(c, "\x00")
	free_inodes, err := strconv.Atoi(c)
	if err != nil {
		fmt.Println("Error, no se pudo obtener free blocks, cic")
		return 0
	}
	copy((*sb).S_free_inodes_count[:], strconv.Itoa(free_inodes-1))

	copy(bi.I_block[0][:], strconv.Itoa(pos_first_block))
	//se insertan el inodo y bloque de la nueva carpeta
	flag := write_inodo(bi, pos_inode_start+(int(unsafe.Sizeof(Inodo{}))*pos_first_inodo))
	if !flag {
		fmt.Println("Error, no se pudo escribir el inodo, cic")
		return 0
	}
	flag = write_b_carpeta(bc, pos_block_start+(int(unsafe.Sizeof(B_Carpeta{}))*(pos_first_block)))
	if !flag {
		fmt.Println("Error, no se pudo escribir el bloque, cic")
		return 0
	}
	//se escribe el bitmap de inodos y bloques
	flag = write_bitmap_inodos()
	if !flag {
		fmt.Println("Error, no se pudo escribir el bitmap de inodos, cic")
		return 0
	}
	flag = write_bitmap_bloques()
	if !flag {
		fmt.Println("Error, no se pudo escribir el bitmap de bloques, cic")
		return 0
	}
	//se actualiza el superbloque
	copy((*sb).S_first_blo[:], strconv.Itoa(pos_first_block+1))
	copy((*sb).S_first_ino[:], strconv.Itoa(pos_first_inodo+1))
	//se escribe el superbloque
	c = string(ItemLogin.LoginItem.Part.Part_start[:])
	c = strings.TrimRight(c, "\x00")
	posicion_superbloque, err1 := strconv.Atoi(c)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero superbloque, cic")
		return 0
	}
	flag = write_sb((*sb), ItemLogin.LoginItem.Path, posicion_superbloque)
	if !flag {
		fmt.Println("Error, no se pudo escribir el superbloque, cic")
		return 0
	}
	return pos_first_inodo

}

func createFileB(array_path []string, file string, r bool) {
	//se recupera el superbloque
	sb, flag := read_sb(ItemLogin.LoginItem.Path, ItemLogin.LoginItem.Part.Part_start)
	if !flag {
		fmt.Println("Error, no se pudo recuperar el superbloque, cfw")
		return
	}
	//se recupera el nodo raiz
	c := string(sb.S_inode_start[:])
	c = strings.TrimRight(c, "\x00")
	pos_inode_start, err := strconv.Atoi(c) //inicio de inodos
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, cfw")
		return
	}
	c = string(sb.S_block_start[:])
	c = strings.TrimRight(c, "\x00")
	pos_block_start, err := strconv.Atoi(c) //inicio de bloques
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, cfw")
		return
	}

	//se lee el inodo raiz
	inodo, flag1 := read_inodo(ItemLogin.LoginItem.Path, pos_inode_start)
	if !flag1 {
		fmt.Println("Error, no se pudo recuperar el inodo, createdir")
		return
	}
	//se recorre el array para encontrar la carpeta o crearlas
	temp_inodo := inodo
	pos_inodo := pos_inode_start
	for i := 0; i < len(array_path); i++ {
		inodo = temp_inodo
		encounter := encounterInode(inodo, sb, array_path[i])
		if encounter == 0 {
			if r {
				//se crea el inodo padre
				pos_inodo = createInodeCarpeta(inodo, &sb, pos_inodo, array_path[i])
				if pos_inodo == 0 {
					//fmt.Println("Error, r, no se creo el inodo, mkdir 1")
					return
				}
				temp_inodo, flag1 = read_inodo(ItemLogin.LoginItem.Path, pos_inode_start+(int(unsafe.Sizeof(Inodo{}))*pos_inodo))
				if !flag1 {
					fmt.Println("Error, no se pudo recuperar el inodo, cfw 1")
					return
				}

			} else {
				fmt.Println("Error, la carpeta no existe, esta no se puede crear")
				cadRespuesta += "Error, la carpeta no existe, esta no se puede crear\n"
				return
			}
		} else {
			//se recupera el inodo
			temp_inodo, flag1 = read_inodo(ItemLogin.LoginItem.Path, pos_inode_start+(int(unsafe.Sizeof(Inodo{}))*encounter))
			if !flag1 {
				fmt.Println("Error, no se pudo recuperar el inodo, cfw 3")
				return
			}
			pos_inodo = encounter
		}
	}
	//temp_inodo tiene la carpeta donde se creara el archivo

	//se lee el contenido del archivo

	file_content := ""

	//se busca en el inodo el primer espacio donde se pueda insertar el archivo
	for i := 0; i < 16; i++ {
		if temp_inodo.I_block[i][0] == 0 {
			createInodeArchivo(temp_inodo, &sb, pos_inodo, file_content, file)
			fmt.Println("El archivo ha sido creado.")
			cadRespuesta += "El archivo ha sido creado.\n"
			return
		} else {
			//se recupera el bloque
			c := string(temp_inodo.I_block[i][:])
			c = strings.TrimRight(c, "\x00")
			pos_block, err := strconv.Atoi(c) //inicio de inodos
			if err != nil {
				fmt.Println("Error, no se pudo leer la pos bloque, cfw")
				return
			}
			bc, flag := read_b_carpeta(ItemLogin.LoginItem.Path, pos_block_start+(int(unsafe.Sizeof(B_Carpeta{}))*pos_block))
			if !flag {
				fmt.Println("Error, no se pudo recuperar el bloque, cfw 4")
				return
			}
			//se recorre el bloque para encontrar el primer espacio libre
			for j := 0; j < 4; j++ {
				if bc.B_content[j].B_inodo[0] != 0 {
					c := string(bc.B_content[j].B_name[:])
					c = strings.TrimRight(c, "\x00")
					if c == file {
						/* fmt.Println("Desea sobreescribir el archivo, Y/N")
						var confirm string
						fmt.Scanln(&confirm)
						if confirm != "Y" {
							fmt.Println("No sobreescribira el archivo ")
							return
						} */
						c = string(bc.B_content[j].B_inodo[:])
						c = strings.TrimRight(c, "\x00")
						poss, err := strconv.Atoi(c) //inicio de inodos
						if err != nil {
							fmt.Println("Error, no se pudo leer la pos bloque, cfw")
							return
						}
						reWriteArchivo(temp_inodo, &sb, poss, file_content, file)
					}
				} else {
					createInodeArchivo(temp_inodo, &sb, pos_inodo, file_content, file)
					fmt.Println("El archivo ha sido creado.")
					cadRespuesta += "El archivo ha sido creado.\n"
					i = 16
					break
				}
			}
		}
	}
}

func createFileS(array_path []string, file string, r bool, size int) {
	//se recupera el superbloque
	sb, flag := read_sb(ItemLogin.LoginItem.Path, ItemLogin.LoginItem.Part.Part_start)
	if !flag {
		fmt.Println("Error, no se pudo recuperar el superbloque, cfw")
		return
	}
	//se recupera el nodo raiz
	c := string(sb.S_inode_start[:])
	c = strings.TrimRight(c, "\x00")
	pos_inode_start, err := strconv.Atoi(c) //inicio de inodos
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, cfw")
		return
	}
	c = string(sb.S_block_start[:])
	c = strings.TrimRight(c, "\x00")
	pos_block_start, err := strconv.Atoi(c) //inicio de bloques
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, cfw")
		return
	}

	//se lee el inodo raiz
	inodo, flag1 := read_inodo(ItemLogin.LoginItem.Path, pos_inode_start)
	if !flag1 {
		fmt.Println("Error, no se pudo recuperar el inodo, createdir")
		return
	}
	//se recorre el array para encontrar la carpeta o crearlas
	temp_inodo := inodo
	pos_inodo := pos_inode_start
	for i := 0; i < len(array_path); i++ {
		inodo = temp_inodo
		encounter := encounterInode(inodo, sb, array_path[i])
		if encounter == 0 {
			if r {
				//se crea el inodo padre
				pos_inodo = createInodeCarpeta(inodo, &sb, pos_inodo, array_path[i])
				if pos_inodo == 0 {
					//fmt.Println("Error, r, no se creo el inodo, mkdir 1")
					return
				}
				temp_inodo, flag1 = read_inodo(ItemLogin.LoginItem.Path, pos_inode_start+(int(unsafe.Sizeof(Inodo{}))*pos_inodo))
				if !flag1 {
					fmt.Println("Error, no se pudo recuperar el inodo, cfw 1")
					return
				}

			} else {
				fmt.Println("Error, la carpeta no existe, esta no se puede crear")
				cadRespuesta += "Error, la carpeta no existe, esta no se puede crear\n"
				return
			}
		} else {
			//se recupera el inodo
			temp_inodo, flag1 = read_inodo(ItemLogin.LoginItem.Path, pos_inode_start+(int(unsafe.Sizeof(Inodo{}))*encounter))
			if !flag1 {
				fmt.Println("Error, no se pudo recuperar el inodo, cfw 3")
				return
			}
			pos_inodo = encounter
		}
	}
	//temp_inodo tiene la carpeta donde se creara el archivo

	//se crea el contenido del archivo
	str := ""
	for i := 0; i < size; i++ {
		str += strconv.Itoa(i % 10)
	}
	file_content := str

	//se busca en el inodo el primer espacio donde se pueda insertar el archivo
	for i := 0; i < 16; i++ {
		if temp_inodo.I_block[i][0] == 0 {
			createInodeArchivo(temp_inodo, &sb, pos_inodo, file_content, file)
			fmt.Println("El archivo ha sido creado.")
			cadRespuesta += "El archivo ha sido creado.\n"
			return
		} else {
			//se recupera el bloque
			c := string(temp_inodo.I_block[i][:])
			c = strings.TrimRight(c, "\x00")
			pos_block, err := strconv.Atoi(c) //inicio de inodos
			if err != nil {
				fmt.Println("Error, no se pudo leer la pos bloque, cfw")
				return
			}
			bc, flag := read_b_carpeta(ItemLogin.LoginItem.Path, pos_block_start+(int(unsafe.Sizeof(B_Carpeta{}))*pos_block))
			if !flag {
				fmt.Println("Error, no se pudo recuperar el bloque, cfw 4")
				return
			}
			//se recorre el bloque para encontrar el primer espacio libre
			for j := 0; j < 4; j++ {
				if bc.B_content[j].B_inodo[0] != 0 {
					c := string(bc.B_content[j].B_name[:])
					c = strings.TrimRight(c, "\x00")
					if c == file {
						/* fmt.Println("Desea sobreescribir el archivo, Y/N")
						var confirm string
						fmt.Scanln(&confirm)
						if confirm != "Y" {
							fmt.Println("No sobreescribira el archivo ")
							return
						} */
						c = string(bc.B_content[j].B_inodo[:])
						c = strings.TrimRight(c, "\x00")
						poss, err := strconv.Atoi(c) //inicio de inodos
						if err != nil {
							fmt.Println("Error, no se pudo leer la pos bloque, cfw")
							return
						}
						reWriteArchivo(temp_inodo, &sb, poss, file_content, file)
					}
				} else {
					createInodeArchivo(temp_inodo, &sb, pos_inodo, file_content, file)
					fmt.Println("El archivo ha sido creado.")
					cadRespuesta += "El archivo ha sido creado.\n"
					i = 16
					break
				}
			}
		}
	}
}

func createFileW(array_path []string, file string, r bool, content string) {

	//se recupera el superbloque
	sb, flag := read_sb(ItemLogin.LoginItem.Path, ItemLogin.LoginItem.Part.Part_start)
	if !flag {
		fmt.Println("Error, no se pudo recuperar el superbloque, cfw")
		return
	}
	//se recupera el nodo raiz
	c := string(sb.S_inode_start[:])
	c = strings.TrimRight(c, "\x00")
	pos_inode_start, err := strconv.Atoi(c) //inicio de inodos
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, cfw")
		return
	}
	c = string(sb.S_block_start[:])
	c = strings.TrimRight(c, "\x00")
	pos_block_start, err := strconv.Atoi(c) //inicio de bloques
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, cfw")
		return
	}

	//se lee el inodo raiz
	inodo, flag1 := read_inodo(ItemLogin.LoginItem.Path, pos_inode_start)
	if !flag1 {
		fmt.Println("Error, no se pudo recuperar el inodo, createdir")
		return
	}
	//se recorre el array para encontrar la carpeta o crearlas
	temp_inodo := inodo
	pos_inodo := pos_inode_start
	for i := 0; i < len(array_path); i++ {
		inodo = temp_inodo
		encounter := encounterInode(inodo, sb, array_path[i])
		if encounter == 0 {
			if r {
				//se crea el inodo padre
				pos_inodo = createInodeCarpeta(inodo, &sb, pos_inodo, array_path[i])
				if pos_inodo == 0 {
					//fmt.Println("Error, r, no se creo el inodo, mkdir 1")
					return
				}
				temp_inodo, flag1 = read_inodo(ItemLogin.LoginItem.Path, pos_inode_start+(int(unsafe.Sizeof(Inodo{}))*pos_inodo))
				if !flag1 {
					fmt.Println("Error, no se pudo recuperar el inodo, cfw 1")
					return
				}

			} else {
				fmt.Println("Error, la carpeta no existe, esta no se puede crear")
				cadRespuesta += "Error, la carpeta no existe, esta no se puede crear\n"
				return
			}
		} else {
			//se recupera el inodo
			temp_inodo, flag1 = read_inodo(ItemLogin.LoginItem.Path, pos_inode_start+(int(unsafe.Sizeof(Inodo{}))*encounter))
			if !flag1 {
				fmt.Println("Error, no se pudo recuperar el inodo, cfw 3")
				return
			}
			pos_inodo = encounter
		}
	}
	//temp_inodo tiene la carpeta donde se creara el archivo

	//se lee el contenido del archivo
	str, err := os.ReadFile(content)
	if err != nil {
		fmt.Println("Error, no se puede abrir el archivo mkfile", err)
		return
	}
	file_content := string(str)

	//se busca en el inodo el primer espacio donde se pueda insertar el archivo
	for i := 0; i < 16; i++ {
		if temp_inodo.I_block[i][0] == 0 {
			createInodeArchivo(temp_inodo, &sb, pos_inodo, file_content, file)
			fmt.Println("El archivo ha sido creado.")
			cadRespuesta += "El archivo ha sido creado.\n"
			return
		} else {
			//se recupera el bloque
			c := string(temp_inodo.I_block[i][:])
			c = strings.TrimRight(c, "\x00")
			pos_block, err := strconv.Atoi(c) //inicio de inodos
			if err != nil {
				fmt.Println("Error, no se pudo leer la pos bloque, cfw")
				return
			}
			bc, flag := read_b_carpeta(ItemLogin.LoginItem.Path, pos_block_start+(int(unsafe.Sizeof(B_Carpeta{}))*pos_block))
			if !flag {
				fmt.Println("Error, no se pudo recuperar el bloque, cfw 4")
				return
			}
			//se recorre el bloque para encontrar el primer espacio libre
			for j := 0; j < 4; j++ {
				if bc.B_content[j].B_inodo[0] != 0 {
					c := string(bc.B_content[j].B_name[:])
					c = strings.TrimRight(c, "\x00")
					if c == file {
						/* fmt.Println("Desea sobreescribir el archivo, Y/N")
						var confirm string
						fmt.Scanln(&confirm)
						if confirm != "Y" {
							fmt.Println("No sobreescribira el archivo ")
							return
						} */
						c = string(bc.B_content[j].B_inodo[:])
						c = strings.TrimRight(c, "\x00")
						poss, err := strconv.Atoi(c) //inicio de inodos
						if err != nil {
							fmt.Println("Error, no se pudo leer la pos bloque, cfw")
							return
						}
						reWriteArchivo(temp_inodo, &sb, poss, file_content, file)
					}
				} else {
					createInodeArchivo(temp_inodo, &sb, pos_inodo, file_content, file)
					fmt.Println("El archivo ha sido creado.")
					cadRespuesta += "El archivo ha sido creado.\n"
					i = 16
					break
				}
			}
		}
	}
}

func reWriteArchivo(inodo Inodo, sb *SuperBloque, pos_inodo int, content string, name string) {

	//se validan los permisos
	gid := string(inodo.I_gid[:])
	gid = strings.TrimRight(gid, "\x00") //id de grupo

	uid := string(inodo.I_uid[:])
	uid = strings.TrimRight(uid, "\x00") //id de usuario
	if ItemLogin.Grupo_id != gid && ItemLogin.User_id != uid {
		fmt.Println("Error, no se tienen los permisos para crear la carpeta")
		cadRespuesta += "Error, no se tienen los permisos para crear la carpeta\n"
		return
	}

	c := string((*sb).S_inode_start[:])
	c = strings.TrimRight(c, "\x00")
	pos_inode_start, err := strconv.Atoi(c) //inicio de inodos
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, cia")
		return
	}
	c = string((*sb).S_block_start[:])
	c = strings.TrimRight(c, "\x00")
	pos_block_start, err := strconv.Atoi(c) //inicio de bloques
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, cia")
		return
	}

	//se recupera el primer bloque libre
	c = string((*sb).S_first_blo[:])
	c = strings.TrimRight(c, "\x00")
	pos_first_block, err := strconv.Atoi(c) //inicio de bloques
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, cia")
		return
	}
	//se crea el inodo
	date := time.Now()
	formatted := date.Format("02/01/2006 15:04:05")
	bi := Inodo{}
	bi.I_uid[0] = ItemLogin.User_id[0]
	bi.I_gid[0] = ItemLogin.Grupo_id[0]
	copy(bi.I_atime[:], formatted)
	copy(bi.I_ctime[:], formatted)
	copy(bi.I_mtime[:], formatted)
	bi.I_type[0] = '1'
	bi.I_perm[0] = '6'
	bi.I_perm[1] = '6'
	bi.I_perm[2] = '4'

	//verificamos el contenido
	if len(content) == 0 {
		bi.I_size[0] = '0'
		//ingresamos el inodo
		flag := write_inodo(bi, pos_inode_start+(int(unsafe.Sizeof(Inodo{}))*pos_inodo))
		if !flag {
			fmt.Println("Error, no se pudo escribir el inodo, cia")
			return
		}

		cadRespuesta += "El archivo ha sido sobreescrito correctamente\n"
		return

	}
	pos := 0
	cont := 0
	for {
		if cont == 16 {
			break
		}
		content_block := ""
		if pos+64 > len(content) {
			content_block = content[pos:]
		} else {
			content_block = content[pos : pos+64]
		}

		ba := B_Archivo{}
		copy(ba.B_content[:], content_block)
		flag := write_b_archivo(ba, pos_block_start+(int(unsafe.Sizeof(B_Archivo{}))*pos_first_block))
		if !flag {
			fmt.Println("Error, no se pudo escribir el bloque, cia")
			return
		}
		//se actualiza el bitmap de bloques
		flag = write_bitmap_bloques()
		if !flag {
			fmt.Println("Error, no se pudo escribir el bitmap de bloques, cia")
			return
		}
		c = strconv.Itoa(pos_first_block)
		copy(bi.I_block[cont][:], c)
		pos_first_block++
		c = string((*sb).S_free_blocks_count[:])
		c = strings.TrimRight(c, "\x00")
		free_blocks, err := strconv.Atoi(c) //posicion del inicio de bloques
		if err != nil {
			fmt.Println("Error, no se pudo obtener free inodes, cia")
			return
		}
		copy((*sb).S_free_blocks_count[:], strconv.Itoa(free_blocks-1))
		pos += 64
		cont++
		if pos > len(content) {
			break
		}
	}

	//se escribe el inodo
	flag := write_inodo(bi, pos_inode_start+(int(unsafe.Sizeof(Inodo{}))*pos_inodo))
	if !flag {
		fmt.Println("Error, no se pudo escribir el inodo, cia")
		return
	}

	//se actualizan el superbloque
	copy((*sb).S_first_blo[:], strconv.Itoa(pos_first_block))
	c = string(ItemLogin.LoginItem.Part.Part_start[:])
	c = strings.TrimRight(c, "\x00")
	posicion_superbloque, err1 := strconv.Atoi(c)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero superbloque, cia")
		return
	}
	flag = write_sb((*sb), ItemLogin.LoginItem.Path, posicion_superbloque)
	if !flag {
		fmt.Println("Error, no se pudo escribir el superbloque, cia")
		return
	}
	fmt.Println("todo nice")
	cadRespuesta += "El archivo ha sido sobreescrito correctamente\n"

}

func createInodeArchivo(inodo Inodo, sb *SuperBloque, pos_inodo int, content string, name string) {

	//se validan los permisos
	gid := string(inodo.I_gid[:])
	gid = strings.TrimRight(gid, "\x00") //id de grupo

	uid := string(inodo.I_uid[:])
	uid = strings.TrimRight(uid, "\x00") //id de usuario
	if ItemLogin.Grupo_id != gid && ItemLogin.User_id != uid {
		fmt.Println("Error, no se tienen los permisos para crear la carpeta")
		cadRespuesta += "Error, no se tienen los permisos para crear la carpeta\n"
		return
	}

	c := string((*sb).S_inode_start[:])
	c = strings.TrimRight(c, "\x00")
	pos_inode_start, err := strconv.Atoi(c) //inicio de inodos
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, cia")
		return
	}
	c = string((*sb).S_block_start[:])
	c = strings.TrimRight(c, "\x00")
	pos_block_start, err := strconv.Atoi(c) //inicio de bloques
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, cia")
		return
	}

	//se recupera el primer inodo libre
	c = string((*sb).S_first_ino[:])
	c = strings.TrimRight(c, "\x00")
	pos_first_inodo, err := strconv.Atoi(c) //inicio de inodos
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, cia")
		return
	}
	//se recupera el primer bloque libre
	c = string((*sb).S_first_blo[:])
	c = strings.TrimRight(c, "\x00")
	pos_first_block, err := strconv.Atoi(c) //inicio de bloques
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, cia")
		return
	}

	//asignamos al inodo anterior el bloque y al bloque el nuevo inodo
	for i := 0; i < 16; i++ {
		if inodo.I_block[i][0] == 0 {
			bc := B_Carpeta{}
			//se asigna el bloque al inodo
			c = strconv.Itoa(pos_first_block)
			copy(inodo.I_block[i][:], c)
			//se asigna el inodo al bloque
			c = strconv.Itoa(pos_first_inodo)
			copy(bc.B_content[0].B_inodo[:], c)
			//se asigna el nombre al bloque
			copy(bc.B_content[0].B_name[:], name)
			//se escribe el bloque
			flag := write_b_carpeta(bc, pos_block_start+(int(unsafe.Sizeof(B_Carpeta{}))*pos_first_block))
			if !flag {
				fmt.Println("Error, no se pudo escribir el bloque, cia")
				return
			}
			//se escribe el inodo
			flag = write_inodo(inodo, pos_inode_start+(int(unsafe.Sizeof(Inodo{}))*pos_inodo))
			if !flag {
				fmt.Println("Error, no se pudo escribir el inodo, cia")
				return
			}
			//se escribe el bloque en el bitmap
			flag = write_bitmap_bloques()
			if !flag {
				fmt.Println("Error, no se pudo escribir el bitmap de bloques, cia")
				return
			}
			//se actualizan los valores del superbloque
			copy((*sb).S_first_blo[:], strconv.Itoa(pos_first_block+1))
			pos_first_block++
			c = string((*sb).S_free_blocks_count[:])
			c = strings.TrimRight(c, "\x00")
			free_blocks, err := strconv.Atoi(c) //posicion del inicio de bloques
			if err != nil {
				fmt.Println("Error, no se pudo obtener free blocks, cia")
				return
			}
			copy((*sb).S_free_blocks_count[:], strconv.Itoa(free_blocks-1))
			break
		} else {
			c := string(inodo.I_block[i][:])
			c = strings.TrimRight(c, "\x00")
			pos_block, err := strconv.Atoi(c) //pos block
			if err != nil {
				fmt.Println("Error, no se pudo obtener pos bloque, cia")
				return
			}
			bc, flag := read_b_carpeta(ItemLogin.LoginItem.Path, pos_block_start+(int(unsafe.Sizeof(B_Carpeta{}))*pos_block))
			if !flag {
				fmt.Println("Error, no se pudo leer el bloque, cia")
				return
			}
			flag_e := false
			for j := 0; j < 4; j++ {
				if bc.B_content[j].B_name[0] == 0 {
					//se asigna el inodo al bloque
					c = strconv.Itoa(pos_first_inodo)
					copy(bc.B_content[j].B_inodo[:], c)
					//se asigna el nombre al bloque
					copy(bc.B_content[j].B_name[:], name)
					//se escribe el bloque
					flag := write_b_carpeta(bc, pos_block_start+(int(unsafe.Sizeof(B_Carpeta{}))*pos_block))
					if !flag {
						fmt.Println("Error, no se pudo escribir el bloque, cic")
						return
					}
					flag_e = true
					break
				}
			}
			if flag_e {
				break
			}
		}
	}

	//se crea el inodo
	date := time.Now()
	formatted := date.Format("02/01/2006 15:04:05")
	bi := Inodo{}
	bi.I_uid[0] = ItemLogin.User_id[0]
	bi.I_gid[0] = ItemLogin.Grupo_id[0]
	copy(bi.I_atime[:], formatted)
	copy(bi.I_ctime[:], formatted)
	copy(bi.I_mtime[:], formatted)
	bi.I_type[0] = '1'
	bi.I_perm[0] = '6'
	bi.I_perm[1] = '6'
	bi.I_perm[2] = '4'

	//verificamos el contenido
	if len(content) == 0 {
		bi.I_size[0] = '0'
		//ingresamos el inodo
		flag := write_inodo(bi, pos_inode_start+(int(unsafe.Sizeof(Inodo{}))*pos_first_inodo))
		if !flag {
			fmt.Println("Error, no se pudo escribir el inodo, cia")
			return
		}
		//se actualiza el bitmap de inodos
		flag = write_bitmap_inodos()
		if !flag {
			fmt.Println("Error, no se pudo escribir el bitmap de inodos, cia")
			return
		}
		//se actualizan el superbloque
		copy((*sb).S_first_ino[:], strconv.Itoa(pos_first_inodo+1))
		pos_first_inodo++
		c = string((*sb).S_free_inodes_count[:])
		c = strings.TrimRight(c, "\x00")
		free_inodes, err := strconv.Atoi(c) //posicion del inicio de bloques
		if err != nil {
			fmt.Println("Error, no se pudo obtener free inodes, cia")
			return
		}
		copy((*sb).S_free_inodes_count[:], strconv.Itoa(free_inodes-1))
		c = string(ItemLogin.LoginItem.Part.Part_start[:])
		c = strings.TrimRight(c, "\x00")
		posicion_superbloque, err1 := strconv.Atoi(c)
		if err1 != nil {
			fmt.Println("Error posicionando el puntero superbloque, cia")
			return
		}
		flag = write_sb((*sb), ItemLogin.LoginItem.Path, posicion_superbloque)
		if !flag {
			fmt.Println("Error, no se pudo escribir el superbloque, cia")
			return
		}
		fmt.Println("todo nice")
		return

	}
	pos := 0
	cont := 0
	for {
		if cont == 16 {
			break
		}
		content_block := ""
		if pos+64 > len(content) {
			content_block = content[pos:]
		} else {
			content_block = content[pos : pos+64]
		}

		ba := B_Archivo{}
		copy(ba.B_content[:], content_block)
		flag := write_b_archivo(ba, pos_block_start+(int(unsafe.Sizeof(B_Archivo{}))*pos_first_block))
		if !flag {
			fmt.Println("Error, no se pudo escribir el bloque, cia")
			return
		}
		//se actualiza el bitmap de bloques
		flag = write_bitmap_bloques()
		if !flag {
			fmt.Println("Error, no se pudo escribir el bitmap de bloques, cia")
			return
		}
		c = strconv.Itoa(pos_first_block)
		copy(bi.I_block[cont][:], c)
		pos_first_block++
		c = string((*sb).S_free_blocks_count[:])
		c = strings.TrimRight(c, "\x00")
		free_blocks, err := strconv.Atoi(c) //posicion del inicio de bloques
		if err != nil {
			fmt.Println("Error, no se pudo obtener free inodes, cia")
			return
		}
		copy((*sb).S_free_blocks_count[:], strconv.Itoa(free_blocks-1))
		pos += 64
		cont++
		if pos > len(content) {
			break
		}
	}

	//se escribe el inodo
	flag := write_inodo(bi, pos_inode_start+(int(unsafe.Sizeof(Inodo{}))*pos_first_inodo))
	if !flag {
		fmt.Println("Error, no se pudo escribir el inodo, cia")
		return
	}
	//se actualiza el bitmap de inodos
	flag = write_bitmap_inodos()
	if !flag {
		fmt.Println("Error, no se pudo escribir el bitmap de inodos, cia")
		return
	}
	//se actualizan el superbloque
	copy((*sb).S_first_ino[:], strconv.Itoa(pos_first_inodo+1))
	copy((*sb).S_first_blo[:], strconv.Itoa(pos_first_block))
	c = string((*sb).S_free_inodes_count[:])
	c = strings.TrimRight(c, "\x00")
	free_inodes, err := strconv.Atoi(c) //posicion del inicio de bloques
	if err != nil {
		fmt.Println("Error, no se pudo obtener free inodes, cia")
		return
	}
	copy((*sb).S_free_inodes_count[:], strconv.Itoa(free_inodes-1))
	c = string(ItemLogin.LoginItem.Part.Part_start[:])
	c = strings.TrimRight(c, "\x00")
	posicion_superbloque, err1 := strconv.Atoi(c)
	if err1 != nil {
		fmt.Println("Error posicionando el puntero superbloque, cia")
		return
	}
	flag = write_sb((*sb), ItemLogin.LoginItem.Path, posicion_superbloque)
	if !flag {
		fmt.Println("Error, no se pudo escribir el superbloque, cia")
		return
	}
	fmt.Println("todo nice")
}
