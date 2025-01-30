package functions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type WebhookCreatedResp struct {
	Data struct {
		Id         string `json:"id"`
		Attributes struct {
			ApplicationId      int    `json:"application_id"`
			AuthenticitySecret string `json:"authenticity_secret"`
		} `json:"attributes"`
	} `json:"data"`
}

type CreateWebhookReturn struct {
	Id                 int    `json:"id"`
	AuthenticitySecret string `json:"authenticity_secret"`
}

func CreateWebhook(name string, url string, accessToken string) (CreateWebhookReturn, error) {
	data := map[string]any{
		"data": map[string]any{
			"attributes": map[string]any{
				"name":   name,
				"url":    url,
				"active": true,
			},
		},
	}

	// Marshal the JSON data
	jsonData, err := json.Marshal(data)
	if err != nil {
		return CreateWebhookReturn{}, errors.New("failed to create the json body that's required to send to the planning center api")
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.planningcenteronline.com/webhooks/v2/subscriptions",
		bytes.NewBuffer(jsonData),
	)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return CreateWebhookReturn{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return CreateWebhookReturn{}, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return CreateWebhookReturn{}, err
	}

	bodyStruct := WebhookCreatedResp{}
	err = json.Unmarshal(resBody, &bodyStruct)
	if err != nil {
		return CreateWebhookReturn{}, err
	}

	respon := CreateWebhookReturn{
		Id:                 bodyStruct.Data.Attributes.ApplicationId,
		AuthenticitySecret: bodyStruct.Data.Attributes.AuthenticitySecret,
	}

	fmt.Println(respon)
	return respon, nil
}
