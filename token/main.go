package token

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase"
)

type NewAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int    `json:"created_at"`
}

// This function will refresh the auth token. It will return the auth token and store the new auth and refresh tokens in the user db.
func RefreshTheAuthToken(userId string, base *pocketbase.PocketBase) (string, error) {
	record, err := base.FindRecordById("users", userId)
	if err != nil {
		return "", errors.New("user doesnt exist")
	}

	data := map[string]string{
		"client_id":     "a43d0cc3974d32cc7e2181849b393ee8e127eebedf72c1056e54e13c95af568f",
		"client_secret": "pco_app_c479e2122b4a181b535e3a9dec7ad309b668b1ccd56c89371616918f04d75b9c4c6fa23e",
		"refresh_token": record.GetString("refreshToken"),
		"grant_type":    "refresh_token",
	}

	// Marshal the JSON data
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", errors.New("failed to jsonify the request body")
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.planningcenteronline.com/oauth/token",
		bytes.NewBuffer(jsonData),
	)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return "", errors.New("failed to create the http request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.New("failed to send the http request")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("failed to read the http response body")
	}

	bodyStruct := NewAuthTokenResponse{}
	fmt.Println(string(resBody[:]))
	err = json.Unmarshal(resBody, &bodyStruct)
	if err != nil {
		return "", errors.New("failed to put http response body into a struct")
	}

	if len(bodyStruct.AccessToken) == 0 || len(bodyStruct.RefreshToken) == 0 {
		return "", errors.New("failed access the tokens from the struct")
	}

	date := time.Now()
	futureDate := date.AddDate(0, 0, 89)

	record.Set("refreshToken", bodyStruct.RefreshToken)
	record.Set("authToken", bodyStruct.AccessToken)
	record.Set("refreshTokenExpires", futureDate.Format("2006-01-02 15:04:05Z"))

	err = base.Save(record)
	if err != nil {
		return "", err
	}

	fmt.Println("Donw", bodyStruct.AccessToken, bodyStruct.RefreshToken)

	return bodyStruct.AccessToken, nil
}
