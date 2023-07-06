package go5paisa

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	// "log"
)

// OrderForStatus represents details of an order in the OrderBook
type OrderForStatus struct {
	Exchange      string `json:"Exch"`
	ExchangeType  string `json:"ExchType"`
	ScripCode     int    `json:"ScripCode"`
	RemoteOrderID string `json:"RemoteOrderID"`
}

type OrderResponseList struct {
	OrderList []OrderResponse `json:"OrdStatusResLst"`
}

// OrderResponse contains the order status details
type OrderResponse struct {
	Exchange          string  `json:"Exch"`
	ExchangeType      string  `json:"ExchType"`
	ScripCode         int     `json:"ScripCode"`
	ExchangeOrderID   int     `json:"ExchOrderID"`
	ExchangeOrderTime string  `json:"ExchOrderTime"`
	OrderQty          int     `json:"OrderQty"`
	OrderRate         float32 `json:"OrderRate"`
	PendingQty        int     `json:"PendingQty"`
	Status            string  `json:"Status"`
	Symbol            string  `json:"Symbol"`
	TradedQty         int     `json:"TradedQty"`
}

// OrderStatusReqList contains a order status request list
type OrderStatusReqList struct {
	OrderStatusList []OrderForStatus `json:"OrdStatusReqList"`
}

// GetOrderStatus fetches order book of the user
func (c *Client) GetOrderStatus(orderList OrderStatusReqList) (OrderResponseList, error) {
	var orderStatus OrderResponseList
	//c.appConfig.head.RequestCode =

	head := c.buildHeader(orderStatusRequestCode)
	payloadBody := orderStatusPayloadBody{
		ClientCode: c.clientCode,
		OrdList:    orderList.OrderStatusList,
	}
	payload := orderStatusPayload{
		Head: &head,
		Body: payloadBody,
	}
	jsonValue, _ := json.Marshal(payload)
	res, err := c.connection.Post(baseURL+orderStatusRoute, contentType, bytes.NewBuffer(jsonValue))
	if err != nil {
		return orderStatus, err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return orderStatus, err
	}
	parseResBody(resBody, &orderStatus)
	return orderStatus, nil
}
