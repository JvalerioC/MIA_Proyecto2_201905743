package main

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

func login2(user string, password string, id string) bool {
	//se valida que exista una particion Montada
	if len(PartMount) == 0 {
		fmt.Println("Error, no hay particiones montadas")
		cadRespuesta += "Error, no hay particiones montadas\n"
		return false
	}

	//se obtienen los parametros
	/* for i := 0; i < len(params); i++ {
		array := strings.Split(params[i], "=")
		param := strings.ToLower(array[0])
		if param == ">user" {
			user = array[1]
		} else if param == ">pwd" {
			password = array[1]
		} else if param == ">id" {
			id = array[1]
		} else {
			fmt.Println("El parametro ingresado no es valido")
			return
		}
	} */

	//se validan los parametros obligatorios
	if user == "" || password == "" || id == "" {
		fmt.Println("Error, parametro obligatorio vacio")
		cadRespuesta += "Error, parametro obligatorio vacio\n"
		return false
	}
	//se valida que no haya una sesion iniciada
	if ItemLogin.Iniciado {
		fmt.Println("Error, ya hay una sesion iniciada")
		cadRespuesta += "Error, ya hay una sesion iniciada\n"
		return false
	}
	//se busca el id en las particiones montadas
	flag := false
	item := itemMount{}
	for i := 0; i < len(PartMount); i++ {
		if PartMount[i].Id == id {
			flag = true
			item = PartMount[i]
			break
		}
	}
	if !flag {
		fmt.Println("Error, el id dela particion montada no existe")
		cadRespuesta += "Error, el id dela particion montada no existe\n"
		return false
	}
	//vamos a verificar si es el root
	if user == "root" && password == "123" {
		ItemLogin.Iniciado = true
		ItemLogin.Admin = true
		ItemLogin.User = user
		ItemLogin.LoginItem = item
		ItemLogin.Grupo = "root"
		ItemLogin.Grupo_id = "1"
		ItemLogin.User_id = "1"
		fmt.Println("Inicion de sesion exitoso, Bienvenido root")
		cadRespuesta += "Inicion de sesion exitoso, Bienvenido root\n"
		return true
	} else {

		file_content := contenidoUsers(item)
		if len(file_content) == 0 {
			fmt.Println("Error, no se pudo recuperar el contenido del archivo users.txt")
			return false
		}
		//se separa por saltos de linea y luego por comas
		encounter := false
		array_users := strings.Split(file_content, "\n")

		for i := 0; i < len(array_users); i++ {
			array := strings.Split(array_users[i], ",")
			if array[0] != "0" {
				if len(array) == 5 {
					if array[3] == user && array[4] == password {
						encounter = true
						ItemLogin.Iniciado = true
						ItemLogin.Admin = false
						ItemLogin.User = user
						ItemLogin.LoginItem = item
						ItemLogin.Grupo = array[3]
						ItemLogin.Grupo_id = ""
						ItemLogin.User_id = array[0]
						break
					}
				}
			}
		}
		if encounter {
			//se hara lo mismo pero para encontrar el id del grupo, no se si ha de servir de algo pero por si las moscas
			for i := 0; i < len(array_users); i++ {
				array := strings.Split(array_users[i], ",")
				if array[0] != "0" {
					if len(array) == 3 {
						if ItemLogin.Grupo == array[3] {
							ItemLogin.Grupo_id = array[0]
							break
						}
					}
				}
			}
			fmt.Println("Inicio de sesion exitoso, Bienvenido ", user)
			cadRespuesta += "Inicio de sesion exitoso, Bienvenido " + user
			return true
		} else {
			fmt.Println("Error, usuario o contrasena incorrecto")
			cadRespuesta += "Error, usuario o contrasena incorrecto"
			return false
		}
	}
}

