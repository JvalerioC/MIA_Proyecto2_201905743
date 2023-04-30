package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"
)

func rep(params []string) {
	var name string
	var path string
	var ruta string
	var id string
	for i := 0; i < len(params); i++ {
		array := strings.Split(params[i], "=")
		param := strings.ToLower(array[0])
		if param == ">name" {
			name = array[1]
		} else if param == ">path" {
			path = array[1]
		} else if param == ">id" {
			id = array[1]
		} else if param == ">ruta" {
			ruta = array[1]
		} else {
			fmt.Println("Error, el parametro ingresado no es valido")
			cadRespuesta += "Error, el parametro ingresado no es valido\n"
			return
		}
	}
	//validar parametros obligatorios
	if name == "" || path == "" || id == "" {
		fmt.Println("Error, parametro obligatorio vacio")
		cadRespuesta += "Error, parametro obligatorio vacio\n"
		return
	}
	//se busca el id
	item := itemMount{}
	for i := 0; i < len(PartMount); i++ {
		if PartMount[i].Id == id {
			item = PartMount[i]
			break
		}
	}
	//se valida que exista el id
	if item.Id == "" {
		fmt.Println("Error, el parametro id ingresado para el reporte no existe")
		cadRespuesta += "Error, el parametro id ingresado para el reporte no existe\n"
		return
	}

	//se crean los directorios para el reporte
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)

	if err != nil {
		fmt.Println("Error, creando el directorio, rep")
		return
	}
	/* //verificamos si el archivo existe o no, en este caso el reporte
	_, err1 := os.Stat(path)
	if os.IsNotExist(err1) {
		//fmt.Println("Error, ")
	} else {
		fmt.Println("Error, el archivo ya existe")
		return
	} */

	//se verifica el nombre del reporte para saber que reporte crear
	if name == "disk" {
		reportDisk(item, path)
	} else if name == "sb" {
		reportSB(item, path)
	} else if name == "file" {
		reportFile(item, path, ruta)
	} else if name == "tree" {
		reportTree(item, path)
	} else {
		fmt.Println("Error, el nombre del reporte no es valido")
		cadRespuesta += "Error, el nombre del reporte no es valido\n"
		return
	}
}

