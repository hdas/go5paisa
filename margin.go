package go5paisa

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
)

// Margin represents a margin object returned in the JSON response
type Margin struct {
	ALB             float32 `json:"ALB"`
	Adhoc           float32 `json:"Adhoc"`
	AvailableMargin float32 `json:"AvailableMargin"`
	GHV             float32 `json:"GHV"`
	GHVPer          float32 `json:"GHVPer"`
	GrossMargin     float32 `json:"GrossMargin"`
	Mgn4PendOrd     float32 `json:"Mgn4PendOrd"`
	Mgn4Position    float32 `json:"Mgn4Position"`
	OptionsMtoMLoss float32 `json:"OptionsMtoMLoss"`
	PDHV            float32 `json:"PDHV"`
	Payments        float32 `json:"Payments"`
	Receipts        float32 `json:"Receipts"`
	THV             float32 `json:"THV"`
}

type marginResponseData struct {
	Head interface{}  `json:"head"`
	Body EquityMargin `json:"body"`
}

// EquityMargin contains Margin objects
type EquityMargin struct {
	EquityMargin []Margin `json:"EquityMargin"`
}

func parseMarginResponseBody(resBody []byte, obj EquityMargin) {
	var body marginResponseData
	body.Body = obj
	if err := json.Unmarshal(resBody, &body); err != nil {
		log.Fatal("Error parsing JSON response:", err)
	}
}

// GetMargin fetches margins of the user
func (c *Client) GetMargin() (EquityMargin, error) {
	var margin EquityMargin
	//c.appConfig.head.RequestCode = marginRequestCode
	payloadBody := genericPayloadBody{
		ClientCode: c.clientCode,
	}

	head := c.buildHeader("5PMarginV3")

	payload := genericPayload{
		Head: &head,
		Body: payloadBody,
	}
	jsonValue, _ := json.Marshal(payload)
	res, err := c.connection.Post(baseURL+marginRoute, contentType, bytes.NewBuffer(jsonValue))
	if err != nil {
		return margin, err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return margin, err
	}

	parseResBody(resBody, &margin)
	return margin, nil
}
