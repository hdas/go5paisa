package go5paisa

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
)

// Position represents a position object returned in the JSON response
type Position struct {
	BodQty        int32   `json:"BodQty"`
	BookedPL      float32 `json:"BookedPL"`
	BuyAvgRate    float32 `json:"BuyAvgRate"`
	BuyValue      float32 `json:"BuyValue"`
	Exch          string  `json:"Exch"`
	ExchType      string  `json:"ExchType"`
	LTP           float32 `json:"LTP"`
	MTOM          float32 `json:"MTOM"`
	Multiplier    float32 `json:"Multiplier"`
	NetQty        int32   `json:"NetQty"`
	OrderFor      string  `json:"OrderFor"`
	PreviousClose float32 `json:"PreviousClose"`
	ScripCode     int32   `json:"ScripCode"`
	ScripName     string  `json:"ScripName"`
	SellAvgRate   float32 `json:"SellAvgRate"`
	SellQty       int32   `json:"SellQty"`
	SellValue     float32 `json:"SellValue"`
}

type positionResponseData struct {
	Head interface{} `json:"head"`
	Body Positions   `json:"body"`
}

// Positions contains Position objects
type Positions struct {
	PositionDetail []Position `json:"NetPositionDetail"`
}

func parseResponseBody(resBody []byte, obj Positions) {
	var body positionResponseData
	body.Body = obj
	if err := json.Unmarshal(resBody, &body); err != nil {
		log.Fatal("Error parsing JSON response:", err)
	}
}

// GetPositions fetches margins of the user
func (c *Client) GetPositions() (Positions, error) {
	var positions Positions
	// c.appConfig.head.RequestCode =
	head := c.buildHeader(positionsRequestCode)
	payloadBody := genericPayloadBody{
		ClientCode: c.clientCode,
	}
	payload := genericPayload{
		Head: &head,
		Body: payloadBody,
	}
	jsonValue, _ := json.Marshal(payload)
	res, err := c.connection.Post(baseURL+positionsRoute, contentType, bytes.NewBuffer(jsonValue))
	if err != nil {
		return positions, err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return positions, err
	}
	parseResBody(resBody, &positions)
	return positions, nil
}
