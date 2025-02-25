package event

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

type ResourcesFetchResponse struct {
	Data []ResourceJsonType `json:"data"`
}

func FetchResources(authorizationHeader string) []ResourceJsonType {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.planningcenteronline.com/calendar/v2/resources",
		nil,
	)
	req.Header.Add("Authorization", authorizationHeader)

	if err != nil {
		fmt.Println(err)
		return []ResourceJsonType{}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return []ResourceJsonType{}
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return []ResourceJsonType{}
	}

	var responseJson ResourcesFetchResponse

	err = json.Unmarshal([]byte(resBody), &responseJson)

	if err != nil {
		fmt.Println(err)
		return []ResourceJsonType{}
	}

	return responseJson.Data
}
