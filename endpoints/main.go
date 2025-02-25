package endpoints

import (
	"encoding/json"
	"fmt"
	"io"

	"time"

	"github.com/Stephen10121/planningcenterbackend/event"
	"github.com/Stephen10121/planningcenterbackend/functions"
	"github.com/Stephen10121/planningcenterbackend/initializers"
	"github.com/Stephen10121/planningcenterbackend/token"
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
		newResp, err := event.EventFetcher("Basic " + initializers.Credentials + "==")
		if err != nil {
			fmt.Println(err)
			return c.JSON(500, map[string]string{
				"msg":  "Not all Good",
				"resp": err.Error(),
			})
		}

		newRespJSON, err := json.Marshal(newResp)
		if err != nil {
			fmt.Println(err)
			return c.JSON(500, map[string]string{
				"msg":  "Not all Good",
				"resp": err.Error(),
			})
		}

		return c.JSON(200, map[string]string{
			"msg":     "All Good",
			"newresp": string(newRespJSON),
		})
	})
}

func GetEvents(e *core.ServeEvent, base *pocketbase.PocketBase) {
	e.Router.GET("/events", func(c *core.RequestEvent) error {
		auth := c.Request.Header.Get("Authorization")

		if auth != "Basic "+initializers.Credentials {
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

		events := []event.Event{}

		for i := 0; i < len(records); i++ {
			var times []event.SpecificEventTimes
			var resources []event.Resource
			var tags []event.EventTag

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

			events = append(events, event.Event{
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

type UserSubscribedBody struct {
	Id          string `json:"id"`
	AccessToken string `json:"accessToken"`
}

func UserHasSubscribed(e *core.ServeEvent, base *pocketbase.PocketBase) {
	e.Router.POST("/userSubscribed", func(c *core.RequestEvent) error {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return c.JSON(500, map[string]string{
				"msg": "Failed to read the body!",
			})
		}

		bodyStruct := UserSubscribedBody{}
		err = json.Unmarshal(body, &bodyStruct)
		if err != nil {
			return c.JSON(500, map[string]string{
				"msg": "Failed to read the body!",
			})
		}

		if len(bodyStruct.Id) == 0 || len(bodyStruct.AccessToken) == 0 {
			return c.JSON(400, map[string]string{
				"msg": "Missing data in the body!",
			})
		}

		record, err := base.FindRecordById("users", bodyStruct.Id)
		if err != nil {
			return c.JSON(401, map[string]string{
				"msg": "Unauthorized!",
			})
		}

		if record.GetString("authToken") != bodyStruct.AccessToken {
			return c.JSON(401, map[string]string{
				"msg": "Unauthorized!",
			})
		}

		tok, err := token.RefreshTheAuthToken(record.Id, base)

		if err != nil {
			fmt.Println(err)
			return c.JSON(400, map[string]string{
				"msg": "Failed to refresh the token!",
			})
		}

		webhooksNeeded := []string{
			"calendar.v2.events.event_request.approved",
			"calendar.v2.events.event_request.updated",
			"calendar.v2.events.event_request.created",
			"calendar.v2.events.event_request.destroyed",
		}

		createdWebhooks := []functions.CreateWebhookReturn{}

		for _, webhook := range webhooksNeeded {
			resp, err := functions.CreateWebhook(
				webhook,
				"https://calapi.stephengruzin.dev/tester",
				tok,
			)

			if err != nil {
				fmt.Println(err)
				return c.JSON(400, map[string]string{
					"msg": "Failed to add a webhook!",
				})
			}

			createdWebhooks = append(createdWebhooks, resp)
		}

		createdWebhooksJson, err := json.Marshal(createdWebhooks)
		if err != nil {
			return c.JSON(400, map[string]string{
				"msg": "Failed to jsonify the webhooks!",
			})
		}

		record.Set("webhooks", createdWebhooksJson)

		if err := base.Save(record); err != nil {
			return c.JSON(400, map[string]string{
				"msg": "Failed to save the webhooks in the database.",
			})
		}

		return c.JSON(200, map[string]string{
			"msg": "ok",
		})
	})
}
