package event

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Stephen10121/planningcenterbackend/initializers"
)

type EventTimes struct {
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type ResourceType struct {
	Id       string `json:"id"`
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	PathName string `json:"path_name"`
}

type TagsType struct {
	Color string `json:"color"`
	Name  string `json:"name"`
	Id    string `json:"id"`
}

type EventType struct {
	InstanceId string         `json:"instanceId"`
	StartTime  string         `json:"startTime"`
	EndTime    string         `json:"endTime"`
	Name       string         `json:"name"`
	Location   string         `json:"location"`
	Times      []EventTimes   `json:"times"`
	Resources  []ResourceType `json:"resources"`
	Tags       []TagsType     `json:"tags"`
}

func str(num int) string {
	return strconv.Itoa(num)
}

type EventItselfResponseType struct {
	Data struct {
		Type       string `json:"type"`
		Id         string `json:"id"`
		Attributes struct {
			ApprovalStatus       string `json:"approval_status"`
			CreatedAt            string `json:"created_at"`
			Description          any    `json:"description"`
			Featured             bool   `json:"featured"`
			ImageUrl             any    `json:"image_url"`
			Name                 string `json:"name"`
			PercentApproved      int    `json:"percent_approved"`
			PercentRejected      int    `json:"percent_rejected"`
			RegistrationURL      any    `json:"registration_url"`
			Summary              any    `json:"summary"`
			UpdatedAt            string `json:"updated_at"`
			VisibleInChuchCenter bool   `json:"visible_in_church_center"`
		} `json:"attributes"`
		Relationships map[string]any `json:"relationships"`
		Links         map[string]any `json:"links"`
	} `json:"data"`
	Included []any `json:"included"`
	Meta     map[string]any
}

func FetchEventItself(id string) (*EventItselfResponseType, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.planningcenteronline.com/calendar/v2/event_instances/"+id+"/event",
		nil,
	)
	req.Header.Add("Authorization", "Basic "+initializers.Credentials+"==")

	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	responseJson := new(EventItselfResponseType)

	err = json.Unmarshal([]byte(resBody), &responseJson)

	if err != nil {
		return nil, err
	}

	return responseJson, nil
}

type EventTimeResponseType struct {
	Data []struct {
		Type       string `json:"type"`
		Id         string `json:"id"`
		Attributes struct {
			EndsAt                 string `json:"ends_at"`
			Name                   string `json:"name"`
			StartsAt               string `json:"starts_at"`
			VisibleOnKiosks        bool   `json:"visible_on_kiosks"`
			VisibleOnWidgetAndIcal bool   `json:"visible_on_widget_and_ical"`
		} `json:"attributes"`
		Relationships map[string]any `json:"relationships"`
		Links         map[string]any `json:"links"`
	} `json:"data"`
	Links    map[string]any `json:"links"`
	Included []any          `json:"included"`
	Meta     map[string]any
}

func FetchEventTime(id string) ([]EventTimes, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.planningcenteronline.com/calendar/v2/event_instances/"+id+"/event_times",
		nil,
	)
	req.Header.Add("Authorization", "Basic "+initializers.Credentials+"==")

	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	responseJson := new(EventTimeResponseType)

	err = json.Unmarshal([]byte(resBody), &responseJson)

	if err != nil {
		return nil, err
	}

	eventTimes := []EventTimes{}

	for i := 0; i < len(responseJson.Data); i++ {
		eventTimes = append(eventTimes, EventTimes{
			Name:      responseJson.Data[i].Attributes.Name,
			StartTime: responseJson.Data[i].Attributes.StartsAt,
			EndTime:   responseJson.Data[i].Attributes.EndsAt,
		})
	}

	return eventTimes, nil
}

type ResourceBookingsResponseType struct {
	Data []struct {
		Type       string `json:"type"`
		Id         string `json:"id"`
		Attributes struct {
			CreatedAt string `json:"created_at"`
			EndsAt    string `json:"ends_at"`
			StartsAt  string `json:"starts_at"`
			UpdatedAt string `json:"updated_at"`
			Quantity  int    `json:"quantity"`
		} `json:"attributes"`
		Relationships struct {
			Resource struct {
				Data struct {
					Type string `json:"type"`
					Id   string `json:"Id"`
				} `json:"data"`
			} `json:"resource"`
		} `json:"relationships"`
		Links map[string]any `json:"links"`
	} `json:"data"`
	Links    map[string]any `json:"links"`
	Included []any          `json:"included"`
	Meta     map[string]any
}

func FetchResourceBookings(id string, resourcesMap []ResourceJsonType) ([]ResourceType, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.planningcenteronline.com/calendar/v2/event_instances/"+id+"/resource_bookings",
		nil,
	)
	req.Header.Add("Authorization", "Basic "+initializers.Credentials+"==")

	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	responseJson := new(ResourceBookingsResponseType)

	err = json.Unmarshal([]byte(resBody), &responseJson)

	if err != nil {
		return nil, err
	}

	resources := []ResourceType{}

	for a := 0; a < len(responseJson.Data); a++ {
		for b := 0; b < len(resourcesMap); b++ {
			if responseJson.Data[a].Relationships.Resource.Data.Id == resourcesMap[b].Id {
				resources = append(resources, ResourceType{
					Id:       responseJson.Data[a].Relationships.Resource.Data.Id,
					Name:     resourcesMap[b].Attributes.Name,
					PathName: resourcesMap[b].Attributes.PathName,
					Kind:     resourcesMap[b].Attributes.Kind,
				})
			}
		}
	}

	return resources, nil
}

