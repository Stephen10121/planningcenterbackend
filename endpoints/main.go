package endpoints

import (
	"net/http"
	"time"

	"github.com/Stephen10121/planningcenterbackend/event"
	"github.com/Stephen10121/planningcenterbackend/initializers"
	"github.com/labstack/echo/v5"
)

func LoginGetEndpoint(c echo.Context) error {
	password, err := c.Cookie("password")

	if err != nil {
		return c.Render(http.StatusOK, "login", map[string]interface{}{})
	}

	if password.Value == initializers.Password {
		return c.Redirect(301, "/")
	}

	return c.Render(http.StatusOK, "login", map[string]interface{}{})
}

func LoginPostEndpoint(c echo.Context) error {
	password := c.FormValue("password")
	if password != initializers.Password {
		return c.Render(http.StatusOK, "login", map[string]interface{}{
			"Error": "Invalid Password!",
		})
	}

	cookie := new(http.Cookie)
	cookie.Name = "password"
	cookie.Value = password
	cookie.Expires = time.Now().Add(24 * time.Hour)

	c.SetCookie(cookie)
	return c.Redirect(301, "/")
}

func LogoutGetEndpoint(c echo.Context) error {
	cookie := &http.Cookie{
		Name:    "password",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	c.SetCookie(cookie)
	return c.Redirect(301, "/login")
}

func MainGetEndpoint(c echo.Context) error {
	password, err := c.Cookie("password")

	if err != nil {
		return c.Redirect(301, "/login")
	}

	if password.Value != initializers.Password {
		return c.Redirect(301, "/login")
	}

	event.FetchEvents()
	//Test data. Will change later.
	events := []event.EventType{
		{
			InstanceId: "123",
			StartTime:  "string",
			EndTime:    "string",
			Name:       "string",
			Location:   "string",
			Times:      []event.EventTimes{},
			Resources:  []event.ResourceType{},
			Tags:       []event.TagsType{},
		},
		{
			InstanceId: "124",
			StartTime:  "string",
			EndTime:    "string",
			Name:       "string",
			Location:   "string",
			Times:      []event.EventTimes{},
			Resources:  []event.ResourceType{},
			Tags:       []event.TagsType{},
		},
	}

	return c.Render(http.StatusOK, "index", map[string]interface{}{
		"Events": events,
	})
}