func login(params []string) {
	var user string
	var password string
	var id string
	//se valida que exista una particion Montada
	if len(PartMount) == 0 {
		fmt.Println("Error, no hay particiones montadas")
		cadRespuesta += "Error, no hay particiones montadas\n"
		return
	}

	//se obtienen los parametros
	for i := 0; i < len(params); i++ {
		array := strings.Split(params[i], "=")
		param := strings.ToLower(array[0])
		if param == ">user" {
			user = array[1]
		} else if param == ">pwd" {
			password = array[1]
		} else if param == ">id" {
			id = array[1]
		} else {
			fmt.Println("El parametro ingresado no es valido")
			cadRespuesta += "El parametro ingresado no es valido\n"
			return
		}
	}

	//se validan los parametros obligatorios
	if user == "" || password == "" || id == "" {
		fmt.Println("Error, parametro obligatorio vacio")
		cadRespuesta += "Error, parametro obligatorio vacio\n"
		return
	}
	//se valida que no haya una sesion iniciada
	if ItemLogin.Iniciado {
		fmt.Println("Error, ya hay una sesion iniciada")
		cadRespuesta += "Error, ya hay una sesion iniciada\n"
		return
	}
	//se busca el id en las particiones montadas
	flag := false
	item := itemMount{}
	for i := 0; i < len(PartMount); i++ {
		if PartMount[i].Id == id {
			flag = true
			item = PartMount[i]
			break
		}
	}
	if !flag {
		fmt.Println("Error, el id dela particion montada no existe")
		cadRespuesta += "Error, el id de la particion montada no existe\n"
		return
	}
	//vamos a verificar si es el root
	if user == "root" && password == "123" {
		ItemLogin.Iniciado = true
		ItemLogin.Admin = true
		ItemLogin.User = user
		ItemLogin.LoginItem = item
		ItemLogin.Grupo = "root"
		ItemLogin.Grupo_id = "1"
		ItemLogin.User_id = "1"
		fmt.Println("Inicion de sesion exitoso, Bienvenido root")
		cadRespuesta += "Inicion de sesion exitoso, Bienvenido root\n"
		return
	} else {

		file_content := contenidoUsers(item)
		if len(file_content) == 0 {
			fmt.Println("Error, no se pudo recuperar el contenido del archivo users.txt")
			cadRespuesta += "Inicion de sesion exitoso, Bienvenido root\n"
			return
		}
		//se separa por saltos de linea y luego por comas
		encounter := false
		array_users := strings.Split(file_content, "\n")

		for i := 0; i < len(array_users); i++ {
			array := strings.Split(array_users[i], ",")
			if array[0] != "0" {
				if len(array) == 5 {
					if array[3] == user && array[4] == password {
						encounter = true
						ItemLogin.Iniciado = true
						ItemLogin.Admin = false
						ItemLogin.User = user
						ItemLogin.LoginItem = item
						ItemLogin.Grupo = array[3]
						ItemLogin.Grupo_id = ""
						ItemLogin.User_id = array[0]
						break
					}
				}
			}
		}
		if encounter {
			//se hara lo mismo pero para encontrar el id del grupo, no se si ha de servir de algo pero por si las moscas
			for i := 0; i < len(array_users); i++ {
				array := strings.Split(array_users[i], ",")
				if array[0] != "0" {
					if len(array) == 3 {
						if ItemLogin.Grupo == array[3] {
							ItemLogin.Grupo_id = array[0]
							break
						}
					}
				}
			}
			fmt.Println("Inicio de sesion exitoso, Bienvenido ", user)
		} else {
			fmt.Println("Error, usuario o contrasena incorrecto")
		}

	}
}

func logout() {
	if ItemLogin.Iniciado {
		ItemLogin.Iniciado = false
		ItemLogin.Admin = false
		ItemLogin.User = ""
		ItemLogin.LoginItem = itemMount{}
		ItemLogin.Grupo = ""
		ItemLogin.Grupo_id = ""
		ItemLogin.User_id = ""
		fmt.Println("Cierre de sesion exitosa")
		cadRespuesta += "Cierre de sesion exitosa\n"
	} else {
		fmt.Println("Error, no hay una sesion iniciada")
		cadRespuesta += "Error, no hay una sesion iniciada\n"
	}
}