func reportDisk(item itemMount, path string) {
	//se crea el archivo
	file, err := os.Create("dot/disk.dot")
	if err != nil {
		fmt.Println("Error, creando el archivo, rep disk")
		return
	}
	dot := ""
	dot += "digraph disk {\n"
	dot += "	tbl [\n"
	dot += "		shape=plaintext\n"
	dot += "		label=<\n"
	dot += "			<table border='0' cellborder='1' cellspacing='1' color='grey'>\n"
	dot += "				<tr>\n"
	dot += "					<td border='0' colspan='3'><b>" + filepath.Base(item.Path) + "</b></td>\n"
	dot += "				</tr>\n"
	dot += "				<tr>\n"
	dot += "        <td border='0'>\n"
	dot += "          <table color='gray' cellspacing='3' cellpadding='5'>\n"
	dot += "            <tr>\n"
	dot += "              <td>MBR</td>"

	mbr, flag := read_MBR(item.Path)
	if !flag {
		fmt.Println("Error, leyendo el mbr, rep disk")
		return
	}

	inicio := int(unsafe.Sizeof(mbr))
	c := string(mbr.Mbr_tamano[:])
	c = strings.TrimRight(c, "\x00")
	disk_size, err := strconv.Atoi(c)
	if err != nil {
		fmt.Println("Error, convirtiendo el tamaño del disco, rep disk")
		return
	}
	var parts []Partition
	//se recorren las particiones
	for i := 0; i < 4; i++ {
		if mbr.Mbr_Partition[i].Part_size[0] != 0 {
			parts = append(parts, mbr.Mbr_Partition[i])
		}
	}

	for i := 0; i < len(parts); i++ {
		c1 := string(mbr.Mbr_Partition[i].Part_start[:])
		c1 = strings.TrimRight(c1, "\x00")
		start_partition, err2 := strconv.Atoi(c1) //variable
		if err2 != nil {
			fmt.Println("Error, convirtiendo el inicio de la particion, rep disk")
			return
		}
		c1 = string(mbr.Mbr_Partition[i].Part_size[:])
		c1 = strings.TrimRight(c1, "\x00")
		size_partition, err := strconv.Atoi(c1) //variable
		if err != nil {
			fmt.Println("Error, convirtiendo el tamano de la particion, rep disk")
			return
		}
		if inicio == start_partition {
			if mbr.Mbr_Partition[i].Part_type[0] == 'E' {
				dot += "              <td border='0' cellpadding='0'>\n"
				dot += "                <table cellborder='0' color='gray' cellspacing='0' cellpadding='5'>\n"
				dot += "                  <tr>\n"
				dot += "                    <td border='1' colspan='8'>Extendida</td>\n"
				dot += "                  </tr>\n"
				dot += "                  <tr>\n"
				//recuperamos el ebr para el reporte
				ebr := EBR{}
				flag := false
				temporal_size := 0
				temporal := mbr.Mbr_Partition[i].Part_start
				for {
					ebr, flag = read_ebr(item.Path, temporal)
					if !flag {
						fmt.Println("Error, leyendo el ebr, rep disk")
						return
					}
					c1 = string(ebr.Part_size[:])
					c1 = strings.TrimRight(c1, "\x00")
					size_ebr, err1 := strconv.Atoi(c1) //variable
					if err1 != nil {
						fmt.Println("Error, convirtiendo el tamaño del ebr, rep disk")
						return
					}
					if size_ebr == 0 {
						espace := (float64(size_partition * 100.0)) / float64(disk_size)
						dot += "                    <td border='1'>Libre"
						dot += "<br/>" + strconv.Itoa(int(math.Round(espace))) + "%</td>\n"
						temporal_size = size_partition
					} else {
						dot += "                    <td border='1'>EBR</td>\n"
						espace := (float64(size_ebr * 100.0)) / float64(disk_size)
						dot += "                    <td border='1'>Logica"
						dot += "<br/>" + strconv.Itoa(int(math.Round(espace))) + "%</td>\n"
						temporal_size += size_ebr
					}
					temporal = ebr.Part_next
					if temporal[0] == 0 {
						break
					}
				}
				if size_partition > temporal_size {
					espace := (float64((size_partition - temporal_size) * 100.0)) / float64(disk_size)
					dot += "                    <td border='1'>Libre"
					dot += "<br/>" + strconv.Itoa(int(math.Round(espace))) + "%</td>\n"
				}
				dot += "                  </tr>\n"
				dot += "                </table>\n"
				dot += "              </td>\n"
			} else {
				espace := (float64(size_partition * 100.0)) / float64(disk_size)
				dot += "              <td border='1'>Particion " + strconv.Itoa(i+1)
				dot += "<br/>" + strconv.Itoa(int(math.Round(espace))) + "%</td>\n"
			}
			inicio = inicio + size_partition
		} else {
			espace := (float64((start_partition - inicio) * 100.0)) / float64(disk_size)
			dot += "              <td>Libre"
			dot += "<br/>" + strconv.Itoa(int(math.Round(espace))) + "%</td>\n"
			if mbr.Mbr_Partition[i].Part_type[0] == 'E' {
				dot += "              <td border='0' cellpadding='0'>\n"
				dot += "                <table cellborder='0' color='gray' cellspacing='0' cellpadding='5'>\n"
				dot += "                  <tr>\n"
				dot += "                    <td border='1' colspan='8'>Extendida</td>\n"
				dot += "                  </tr>\n"
				dot += "                  <tr>\n"
				//recuperamos el ebr para el reporte
				ebr := EBR{}
				flag := false
				temporal_size := 0
				temporal := mbr.Mbr_Partition[i].Part_start
				for {
					ebr, flag = read_ebr(item.Path, temporal)
					if !flag {
						fmt.Println("Error, leyendo el ebr, rep disk")
						return
					}
					c1 = string(ebr.Part_size[:])
					c1 = strings.TrimRight(c1, "\x00")
					size_ebr, err1 := strconv.Atoi(c1) //variable
					if err1 != nil {
						fmt.Println("Error, convirtiendo el tamaño del ebr, rep disk")
						return
					}
					if size_ebr == 0 {
						espace := (float64(size_partition * 100.0)) / float64(disk_size)
						dot += "                    <td border='1'>Libre"
						dot += "<br/>" + strconv.Itoa(int(math.Round(espace))) + "%</td>\n"
						temporal_size = size_partition
					} else {
						dot += "                    <td border='1'>EBR</td>\n"
						espace := (float64(size_ebr * 100.0)) / float64(disk_size)
						dot += "                    <td border='1'>Logica"
						dot += "<br/>"
						dot += strconv.Itoa(int(math.Round(espace))) + "%</td>\n"
						temporal_size += size_ebr
					}
					temporal = ebr.Part_next
					if temporal[0] == 0 {
						break
					}
				}
				if size_partition > temporal_size {
					espace := (float64((size_partition - temporal_size) * 100.0)) / float64(disk_size)
					dot += "                    <td border='1'>Libre"
					dot += "<br/>" + strconv.Itoa(int(math.Round(espace))) + "%</td>\n"
				}
				dot += "                  </tr>\n"
				dot += "                </table>\n"
				dot += "              </td>\n"
			} else {
				espace := (float64(size_partition * 100.0)) / float64(disk_size)
				dot += "              <td border='1'>Particion " + strconv.Itoa(i+1)
				dot += "<br/>" + strconv.Itoa(int(math.Round(espace))) + "%</td>\n"
			}
			inicio = start_partition + size_partition
		}
	}
	if inicio < disk_size {
		espace := (float64((disk_size - inicio) * 100.0)) / float64(disk_size)
		dot += "              <td>Libre"
		dot += "<br/>" + strconv.Itoa(int(math.Round(espace))) + "%</td>\n"
	}
	dot += "            </tr>\n"
	dot += "          </table>\n"
	dot += "        </td>\n"
	dot += "      </tr>\n"
	dot += "    </table>\n"
	dot += "  >];\n"
	dot += "}\n"

	_, err = io.WriteString(file, dot)
	if err != nil {
		fmt.Println("error al escribir el dot")
		return
	}
	file.Close()
	flag = createReport(dot, "disk", path)
	if flag {
		fmt.Println("Reporte creado con exito")
		cadRespuesta += "Reporte creado con exito\n"
		reportes = append(reportes, path)
	} else {
		fmt.Println("Error, creando el reporte")
		cadRespuesta += "Error, creando el reporte\n"
	}
}

