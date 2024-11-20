package main

import (
	"log"
	"net/http"

	"github.com/Stephen10121/planningcenterbackend/initializers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
)

func main() {
	initializers.SetupEnv()

	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		password := c.Cookies("password")
		if password != initializers.Password {
			c.Redirect("/login")
		}
		return c.SendString("Hello, World!")
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{
			"Title": "Hello, World! Test",
		})
	})

	app.Use("/static", filesystem.New(filesystem.Config{
		Root: http.Dir("./assets"),
	}))

	log.Fatal(app.Listen(":3000"))
}