func mkgrp(params []string) {
	//se hacen las validaciones para poder ejecutar este comando
	if len(PartMount) == 0 {
		fmt.Println("Error, no hay particion montada")
		cadRespuesta += "Error, no hay particion montada \n"
		return
	}
	if !ItemLogin.Iniciado {
		fmt.Println("Error, no hay una sesion iniciada")
		cadRespuesta += "Error, no hay una sesion iniciada \n"
		return
	}
	if !ItemLogin.Admin {
		fmt.Println("Error, el usuario no puede ejecutar este comando")
		cadRespuesta += "Error, el usuario no puede ejecutar este comando \n"
		return
	}

	//se obtienen los parametros
	var name string
	for i := 0; i < len(params); i++ {
		array := strings.Split(params[i], "=")
		param := strings.ToLower(array[0])
		if param == ">name" {
			name = array[1]
		} else {
			fmt.Println("El parametro ingresado no es valido")
			cadRespuesta += "El parametro ingresado no es valido \n"
			return
		}
	}
	//se validan los parametros obligatorios
	if name == "" {
		fmt.Println("Error, parametro obligatorio vacio")
		cadRespuesta += "Error, parametro obligatorio vacio \n"
		return
	}
	if len(name) > 10 {
		fmt.Println("Error, el nombre del grupo no puede ser mayor a 10 caracteres")
		cadRespuesta += "Error, el nombre del grupo no puede ser mayor a 10 caracteres \n"
		return
	}

	//se obtiene el contenido del archivo users.txt
	file_content := contenidoUsers(ItemLogin.LoginItem)
	if len(file_content) == 0 {
		fmt.Println("Error, no se pudo recuperar el contenido del archivo users.txt")

		return
	}
	//se separa por saltos de linea y luego por comas
	encounter := false
	contador_grupos := 1
	array_users := strings.Split(file_content, "\n")
	array_users = array_users[:len(array_users)-1]
	for i := 0; i < len(array_users); i++ {
		array := strings.Split(array_users[i], ",")
		if array[0] != "0" {
			if len(array) == 3 || array[1] == "G" {
				contador_grupos++
				if array[2] == name {
					encounter = true
					break
				}
			}
		} else {
			if len(array) == 3 || array[1] == "G" {
				contador_grupos++
			}
		}
	}
	if encounter {
		fmt.Println("Error, el nombre de grupo ya existe")
		cadRespuesta += "Error, el nombre de grupo ya existe \n"
		return
	}
	//se agrega el nuevo grupo al archivo users.txt
	//se obtiene el superbloque
	sb, flag := read_sb(ItemLogin.LoginItem.Path, ItemLogin.LoginItem.Part.Part_start)
	if !flag {
		fmt.Println("Error, no se pudo leer el superbloque, mkgrp")
		return
	}
	nuevo_grupo := strconv.Itoa(contador_grupos)
	nuevo_grupo += ",G," + name + "\n"
	//se obtiene la posicion del inodo del archivo users
	c := string(sb.S_inode_start[:])
	c = strings.TrimRight(c, "\x00")
	posicion_inodo, err := strconv.Atoi(c)
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, mkgrp")
		return
	}
	//se sabe que el inodo de users.txt es el uno
	posicion_inodo += int(unsafe.Sizeof(Inodo{}))

	//se obtiene el inodo del archivo users.txt
	inodo, flag := read_inodo(ItemLogin.LoginItem.Path, posicion_inodo)
	if !flag {
		fmt.Println("Error, no se pudo leer el inodo de users.txt, mkgrp")
		return
	}
	//se obtiene la posicion inicial de bloques
	c = string(sb.S_block_start[:])
	c = strings.TrimRight(c, "\x00")
	posicion_bloque_inicio, err := strconv.Atoi(c)
	if err != nil {
		fmt.Println("Error, no se pudo leer la posicion inicial de bloques, mkgrp")
		return
	}
	//se obtiene el bloque del archivo users.txt
	for i := 0; i < 16; i++ {
		if inodo.I_block[i][0] == 0 {
			break
		}
		d := string(inodo.I_block[i][:])
		d = strings.TrimRight(d, "\x00")
		posicion_bloque, err := strconv.Atoi(d)
		if err != nil {
			fmt.Println("Error, no se pudo obtener la posicion del bloque, mkgrp")
			return
		}
		//se obtiene el bloque archivo
		bloque, flag := read_b_archivo(ItemLogin.LoginItem.Path, posicion_bloque_inicio+(int(unsafe.Sizeof(B_Archivo{}))*posicion_bloque))
		if !flag {
			fmt.Println("Error, no se pudo leer el bloque del archivo users.txt, mkgrp")
			return
		}
		temporal_content := string(bloque.B_content[:])
		temporal_content = strings.TrimRight(temporal_content, "\x00")
		if i != 15 {

			if inodo.I_block[i+1][0] == 0 {
				//se lee el contenido del bloque
				if len(temporal_content)+len(nuevo_grupo) >= 64 {
					//se crea un nuevo bloque
					new_bloque := B_Archivo{}
					copy(new_bloque.B_content[:], nuevo_grupo)
					//se escribe el nuevo bloque
					c1 := string(sb.S_first_blo[:])
					c1 = strings.TrimRight(c1, "\x00")
					posicion_free_block, err1 := strconv.Atoi(c1)
					if err1 != nil {
						fmt.Println("Error posicionando el puntero, ba")
						return
					}
					flag := write_b_archivo(new_bloque, posicion_bloque_inicio+(int(unsafe.Sizeof(B_Archivo{}))*posicion_free_block))
					if !flag {
						fmt.Println("Error, no se pudo escribir el nuevo bloque, mkgrp")
						return
					}
					//se escribe en el bm de bloques y actualiza el superbloque
					flag = write_bitmap_bloques()
					if !flag {
						fmt.Println("Error, no se pudo escribir el bm de bloques, mkgrp")
						return
					}

					//se actualiza el inodo
					copy(inodo.I_block[i+1][:], strconv.Itoa(posicion_free_block))
					c1 = strconv.Itoa(len(file_content + nuevo_grupo))
					copy(inodo.I_size[:], c1)
					flag = write_inodo(inodo, posicion_inodo)
					if !flag {
						fmt.Println("Error, no se pudo escribir el inodo, mkgrp")
						return
					}
					fmt.Println("Grupo ingresado con exito")
					cadRespuesta += "Grupo ingresado con exito \n"
					return
				} else {
					//se escribe el nuevo grupo en el bloque
					copy(bloque.B_content[:], temporal_content+nuevo_grupo)
					//se escribe el bloque
					flag := write_b_archivo(bloque, posicion_bloque_inicio+(int(unsafe.Sizeof(B_Archivo{}))*posicion_bloque))
					if !flag {
						fmt.Println("Error, no se pudo escribir el bloque, mkgrp")
						return
					}
					//se actualiza el inodo
					c1 := strconv.Itoa(len(file_content + nuevo_grupo))
					copy(inodo.I_size[:], c1)
					flag = write_inodo(inodo, posicion_inodo)
					if !flag {
						fmt.Println("Error, no se pudo escribir el inodo, mkgrp")
						return
					}
					fmt.Println("Grupo ingresado con exito")
					cadRespuesta += "Grupo ingresado con exito\n"
					return
				}
			}
		} else {
			if len(temporal_content)+len(nuevo_grupo) >= 64 {
				fmt.Println("Error, no se pudo ingresar el grupo, mkgrp")
				return
			} else {
				//se escribe el nuevo grupo en el bloque
				copy(bloque.B_content[:], temporal_content+nuevo_grupo)
				//se escribe el bloque
				flag := write_b_archivo(bloque, posicion_bloque_inicio+(int(unsafe.Sizeof(B_Archivo{}))*posicion_bloque))
				if !flag {
					fmt.Println("Error, no se pudo escribir el bloque, mkgrp")
					return
				}
				//se actualiza el inodo
				c1 := strconv.Itoa(len(file_content + nuevo_grupo))
				copy(inodo.I_size[:], c1)
				flag = write_inodo(inodo, posicion_inodo)
				if !flag {
					fmt.Println("Error, no se pudo escribir el inodo, mkgrp")
					return
				}
				fmt.Println("Grupo ingresado con exito")
				cadRespuesta += "Grupo ingresado con exito\n"
				return
			}
		}
	}
}