func reportSB(item itemMount, path string) {
	//se crea el archivo
	file, err := os.Create("dot/sb.dot")
	if err != nil {
		fmt.Println("Error, creando el archivo, rep sb")
		return
	}
	//se obtiene el superbloque
	sb, flag := read_sb(item.Path, item.Part.Part_start)
	if !flag {
		fmt.Println("Error, leyendo el superbloque, rep sb")
		return
	}
	dot := ""
	dot += "digraph sb {\n"
	dot += "  some_node [\n"
	dot += "    shape=plaintext\n"
	dot += "    label=<\n"
	dot += "      <table cellpadding='4' cellborder='1' color='grey' cellspacing='1'>\n"
	dot += "        <tr>\n"
	dot += "          <td bgcolor='grey' colspan='2'>Reporte SuperBloque</td>\n"
	dot += "        </tr>\n"
	dot += "        <tr>\n"
	dot += "					<td>s_filesystem_type</td>\n"
	dot += "					<td>EXT2</td>\n"
	dot += "				</tr>\n"
	dot += "				<tr>\n"
	dot += "					<td>s_inodes_count</td>\n"
	dot += "					<td>" + strings.TrimRight(string(sb.S_inodes_count[:]), "\x00") + "</td>\n"
	dot += "				</tr>\n"
	dot += "				<tr>\n"
	dot += "					<td>s_blocks_count</td>\n"
	dot += "					<td>" + strings.TrimRight(string(sb.S_blocks_count[:]), "\x00") + "</td>\n"
	dot += "				</tr>\n"
	dot += "				<tr>\n"
	dot += "					<td>s_free_inodes_count</td>\n"
	dot += "					<td>" + strings.TrimRight(string(sb.S_free_inodes_count[:]), "\x00") + "</td>\n"
	dot += "				</tr>\n"
	dot += "				<tr>\n"
	dot += "					<td>s_free_blocks_count</td>\n"
	dot += "					<td>" + strings.TrimRight(string(sb.S_free_blocks_count[:]), "\x00") + "</td>\n"
	dot += "				</tr>\n"
	dot += "				<tr>\n"
	dot += "					<td>s_mtime</td>\n"
	dot += "					<td>" + strings.TrimRight(string(sb.S_mtime[:]), "\x00") + "</td>\n"
	dot += "				</tr>\n"
	dot += "				<tr>\n"
	dot += "					<td>s_mnt_count</td>\n"
	dot += "					<td>" + strings.TrimRight(string(sb.S_mnt_count[:]), "\x00") + "</td>\n"
	dot += "				</tr>\n"
	dot += "				<tr>\n"
	dot += "					<td>s_magic</td>\n"
	dot += "					<td>" + strings.TrimRight(string(sb.S_magic[:]), "\x00") + "</td>\n"
	dot += "				</tr>\n"
	dot += "				<tr>\n"
	dot += "					<td>s_inode_size</td>\n"
	dot += "					<td>" + strings.TrimRight(string(sb.S_inode_size[:]), "\x00") + "</td>\n"
	dot += "				</tr>\n"
	dot += "				<tr>\n"
	dot += "					<td>s_block_size</td>\n"
	dot += "					<td>" + strings.TrimRight(string(sb.S_block_size[:]), "\x00") + "</td>\n"
	dot += "				</tr>\n"
	dot += "				<tr>\n"
	dot += "					<td>s_first_ino</td>\n"
	dot += "					<td>" + strings.TrimRight(string(sb.S_first_ino[:]), "\x00") + "</td>\n"
	dot += "				</tr>\n"
	dot += "				<tr>\n"
	dot += "					<td>s_first_blo</td>\n"
	dot += "					<td>" + strings.TrimRight(string(sb.S_first_blo[:]), "\x00") + "</td>\n"
	dot += "				</tr>\n"
	dot += "				<tr>\n"
	dot += "					<td>s_bm_inode_start</td>\n"
	dot += "					<td>" + strings.TrimRight(string(sb.S_bm_inode_start[:]), "\x00") + "</td>\n"
	dot += "				</tr>\n"
	dot += "				<tr>\n"
	dot += "					<td>s_bm_block_start</td>\n"
	dot += "					<td>" + strings.TrimRight(string(sb.S_bm_block_start[:]), "\x00") + "</td>\n"
	dot += "				</tr>\n"
	dot += "				<tr>\n"
	dot += "					<td>s_inode_start</td>\n"
	dot += "					<td>" + strings.TrimRight(string(sb.S_inode_start[:]), "\x00") + "</td>\n"
	dot += "				</tr>\n"
	dot += "				<tr>\n"
	dot += "					<td>s_block_start</td>\n"
	dot += "					<td>" + strings.TrimRight(string(sb.S_block_start[:]), "\x00") + "</td>\n"
	dot += "				</tr>\n"
	dot += "      </table>\n"
	dot += "    >];\n"
	dot += "}"

	//se escribe el dot al archivo
	_, err = io.WriteString(file, dot)
	if err != nil {
		fmt.Println("Error, al escribir el dot")
		return
	}
	file.Close()
	flag = createReport(dot, "sb", path)
	if flag {
		fmt.Println("Reporte creado con exito")
		cadRespuesta += "Reporte creado con exito\n"
		reportes = append(reportes, path)
	} else {
		fmt.Println("Error, creando el reporte")
		cadRespuesta += "Error, creando el reporte\n"
	}
}

