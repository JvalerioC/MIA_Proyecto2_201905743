package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	flag := true
	fmt.Println("**********************************************************")
	fmt.Println("**                                                      **")
	fmt.Println("**               201905743  -  PROYECTO 2               **")
	fmt.Println("**                                                      **")
	fmt.Println("**********************************************************")

	for flag {
		fmt.Print("201905743@P2:~$ ")
		reader := bufio.NewReader(os.Stdin)
		entrada, _ := reader.ReadString('\n')
		entrada = strings.TrimRight(entrada, "\r\n")
		if entrada == "exit" {
			flag = false
			fmt.Println("Saliendo de la aplicacion...")
			continue
		} else if entrada[0] == '#' {
			fmt.Println(entrada)
			continue
		}
		analizador(entrada)
	}
}