func mkusr(params []string) {
	//se hacen las validaciones para poder ejecutar este comando
	if len(PartMount) == 0 {
		fmt.Println("Error, no hay particion montada")
		cadRespuesta += "Error, no hay particion montada\n"
		return
	}
	if !ItemLogin.Iniciado {
		fmt.Println("Error, no hay una sesion iniciada")
		cadRespuesta += "Error, no hay una sesion iniciada\n"
		return
	}
	if !ItemLogin.Admin {
		fmt.Println("Error, el usuario no puede ejecutar este comando")
		cadRespuesta += "Error, el usuario no puede ejecutar este comando\n"
		return
	}

	//se obtienen los parametros
	var username string
	var password string
	var group string
	for i := 0; i < len(params); i++ {
		array := strings.Split(params[i], "=")
		param := strings.ToLower(array[0])
		if param == ">user" {
			username = array[1]
		} else if param == ">pwd" {
			password = array[1]
		} else if param == ">grp" {
			group = array[1]
		} else {
			fmt.Println("El parametro ingresado no es valido")
			cadRespuesta += "El parametro ingresado no es valido\n"
			return
		}
	}
	//se validan los parametros obligatorios
	if username == "" || password == "" || group == "" {
		fmt.Println("Error, parametro obligatorio vacio")
		cadRespuesta += "Error, parametro obligatorio vacio\n"
		return
	}
	if len(username) > 10 || len(password) > 10 {
		fmt.Println("Error, el nombre del grupo no puede ser mayor a 10 caracteres")
		cadRespuesta += "Error, el nombre del grupo no puede ser mayor a 10 caracteres\n"
		return
	}

	//se obtiene el contenido del archivo users.txt
	file_content := contenidoUsers(ItemLogin.LoginItem)
	if len(file_content) == 0 {
		fmt.Println("Error, no se pudo recuperar el contenido del archivo users.txt")
		return
	}
	//se separa por saltos de linea y luego por comas
	encounter := false
	encounter2 := false
	contador_usuarios := 1
	array_users := strings.Split(file_content, "\n")
	array_users = array_users[:len(array_users)-1]
	for i := 0; i < len(array_users); i++ {
		array := strings.Split(array_users[i], ",")
		if array[0] != "0" {
			if len(array) == 3 || array[1] == "G" {
				if array[2] == group {
					encounter = true
				}
			} else {
				contador_usuarios++
				if array[3] == username {
					encounter2 = true
					break
				}
			}
		} else {
			if len(array) == 5 || array[1] == "U" {
				contador_usuarios++
			}
		}
	}
	if !encounter {
		fmt.Println("Error, el nombre de grupo no existe")
		cadRespuesta += "Error, el nombre de grupo no existe\n"
		return
	}
	if encounter2 {
		fmt.Println("Error, el nombre de usuario ya existe")
		cadRespuesta += "Error, el nombre de usuario ya existe\n"
		return
	}
	//se agrega el nuevo grupo al archivo users.txt
	//se obtiene el superbloque
	sb, flag := read_sb(ItemLogin.LoginItem.Path, ItemLogin.LoginItem.Part.Part_start)
	if !flag {
		fmt.Println("Error, no se pudo leer el superbloque, mkuser")
		return
	}
	nuevo_usuario := strconv.Itoa(contador_usuarios)
	nuevo_usuario += ",U," + group + "," + username + "," + password + "\n"
	//se obtiene la posicion del inodo del archivo users
	c := string(sb.S_inode_start[:])
	c = strings.TrimRight(c, "\x00")
	posicion_inodo, err := strconv.Atoi(c)
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, mkuser")
		return
	}
	//se sabe que el inodo de users.txt es el uno
	posicion_inodo += int(unsafe.Sizeof(Inodo{}))

	//se obtiene el inodo del archivo users.txt
	inodo, flag := read_inodo(ItemLogin.LoginItem.Path, posicion_inodo)
	if !flag {
		fmt.Println("Error, no se pudo leer el inodo de users.txt, mkuser")
		return
	}
	//se obtiene la posicion inicial de bloques
	c = string(sb.S_block_start[:])
	c = strings.TrimRight(c, "\x00")
	posicion_bloque_inicio, err := strconv.Atoi(c)
	if err != nil {
		fmt.Println("Error, no se pudo leer la posicion inicial de bloques, mkuser")
		return
	}
	//se obtiene el bloque del archivo users.txt
	for i := 0; i < 16; i++ {
		if inodo.I_block[i][0] == 0 {
			break
		}
		d := string(inodo.I_block[i][:])
		d = strings.TrimRight(d, "\x00")
		posicion_bloque, err := strconv.Atoi(d)
		if err != nil {
			fmt.Println("Error, no se pudo obtener la posicion del bloque, mkuser")
			return
		}
		//se obtiene el bloque archivo
		bloque, flag := read_b_archivo(ItemLogin.LoginItem.Path, posicion_bloque_inicio+(int(unsafe.Sizeof(B_Archivo{}))*posicion_bloque))
		if !flag {
			fmt.Println("Error, no se pudo leer el bloque del archivo users.txt, mkuser")
			return
		}
		temporal_content := string(bloque.B_content[:])
		temporal_content = strings.TrimRight(temporal_content, "\x00")
		if i != 15 {

			if inodo.I_block[i+1][0] == 0 {
				//se lee el contenido del bloque
				if len(temporal_content)+len(nuevo_usuario) >= 64 {
					//se crea un nuevo bloque
					new_bloque := B_Archivo{}
					copy(new_bloque.B_content[:], nuevo_usuario)
					//se escribe el nuevo bloque
					c1 := string(sb.S_first_blo[:])
					c1 = strings.TrimRight(c1, "\x00")
					posicion_free_block, err1 := strconv.Atoi(c1)
					if err1 != nil {
						fmt.Println("Error posicionando el puntero, ba")
						return
					}
					flag := write_b_archivo(new_bloque, posicion_bloque_inicio+(int(unsafe.Sizeof(B_Archivo{}))*posicion_free_block))
					if !flag {
						fmt.Println("Error, no se pudo escribir el nuevo bloque, mkuser")
						return
					}
					//se escribe en el bm de bloques y actualiza el superbloque
					flag = write_bitmap_bloques()
					if !flag {
						fmt.Println("Error, no se pudo escribir el bm de bloques, mkuser")
						return
					}

					//se actualiza el inodo
					copy(inodo.I_block[i+1][:], strconv.Itoa(posicion_free_block))
					c1 = strconv.Itoa(len(file_content + nuevo_usuario))
					copy(inodo.I_size[:], c1)
					flag = write_inodo(inodo, posicion_inodo)
					if !flag {
						fmt.Println("Error, no se pudo escribir el inodo, mkuser")
						return
					}
					fmt.Println("Usuario ingresado con exito")
					cadRespuesta += "Usuario ingresado con exito\n"
					return
				} else {
					//se escribe el nuevo grupo en el bloque
					copy(bloque.B_content[:], temporal_content+nuevo_usuario)
					//se escribe el bloque
					flag := write_b_archivo(bloque, posicion_bloque_inicio+(int(unsafe.Sizeof(B_Archivo{}))*posicion_bloque))
					if !flag {
						fmt.Println("Error, no se pudo escribir el bloque, mkuser")
						return
					}
					//se actualiza el inodo
					c1 := strconv.Itoa(len(file_content + nuevo_usuario))
					copy(inodo.I_size[:], c1)
					flag = write_inodo(inodo, posicion_inodo)
					if !flag {
						fmt.Println("Error, no se pudo escribir el inodo, mkgrp")
						return
					}
					fmt.Println("Usuario ingresado con exito")
					cadRespuesta += "Usuario ingresado con exito\n"
					return
				}
			}
		} else {
			if len(temporal_content)+len(nuevo_usuario) >= 64 {
				fmt.Println("Error, no se pudo ingresar el grupo, mkuser")
				return
			} else {
				//se escribe el nuevo grupo en el bloque
				copy(bloque.B_content[:], temporal_content+nuevo_usuario)
				//se escribe el bloque
				flag := write_b_archivo(bloque, posicion_bloque_inicio+(int(unsafe.Sizeof(B_Archivo{}))*posicion_bloque))
				if !flag {
					fmt.Println("Error, no se pudo escribir el bloque, mkuser")
					return
				}
				//se actualiza el inodo
				c1 := strconv.Itoa(len(file_content + nuevo_usuario))
				copy(inodo.I_size[:], c1)
				flag = write_inodo(inodo, posicion_inodo)
				if !flag {
					fmt.Println("Error, no se pudo escribir el inodo, mkgrp")
					return
				}
				fmt.Println("Usuario ingresado con exito")
				cadRespuesta += "Usuario ingresado con exito\n"
				return
			}
		}
	}
}