func reportFile(item itemMount, path string, ruta string) {

	//se lee el superbloque
	sb, flag := read_sb(item.Path, item.Part.Part_start)
	if !flag {
		fmt.Println("Error, leyendo el superbloque, rep file")
		return
	}
	//se crea el array de carpetas a buscar
	array := strings.Split(ruta, "/")
	copy(array[:], array[1:])
	array = array[:len(array)-1]
	//se lee el inodo, desde la raiz
	inodo := Inodo{} // se declara un inodo temporal
	c1 := string(sb.S_inode_start[:])
	c1 = strings.TrimRight(c1, "\x00")
	pos_inodo, err2 := strconv.Atoi(c1) //variable
	if err2 != nil {
		fmt.Println("Error, convirtiendo la posicion del inodo, rep disk")
		return
	}
	pos_start := pos_inodo
	c2 := string(sb.S_block_start[:])
	c2 = strings.TrimRight(c2, "\x00")
	pos_block_start, err1 := strconv.Atoi(c2) //variable
	if err1 != nil {
		fmt.Println("Error, convirtiendo la posicion del inodo, rep disk")
		return
	}
	encounter := false
	temp_inodo, flag := read_inodo(item.Path, pos_inodo)
	if !flag {
		fmt.Println("Error, leyendo el inodo, rep file")
		return
	}
	//se recorre el array de carpetas
	for i := 0; i < len(array); i++ {
		inodo = temp_inodo
		if inodo.I_type[0] != '1' {
			//se recorren los apuntadores directos
			bc := B_Carpeta{}
			flag := false
			for j := 0; j < 16; j++ {
				if inodo.I_block[j][0] == 0 {
					break
				}
				//se obtiene el apuntador
				c1 := string(inodo.I_block[j][:])
				c1 = strings.TrimRight(c1, "\x00")
				pos_block, err2 := strconv.Atoi(c1) //variable
				new_pos_block := pos_block_start + (int(unsafe.Sizeof(B_Carpeta{})) * pos_block)
				if err2 != nil {
					fmt.Println("Error, convirtiendo la posicion del bloque, rep disk")
					return
				}
				//se lee el bloque
				bc, flag = read_b_carpeta(item.Path, new_pos_block)
				if !flag {
					fmt.Println("Error, leyendo el bloque, rep file")
					return
				}
				//se recorren los apuntadores directos
				for k := 0; k < 4; k++ {
					if bc.B_content[k].B_inodo[0] == 0 {
						break
					}

					//se obtiene el apuntador
					c1 := string(bc.B_content[k].B_name[:])
					c1 = strings.TrimRight(c1, "\x00")
					if c1 == array[i] {
						//se obtiene el inodo
						c1 := string(bc.B_content[k].B_inodo[:])
						c1 = strings.TrimRight(c1, "\x00")
						pos_temp, err2 := strconv.Atoi(c1) //variable
						if err2 != nil {
							fmt.Println("Error, convirtiendo el apuntador al inodo, rep disk")
							return
						}
						pos_inodo = pos_start + (int(unsafe.Sizeof(Inodo{})) * pos_temp)
						//se lee el inodo
						temp_inodo, flag = read_inodo(item.Path, pos_inodo)
						if !flag {
							fmt.Println("Error, leyendo el inodo, rep file")
							return
						}
						if temp_inodo.I_type[0] == '1' {
							encounter = true
						}
					}
				}
			}

		}
	}
	//se valida si el archivo fue encontrado, este estara en temp_inodo
	if !encounter {
		fmt.Println("Error, el archivo no fue encontrado, rep file")
		return
	}
	content := "Reporte de archivo: " + array[len(array)-1] + "\n"
	//se recorre los apuntadores para recuperar el contenido del archivo
	for i := 0; i < 16; i++ {
		if temp_inodo.I_block[i][0] == 0 {
			break
		}
		c1 := string(temp_inodo.I_block[i][:])
		c1 = strings.TrimRight(c1, "\x00")
		pos_temp, err2 := strconv.Atoi(c1) //variable
		if err2 != nil {
			fmt.Println("Error, convirtiendo el inicio de la particion, rep disk")
			return
		}
		//se obtiene el bloque archivo
		new_position := pos_block_start + (int(unsafe.Sizeof(B_Archivo{})) * pos_temp)
		arch, flag := read_b_archivo(item.Path, new_position)
		if !flag {
			fmt.Println("Error, leyendo el bloque archivo, rep file")
			return
		}
		//se obtiene el contenido
		c1 = string(arch.B_content[:])
		c1 = strings.TrimRight(c1, "\x00")
		content += c1
	}
	//se crea el archivo
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error, creando el archivo, rep file")
		return
	}
	_, err = io.WriteString(file, content)
	if err != nil {
		fmt.Println("error al escribir el contenido al archivo")
		return
	}
	file.Close()
	fmt.Println("Reporte creado con exito")
	cadRespuesta += "Reporte creado con exito\n"
	reportes = append(reportes, path)
}

