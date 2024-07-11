package go5paisa

import (
	"encoding/json"
)

// OrderDetail represents details of an order in the OrderBook
type OrderDetail struct {
	AHProcess          string  `json:"AHProcess"`
	AfterHours         string  `json:"AfterHours"`
	AtMarket           string  `json:"AtMarket"`
	BrokerOrderID      int     `json:"BrokerOrderId"`
	BrokerOrderTime    string  `json:"BrokerOrderTime"`
	BuySell            string  `json:"BuySell"`
	DelvIntra          string  `json:"DelvIntra"`
	DisClosedQty       int     `json:"DisClosedQty"`
	Exchange           string  `json:"Exch"`
	ExchOrderID        string  `json:"ExchOrderID"`
	ExchOrderTime      string  `json:"ExchOrderTime"`
	ExchangeType       string  `json:"ExchType"`
	MarketLot          int     `json:"MarketLot"`
	OldorderQty        int     `json:"OldorderQty"`
	OrderRequesterCode string  `json:"OrderRequesterCode"`
	OrderStatus        string  `json:"OrderStatus"`
	OrderValidUpto     string  `json:"OrderValidUpto"`
	OrderValidity      int     `json:"OrderValidity"`
	PendingQty         int     `json:"PendingQty"`
	Qty                int     `json:"Qty"`
	Rate               float32 `json:"Rate"`
	Reason             string  `json:"Reason"`
	RequestType        string  `json:"RequestType"`
	SLTriggerRate      float32 `json:"SLTriggerRate"`
	SLTriggered        string  `json:"SLTriggered"`
	SMOProfitRate      float32 `json:"SMOProfitRate"`
	SMOSLLimitRate     float32 `json:"SMOSLLimitRate"`
	SMOSLTriggerRate   float32 `json:"SMOSLTriggerRate"`
	SMOTrailingSL      float32 `json:"SMOTrailingSL"`
	ScripCode          int     `json:"ScripCode"`
	ScripName          string  `json:"ScripName"`
	TerminalID         int     `json:"TerminalId"`
	TradedQty          int     `json:"TradedQty"`
	WithSL             string  `json:"WithSL"`
}

type orderBookResponseData struct {
	Head interface{} `json:"head"`
	Body OrderBook   `json:"body"`
}

// OrderBook contains details for orders
type OrderBook struct {
	OrderBookDetail []OrderDetail `json:"OrderBookDetail"`
}

func parseOrderBookResponseBody(resBody []byte, obj OrderBook) {
	var body orderBookResponseData
	body.Body = obj
	if err := json.Unmarshal(resBody, &body); err != nil {
		panic(err)
	}
}

// GetOrderBook fetches order book of the user
func (c *Client) GetOrderBook() (OrderBook, error) {
	var orderBook OrderBook

	head := c.buildHeader(orderBookRequestCode)

	payloadBody := genericPayloadBody{
		ClientCode: c.clientCode,
	}
	payload := genericPayload{
		Head: &head,
		Body: payloadBody,
	}

	resBody, err := c.postRequest(payload, orderBookRoute)
	//fmt.Println(string(resBody))
	if err != nil {
		return orderBook, err
	}

	parseResBody(resBody, &orderBook)

	return orderBook, nil
}
