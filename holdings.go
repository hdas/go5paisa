package go5paisa

import (
	"encoding/json"
)

// Holding represents a single holding
type Holding struct {
	BseCode         int     `json:"BseCode"`
	CurrentPrice    float32 `json:"CurrentPrice"`
	DPQty           int     `json:"DPQty"`
	Exchange        string  `json:"Exch"`
	ExchangeType    string  `json:"ExchType"`
	Name            string  `json:"FullName"`
	NseCode         int     `json:"NseCode"`
	POASigned       string  `json:"POASigned"`
	PoolQty         int     `json:"PoolQty"`
	Quantity        int     `json:"Quantity"`
	ScripMultiplier int     `json:"ScripMultiplier"`
	Symbol          string  `json:"Symbol"`
}

// Data has all holdings for a user
type responseData struct {
	Head interface{} `json:"head"`
	Body Holdings    `json:"body"`
}

type Holdings struct {
	Data []Holding `json:"Data"`
}

func parsHoldingsResponse(resBody []byte, obj Holdings) {
	var body responseData
	body.Body = obj
	if err := json.Unmarshal(resBody, &body); err != nil {
		panic(err)
	}
}

// GetHoldings fetches holdings of the user
func (c *Client) GetHoldings() (Holdings, error) {
	var holdings Holdings
	payloadBody := genericPayloadBody{
		ClientCode: c.clientCode,
	}

	head := c.buildHeader(holdingsRequestCode)

	payload := genericPayload{
		Head: &head,
		Body: payloadBody,
	}

	resBody, err := c.postRequest(payload, holdingsRoute)
	if err != nil {
		return holdings, err
	}

	parseResBody(resBody, &holdings)
	return holdings, nil
}