func reportTree(item itemMount, path string) {
	//se obtiene el superbloque
	sb, flag := read_sb(item.Path, item.Part.Part_start)
	if !flag {
		fmt.Println("Error, leyendo el superbloque, rep tree")
		return
	}
	//se convierten las posiciones iniciales de bloques e inodos
	c1 := string(sb.S_inode_start[:])
	c1 = strings.TrimRight(c1, "\x00")
	pos_inodo_start, err := strconv.Atoi(c1) //variable
	if err != nil {
		fmt.Println("Error, convirtiendo la posicion del inodo, rep tree")
		return
	}

	//se obtiene el inodo raiz
	inodo, flag1 := read_inodo(item.Path, pos_inodo_start)
	if !flag1 {
		fmt.Println("Error, leyendo el inodo raiz, rep tree")
		return
	}
	dot := "digraph tree {\n"
	dot += "	rankdir=LR;\n"
	dot += "	inode0 [\n"
	dot += "		shape=plaintext\n"
	dot += "		label=<\n"
	dot += "      <table color='grey' cellborder='1' cellspacing='1' cellpadding='4'>\n"
	dot += "        <tr>\n"
	dot += "          <td colspan='2' bgcolor='gray'>Inodo / </td>\n"
	dot += "        </tr>\n"
	dot += "        <tr>\n"
	dot += "          <td>i_uid</td>\n"
	dot += "          <td>" + strings.TrimRight(string(inodo.I_uid[:]), "\x00") + "</td>\n"
	dot += "        </tr>\n"
	dot += "        <tr>\n"
	dot += "          <td>i_gid</td>\n"
	dot += "          <td>" + strings.TrimRight(string(inodo.I_gid[:]), "\x00") + "</td>\n"
	dot += "        </tr>\n"
	dot += "        <tr>\n"
	dot += "          <td>i_size</td>\n"
	dot += "          <td>" + strings.TrimRight(string(inodo.I_size[:]), "\x00") + "</td>\n"
	dot += "        </tr>\n"
	dot += "        <tr>\n"
	dot += "          <td>i_atime</td>\n"
	dot += "          <td>" + strings.TrimRight(string(inodo.I_atime[:]), "\x00") + "</td>\n"
	dot += "        </tr>\n"
	dot += "        <tr>\n"
	dot += "          <td>i_ctime</td>\n"
	dot += "          <td>" + strings.TrimRight(string(inodo.I_ctime[:]), "\x00") + "</td>\n"
	dot += "        </tr>\n"
	dot += "        <tr>\n"
	dot += "          <td>i_mtime</td>\n"
	dot += "          <td>" + strings.TrimRight(string(inodo.I_mtime[:]), "\x00") + "</td>\n"
	dot += "        </tr>\n"
	dot += "        <tr>\n"
	dot += "          <td>i_type</td>\n"
	dot += "          <td>" + strings.TrimRight(string(inodo.I_type[:]), "\x00") + "</td>\n"
	dot += "        </tr>\n"
	dot += "        <tr>\n"
	dot += "          <td>i_perm</td>\n"
	dot += "          <td>" + strings.TrimRight(string(inodo.I_perm[:]), "\x00") + "</td>\n"
	dot += "        </tr>\n"
	for i := 0; i < 16; i++ {
		dot += "        <tr>\n"
		dot += "          <td>ap" + strconv.Itoa(i) + "</td>\n"
		dot += "          <td port='b" + strconv.Itoa(i) + "' >" + strings.TrimRight(string(inodo.I_block[i][:]), "\x00") + "</td>\n"
		dot += "        </tr>\n"
	}
	dot += "      </table>\n"
	dot += "    >\n"
	dot += "	];\n"
	name_inodo := "inode0"
	counter := 0
	for i := 0; i < 16; i++ {
		if inodo.I_block[i][0] == 0 {
			continue
		} else {
			c := string(inodo.I_block[i][:])
			c = strings.TrimRight(c, "\x00")
			poss, err1 := strconv.Atoi(c) //variable
			if err1 != nil {
				fmt.Println("Error, convirtiendo la posicion del inodo, rep tree")
				return
			}
			if inodo.I_type[0] == '0' {
				dot += printBloqueContenido(sb, item.Path, poss, name_inodo+":b"+strconv.Itoa(i), &counter)
			} else {
				dot += printBloqueArchivo(sb, item.Path, poss, name_inodo+":b"+strconv.Itoa(i), &counter)
			}
		}
	}
	dot += "}"

	//se crea el archivo
	file, err := os.Create("dot/tree.dot")
	if err != nil {
		fmt.Println("Error, creando el archivo, rep tree")
		return
	}
	//se escribe el archivo
	_, err = io.WriteString(file, dot)
	if err != nil {
		fmt.Println("error al escribir el dot")
		return
	}
	file.Close()
	flag = createReport(dot, "tree", path)
	if flag {
		fmt.Println("Reporte creado con exito")
		cadRespuesta += "Reporte creado con exito\n"
		reportes = append(reportes, path)
	} else {
		fmt.Println("Error, creando el reporte")
		cadRespuesta += "Error, creando el reporte\n"
	}
}