type TagsResponseType struct {
	Data []struct {
		Type       string `json:"type"`
		Id         string `json:"id"`
		Attributes struct {
			CreatedAt            string `json:"created_at"`
			UpdatedAt            string `json:"updated_at"`
			Position             any    `json:"position"`
			Name                 string `json:"name"`
			Color                string `json:"color"`
			ChurchCenterCategory bool   `json:"church_center_category"`
		} `json:"attributes"`
		Links map[string]any `json:"links"`
	} `json:"data"`
	Links    map[string]any `json:"links"`
	Included []any          `json:"included"`
	Meta     map[string]any
}

func FetchTags(id string) ([]TagsType, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.planningcenteronline.com/calendar/v2/event_instances/"+id+"/tags",
		nil,
	)
	req.Header.Add("Authorization", "Basic "+initializers.Credentials+"==")

	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	responseJson := new(TagsResponseType)

	err = json.Unmarshal([]byte(resBody), &responseJson)

	if err != nil {
		return nil, err
	}

	tags := []TagsType{}

	for a := 0; a < len(responseJson.Data); a++ {
		tags = append(tags, TagsType{
			Id:    responseJson.Data[a].Id,
			Color: responseJson.Data[a].Attributes.Color,
			Name:  responseJson.Data[a].Attributes.Name,
		})
	}

	return tags, nil
}

type EventInstancesResponseType struct {
	Links map[string]string `json:"links"`
	Data  []struct {
		Type       string `json:"type"`
		Id         string `json:"id"`
		Attributes struct {
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
		} `json:"attributes"`
		Relationships map[string]any `json:"relationships"`
		Links         map[string]any `json:"links"`
	} `json:"data"`
	Included []any          `json:"included"`
	Meta     map[string]any `json:"meta"`
}

type ResourceJsonType struct {
	Type       string `json:"type"`
	Id         string `json:"id"`
	Attributes struct {
		CreatedAt    string         `json:"created_at"`
		Description  any            `json:"description"`
		ExpiresAt    any            `json:"expires_at"`
		HomeLocation any            `json:"home_location"`
		Image        map[string]any `json:"image"`
		Kind         string         `json:"kind"`
		Name         string         `json:"name"`
		PathName     string         `json:"path_name"`
		Quantity     int            `json:"quantity"`
		SerialNumber any            `json:"serial_number"`
		UpdatedAt    string         `json:"updated_at"`
	} `json:"attributes"`
	Links map[string]any `json:"links"`
}

func FetchEvents() ([]EventType, error) {
	plan, _ := os.ReadFile("./resources.json")
	var resourcesMap []ResourceJsonType
	err := json.Unmarshal(plan, &resourcesMap)

	if err != nil {
		return nil, err
	}

	year, month, day := time.Now().Add(-72 * time.Hour).Date()
	thirdYear, thirdMonth, thirdDay := time.Now().Add(72 * time.Hour).Date()

	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.planningcenteronline.com/calendar/v2/event_instances?where[starts_at][gt]="+str(year)+"-"+str(int(month))+"-"+str(day),
		nil,
	)
	req.Header.Add("Authorization", "Basic "+initializers.Credentials+"==")

	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	responseJson := new(EventInstancesResponseType)

	err = json.Unmarshal([]byte(resBody), &responseJson)

	if err != nil {
		return nil, err
	}

	var fetchedEvents []EventType

	for i := 0; i < len(responseJson.Data); i++ {
		date, err := time.Parse(time.RFC3339, responseJson.Data[i].Attributes.StartsAt)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if date.Year() <= thirdYear && date.Month() <= thirdMonth && date.Day() <= thirdDay {
			eventItself, err := FetchEventItself(responseJson.Data[i].Id)

			if err != nil {
				fmt.Println(err)
				continue
			}

			eventTime, err := FetchEventTime(responseJson.Data[i].Id)

			if err != nil {
				fmt.Println(err)
				continue
			}

			resources, err := FetchResourceBookings(responseJson.Data[i].Id, resourcesMap)

			if err != nil {
				fmt.Println(err)
				continue
			}

			tags, err := FetchTags(responseJson.Data[i].Id)

			if err != nil {
				fmt.Println(err)
				continue
			}

			fetchedEvents = append(fetchedEvents, EventType{
				InstanceId: responseJson.Data[i].Id,
				StartTime:  responseJson.Data[i].Attributes.StartsAt,
				EndTime:    eventTime[len(eventTime)-1].EndTime,
				Name:       eventItself.Data.Attributes.Name,
				Location:   responseJson.Data[i].Attributes.Location,
				Times:      eventTime,
				Resources:  resources,
				Tags:       tags,
			})
		} else {
			continue
		}
	}

	return fetchedEvents, nil
}
