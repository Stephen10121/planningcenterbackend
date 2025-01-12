package main

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Stephen10121/planningcenterbackend/event"
	"github.com/Stephen10121/planningcenterbackend/initializers"

	// "github.com/Stephen10121/planningcenterbackend/initializers"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	// "github.com/pocketbase/pocketbase/models"
)

func main() {
	initializers.SetupEnv()
	base := pocketbase.New()

	base.OnServe().BindFunc(func(e *core.ServeEvent) error {
		e.Router.POST("/webhook", func(c *core.RequestEvent) error {
			fmt.Println("Endpoint hit")
			fmt.Println(c.Request.Body)
			return c.JSON(200, map[string]string{
				"msg": "All Good!",
			})
		})
		e.Router.GET("/test", func(c *core.RequestEvent) error {
			event.NewEventFetcher()
			return c.JSON(200, map[string]string{
				"msg": "All Good",
			})
		})
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

		return e.Next()
	})

	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("[server] Fetching events from planning center.")

				events, err := event.FetchEvents()

				if err != nil {
					base.Logger().Error(
						"Failed to fetch data from the planning center api!",
						"id", 123,
						"error", err,
					)
				}

				collection, err := base.FindCollectionByNameOrId("events")
				if err != nil {
					base.Logger().Warn("Create the users collection to save the data fetched from the planning center api!")
					continue
				}

				for i := 0; i < len(events); i++ {
					times, err := json.Marshal(events[i].Times)
					if err != nil {
						continue
					}

					resources, err := json.Marshal(events[i].Resources)
					if err != nil {
						continue
					}

					tags, err := json.Marshal(events[i].Tags)
					if err != nil {
						continue
					}

					existingRecord, err := base.FindRecordById("events", events[i].InstanceId)

					if err != nil {
						record := core.NewRecord(collection)
						record.Set("id", events[i].InstanceId)
						record.Set("startTime", events[i].StartTime)
						record.Set("endTime", events[i].EndTime)
						record.Set("name", events[i].Name)
						record.Set("location", events[i].Location)
						record.Set("times", times)
						record.Set("resources", resources)
						record.Set("tags", tags)

						if err := base.Save(record); err != nil {
							fmt.Println(err)
							continue
						}
					} else {
						existingRecord.Set("startTime", events[i].StartTime)
						existingRecord.Set("endTime", events[i].EndTime)
						existingRecord.Set("name", events[i].Name)
						existingRecord.Set("location", events[i].Location)
						existingRecord.Set("times", times)
						existingRecord.Set("resources", resources)
						existingRecord.Set("tags", tags)

						if err := base.Save(existingRecord); err != nil {
							fmt.Println(err)
							continue
						}
					}
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	defer close(quit)

	if err := base.Start(); err != nil {
		log.Fatal(err)
	}
}
