package main

import (
	"html/template"
	"log"
	"os"

	"github.com/Stephen10121/planningcenterbackend/endpoints"
	"github.com/Stephen10121/planningcenterbackend/initializers"
	"github.com/Stephen10121/planningcenterbackend/setup"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	initializers.SetupEnv()

	pocketbase := pocketbase.New()

	pocketbase.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.Renderer = &setup.Template{
			Templates: template.Must(template.ParseGlob("templates/*.html")),
		}

		e.Router.GET("/", endpoints.MainGetEndpoint)
		e.Router.GET("/login", endpoints.LoginGetEndpoint)
		e.Router.POST("/login", endpoints.LoginPostEndpoint)
		e.Router.GET("/logout", endpoints.LogoutGetEndpoint)
		e.Router.GET("/static/*", apis.StaticDirectoryHandler(os.DirFS("./assets"), false))

		return nil
	})

	log.Println("Server started on: http://localhost:8090")
	if err := pocketbase.Start(); err != nil {
		log.Fatal(err)
	}
}
