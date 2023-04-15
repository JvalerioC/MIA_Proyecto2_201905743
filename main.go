package main

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var PartMount []itemMount
var disc_counter int
var ItemLogin Usuario
var cadRespuesta string

func main() {

	app := fiber.New()
	app.Use(cors.New())

	fmt.Println("**********************************************************")
	fmt.Println("**                                                      **")
	fmt.Println("**               201905743  -  PROYECTO 2               **")
	fmt.Println("**                                                      **")
	fmt.Println("**********************************************************")

	/* app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	})) */

	app.Post("/execute", func(c *fiber.Ctx) error {
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
			"message": "Exito",
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var requestBody map[string]interface{}
		if err := c.BodyParser(&requestBody); err != nil {
			return err
		}
		user := requestBody["username"].(string)
		pass := requestBody["password"].(string)
		id := requestBody["id"].(string)
		fmt.Println(user, pass, id)
		login2(user, pass, id)
		return c.JSON(&fiber.Map{
			"message": cadRespuesta,
		})
	})

	app.Listen(":3000")
	fmt.Println("Server on port 3000")

	/*  */
}
