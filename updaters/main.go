package updaters

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Stephen10121/planningcenterbackend/email"
	"github.com/Stephen10121/planningcenterbackend/event"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func UpdateEvents(base *pocketbase.PocketBase) bool {
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
		return true
	}

	for i := 0; i < len(events); i++ {
		times, err := json.Marshal(events[i].Times)
		if err != nil {
			return true
		}

		resources, err := json.Marshal(events[i].Resources)
		if err != nil {
			return true
		}

		tags, err := json.Marshal(events[i].Tags)
		if err != nil {
			return true
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
				return true
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

	return false
}

func sendWarning(base *pocketbase.PocketBase, daysLeft int, recipient string) {
	_, err := email.TokenExpireWarning(daysLeft, recipient)
	if err != nil {
		base.Logger().Error(
			"Failed to send warning email!",
			"id", 124,
			"error", err,
		)
	}
}

func absInt(x int) int {
	return absDiffInt(x, 0)
}

func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

func CheckForExpiredTokens(base *pocketbase.PocketBase) bool {
	records, err := base.FindAllRecords("users")

	if err != nil {
		return true
	}

	for i := 0; i < len(records); i++ {
		tokenExpiresStr := records[i].GetString("refreshTokenExpires")
		userEmailStr := records[i].GetString("subscriptionEmail")

		if len(tokenExpiresStr) == 0 || len(userEmailStr) == 0 {
			continue
		}

		date := time.Now()

		format := "2006-01-02 15:04:05Z"
		then, _ := time.Parse(format, tokenExpiresStr)

		diff := date.Sub(then)
		diffDay := int(diff.Hours() / 24)

		if diffDay > -11 {
			daysLeftAbs := absInt(diffDay)
			sendWarning(base, daysLeftAbs, userEmailStr)
			return false
		}
	}

	return false
}
