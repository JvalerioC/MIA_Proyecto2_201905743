package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var PartMount []itemMount
var disc_counter int
var ItemLogin Usuario
var cadRespuesta string
var reportes []string

func ReadFileToBase64(filePath string) (string, error) {
	file, err := os.Open(filePath)
	fmt.Println(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a buffer to store the file content
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)

	// Read the file content into the buffer
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	// Convert the file content to Base64
	encoded := base64.StdEncoding.EncodeToString(buffer)
	return encoded, nil
}

func main() {

	app := fiber.New()
	app.Use(cors.New())

	/* fmt.Println("**********************************************************")
	fmt.Println("**                                                      **")
	fmt.Println("**               201905743  -  PROYECTO 2               **")
	fmt.Println("**                                                      **")
	fmt.Println("**********************************************************")
	*/
	/* app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	})) */

	app.Post("/execute", func(c *fiber.Ctx) error {
		cadRespuesta = ""
		var requestBody map[string]interface{}
		if err := c.BodyParser(&requestBody); err != nil {
			return err
		}
		cadena := requestBody["fileContent"].(string)
		//flag := false
		/*
			fmt.Println("**********************************************************")
			fmt.Println("**                                                      **")
			fmt.Println("**               201905743  -  PROYECTO 2               **")
			fmt.Println("**                                                      **")
			fmt.Println("**********************************************************")
		*/
		//se hace split por salto de linea
		array := strings.Split(cadena, "\n")
		//se recorre el array
		for _, line := range array {
			/* fmt.Print("201905743@P2:~$ ")
			reader := bufio.NewReader(os.Stdin)
			entrada, _ := reader.ReadString('\n')
			entrada = strings.TrimRight(entrada, "\r\n") */
			entrada := line
			entrada = strings.TrimRight(entrada, "\r\n")
			entrada = strings.TrimRight(entrada, " ")
			if entrada == "" {
				continue
			}
			if entrada == "exit" {
				fmt.Println("Saliendo de la aplicacion...")
				break
			} else if entrada[0] == '#' {
				fmt.Println(entrada)
				continue
			} else if entrada == "pause" {
				fmt.Println("Presione enter para continuar...")
				fmt.Scanln()
				continue
			}
			analizador(entrada)
		}
		return c.JSON(&fiber.Map{
			"message": cadRespuesta,
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		cadRespuesta = ""
		var requestBody map[string]interface{}
		if err := c.BodyParser(&requestBody); err != nil {
			return err
		}
		user := requestBody["username"].(string)
		pass := requestBody["password"].(string)
		id := requestBody["id"].(string)
		fmt.Println(user, pass, id)
		status := login2(user, pass, id)
		return c.JSON(&fiber.Map{
			"message": cadRespuesta,
			"status":  status,
		})
	})

	app.Get("/reports", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{
			"message": reportes,
		})
	})

	app.Get("/reports/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		fmt.Println(id)
		for i := 0; i < len(reportes); i++ {
			d := strconv.Itoa(i)
			if d == id {
				fmt.Println("si son iguales")
				encodedContent, err := ReadFileToBase64(reportes[i])
				if err != nil {
					return c.JSON(&fiber.Map{
						"message": "Error",
					})
				}
				return c.JSON(&fiber.Map{
					"message": encodedContent,
				})
			}
		}
		return c.JSON(&fiber.Map{
			"message": "No se encontro el reporte",
		})
	})

	app.Listen(":3000")
	fmt.Println("Server on port 3000")

	/*  */
}
