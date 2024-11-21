package event

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Stephen10121/planningcenterbackend/initializers"
)

type EventTimes struct {
	Name      string
	StartTime string
	EndTime   string
}

type ResourceType struct {
	Id       string
	Kind     string
	Name     string
	PathName string
}

type TagsType struct {
	Color string
	Name  string
	Id    string
}

type EventType struct {
	InstanceId string
	StartTime  string
	EndTime    string
	Name       string
	Location   string
	Times      []EventTimes
	Resources  []ResourceType
	Tags       []TagsType
}

func str(num int) string {
	return strconv.Itoa(num)
}

type EventInstanceAttributes struct {
	AllDayEvent                  bool   `json:"all_day_event"`
	ChurchCenterURL              string `json:"church_center_url"`
	CompactRecurrenceDescription string `json:"compact_recurrence_description"`
	CreatedAt                    string `json:"created_at"`
	EndsAt                       string `json:"ends_at"`
	Location                     string `json:"location"`
	PublishedEndsAt              string `json:"published_ends_at"`
	PublishedStartsAt            string `json:"published_starts_at"`
	Recurrence                   string `json:"recurrence"`
	RecurrenceDescription        string `json:"recurrence_description"`
	StartsAt                     string `json:"starts_at"`
	UpdatedAt                    string `json:"updated_at"`
}

type EventInstanceActual struct {
	Type          string                  `json:"type"`
	Id            string                  `json:"id"`
	Attributes    EventInstanceAttributes `json:"attributes"`
	Relationships map[string]any          `json:"relationships"`
	Links         map[string]any          `json:"links"`
}

type EventInstancesResponseType struct {
	Links    map[string]string     `json:"links"`
	Data     []EventInstanceActual `json:"data"`
	Included []any                 `json:"included"`
	Meta     map[string]any        `json:"meta"`
}

func FetchEvents() error {
	year, month, day := time.Now().Add(-72 * time.Hour).Date()

	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.planningcenteronline.com/calendar/v2/event_instances?where[starts_at][gt]="+str(year)+"-"+str(int(month))+"-"+str(day),
		nil,
	)
	req.Header.Add("Authorization", "Basic "+initializers.Credentials+"==")

	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return err
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return err
	}

	responseJson := new(EventInstancesResponseType)

	err = json.Unmarshal([]byte(resBody), &responseJson)

	if err != nil {
		return err
	}

	for i := 0; i < len(responseJson.Data); i++ {
		fmt.Println("Current Data Response: ", responseJson.Data[i].Id, responseJson.Data[i].Attributes.StartsAt)
	}

	return nil
}