// funcion para agregar al dot un bloque de contenido
func printBloqueContenido(sb SuperBloque, path string, pos int, name_padre string, counter *int) string {
	c2 := string(sb.S_block_start[:])
	c2 = strings.TrimRight(c2, "\x00")
	pos_block_start, err1 := strconv.Atoi(c2) //variable
	if err1 != nil {
		fmt.Println("Error, convirtiendo la posicion del inicio de bloques, rep tree")
		return ""
	}
	dot := ""
	//se lee el bloque de contenido
	bloqueContenido, flag := read_b_carpeta(path, pos_block_start+(int(unsafe.Sizeof(B_Carpeta{}))*pos))
	if !flag {
		fmt.Println("Error, leyendo el archivo, rep tree, pbc")
		return ""
	}
	//se agrega al dot la informacion del bloque
	dot += "	bloque" + strconv.Itoa(*counter) + " [\n"
	dot += "		shape=plaintext\n"
	dot += "		label=<\n"
	dot += "      <table color='grey' cellborder='1' cellspacing='1' cellpadding='4'>\n"
	dot += "        <tr>\n"
	dot += "          <td colspan='2' bgcolor='gray'>Bloque Carpeta</td>\n"
	dot += "        </tr>\n"
	for i := 0; i < 4; i++ {
		dot += "        <tr>\n"
		dot += "          <td>" + strings.TrimRight(string(bloqueContenido.B_content[i].B_name[:]), "\x00") + "</td>\n"
		dot += "          <td port='b" + strconv.Itoa(i) + "' >" + strings.TrimRight(string(bloqueContenido.B_content[i].B_inodo[:]), "\x00") + "</td>\n"
		dot += "        </tr>\n"
	}
	dot += "      </table>\n"
	dot += "    >\n"
	dot += "	];\n"
	//se crea la conexion con el inodo
	nuevo_padre := "bloque" + strconv.Itoa(*counter)
	dot += name_padre + "->" + nuevo_padre + ";\n"
	*counter++
	//se recorren los apuntadores
	for i := 0; i < 4; i++ {
		if bloqueContenido.B_content[i].B_inodo[0] == 0 {
			continue
		} else if bloqueContenido.B_content[i].B_inodo[0] == '0' {
			continue
		} else {
			c := string(bloqueContenido.B_content[i].B_inodo[:])
			c = strings.TrimRight(c, "\x00")
			poss, err1 := strconv.Atoi(c) //variable
			if err1 != nil {
				fmt.Println("Error, convirtiendo la posicion del inodo, rep tree")
				return ""
			}
			dot += printInodo(sb, path, poss, nuevo_padre+":b"+strconv.Itoa(i), counter, strings.TrimRight(string(bloqueContenido.B_content[i].B_name[:]), "\x00"))
		}
	}
	return dot
}