func rmusr(params []string) {
	if len(PartMount) == 0 {
		fmt.Println("Error, no hay particiones montadas, rmusr")
		cadRespuesta += "Error, no hay particiones montadas\n"
		return
	}
	if !ItemLogin.Iniciado {
		fmt.Println("Error, no hay sesion iniciada, rmusr")
		cadRespuesta += "Error, no hay sesion iniciada\n"
		return
	}
	if !ItemLogin.Admin {
		fmt.Println("Error, el usuario no puede ejecutar este comando, rmusr")
		cadRespuesta += "Error, el usuario no puede ejecutar este comando\n"
		return
	}

	var username string
	for i := 0; i < len(params); i++ {
		array := strings.Split(params[i], "=")
		param := strings.ToLower(array[0])
		if param == ">user" {
			username = array[1]
		} else {
			fmt.Println("El parametro ingresdo no es valido")
			cadRespuesta += "El parametro ingresado no es valido\n"
			return
		}
	}
	if username == "" {
		fmt.Println("Error parametro obligatorio vacio")
		cadRespuesta += "Error parametro obligatorio vacio\n"
		return
	}

	//se obtiene el superbloque
	sb, flag := read_sb(ItemLogin.LoginItem.Path, ItemLogin.LoginItem.Part.Part_start)
	if !flag {
		fmt.Println("Error, no se pudo leer el superbloque, rmusr")
		return
	}
	//se obtiene la posicion del inodo del archivo users
	c := string(sb.S_inode_start[:])
	c = strings.TrimRight(c, "\x00")
	posicion_inodo, err := strconv.Atoi(c)
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, rmusr")
		return
	}
	//se sabe que el inodo de users.txt es el uno
	posicion_inodo += int(unsafe.Sizeof(Inodo{}))
	//se obtiene la posicion del comienzo de los bloques
	c = string(sb.S_block_start[:])
	c = strings.TrimRight(c, "\x00")
	posicion_inicial_bloque, err := strconv.Atoi(c)
	if err != nil {
		fmt.Println("Error, no se pudo leer la posicion inicial de bloques, rmusr")
		return
	}

	//se obtiene el inodo del archivo users.txt
	inodo, flag := read_inodo(ItemLogin.LoginItem.Path, posicion_inodo)
	if !flag {
		fmt.Println("Error, no se pudo leer el inodo de users.txt, rmusr")
		return
	}
	//se va a buscar de bloque en bloque
	encontrado := false
	for i := 0; i < 16; i++ {
		if inodo.I_block[i][0] == 0 {
			break
		}
		//se obtiene la posicion del bloque
		d := string(inodo.I_block[i][:])
		d = strings.TrimRight(d, "\x00")
		posicion_bloque, err := strconv.Atoi(d)
		if err != nil {
			fmt.Println("Error, no se pudo obtener la posicion del bloque, rmusr")
			return
		}
		//se obtiene el bloque archivo
		bloque, flag := read_b_archivo(ItemLogin.LoginItem.Path, posicion_inicial_bloque+(int(unsafe.Sizeof(B_Archivo{}))*posicion_bloque))
		if !flag {
			fmt.Println("Error, no se pudo leer el bloque del archivo users.txt, rmusr")
			return
		}
		temporal_content := string(bloque.B_content[:])
		temporal_content = strings.TrimRight(temporal_content, "\x00")

		encontrado = strings.Contains(temporal_content, ("," + username + ","))
		if encontrado {
			index := strings.Index(temporal_content, ("," + username + ","))
			lastNewline := strings.LastIndex(temporal_content[:index], "\n")
			// Convertir la cadena en un slice de bytes
			if lastNewline == -1 {
				b := []byte(temporal_content)
				// Modificar el elemento en el índice correspondiente
				b[0] = '0'
				// Convertir de nuevo a cadena
				temporal_content = string(b)
				//se actualiza el bloque
				copy(bloque.B_content[:], temporal_content)
				flag = write_b_archivo(bloque, posicion_inicial_bloque+(int(unsafe.Sizeof(B_Archivo{}))*posicion_bloque))
				if !flag {
					fmt.Println("Error, no se pudo escribir el bloque, rmusr")
					return
				}

			} else {
				b := []byte(temporal_content)
				// Modificar el elemento en el índice correspondiente
				b[lastNewline+1] = '0'
				// Convertir de nuevo a cadena
				temporal_content = string(b)
				//se actualiza el bloque
				copy(bloque.B_content[:], temporal_content)
				flag = write_b_archivo(bloque, posicion_inicial_bloque+(int(unsafe.Sizeof(B_Archivo{}))*posicion_bloque))
				if !flag {
					fmt.Println("Error, no se pudo escribir el bloque, rmusr")
					return
				}
			}

			fmt.Println("Usuario eliminado con exito")
			cadRespuesta += "Usuario eliminado con exito\n"
			break
		}

	}
	if !encontrado {
		fmt.Println("Error, el usuario a eliminar no esta registrado")
		cadRespuesta += "Error, el usuario a eliminar no esta registrado\n"
		return
	}
}

