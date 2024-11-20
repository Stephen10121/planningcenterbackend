package main

import (
	"log"
	"net/http"
	"time"

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
			return c.Redirect("/login")
		}
		return c.SendString("Hello, World!")
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		password := c.Cookies("password")
		if password == initializers.Password {
			return c.Redirect("/")
		}
		return c.Render("login", fiber.Map{})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		password := c.FormValue("password")
		if password != initializers.Password {
			return c.Render("login", fiber.Map{
				"Error": "Invalid Password!",
			})
		}

		cookie := new(fiber.Cookie)
		cookie.Name = "password"
		cookie.Value = password
		cookie.Expires = time.Now().Add(24 * time.Hour)

		c.Cookie(cookie)
		return c.Redirect("/")
	})

	app.Get("/logout", func(c *fiber.Ctx) error {
		c.ClearCookie("password")
		return c.Redirect("/login")
	})

	app.Use("/static", filesystem.New(filesystem.Config{
		Root: http.Dir("./assets"),
	}))

	log.Fatal(app.Listen(":3000"))
}