// funcion para agregar al dot un bloque de archivo
func printBloqueArchivo(sb SuperBloque, path string, pos int, name_padre string, counter *int) string {
	c2 := string(sb.S_block_start[:])
	c2 = strings.TrimRight(c2, "\x00")
	pos_block_start, err1 := strconv.Atoi(c2) //variable
	if err1 != nil {
		fmt.Println("Error, convirtiendo la posicion del inicio de bloques, rep tree")
		return ""
	}
	dot := ""
	//se lee el bloque de contenido
	bloqueArchivo, flag := read_b_archivo(path, pos_block_start+(int(unsafe.Sizeof(B_Archivo{}))*pos))
	if !flag {
		fmt.Println("Error, leyendo el archivo, rep tree, pba")
		return ""
	}
	//se agrega al dot la informacion del bloque
	dot += "	bloque" + strconv.Itoa(*counter) + " [\n"
	dot += "		shape=plaintext\n"
	dot += "		label=<\n"
	dot += "      <table color='grey' cellborder='1' cellspacing='1' cellpadding='4'>\n"
	dot += "        <tr>\n"
	dot += "          <td colspan='2' bgcolor='gray'>Bloque Archivo</td>\n"
	dot += "        </tr>\n"
	dot += "        <tr>\n"
	dot += "          <td colspan='2'>" + strings.TrimRight(string(bloqueArchivo.B_content[:]), "\x00") + "</td>\n"
	dot += "        </tr>\n"
	dot += "      </table>\n"
	dot += "    >\n"
	dot += "	];\n"
	//se crea la conexion con el inodo
	nuevo_padre := "bloque" + strconv.Itoa(*counter)
	dot += name_padre + "->" + nuevo_padre + ";\n"
	*counter++
	return dot
}