func rmgrp(params []string) {
	if len(PartMount) == 0 {
		fmt.Println("Error, no hay particiones montadas, rmgrp")
		cadRespuesta += "Error, no hay particiones montadas\n"
		return
	}
	if !ItemLogin.Iniciado {
		fmt.Println("Error, no hay sesion iniciada, rmgrp")
		cadRespuesta += "Error, no hay sesion iniciada\n"
		return
	}
	if !ItemLogin.Admin {
		fmt.Println("Error, el usuario no puede ejecutar este comando, rmgrp")
		cadRespuesta += "Error, el usuario no puede ejecutar este comando\n"
		return
	}

	var group string
	for i := 0; i < len(params); i++ {
		array := strings.Split(params[i], "=")
		param := strings.ToLower(array[0])
		if param == ">name" {
			group = array[1]
		} else {
			fmt.Println("El parametro ingresado no es valido")
			cadRespuesta += "El parametro ingresado no es valido\n"
			return
		}
	}
	if group == "" {
		fmt.Println("Error parametro obligatorio vacio")
		cadRespuesta += "Error parametro obligatorio vacio\n"
		return
	}

	//se obtiene el superbloque
	sb, flag := read_sb(ItemLogin.LoginItem.Path, ItemLogin.LoginItem.Part.Part_start)
	if !flag {
		fmt.Println("Error, no se pudo leer el superbloque, rmgrp")
		return
	}
	//se obtiene la posicion del inodo del archivo users
	c := string(sb.S_inode_start[:])
	c = strings.TrimRight(c, "\x00")
	posicion_inodo, err := strconv.Atoi(c)
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, rmgrp")
		return
	}
	//se sabe que el inodo de users.txt es el uno
	posicion_inodo += int(unsafe.Sizeof(Inodo{}))
	//se obtiene la posicion del comienzo de los bloques
	c = string(sb.S_block_start[:])
	c = strings.TrimRight(c, "\x00")
	posicion_inicial_bloque, err := strconv.Atoi(c)
	if err != nil {
		fmt.Println("Error, no se pudo leer la posicion inicial de bloques, rmgrp")
		return
	}

	//se obtiene el inodo del archivo users.txt
	inodo, flag := read_inodo(ItemLogin.LoginItem.Path, posicion_inodo)
	if !flag {
		fmt.Println("Error, no se pudo leer el inodo de users.txt, rmgrp")
		return
	}
	//se va a buscar de bloque en bloque
	encontrado := false
	for i := 0; i < 16; i++ {
		if inodo.I_block[i][0] == 0 {
			break
		}
		//se obtiene la posicion del bloque
		d := string(inodo.I_block[i][:])
		d = strings.TrimRight(d, "\x00")
		posicion_bloque, err := strconv.Atoi(d)
		if err != nil {
			fmt.Println("Error, no se pudo obtener la posicion del bloque, rmgrp")
			return
		}
		//se obtiene el bloque archivo
		bloque, flag := read_b_archivo(ItemLogin.LoginItem.Path, posicion_inicial_bloque+(int(unsafe.Sizeof(B_Archivo{}))*posicion_bloque))
		if !flag {
			fmt.Println("Error, no se pudo leer el bloque del archivo users.txt, rmgrp")
			return
		}
		temporal_content := string(bloque.B_content[:])
		temporal_content = strings.TrimRight(temporal_content, "\x00")

		encontrado = strings.Contains(temporal_content, ("G," + group))
		if encontrado {
			index := strings.Index(temporal_content, ("G," + group))
			// Convertir la cadena en un slice de bytes
			b := []byte(temporal_content)
			// Modificar el elemento en el índice correspondiente
			b[index-2] = '0'
			// Convertir de nuevo a cadena
			temporal_content = string(b)
			//se actualiza el bloque
			copy(bloque.B_content[:], temporal_content)
			flag = write_b_archivo(bloque, posicion_inicial_bloque+(int(unsafe.Sizeof(B_Archivo{}))*posicion_bloque))
			if !flag {
				fmt.Println("Error, no se pudo escribir el bloque, rmgrp")
				return
			}
			fmt.Println("Grupo eliminado con exito")
			cadRespuesta += "Grupo eliminado con exito\n"
			break
		}

	}
	if !encontrado {
		fmt.Println("Error, el grupo a eliminar no esta registrado")
		cadRespuesta += "Error, el grupo a eliminar no esta registrado\n"
		return
	}

}

