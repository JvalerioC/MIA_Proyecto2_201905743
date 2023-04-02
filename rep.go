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
			return
		}
	}
	//validar parametros obligatorios
	if name == "" || path == "" || id == "" {
		fmt.Println("Error, parametro obligatorio vacio")
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
		return
	}

	//se crean los directorios para el reporte
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)

	if err != nil {
		fmt.Println("Error, creando el directorio, rep")
		return
	}
	//verificamos si el archivo existe o no, en este caso el reporte
	_, err1 := os.Stat(path)
	if os.IsNotExist(err1) {
		//fmt.Println("Error, ")
	} else {
		fmt.Println("Error, el archivo ya existe")
		return
	}

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
						temporal_size = size_ebr
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
						temporal_size = size_ebr
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
	} else {
		fmt.Println("Error, creando el reporte")
	}
}

func reportSB(item itemMount, path string) {
	//se crea el archivo
	fmt.Println("reporte SB")
}

func reportFile(item itemMount, path string, ruta string) {
	//se crea el archivo
	fmt.Println("reporte file")
}

func reportTree(item itemMount, path string) {
	//se crea el archivo
	fmt.Println("reporte tree")
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