// funcion para agregar al dot un inodo
func printInodo(sb SuperBloque, path string, pos int, name_padre string, counter *int, inode_name string) string {
	c := string(sb.S_inode_start[:])
	c = strings.TrimRight(c, "\x00")
	pos_inode_start, err1 := strconv.Atoi(c) //variable
	if err1 != nil {
		fmt.Println("Error, convirtiendo la posicion del inicio de inodos, rep tree")
		return ""
	}
	dot := ""
	//se lee el inodo
	inodo, flag := read_inodo(path, pos_inode_start+(int(unsafe.Sizeof(Inodo{}))*pos))
	if !flag {
		fmt.Println("Error, leyendo el archivo, rep tree, pi")
		return ""
	}
	//se agrega al dot la informacion del inodo
	dot += "	inodo" + strconv.Itoa(*counter) + " [\n"
	dot += "		shape=plaintext\n"
	dot += "		label=<\n"
	dot += "      <table color='grey' cellborder='1' cellspacing='1' cellpadding='4'>\n"
	dot += "        <tr>\n"
	dot += "          <td colspan='2' bgcolor='gray'>Inodo " + inode_name + "</td>\n"
	dot += "        </tr>\n"
	dot += "        <tr>\n"
	dot += "          <td>i_uid</td>\n"
	dot += "          <td>" + strings.TrimRight(string(inodo.I_uid[:]), "\x00") + "</td>\n"
	dot += "        </tr>\n"
	dot += "        <tr>\n"
	dot += "          <td>i_gid</td>\n"
	dot += "          <td>" + strings.TrimRight(string(inodo.I_gid[:]), "\x00") + "</td>\n"
	dot += "        </tr>\n"
	dot += "        <tr>\n"
	dot += "          <td>i_size</td>\n"
	dot += "          <td>" + strings.TrimRight(string(inodo.I_size[:]), "\x00") + "</td>\n"
	dot += "        </tr>\n"
	dot += "        <tr>\n"
	dot += "          <td>i_atime</td>\n"
	dot += "          <td>" + strings.TrimRight(string(inodo.I_atime[:]), "\x00") + "</td>\n"
	dot += "        </tr>\n"
	dot += "        <tr>\n"
	dot += "          <td>i_ctime</td>\n"
	dot += "          <td>" + strings.TrimRight(string(inodo.I_ctime[:]), "\x00") + "</td>\n"
	dot += "        </tr>\n"
	dot += "        <tr>\n"
	dot += "          <td>i_mtime</td>\n"
	dot += "          <td>" + strings.TrimRight(string(inodo.I_mtime[:]), "\x00") + "</td>\n"
	dot += "        </tr>\n"
	dot += "        <tr>\n"
	dot += "          <td>i_type</td>\n"
	dot += "          <td>" + strings.TrimRight(string(inodo.I_type[:]), "\x00") + "</td>\n"
	dot += "        </tr>\n"
	dot += "        <tr>\n"
	dot += "          <td>i_perm</td>\n"
	dot += "          <td>" + strings.TrimRight(string(inodo.I_perm[:]), "\x00") + "</td>\n"
	dot += "        </tr>\n"
	for i := 0; i < 16; i++ {
		dot += "        <tr>\n"
		dot += "          <td>ap" + strconv.Itoa(i) + "</td>\n"
		dot += "          <td port='b" + strconv.Itoa(i) + "' >" + strings.TrimRight(string(inodo.I_block[i][:]), "\x00") + "</td>\n"
		dot += "        </tr>\n"
	}
	dot += "      </table>\n"
	dot += "    >\n"
	dot += "	];\n"
	//se crea la conexion con el padre
	nuevo_padre := "inodo" + strconv.Itoa(*counter)
	dot += name_padre + "->" + nuevo_padre + ";\n"
	*counter++
	//se recorren los apuntadores
	for i := 0; i < 16; i++ {
		if inodo.I_block[i][0] == 0 {
			continue
		} else {
			c := string(inodo.I_block[i][:])
			c = strings.TrimRight(c, "\x00")
			poss, err1 := strconv.Atoi(c) //variable
			if err1 != nil {
				fmt.Println("Error, convirtiendo la posicion del inodo, rep tree")
				return ""
			}
			if inodo.I_type[0] == '0' {
				dot += printBloqueContenido(sb, path, poss, nuevo_padre+":b"+strconv.Itoa(i), counter)
			} else {
				dot += printBloqueArchivo(sb, path, poss, nuevo_padre+":b"+strconv.Itoa(i), counter)
			}
		}
	}
	return dot
}

// funcion para generar el la imagen, pdf o jpg
func createReport(dot string, name string, path string) bool {
	//se busca el nombre del dot
	namedot := "dot/" + name + ".dot"
	//se obtiene la extension del archivo
	extension := filepath.Ext(path)
	//se obtiene el contenido del dot
	dotfile, err := ioutil.ReadFile(namedot)
	if err != nil {
		fmt.Println("Error, al leer el archivo dot")
		return false
	}
	formato := ""
	if extension == ".png" {
		formato = "-Tpng"
	} else if extension == ".jpg" {
		formato = "-Tjpg"
	} else if extension == ".pdf" {
		formato = "-Tpdf"
	} else {
		fmt.Println("Error, el formato del reporte no es valido")
		return false
	}

	cmd := exec.Command("dot", formato)
	cmd.Stdin = bytes.NewReader(dotfile)
	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error, al ejecutar el comando dot")
		return false
	}
	//se crea el archivo
	err = ioutil.WriteFile(path, out.Bytes(), os.ModePerm)
	if err != nil {
		fmt.Println("Error, al crear el archivo")
		return false
	}
	return true

}
