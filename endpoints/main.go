package endpoints

import (
	"fmt"
	"time"

	"github.com/Stephen10121/planningcenterbackend/event"
	"github.com/Stephen10121/planningcenterbackend/initializers"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func WebhookTest(e *core.ServeEvent) {
	e.Router.POST("/webhook", func(c *core.RequestEvent) error {
		fmt.Println("Endpoint hit")
		fmt.Println(c.Request.Body)
		return c.JSON(200, map[string]string{
			"msg": "All Good!",
		})
	})
}

func TestEndpoint(e *core.ServeEvent) {
	e.Router.GET("/test", func(c *core.RequestEvent) error {
		event.NewEventFetcher()
		return c.JSON(200, map[string]string{
			"msg": "All Good",
		})
	})
}

func GetEvents(e *core.ServeEvent, base *pocketbase.PocketBase) {
	e.Router.GET("/events", func(c *core.RequestEvent) error {
		auth := c.Request.Header.Get("Authorization")

		if auth != "Basic "+initializers.Credentials+"==" {
			return c.JSON(401, map[string]string{
				"error": "Invalid Credentials",
			})
		}
		d := time.Now().Add(-72 * time.Hour)

		records, err := base.FindAllRecords("events",
			dbx.NewExp("startTime >= {:filterDate}", dbx.Params{"filterDate": d}),
		)
		// records, err := base.Dao().FindRecordsByExpr("events",
		// 	dbx.NewExp("startTime >= {:filterDate}", dbx.Params{"filterDate": d}),
		// )

		if err != nil {
			base.Logger().Error(
				"Failed to fetch events!",
				"id", 123,
				"error", err,
			)
			return c.String(500, "Failed to fetch events from the database.")
		}

		events := []event.EventType{}

		for i := 0; i < len(records); i++ {
			var times []event.EventTimes
			var resources []event.ResourceType
			var tags []event.TagsType

			err := records[i].UnmarshalJSONField("times", &times)
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = records[i].UnmarshalJSONField("resources", &resources)
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = records[i].UnmarshalJSONField("tags", &tags)
			if err != nil {
				fmt.Println(err)
				continue
			}

			events = append(events, event.EventType{
				InstanceId: records[i].Id,
				StartTime:  records[i].GetString("startTime"),
				EndTime:    records[i].GetString("endTime"),
				Name:       records[i].GetString("name"),
				Location:   records[i].GetString("location"),
				Times:      times,
				Resources:  resources,
				Tags:       tags,
			})
		}

		return c.JSON(200, events)
	})
}
