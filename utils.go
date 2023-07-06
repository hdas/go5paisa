package go5paisa

import (
	"encoding/json"
	"log"
)

type body struct {
	ClientCode   string `json:"ClientCode"`
	Message      string `json:"Message"`
	RequestToken string `json:"RequestToken"`
}

type responseBody struct {
	Head interface{} `json:"head"`
	Body interface{} `json:"body"`
}

type AccessTokenBody struct {
	Message     string `json:"Message"`
	AccessToken string `json:"AccessToken"`
}

func parseResBody(resBody []byte, obj interface{}) {
	var body responseBody
	body.Body = obj
	if err := json.Unmarshal(resBody, &body); err != nil {
		log.Fatal("Error parsing JSON response:", err)
	}
}

func parseAccessTokenResponse(resBody []byte, obj interface{}) {
	var body responseBody
	body.Body = obj
	if err := json.Unmarshal(resBody, &body); err != nil {
		log.Fatal("Error parsing JSON response:", err)
	}
}
