package event

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func str(num int) string {
	return strconv.Itoa(num)
}

type EventInstance struct {
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
	Relationships struct {
		Event            EventRelationship            `json:"event"`
		EventTimes       EventTimesRelationship       `json:"event_times"`
		ResourceBookings ResourceBookingsRelationship `json:"resource_bookings"`
		Tags             TagsRelationship             `json:"tags"`
	} `json:"relationships"`
	Links map[string]any `json:"links"`
}

type IncludedType struct {
	Type          string         `json:"type"`
	Id            string         `json:"id"`
	Attributes    map[string]any `json:"attributes"`
	Relationships map[string]any `json:"relationships"`
	Links         map[string]any `json:"links"`
}

type NewEventInstancesResponseType struct {
	Links    map[string]string `json:"links"`
	Data     []EventInstance   `json:"data"`
	Included []IncludedType    `json:"included"`
	Meta     map[string]any    `json:"meta"`
}

func GetIncludedStructs(included []IncludedType) ([]EventItself, []EventTime, []ResourceBooking, []Tag) {
	eventsitself := []EventItself{}
	eventtimes := []EventTime{}
	resourceBookings := []ResourceBooking{}
	tags := []Tag{}

	for _, aType := range included {
		switch aType.Type {
		case "Event":
			newEvent, ok := RestructureEvent(aType)
			if !ok {
				continue
			}
			eventsitself = append(eventsitself, newEvent)
		case "EventTime":
			newEventTime, ok := RestructureEventTime(aType)
			if !ok {
				continue
			}
			eventtimes = append(eventtimes, newEventTime)
		case "ResourceBooking":
			newResourceBooking, ok := RestructureResourceBooking(aType)
			if !ok {
				continue
			}
			resourceBookings = append(resourceBookings, newResourceBooking)
		case "Tag":
			newTags, ok := RestructureTag(aType)
			if !ok {
				continue
			}
			tags = append(tags, newTags)
		default:
			continue
		}
	}

	return eventsitself, eventtimes, resourceBookings, tags
}

type Event struct {
	InstanceId string               `json:"instanceId"`
	StartTime  string               `json:"startTime"`
	EndTime    string               `json:"endTime"`
	Name       string               `json:"name"`
	Location   string               `json:"location"`
	Times      []SpecificEventTimes `json:"times"`
	Resources  []Resource           `json:"resources"`
	Tags       []EventTag           `json:"tags"`
}

func EventFetcher(authorizationHeader string) ([]Event, error) {
	year, month, day := time.Now().Add(-72 * time.Hour).Date()
	thirdYear, thirdMonth, thirdDay := time.Now().Add(72 * time.Hour).Date()
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.planningcenteronline.com/calendar/v2/event_instances?include=event%2Cevent_times%2Cresource_bookings%2Ctags&order=starts_at&where[starts_at][gt]="+str(year)+"-"+str(int(month))+"-"+str(day),
		nil,
	)
	req.Header.Add("Authorization", authorizationHeader)

	if err != nil {
		return []Event{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []Event{}, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return []Event{}, err
	}

	responseJson := new(NewEventInstancesResponseType)

	err = json.Unmarshal([]byte(resBody), &responseJson)

	if err != nil {
		return []Event{}, err
	}

	events, eventTimes, resourceBookings, allTags := GetIncludedStructs(responseJson.Included)

	resources := FetchResources(authorizationHeader)

	var fetchedEvents []Event

	for i := 0; i < len(responseJson.Data); i++ {
		date, err := time.Parse(time.RFC3339, responseJson.Data[i].Attributes.StartsAt)

		if err != nil {
			fmt.Println(err)
			continue
		}

		// This will only process the events that are happening in the next 3 days. We can probably afford to remove this limitation.
		if date.Year() <= thirdYear && date.Month() <= thirdMonth && date.Day() <= thirdDay {
			eventTime := ParseEventTimes(responseJson.Data[i].Relationships.EventTimes, eventTimes)
			resources := ParseResourceBookings(responseJson.Data[i].Relationships.ResourceBookings, resourceBookings, resources)
			tags := ParseTags(responseJson.Data[i].Relationships.Tags, allTags)
			eventItself, ok := ParseEventItself(responseJson.Data[i].Relationships.Event, events)
			if !ok {
				fmt.Println("No actuall event found for this event instance:", responseJson.Data[i].Id)
				continue
			}

			fetchedEvents = append(fetchedEvents, Event{
				InstanceId: responseJson.Data[i].Id,
				StartTime:  responseJson.Data[i].Attributes.StartsAt,
				EndTime:    eventTime[len(eventTime)-1].EndTime,
				Name:       eventItself.Name,
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