func contenidoUsers(item itemMount) string {
	//se obtiene los datos del disco y particion para obtener los usuarios
	//se obtiene el superbloque
	sb, flag := read_sb(item.Path, item.Part.Part_start)
	if !flag {
		fmt.Println("Error, no se pudo leer el superbloque, cu")
		return ""
	}
	//se obtiene la posicion del inodo del archivo users
	c := string(sb.S_inode_start[:])
	c = strings.TrimRight(c, "\x00")
	posicion_inodo, err := strconv.Atoi(c)
	if err != nil {
		fmt.Println("Error, no se pudo leer el inodo raiz, cu")
		return ""
	}
	//se sabe que el inodo de users.txt es el uno
	posicion_inodo += int(unsafe.Sizeof(Inodo{}))
	c = string(sb.S_block_start[:])
	c = strings.TrimRight(c, "\x00")
	posicion_inicial_bloque, err := strconv.Atoi(c)
	if err != nil {
		fmt.Println("Error, no se pudo leer la posicion inicial de bloques, cu")
		return ""
	}

	//se obtiene el inodo del archivo users.txt
	inodo, flag := read_inodo(item.Path, posicion_inodo)
	if !flag {
		fmt.Println("Error, no se pudo leer el inodo de users.txt, cu")
		return ""
	}
	file_content := ""
	//se obtiene el bloque del archivo users.txt
	for i := 0; i < 16; i++ {
		if inodo.I_block[i][0] == 0 {
			break
		}
		d := string(inodo.I_block[i][:])
		d = strings.TrimRight(d, "\x00")
		posicion_bloque, err := strconv.Atoi(d)
		if err != nil {
			fmt.Println("Error, no se pudo obtener la posicion del bloque, cu")
			return ""
		}
		//se obtiene el bloque archivo
		bloque, flag := read_b_archivo(item.Path, posicion_inicial_bloque+(int(unsafe.Sizeof(B_Archivo{}))*posicion_bloque))
		if !flag {
			fmt.Println("Error, no se pudo leer el bloque del archivo users.txt, cu")
			return ""
		}
		temporal_content := string(bloque.B_content[:])
		temporal_content = strings.TrimRight(temporal_content, "\x00")
		file_content += temporal_content
	}
	return file_content
}
