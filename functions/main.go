package functions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func CreateWebhook(name string, url string, accessToken string) error {
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
		return errors.New("failed to create the json body that's required to send to the planning center api")
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.planningcenteronline.com/webhooks/v2/subscriptions",
		bytes.NewBuffer(jsonData),
	)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(resBody[:]))

	return nil
}
