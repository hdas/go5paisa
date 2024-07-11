package go5paisa

// "log"

type TradeBookDetail struct {
	BuySell           string  `json:"BuySell"`
	DelvIntra         string  `json:"DelvIntra"`
	Exch              string  `json:"Exch"`
	ExchOrderID       string  `json:"ExchOrderID"`
	ExchangeTradeTime string  `json:"ExchangeTradeTime"`
	ExchType          string  `json:"ExchType"`
	ExchangeTradeID   string  `json:"ExchangeTradeID"`
	OrgQty            int     `json:"OrgQty"`
	PendingQty        int     `json:"PendingQty"`
	Qty               int     `json:"Qty"`
	Rate              float64 `json:"Rate"`
	RemoteOrderID     string  `json:"RemoteOrderID"`
	ScripCode         int     `json:"ScripCode"`
	ScripName         string  `json:"ScripName"`
	TradeType         string  `json:"TradeType"`
}

type TradeBookResponse struct {
	Message       string            `json:"Message"`
	Status        int               `json:"Status"`
	TradeBookList []TradeBookDetail `json:"TradeBookDetail"`
}

// GetTradeBook fetches Trade book of the user
func (c *Client) GetTradeBook() (TradeBookResponse, error) {
	var tradeBookResponse TradeBookResponse
	head := c.buildHeader(tradeBookRequestCode)
	payloadBody := genericPayloadBody{
		ClientCode: c.clientCode,
	}
	payload := genericPayload{
		Head: &head,
		Body: payloadBody,
	}

	//res, err := c.connection.Post(baseURL+tradeBookRoute, contentType, bytes.NewBuffer(jsonValue))
	resBuf, err := c.postRequest(payload, tradeBookRoute)
	if err != nil {
		return tradeBookResponse, err
	}

	parseResBody(resBuf, &tradeBookResponse)
	return tradeBookResponse, nil
}
